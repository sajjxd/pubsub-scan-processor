package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"cloud.google.com/go/pubsub"
	"github.com/sajjxd/pubsub-scan-processor/pkg/processing"
	"github.com/sajjxd/pubsub-scan-processor/pkg/storage"
	"google.golang.org/api/option"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	projectID := getEnv("GOOGLE_CLOUD_PROJECT", "test-project")
	subID := getEnv("SUBSCRIPTION_ID", "scan-sub")
	dbPath := getEnv("DATABASE_PATH", "data/scans.db")
	emulatorHost := os.Getenv("PUBSUB_EMULATOR_HOST")

	repo, err := storage.NewRepository(dbPath)
	if err != nil {
		log.Fatalf("Storage init failed: %v", err)
	}
	defer repo.Close()

	var opts []option.ClientOption
	if emulatorHost != "" {
		opts = append(opts, option.WithEndpoint(emulatorHost), option.WithoutAuthentication())
	}
	client, err := pubsub.NewClient(ctx, projectID, opts...)
	if err != nil {
		log.Fatalf("Pub/Sub client creation failed: %v", err)
	}
	defer client.Close()

	handler := processing.NewMessageHandler(repo)
	sub := client.Subscription(subID)

	go func() {
		log.Printf("Listening on subscription [%s]", subID)
		if err := sub.Receive(ctx, handler.HandleMessage); err != nil {
			log.Printf("Pub/Sub Receive error: %v", err)
		}
	}()

	handleShutdown(cancel)
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func handleShutdown(cancel context.CancelFunc) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	log.Println("Shutdown signal received, terminating...")
	cancel()
}
