package domain

import "fmt"

type Job struct{ AssetID, SubscriberID, Operation string }

func Fingerprint(job Job) []string {
	return []string{"creator-commerce", job.Operation}
}

func Process(job Job, capture func(map[string]any) error) (string, error) {
	if job.AssetID == "" || job.SubscriberID == "" {
		return "", fmt.Errorf("asset and subscriber are required")
	}
	if job.Operation == "" {
		job.Operation = "content_processing"
	}
	if err := capture(map[string]any{
		"title":       "backend workflow error",
		"message":     "content processing requires review",
		"level":       "error",
		"fingerprint": Fingerprint(job),
		"exception":   map[string]any{"type": "ContentProcessingError", "asset_id": job.AssetID},
		"context":     map[string]any{"asset_id": job.AssetID, "subscriber_id": job.SubscriberID, "operation": job.Operation},
	}); err != nil {
		return "", err
	}
	return "captured", nil
}
