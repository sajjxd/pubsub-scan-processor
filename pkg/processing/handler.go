package processing

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/sajjxd/pubsub-scan-processor/pkg/storage"
	"github.com/sajjxd/pubsub-scan-processor/pkg/types"
)

type MessageHandler struct {
	repo *storage.Repository
}

func NewMessageHandler(repo *storage.Repository) *MessageHandler {
	return &MessageHandler{repo: repo}
}

func (h *MessageHandler) HandleMessage(ctx context.Context, msg *pubsub.Message) {
	record, err := h.parseMessage(msg.Data)
	if err != nil {
		log.Printf("Error parsing message, nacking: %v", err)
		msg.Nack()
		return
	}

	record.LastScanned = time.Now()

	if err := h.repo.UpsertRecord(*record); err != nil {
		log.Printf("Failed to store record, nacking: %v", err)
		msg.Nack()
		return
	}

	msg.Ack()
}

func (h *MessageHandler) parseMessage(data []byte) (*types.ScanRecord, error) {
	var msg types.Scan
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, fmt.Errorf("unmarshal scan: %w", err)
	}

	var response string

	switch msg.DataVersion {
	case types.V1:
		dataMap, ok := msg.Data.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("v1 data is not a valid map")
		}
		base64Str, ok := dataMap["response_bytes_utf8"].(string)
		if !ok {
			return nil, fmt.Errorf("v1 response_bytes_utf8 is not a string")
		}
		decoded, err := base64.StdEncoding.DecodeString(base64Str)
		if err != nil {
			return nil, fmt.Errorf("failed to decode v1 base64: %w", err)
		}
		response = string(decoded)

	case types.V2:
		dataMap, ok := msg.Data.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("v2 data is not a valid map")
		}
		response, ok = dataMap["response_str"].(string)
		if !ok {
			return nil, fmt.Errorf("v2 response_str is not a string")
		}

	default:
		return nil, fmt.Errorf("unsupported data_version: %d", msg.DataVersion)
	}

	return &types.ScanRecord{
		Ip:       msg.Ip,
		Port:     msg.Port,
		Service:  msg.Service,
		Response: response,
	}, nil
}
