package main

import (
	"context"
	"fmt"
	"os"

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

	fmt.Printf("account: %s", account)
	fmt.Printf("email: %s", email)

	// Set credentials from environment variables
	client.SetSubdomain(account)
	client.SetCredential(zendesk.NewAPITokenCredential(email, token))

	// Result data
	triggers := []zendesk.Trigger{}

	// Dump resources
	pageNum := 1
	for {
		triggersInPage, page, err := client.GetTriggers(context.Background(), &zendesk.TriggerListOptions{
			PageOptions: zendesk.PageOptions{
				Page: pageNum,
			},
		})

		if err != nil {
			fmt.Errorf("[E] %v", err)
			os.Exit(1)
		}

		triggers = append(triggers, triggersInPage...)

		if !page.HasNext() {
			break
		}

		pageNum = pageNum + 1
	}

	// Output
	fmt.Println(triggers)
}
