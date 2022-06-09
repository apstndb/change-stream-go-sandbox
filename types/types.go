package types

import (
	"time"

	"cloud.google.com/go/spanner"
)

type DataChangeRecord struct {
	CommitTimestamp                      time.Time `spanner:"commit_timestamp"`
	RecordSequence                       string    `spanner:"record_sequence"`
	ServerTransactionId                  string    `spanner:"server_transaction_id"`
	IsLastRecordInTransactionInPartition bool      `spanner:"is_last_record_in_transaction_in_partition"`
	TableName                            string    `spanner:"table_name"`
	ValueCaptureType                     string    `spanner:"value_capture_type"`
	ColumnTypes                          []*struct {
		Name            string           `spanner:"name"`
		Type            spanner.NullJSON `spanner:"type"`
		IsPrimaryKey    bool             `spanner:"is_primary_key"`
		OrdinalPosition int64            `spanner:"ordinal_position"`
	} `spanner:"column_types"`
	Mods []*struct {
		Keys      spanner.NullJSON `spanner:"keys"`
		NewValues spanner.NullJSON `spanner:"new_values"`
		OldValues spanner.NullJSON `spanner:"old_values"`
	} `spanner:"mods"`
	ModType                         string `spanner:"mod_type"`
	NumberOfRecordsInTransaction    int64  `spanner:"number_of_records_in_transaction"`
	NumberOfPartitionsInTransaction int64  `spanner:"number_of_partitions_in_transaction"`
}

type HeartbeatRecord struct {
	Timestamp time.Time `spanner:"timestamp"`
}

type ChildPartitionsRecord struct {
	RecordSequence  string    `spanner:"record_sequence"`
	StartTimestamp  time.Time `spanner:"start_timestamp"`
	ChildPartitions []*struct {
		Token                 string   `spanner:"token"`
		ParentPartitionTokens []string `spanner:"parent_partition_tokens"`
	} `spanner:"child_partitions"`
}

type ChangeStreamRecord struct {
	DataChangeRecord      []*DataChangeRecord      `spanner:"data_change_record" json:",omitempty"`
	HeartbeatRecord       []*HeartbeatRecord       `spanner:"heartbeat_record" json:",omitempty"`
	ChildPartitionsRecord []*ChildPartitionsRecord `spanner:"child_partitions_record" json:",omitempty"`
}
