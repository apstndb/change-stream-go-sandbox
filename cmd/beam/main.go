package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/apache/beam/sdks/v2/go/pkg/beam"
	"github.com/apache/beam/sdks/v2/go/pkg/beam/io/avroio"
	"github.com/apache/beam/sdks/v2/go/pkg/beam/x/beamx"
	"github.com/apache/beam/sdks/v2/go/pkg/beam/x/debug"
	"github.com/apstndb/change-stream-go-sandbox/types"
)

var (
	input  = flag.String("input", "./twitter.avro", "input avro file")
	output = flag.String("output", "./output.avro", "output avro file")
)

// Doc type used to unmarshal avro json data
type Doc struct {
	Stamp int64  `json:"timestamp"`
	Tweet string `json:"tweet"`
	User  string `json:"username"`
}

// Note that the schema is only required for Writing avro.
// not Reading.
const schema = `{
	"type": "record",
	"name": "tweet",
	"namespace": "twitter",
	"fields": [
		{ "name": "timestamp", "type": "double" },
		{ "name": "tweet", "type": "string" },
		{ "name": "username", "type": "string" }
	]
}`

func main() {
	fmt.Println("partitionEndTimestamp", time.UnixMicro(253402300799999999))
	fmt.Println("queryStartedAt", time.UnixMicro(-62135596800000000))
	flag.Parse()
	beam.Init()

	p := beam.NewPipeline()
	s := p.Root()

	/*
		// read rows and return JSON string PCollection - PCollection<string>
		rows := avroio.Read(s, *input, reflect.TypeOf(""))
		debug.Print(s, rows)
	*/

	// read rows and return Doc Type PCollection - PCollection<Doc>
	docs := avroio.Read(s, *input, reflect.TypeOf(types.DataChangeRecord{}))
	jsons := beam.ParDo(s, func(record types.DataChangeRecord) (string, error) {
		b, err := json.Marshal(record)
		return string(b), err
	}, docs)

	debug.Print(s, jsons)

	if false {
		// update all values with a single user and tweet.
		format := beam.ParDo(s, func(d Doc, emit func(string)) {
			d.User = "daidokoro"
			d.Tweet = "I was here......"

			b, _ := json.Marshal(d)
			emit(string(b))
		}, docs)

		debug.Print(s, format)

		// write output
		avroio.Write(s, *output, schema, format)
	}

	if err := beamx.Run(context.Background(), p); err != nil {
		log.Fatalf("Failed to execute job: %v", err)
	}
}
