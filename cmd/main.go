package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"

	"transya/config"
	"transya/lib/kafkaProduser"
	yacloudtranslate "transya/lib/yacloudTranslate"
)

type PreparedData struct {
	RequestHash    string `json:"requestHash"`
	Language       string `json:"language"`
	TextHash       string `json:"textHash"`
	Text           string `json:"text"`
	TranslatedText string `json:"translatedText"`
	StatusCode     bool   `json:"statusCode"`
	ErrorText      string `json:"errorText"`
}

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	cfg, err := config.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	nc, err := nats.Connect(cfg.NatsUrl)
	if err != nil {
		log.Fatalf("Error connecting to NATS server: %v", err)
	}

	defer nc.Drain()

	producer, errProd := kafkaProduser.NewKafkaProducer(cfg.KafkaUrl, cfg.KafkaTopic)
	if errProd != nil {
		log.Fatal(errProd)
	}
	defer producer.Producer.Close()

	tr := yacloudtranslate.RestYaTranslate{
		FolderId: cfg.FolderId,
		ApiKey:   cfg.ApiKey,
	}

	subject := "notifications"

	_, err = nc.QueueSubscribe(subject, "workers", func(msg *nats.Msg) {

		fmt.Printf("Received message on subject %s: %s\n", msg.Subject, string(msg.Data))

		var dataOut PreparedData
		errJs := json.Unmarshal(msg.Data, &dataOut)
		if errJs != nil {
			log.Printf("error unmarshalling JSON: %v", errJs)
		}

		result, err := tr.Translate(yacloudtranslate.TranslateRequest{
			TargetLanguageCode: dataOut.Language,
			Texts:              []string{dataOut.Text},
		})
		if err != nil {
			_ = fmt.Sprintf("Error: %s", err.Error())
		}
		dataOut.TranslatedText = result.Translations[0].Text

		jsonData, errJson := json.Marshal(dataOut)
		if errJson != nil {
			log.Printf("error marshalling JSON: %v", errJson)
		}

		errPub := producer.Produce(string(jsonData), dataOut.RequestHash, dataOut.TextHash)
		if errPub != nil {
			fmt.Println("Error:", errPub)
		}
		fmt.Println(result.Translations[0].Text)
	})
	if err != nil {
		log.Fatalf("Error subscribing to subject %s: %v", subject, err)
	}

	fmt.Printf("Subscribed to subject \"%s\". Waiting for messages...\n", subject)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	fmt.Println("\nReceived interrupt signal, draining connection and exiting.")

}
