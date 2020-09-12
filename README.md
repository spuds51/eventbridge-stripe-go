
# eventbridge-stripe-go 

`eventbridge-stripe-go` Builds a serverless architecture to handle Stripe Webhook events 

## Description

Uses AWS CDK to deploy:
* API Gateway endpoint used as the target for a Stripe Webhook event
* GO Lambda function for handling API Gateway request and creation of new customer in DynamoDB
* DynamoDB table for new customers 
* Eventbridge event bus


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

## Tools used
* [Go](https://golang.org/)
* [Python](https://www.python.org/)
* [AWS CDK](https://github.com/aws/aws-cdk)
* [AWS APIGateway](https://aws.amazon.com/api-gateway/)
* [AWS Lambda](https://aws.amazon.com/lambda/)
* [AWS Eventbridge](https://aws.amazon.com/eventbridge/)
* [AWS Dynamodb](https://aws.amazon.com/dynamodb/)
