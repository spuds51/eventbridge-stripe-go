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
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/webhook"
	"os"
)

const (
	EventBusName = "stripeAppEventBus"
	EventSource = "stripeWebHookHandler.lambda"
	StripeCustomerCreatedEvent = "customer.subscription.created"
	SecretName = "dev/StripeApp/stripe/secret"
	SecretVersion = "AWSCURRENT"
)

type Secret struct {
	StripeWebhookEndpointSecret string `json:"stripe-webhook-endpoint-secret"`
}

type EventDetail struct {
	StripeEvent []string `json:"stripeEvent"`
}

func stripeWebhookSecret(cfg aws.Config) Secret{

	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(SecretName),
		VersionStage: aws.String(SecretVersion), // VersionStage defaults to AWSCURRENT if unspecified
	}
	srv := secretsManagerSession(cfg)
	req := srv.GetSecretValueRequest(input)

	resp, err := req.Send(context.TODO())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error %s\n", err)
	}
	secret := Secret{}
	json.Unmarshal([]byte(*resp.SecretString), &secret)
	return secret

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

func defaultConfig() aws.Config {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		panic("unable to load sdk config, " + err.Error())
	}
	cfg.Region = endpoints.EuWest1RegionID
	return cfg
}

func secretsManagerSession(cfg aws.Config) *secretsmanager.Client{
	return secretsmanager.New(cfg)
}

func eventBridgeSession(cfg aws.Config) *eventbridge.Client{
	return eventbridge.New(cfg)
}

func dispatchEvent(eventType string, customerID string, cfg aws.Config){
	srv := eventBridgeSession(cfg)
	details := createEventDetailJSONString(eventType)

	e := eventRequestEntry(details, customerID)
	req := srv.PutEventsRequest(&e)

	_, err := req.Send(context.TODO())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error %s\n", err)
	}
}


func eventHandler(e stripe.Event, cfg aws.Config){
	switch e.Type {
		case StripeCustomerCreatedEvent:
			var customer stripe.Customer
			err := json.Unmarshal(e.Data.Raw, &customer)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error parsing webhook JSON: %v\n", err)
			}
			fmt.Println("Dispatching customer created event..", customer.ID)
			dispatchEvent(e.Type, e.Data.Object["id"].(string), cfg)
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

func verifyWebhookSig(request events.APIGatewayProxyRequest, secret string) bool{

	sig := request.Headers["Stripe-Signature"]
	body := []byte(request.Body)

	_, err := webhook.ConstructEvent(body, sig,	secret)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error verifying webhook signature: %v\n", err)
		return false
	}
	return true
}

func HandleLambdaEvent(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	cfg := defaultConfig()
	secret := stripeWebhookSecret(cfg)

	sigCheck := verifyWebhookSig(request, secret.StripeWebhookEndpointSecret)
	if !sigCheck {
		return events.APIGatewayProxyResponse{Body: fmt.Sprintf("Update failed"), StatusCode: 401}, nil
	}

	event := unmarshalEvent(request)
	eventHandler(event, cfg)
	return events.APIGatewayProxyResponse{Body: fmt.Sprintf("Update succeeded"), StatusCode: 200}, nil



}

func main() {
	lambda.Start(HandleLambdaEvent)
}


