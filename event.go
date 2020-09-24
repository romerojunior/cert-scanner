package main

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

// An `event` represents all collected data and metadata about the sender.
type event struct {
	Sources    sources    `json:"sources"`
	SenderInfo senderInfo `json:"senderInfo"`
}

type senderInfo struct {
	Hostname string `json:"hostname"`
}

type sources struct {
	Certificates []certificate `json:"certificates"`
}

// stringfiedJSON returns the `event` struct as a `JSON` formatted `string` and
// and error in case of failed marshalling.
func (e event) stringfiedJSON() (string, error) {
	j, err := json.Marshal(e)
	if err != nil {
		return "", errors.New("failed to marshall data structure")
	}
	return string(j), err
}

func sendEvent(e event, cfg conf) (err error) {
	sess := session.Must(
		session.NewSessionWithOptions(
			session.Options{
				Profile:           cfg.Destination.Aws.Profile,
				SharedConfigState: session.SharedConfigEnable,
			}),
	)

	svc := sqs.New(sess)
	body, _ := e.stringfiedJSON()

	msg := &sqs.SendMessageInput{
		DelaySeconds: aws.Int64(10),
		MessageAttributes: map[string]*sqs.MessageAttributeValue{
			"ClientVersion": &sqs.MessageAttributeValue{
				DataType:    aws.String("String"),
				StringValue: aws.String(version),
			},
		},
		MessageBody: aws.String(body),
		QueueUrl:    &cfg.Destination.Aws.URL,
	}

	resp, err := svc.SendMessage(msg)

	if err != nil {
		log.Print("error sending event")
		log.Panic(err)
	}

	log.Printf("event sent successfully (id: %v)", *resp.MessageId)
	return err
}
