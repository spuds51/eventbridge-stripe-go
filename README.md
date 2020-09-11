
# eventbridge-stripe-go 

`eventbridge-stripe-go` Builds a serverless architecture to handle Stripe Webhook events 

## Description

Uses AWS CDK to deploy:
* API Gateway endpoint which will be the target for Stripe Webhook events
* GO Lambda function for handling API Gateway request and creation of new customer in DynamoDB
* DynamoDB table for holding new customers 
* Eventbridge event bus for routing events 


## Setup 

Configure Python virtual env
```
python -m venv .env/
source .env/bin/activate
pip install -r requirements.txt
```

Build Lambda handlers
```
cd lambda/stripe-create-customer/
GOOS=linux go build -o createCustomerHandler github.com/cdugga/eventbridge-stripe-go/createCustomer

cd lambda/stripe-create-customer/
GOOS=linux go build -o stripeWebhookHandler github.com/cdugga/eventbridge-stripe-go/stripeWehbookHandler
```

Deploy CDK stack
```
cdk synth
cdk deploy
```
