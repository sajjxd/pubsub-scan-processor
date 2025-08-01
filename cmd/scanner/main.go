package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/sajjxd/pubsub-scan-processor/pkg/types"
)

var (
	services = []string{"HTTP", "SSH", "DNS"}
)

func main() {
	projectId := flag.String("project", "test-project", "GCP Project ID")
	topicId := flag.String("topic", "scan-topic", "GCP PubSub Topic ID")

	ctx := context.Background()

	client, err := pubsub.NewClient(ctx, *projectId)
	if err != nil {
		panic(err)
	}

	topic := client.Topic(*topicId)

	for range time.Tick(time.Second) {

		scan := &types.Scan{
			Ip:        fmt.Sprintf("1.1.1.%d", rand.Intn(255)),
			Port:      uint32(rand.Intn(65535)),
			Service:   services[rand.Intn(len(services))],
			Timestamp: time.Now().Unix(),
		}

		serviceResp := fmt.Sprintf("service response: %d", rand.Intn(100))

		if rand.Intn(2) == 0 {
			scan.DataVersion = types.V1
			scan.Data = &types.V1Data{ResponseBytesUtf8: []byte(serviceResp)}
		} else {
			scan.DataVersion = types.V2
			scan.Data = &types.V2Data{ResponseStr: serviceResp}
		}

		encoded, err := json.Marshal(scan)
		if err != nil {
			panic(err)
		}

		_, err = topic.Publish(ctx, &pubsub.Message{Data: encoded}).Get(ctx)
		if err != nil {
			panic(err)
		}
	}
}
