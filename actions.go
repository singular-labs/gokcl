package kcl

// A general action
type Action struct {
    ActionType          string              `json:"action"`
}

type InitializeAction struct {
    Action
    ShardID             string              `json:"shardId"`
    SequenceNumber      string              `json:"sequenceNumber"`
    SubSequenceNumber   uint64              `json:"subSequenceNumber"`
}

type ProcessRecordsAction struct {
    Action
    Records []Record                        `json:"records"`
    MillisBehindLatest  uint64              `json:"millisBehindLatest"`
}

type ShutdownAction struct {
    Action
    Reason              string              `json:"reason"`
}

type StatusAction struct {
    Action
    ResponseFor         string              `json:"responseFor"`
}

type CheckpointAction struct {
    Action
    SequenceNumber      *string             `json:"sequenceNumber"`
    SubSequenceNumber   *uint64             `json:"subSequenceNumber"`
    Error               *string             `json:"error"`
}

// Record is an individual kinesis record.  Note that the body is always
// base64 encoded.
type Record struct {
    DataB64             string              `json:"data"`
    PartitionKey        string              `json:"partitionKey"`
    SequenceNumber      string              `json:"sequenceNumber"`
    SubSequenceNumber   uint64              `json:"subSequenceNumber"`
    ApproximateArrivalTimestamp struct {
        ID    uint64 `json:"nano"`
        Size  uint64 `json:"epochSecond"`
    }
}