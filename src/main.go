package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/certificate"
)

type Message struct {
	APNS *APNSNotification
}
type MessageTemplate struct {
	NotificationID string `json:"notification_id"`
	DeviceID       string `json:"device_id"`
	Title          string `json:"title"`
	Body           string `json:"body"`
	Silent         bool
	Data           map[string]interface{}
}

type APNSNotification struct {
	Aps        *Aps
	CustomData map[string]interface{}
}

func (apnsNotify *APNSNotification) MarshalJSON() ([]byte, error) {
	m := map[string]interface{}{"aps": apnsNotify.Aps.standardFields()}
	for k, v := range apnsNotify.CustomData {
		m[k] = v
	}
	return json.Marshal(m)
}

type Aps struct {
	Alert            *ApsAlert `json:"alert,omitempty"`
	Badge            *int      `json:"badge,omitempty"`
	Sound            string    `json:"sound,omitempty"`
	ContentAvailable bool      `json:"content-available,omitempty"`
	MutableContent   bool      `json:"mutable-content,omitempty"`
}

type ApsAlert struct {
	Title string `json:"title,omitempty"`
	Body  string `json:"body,omitempty"`
}

func (a *Aps) standardFields() map[string]interface{} {
	m := make(map[string]interface{})
	if a.Alert != nil {
		m["alert"] = a.Alert
	}
	if a.ContentAvailable {
		m["content-available"] = 1
	}
	if a.MutableContent {
		m["mutable-content"] = 1
	}
	if a.Badge != nil {
		m["badge"] = *a.Badge
	}
	if a.Sound != "" {
		m["sound"] = a.Sound
	} else if a.Sound == "" {
		m["sound"] = "default"
	}
	return m
}

func (m *MessageTemplate) Messaging() *Message {
	message := &Message{
		APNS: &APNSNotification{
			Aps: &Aps{
				Alert: &ApsAlert{
					Title: m.Title,
					Body:  m.Body,
				},
				ContentAvailable: m.Silent,
			},
			CustomData: m.Data,
		},
	}

	err := validateMessage(message)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return message
}

func validateMessage(message *Message) error {
	if message.APNS == nil {
		return errors.New("Need message config")
	}

	return nil
}

func NewLocationNotify(c string) *MessageTemplate {
	content := c
	msg := &MessageTemplate{
		Title: "t",
		Body:  c,
		Data: map[string]interface{}{
			"content": content,
		},
	}

	return msg
}

func main() {

	notification := &apns2.Notification{}
	notification.DeviceToken = "token"
	notification.Topic = os.Getenv("TOPIC")
	notification.Payload = NewLocationNotify("t").Messaging()

	client := GetIOSCertification()
	res, err := client.Push(notification)

	if err != nil {
		log.Fatal("Error:", err)
	}

	fmt.Printf("%v %v %v\n", res.StatusCode, res.ApnsID, res.Reason)
}

// GetIOSCertification get APNS config from env
func GetIOSCertification() (client *apns2.Client) {
	apnsConfig := struct {
		PasswdOfPem     string `json:"passwd_of_pem"`
		PushCertificate string `json:"push_certificate"`
		ApnsEnvironment string `json:"apns_environment"`
	}{}
	sDec, err := base64.StdEncoding.DecodeString(os.Getenv("APNS"))
	if err != nil {
		log.Fatal("[GetIOSCertification] cannot decode apns config: ", err)
	}
	err = json.Unmarshal(sDec, &apnsConfig)
	if err != nil {
		log.Fatal("[GetIOSCertification] need apns config: ", err)
	}

	passwd := apnsConfig.PasswdOfPem
	pushCertificate := apnsConfig.PushCertificate

	cert, err := certificate.FromPemBytes([]byte(pushCertificate), passwd)
	if err != nil {
		log.Fatal("Certification Error:", err)
	}

	switch apnsConfig.ApnsEnvironment {
	case "PRODUCTION":
		client = apns2.NewClient(cert).Production()
	case "DEVELOPMENT":
		client = apns2.NewClient(cert).Development()
	}

	return client
}
