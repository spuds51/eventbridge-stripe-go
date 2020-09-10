
# eventbridge-stripe-go 

`eventbridge-stripe-go` Builds a serverless architecture to handle Stripe Webhook events 

## Description

Uses AWS CDK to provision an API Gateway endpoint which acts as a target for a Stripe 'customer.created' webhook event. When the endpoint is invoked a GO Lambda service validates the 
incoming request signature and then pushes an event to AWS EventBridge. Once EventBridge receives the customer.created event it triggers a rule which forwards the event to another GO Lambda service which extracts the customer ID from the webhook event and writes to a DynamoDB table.  




##Usage
Todo