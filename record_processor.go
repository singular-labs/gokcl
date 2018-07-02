package kcl

type RecordProcessor interface {
    // Init is called before record processing with the shardId.
    Init(string, string, uint64) error

    // ProcessRecords is called for each batch of records to be processed.
    ProcessRecords(*[]Record, uint64, *Checkpointer) error

    // Shutdown is called before termination.
    Shutdown(string, *Checkpointer) error

    // ShutdownRequested is called before the KCL decides to terminate us
    ShutdownRequested(*Checkpointer) error
}

