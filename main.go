package main

import (
	"fmt"
	"log"

	"creator-errors/domain"
	"creator-errors/infrai"
)

func main() {
	client, err := infrai.New()
	if err != nil {
		log.Fatal(err)
	}
	job := domain.Job{AssetID: "asset-42", SubscriberID: "sub-7", Operation: "digital_asset_delivery"}
	status, err := domain.Process(job, func(payload map[string]any) error {
		_, err := client.Capture(payload)
		return err
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(status)
}
