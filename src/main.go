package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/certificate"
)

type Alert struct {
	Title string `json:"title,omitempty"`
	Body  string `json:"body,omitempty"`
}

type Aps struct {
	Alert            Alert  `json:"alert,omitempty"`
	Sound            string `json:"sound,omitempty"`
	Badge            int    `json:"badge,omitempty"`
	ContentAvailable int    `json:"content-available,omitempty"`
}

type CustomData struct {
	GroupID  string  `json:"group_id"`
	UserID   string  `json:"user_id"`
	Lat      float64 `json:"lat"`
	Lng      float64 `json:"lng"`
	Accuracy float64 `json:"accuracy"`
}

type Payload struct {
	Aps        *Aps        `json:"aps,omitempty"`
	CustomData *CustomData `json:"data,omitempty"`
}

const a = "13e63da480a4ed41d9498aa9d5c7f2a153e15aff3c1e71c1774061f79939dd27"

func main() {

	err := godotenv.Load("./.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cert, err := certificate.FromPemFile("./config/pushcert.pem", "")
	if err != nil {
		log.Fatal("Cert Error:", err)
	}

	alert := Alert{
		Title: "Day la title",
		Body:  "Day la body",
	}

	aps := Aps{
		Alert: alert,
		Sound: "default",
	}

	payload := Payload{
		Aps: &aps,
	}

	message, err := json.Marshal(payload)

	// log.Fatal(string(message[:len(message)]))

	if err != nil {
		return
	}

	notification := &apns2.Notification{}
	notification.DeviceToken = a
	notification.Topic = os.Getenv("TOPIC")
	notification.Payload = message

	client := apns2.NewClient(cert).Development()
	res, err := client.Push(notification)

	if err != nil {
		log.Fatal("Error:", err)
	}

	fmt.Printf("%v %v %v\n", res.StatusCode, res.ApnsID, res.Reason)
}
