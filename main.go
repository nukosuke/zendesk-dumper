package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/nukosuke/go-zendesk/zendesk"
)

func main() {
	client, err := zendesk.NewClient(nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	account := os.Getenv("ZENDESK_ACCOUNT")
	email := os.Getenv("ZENDESK_EMAIL")
	token := os.Getenv("ZENDESK_TOKEN")

	outputFile := os.Getenv("OUTPUT_FILE")
	if outputFile == "" {
		outputFile = fmt.Sprintf("tickets-%d.json", time.Now().Unix())
	}

	fmt.Printf(`===== Zendesk Account Information =====
account: %s
mail:    %s
`, account, email)

	// Set credentials from environment variables
	client.SetSubdomain(account)
	client.SetCredential(zendesk.NewAPITokenCredential(email, token))

	// Result data
	tickets := []zendesk.Ticket{}

	// Dump resources
	ops := zendesk.NewPaginationOptions()
	it := client.GetTicketsEx(context.Background(), ops)
	for it.HasMore() {
		ticketsInPage, err := it.GetNext()
		if err != nil {
			fmt.Errorf("[E] %v", err)
			os.Exit(1)
		}

		tickets = append(tickets, ticketsInPage...)
	}

	// Output
	jsonBytes, err := json.MarshalIndent(tickets, "", "  ")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = ioutil.WriteFile(outputFile, jsonBytes, 0644)
	if err != nil {
		fmt.Printf("Failed to write file: %v", err)
		os.Exit(1)
	}

	fmt.Println("Tickets have been dumped to", outputFile)
}
