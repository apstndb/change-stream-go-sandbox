package types

import (
	"strconv"
	"time"

	"cloud.google.com/go/spanner"
)

// TimeMills fall back to UnixMicro() if it is not compatible with (*time.Time).Unmarshal()
type TimeMills time.Time

func (t *TimeMills) UnmarshalJSON(b []byte) error {
	if err := ((*time.Time)(t)).UnmarshalJSON(b); err == nil {
		return nil
	}

	i, err := strconv.ParseInt(string(b), 10, 64)
	if err != nil {
		return err
	}
	*t = TimeMills(time.UnixMicro(i).UTC())

	return nil
}

func (t TimeMills) MarshalJSON() ([]byte, error) {
	return time.Time(t).MarshalJSON()
}

// DataChangeRecord is compatible with spanner.Row and avroio.
type DataChangeRecord struct {
	Metadata *struct {
		// https://github.com/apache/beam/blob/master/sdks/java/io/google-cloud-platform/src/main/java/org/apache/beam/sdk/io/gcp/spanner/changestreams/model/ChangeStreamRecordMetadata.java
		ChangeStreamRecordMetadata *struct {
			NumberOfRecordsRead     int64      `json:"numberOfRecordsRead"`
			PartitionCreatedAt      TimeMills  `json:"partitionCreatedAt"`
			PartitionEndTimestamp   TimeMills  `json:"partitionEndTimestamp"`
			PartitionRunningAt      *TimeMills `json:"partitionRunningAt"`
			PartitionScheduledAt    *TimeMills `json:"partitionScheduledAt"`
			PartitionStartTimestamp TimeMills  `json:"partitionStartTimestamp"`
			PartitionToken          string     `json:"partitionToken"`
			QueryStartedAt          TimeMills  `json:"queryStartedAt"`
			RecordReadAt            TimeMills  `json:"recordReadAt"`
			RecordStreamEndedAt     TimeMills  `json:"recordStreamEndedAt"`
			RecordStreamStartedAt   TimeMills  `json:"recordStreamStartedAt"`
			RecordTimestamp         TimeMills  `json:"recordTimestamp"`
			TotalStreamTimeMillis   int64      `json:"totalStreamTimeMillis"`
		} `json:"com.google.cloud.teleport.v2.ChangeStreamRecordMetadata,omitempty"`
	} `spanner:"_" json:"metadata,omitempty"`
	CommitTimestamp                      TimeMills `spanner:"commit_timestamp" json:"commitTimestamp"`
	RecordSequence                       string    `spanner:"record_sequence" json:"recordSequence"`
	ServerTransactionId                  string    `spanner:"server_transaction_id" json:"serverTransactionId"`
	IsLastRecordInTransactionInPartition bool      `spanner:"is_last_record_in_transaction_in_partition" json:"isLastRecordInTransactionInPartition"`
	TableName                            string    `spanner:"table_name" json:"tableName"`
	ValueCaptureType                     string    `spanner:"value_capture_type" json:"valueCaptureType"`
	ColumnTypes                          []*struct {
		Name            string           `spanner:"name" json:"name"`
		Type            spanner.NullJSON `spanner:"type" json:"Type"`
		IsPrimaryKey    bool             `spanner:"is_primary_key" json:"isPrimaryKey"`
		OrdinalPosition int64            `spanner:"ordinal_position" json:"ordinalPosition"`
	} `spanner:"column_types" json:"rowType"`
	Mods []*struct {
		Keys      spanner.NullJSON `spanner:"keys" json:"keysJson"`
		NewValues spanner.NullJSON `spanner:"new_values" json:"newValuesJson"`
		OldValues spanner.NullJSON `spanner:"old_values" json:"oldValuesJson"`
	} `spanner:"mods" json:"mods"`
	ModType                         string `spanner:"mod_type" json:"modType"`
	NumberOfRecordsInTransaction    int64  `spanner:"number_of_records_in_transaction" json:"numberOfRecordsInTransaction"`
	NumberOfPartitionsInTransaction int64  `spanner:"number_of_partitions_in_transaction" json:"numberOfPartitionsInTransaction"`
}

type HeartbeatRecord struct {
	Timestamp TimeMills `spanner:"timestamp"`
}

type ChildPartitionsRecord struct {
	RecordSequence  string    `spanner:"record_sequence"`
	StartTimestamp  TimeMills `spanner:"start_timestamp"`
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
