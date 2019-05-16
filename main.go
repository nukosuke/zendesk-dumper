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
		outputFile = fmt.Sprintf("triggers-%d.json", time.Now().Unix())
	}

	fmt.Printf(`===== Zendesk Account Information =====
account: %s
mail:    %s
`, account, email)

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
	jbytes, err := json.MarshalIndent(triggers, "", "  ")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = ioutil.WriteFile(outputFile, jbytes, 0644)
	if err != nil {
		fmt.Printf("Failed to write file: %v", err)
		os.Exit(1)
	}

	fmt.Println("Triggers have been dumped to", outputFile)
}
