
# eventbridge-stripe-go 

`eventbridge-stripe-go` Builds a serverless architecture to handle [Stripe Webhook](https://stripe.com/docs/api/webhook_endpoints) events 

The AWS CDK stack deploys the following serverless components:
* A single API Gateway endpoint used as the target for the Stripe Webhook [customer.created](https://stripe.com/docs/api/events/types#event_types-customer.created) webhook event
* GO Lambda function for handling API Gateway request and creation of new customer in DynamoDB
* DynamoDB table for new customers 
* Eventbridge event bus


## Setup 

Please refer to [AWS CDK Python workshop](https://cdkworkshop.com/30-python/20-create-project/200-virtualenv.html) for a more detailed set of instructions for initializing and using the python language with the CDK.  

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
