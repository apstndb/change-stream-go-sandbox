package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"cloud.google.com/go/spanner"
	"github.com/apstndb/change-stream-go-sandbox/types"
	"golang.org/x/sync/errgroup"
)

func main() {
	if err := run(context.Background()); err != nil {
		log.Fatalln(err)
	}
}

func run(ctx context.Context) error {
	project := flag.String("project", "", "")
	instance := flag.String("instance", "", "")
	database := flag.String("database", "", "")
	changeStream := flag.String("change-stream", "", "")
	flag.Parse()
	client, err := spanner.NewClient(ctx, fmt.Sprintf("projects/%s/instances/%s/databases/%s", *project, *instance, *database))
	if err != nil {
		return err
	}
	var errgrp errgroup.Group
	errgrp.Go(
		func() error {
			return WatchPartitions(ctx, &errgrp, client, *changeStream, nil, time.Now(), time.Now().Add(10*time.Minute))
		})
	if err != nil {
		return err
	}
	return errgrp.Wait()
}

func WatchPartitions(ctx context.Context, errgrp *errgroup.Group, client *spanner.Client, changeStream string, partition *string, startTimestamp time.Time, endTimestamp time.Time) error {
	iter := client.Single().Query(ctx, spanner.Statement{SQL: fmt.Sprintf(`
SELECT ChangeRecord
    FROM READ_%s (
        @start_timestamp,
        @end_timestamp,
        @partition_token,
        @heartbeat_milliseconds
    )
`, changeStream),
		Params: map[string]interface{}{
			"start_timestamp":        startTimestamp,
			"end_timestamp":          endTimestamp,
			"partition_token":        partition,
			"heartbeat_milliseconds": 300_000,
		},
	})

	enc := json.NewEncoder(os.Stdout)
	err := iter.Do(func(r *spanner.Row) error {
		var records []*types.ChangeStreamRecord
		err := r.ColumnByName("ChangeRecord", &records)
		if err != nil {
			return err
		}
		err = enc.Encode(records)
		if err != nil {
			return err
		}

		for _, record := range records {
			for _, childPartitionsRecord := range record.ChildPartitionsRecord {
				for _, childPartition := range childPartitionsRecord.ChildPartitions {
					childPartitionToken := childPartition.Token
					childStartTimestamp := childPartitionsRecord.StartTimestamp
					errgrp.Go(func() error {
						return WatchPartitions(ctx, errgrp, client, changeStream, &childPartitionToken, childStartTimestamp, endTimestamp)
					})
				}
			}
		}
		return nil
	})
	return err
}
