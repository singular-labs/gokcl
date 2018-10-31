package kcl

type RecordProcessor interface {
    // Init is called before record processing with the shardId.
    Initialize(string, string, uint64) error

    // ProcessRecords is called for each batch of records to be processed.
    ProcessRecords(*[]Record, uint64, *Checkpointer) error

	/**
     * Called when the lease that tied to this record processor has been lost. Once the lease has been lost the record
     * processor can no longer checkpoint.
     * 
     * @param leaseLostInput
     *            access to functions and data related to the loss of the lease. Currently this has no functionality.
     */
    LeaseLost() error

    /**
     * Called when the shard that this record process is handling has been completed. Once a shard has been completed no
     * further records will ever arrive on that shard.
     *
     * When this is called the record processor <b>must</b> call {@link RecordProcessorCheckpointer#checkpoint()},
     * otherwise an exception will be thrown and the all child shards of this shard will not make progress.
     * 
     * @param shardEndedInput
     *            provides access to a checkpointer method for completing processing of the shard.
     */
    ShardEnded(*Checkpointer) error

    /**
     * Called when the Scheduler has been requested to shutdown. This is called while the record processor still holds
     * the lease so checkpointing is possible. Once this method has completed the lease for the record processor is
     * released, and {@link #leaseLost(LeaseLostInput)} will be called at a later time.
     *
     * @param shutdownRequestedInput
     *            provides access to a checkpointer allowing a record processor to checkpoint before the shutdown is
     *            completed.
     */
    ShutdownRequested(*Checkpointer) error
}

