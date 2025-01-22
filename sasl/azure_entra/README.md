# Azure Event Hubs Entra

Provides support for ssl sasl mechanism via [Azure Entra](https://www.microsoft.com/en-us/security/business/identity-access/microsoft-entra-id)

## How to use
You can add this module to your dependencies by running the command below:
```shell
go get github.com/deepshore/kafka-go/sasl/azure_entra
```

To connect to Event Hub with Kafka protocol:
```go
package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/deepshore/kafka-go"
	"github.com/deepshore/kafka-go/sasl/azure_entra"
)

func main() {
	// Create Azure Entra Default Credentials
	cred, err := azidentity.NewDefaultAzureCredential(nil)

	if err != nil {
		fmt.Printf("failed to create Default Azure Credential: %s", err.Error())
		os.Exit(1)
	}

	// Create Azure Entra SASL Mechanism
	entraMechanism := azure_entra.NewMechanism(cred)

	// Reader
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"example.kafka.com:9093"},
		GroupID: "<Arbitrary Consumer Group Id>",
		Topic:   "<Topic>",
		Dialer: &kafka.Dialer{
			SASLMechanism: entraMechanism,
			TLS:           &tls.Config{},
		},
	})

	defer r.Close()

	// Writer
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"<example.kafka.com:9093"},
		Topic:   "<Topic>",
		Dialer: &kafka.Dialer{
			SASLMechanism: entraMechanism,
			TLS:           &tls.Config{},
		},
	})

	defer w.Close()

	err = w.WriteMessages(context.Background(), kafka.Message{
		Value: []byte("test"),
	})

	if err != nil {
		fmt.Printf("failed to write message: %s", err.Error())
		os.Exit(2)
	}

	message, err := r.ReadMessage(context.Background())

	if err != nil {
		fmt.Printf("failed to read message: %s", err.Error())
		os.Exit(3)
	}

	fmt.Printf("received message: %s", string(message.Value))
}

```