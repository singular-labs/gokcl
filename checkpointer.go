package kcl

import (
    "io"
    "fmt"
    "os"
    "encoding/json"
)

// Checkpointer marks a consumers progress.
type Checkpointer struct {
    // Output stream to write the checkpoint action to
    OutputStream io.Writer

    // The corresponding consumer
    Consumer *Consumer
}

// CheckpointAll marks all consumed messages as processed.
func (cp *Checkpointer) CheckpointAll() error {
    checkpointMessage := CheckpointAction{
        Action: Action{"checkpoint"},
        SequenceNumber: nil,
        SubSequenceNumber: nil,
    }

    serializedMessage, err := json.Marshal(checkpointMessage)
    if err != nil {
        fmt.Fprintf(os.Stderr, "CheckpointAll Failed: %s\n", err)
        os.Exit(1)
    }

    return cp.doCheckpoint(string(serializedMessage))
}

// CheckpointSeq marks messages up to sequence number as processed.
func (cp *Checkpointer) CheckpointSeq(seqNum string, subSeqNum uint64) error {
    checkpointMessage := CheckpointAction{
        Action: Action{"checkpoint"},
        SequenceNumber: &seqNum,
        SubSequenceNumber: &subSeqNum,
    }

    serializedMessage, err := json.Marshal(checkpointMessage)
    if err != nil {
        fmt.Fprintf(os.Stderr, "CheckpointSeq Failed: %s\n", err)
        os.Exit(1)
    }

    return cp.doCheckpoint(string(serializedMessage))
}

func (cp *Checkpointer) doCheckpoint(msg string) error {
    // send checkpoint req
    fmt.Println(msg)

    // receive checkpoint ack

    ackBuf := cp.Consumer.readAction()
    var ack CheckpointAction
    err := json.Unmarshal(ackBuf.Bytes(), &ack)

    if err != nil {
        fmt.Fprint(os.Stderr, "Got a bad response for a checkpoint action!\n")
        os.Exit(1)
    } else if ack.ActionType != "checkpoint" {
        fmt.Fprintf(os.Stderr, "Received invalid checkpoint ack: %s\n", ack.Action)
        os.Exit(1)
    } else if ack.Error != nil {
        return fmt.Errorf(*ack.Error)
    }

    // success
    return nil
}

