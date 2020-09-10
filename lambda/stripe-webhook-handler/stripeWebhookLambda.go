package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/webhook"
	"os"
)


const (
	EventBusName = "stripeAppEventBus"
	EventSource = "stripeWebHookHandler.lambda"
	StripeCustomerCreatedEvent = "customer.created"
)

type EventDetail struct {
	StripeEvent []string `json:"stripeEvent"`
}

func verifyWebhookSig(request events.APIGatewayProxyRequest) bool{

	_, err := webhook.ConstructEvent([]byte(request.Body), request.Headers["Stripe-Signature"],
		"some sig")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error verifying webhook signature: %v\n", err)
		return false
	}
	fmt.Print("Succeeded verifying webhook sig", request.Headers["Stripe-Signature"])
	return true

}

func eventRequestEntry(details []byte, customerID string) eventbridge.PutEventsInput{
	return eventbridge.PutEventsInput{Entries: []eventbridge.PutEventsRequestEntry{
		{
			EventBusName: aws.String(EventBusName),
			Detail:       aws.String(string(details)),
			DetailType:   aws.String(customerID),
			Source:       aws.String(EventSource),
		}},
	}
}

func createEventDetailJSONString(t  string) []byte{
	eventDetail := EventDetail{StripeEvent: []string{t}}
	detail, err := json.Marshal(eventDetail)
	if err != nil {
		// do something else
	}
	return detail
}

func eventBridgeSession() *eventbridge.Client{
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		panic("unable to load SDK config, " + err.Error())
	}
	cfg.Region = endpoints.EuWest1RegionID
	return eventbridge.New(cfg)
}

func dispatchEvent(eventType string, customerID string){
	srv := eventBridgeSession()

	details := createEventDetailJSONString(eventType)

	e := eventRequestEntry(details, customerID)

	req := srv.PutEventsRequest(&e)

	_, err := req.Send(context.TODO())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error %s\n", err)
	}
}


func eventHandler(e stripe.Event){
	switch e.Type {
		case StripeCustomerCreatedEvent:
			var customer stripe.Customer
			err := json.Unmarshal(e.Data.Raw, &customer)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error parsing webhook JSON: %v\n", err)
			}
			fmt.Println("Dispatching customer created event..", customer.ID)
			dispatchEvent(e.Type, e.Data.Object["id"].(string))
		default:
			fmt.Fprintf(os.Stderr, "Unexpected event type: %s\n", e.Type)
		}

}

func unmarshalEvent(request events.APIGatewayProxyRequest) stripe.Event{
	event := stripe.Event{}
	if err := json.Unmarshal([]byte(request.Body), &event); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse webhook body json: %v\n", err.Error())
	}
	return event
}

func HandleLambdaEvent(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	verifyWebhookSig(request)

	event := unmarshalEvent(request)
	eventHandler(event)

	return events.APIGatewayProxyResponse{Body: fmt.Sprintf("%s is %d years old!", "Colin", 1000), StatusCode: 200}, nil
}

func main() {
	lambda.Start(HandleLambdaEvent)
}


