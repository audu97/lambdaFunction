package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"net/http"
)

type slackMessage struct {
	Text string `json:"text"`
}

func sendSlackNotification(webhookURL, message string) {
	slackMessage := slackMessage{Text: "Cpu usage is above 50%" + message}
	slackBody, _ := json.Marshal(slackMessage)
	req, err := http.NewRequest(http.MethodPost, webhookURL, bytes.NewBuffer(slackBody))
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending request: %v\n", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error response from slack: %v\n", resp.StatusCode)
	} else {
		fmt.Printf("Successfully sent Slack notification: %v\n", resp.StatusCode)
	}
}

func handleRequest(ctx context.Context, snsEvent events.SNSEvent) error {
	webhookURL := "https://hooks.slack.com/services/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"
	for _, record := range snsEvent.Records {
		snsRecord := record.SNS
		sendSlackNotification(webhookURL, snsRecord.Message)
	}
	return nil
}
func main() {
	lambda.Start(handleRequest)
}
