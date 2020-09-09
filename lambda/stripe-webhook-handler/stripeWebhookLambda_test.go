package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/aws/aws-sdk-go/aws"
	"testing"
)


const (
	Event = "{\"version\": \"0\",\"id\":\"787878787-7441-1213-d0dc-1212121\",\"detail-type\":\"RegistrationRequest\",\"source\":\"stripeHandler.lambda\",\"account\":\"12124545478\",\"region\":\"us-west-1\",\"time\":\"2017-04-11T20:11:04Z\",\"resources\": [],\"detail\": {\"stripeEvent\": [\"customer.created\"]}}"
	EventPattern = "{\"detail\": {\"stripeEvent\": [\"customer.created\"]},\"source\":[\"stripeHandler.lambda\"]}"
	)

func TestEventBridge(t *testing.T) {

	cfg, _ := external.LoadDefaultAWSConfig()
	svc := eventbridge.New(cfg)

	params := &eventbridge.TestEventPatternInput{
		Event:        aws.String(Event),
		EventPattern: aws.String(EventPattern),
	}
	req := svc.TestEventPatternRequest(params)
	resp, err := req.Send(req.Context())

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(resp)
}
