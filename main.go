package main

import (
	"fmt"
	"log"

	"github.com/practice/bank_fcm/fcm"
	"github.com/practice/bank_fcm/model"
)

func main() {
	fmt.Println("Hello World")

	data := map[string]string{
		"msg": "Hello World1",
		"sum": "Happy Day",
	}

	serverKey := ""
	token := "token"

	s := fcm.NewServer(serverKey)

	msg := model.Message{
		Data:             data,
		RegistrationIDs:  []string{token},
		ContentAvailable: true,
		Priority:         model.PriorityHigh,
		Notification: model.Notification{
			Title: "Hello",
			Body:  "World",
		},
	}

	response, err := s.Send(msg)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Status Code   :", response.StatusCode)
	fmt.Println("Success       :", response.Success)
	fmt.Println("Fail          :", response.Fail)
	fmt.Println("Canonical_ids :", response.CanonicalIDs)
	fmt.Println("Topic MsgId   :", response.MsgID)
}
