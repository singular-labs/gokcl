package kcl


import (
    "bufio"
    "bytes"
    "encoding/json"
    "fmt"
    "os"
    "io"
)

type Consumer struct {
    InputStream io.Reader
    OutputStream io.Writer
    Processor RecordProcessor
    Checkpointer *Checkpointer
}


func NewConsumer(processor RecordProcessor) *Consumer {
    consumer := &Consumer{
        InputStream:  os.Stdin,
        OutputStream: os.Stdout,
        Processor:    processor,
        Checkpointer: &Checkpointer{
            OutputStream: os.Stdout,
        },
    }

    consumer.Checkpointer.Consumer = consumer
    return consumer
}

func (consumer *Consumer) Run() {
    for {
        // read next daemon request
        msg := consumer.readAction()
        if msg.Len() == 0 {
            break
        }

        consumer.handleAction(msg)
    }
}

func (consumer *Consumer) readAction() (buffer bytes.Buffer) {
    bio := bufio.NewReader(consumer.InputStream)
    for {
        linePart, hasMoreInLine, err := bio.ReadLine()
        if err != nil {
            panic("Unable to read line from stdin " + err.Error())
        }
        buffer.Write(linePart)
        if hasMoreInLine == false {
            break
        }
    }

    return buffer
}

// handleAction reads a request from the KCL MultiLangDaemon and performs the appropriate action.
func (consumer *Consumer) handleAction(msg bytes.Buffer) {
    var action Action
    err := json.Unmarshal(msg.Bytes(), &action)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Could not understand line read from input: %s\n", msg.String())
        os.Exit(1)
    }

    switch action.ActionType {
    case "initialize":
        err = consumer.handleInitializeAction(&msg)
    case "processRecords":
        err = consumer.handleProcessRecordsAction(&msg)
    case "shutdown":
        err = consumer.handleShutdownAction(&msg)
    case "shutdownRequested":
        consumer.handleShutdownRequestedAction()
    default:
        err = fmt.Errorf("unsupported KCL action: %s", action.ActionType)
    }

    if err != nil {
        fmt.Fprintf(os.Stderr, "%s\n", err)
        os.Exit(1)
    }

    // respond with ack
    status := StatusAction{
        Action: Action{"status"},
        ResponseFor: action.ActionType,
    }
    var serializedStatus []byte
    serializedStatus, err = json.Marshal(status)
    if err != nil {
        fmt.Fprintf(os.Stderr, "%s\n", err)
        os.Exit(1)
    }

    fmt.Println(string(serializedStatus))
}

func (consumer *Consumer) handleInitializeAction(buffer *bytes.Buffer) error {
    var action InitializeAction
    err := json.Unmarshal(buffer.Bytes(), &action)
    if err != nil {
        return err
    }

    consumer.Processor.Init(action.ShardID, action.SequenceNumber, action.SubSequenceNumber)

    return nil
}

func (consumer *Consumer) handleProcessRecordsAction(buffer *bytes.Buffer) error {
    var action ProcessRecordsAction
    err := json.Unmarshal(buffer.Bytes(), &action)
    if err != nil {
        return err
    }

    consumer.Processor.ProcessRecords(&action.Records, action.MillisBehindLatest, consumer.Checkpointer)

    return nil
}

func (consumer *Consumer) handleShutdownAction(buffer *bytes.Buffer) error {
    var action ShutdownAction
    err := json.Unmarshal(buffer.Bytes(), &action)
    if err != nil {
        return err
    }

    consumer.Processor.Shutdown(action.Reason, consumer.Checkpointer)

    return nil
}

func (consumer *Consumer) handleShutdownRequestedAction() {
    consumer.Processor.ShutdownRequested(consumer.Checkpointer)
}