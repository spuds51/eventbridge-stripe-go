# eventbridge-stripe-go 

`eventbridge-stripe-go` Builds a serverless architecture to handle [Stripe Webhook](https://stripe.com/docs/api/webhook_endpoints) events 

The AWS CDK stack deploys the following serverless components:
* A single API Gateway endpoint used as the target for the Stripe Webhook [customer.subscription.created](https://stripe.com/docs/api/events/types#event_types-customer.subscription.created) webhook event
* Go Lambda handler functions:
	* A functions which handles the initial API Gateway request, verifying the Stripe Webhook request signature and subsequently dispatching event to Eventbridge 
	* A second function which is the target of the EventBridge event and writes the newly created customer subscription to a Dynamodb table
* DynamoDB table where new customer subscription d
etails are written
* Eventbridge event bus which orchestrates AWS servcies based on various events

## Prerequisites
The Lambda webhook handler reads a Stripe signing secret from AWS Secrets Manager to verify the incoming request, an assumption is made that this secret has already been created with path: 
```
AWS Secrets Manager
dev/StripeApp/stripe/secret
```
## Deploying the Serverless stack

The fastest way to get started with `eventbridge-stripe-go` is to clone the repo, configure the cdk environment and simply deploy leveraging the already pre-compiled Go Lmabda functions. 

1. Clone repo
```
git clone https://github.com/cdugga/eventbridge-stripe-go.git
```

Configure Python virtual env

2. Configure Python environment
```
python -m venv .env/

source .env/bin/activate

pip install -r requirements.txt
```
Please refer to [AWS CDK Python workshop](https://cdkworkshop.com/30-python/20-create-project/200-virtualenv.html) for a more detailed set of instructions for initializing and using the python language with the CDK.  

3. Deploy stack

```
cdk synth
cdk deploy
```

## (Optional) Compile Lambda functions

Executable versions for each function are included in project source for convenience;
* [createCustomerHandler](https://github.com/cdugga/eventbridge-stripe-go/tree/master/lambda/stripe-create-customer)
* [stripeWebhookHandler](https://github.com/cdugga/eventbridge-stripe-go/tree/master/lambda/stripe-webhook-handler)

Functions can also be compiled from source:

```
cd lambda/stripe-create-customer/
GOOS=linux go build -o createCustomerHandler github.com/cdugga/eventbridge-stripe-go/createCustomer

cd lambda/stripe-create-customer/
GOOS=linux go build -o stripeWebhookHandler github.com/cdugga/eventbridge-stripe-go/stripeWehbookHandler
```
Running programs target the Linux operating system. Use the GOOS runtime value to modify if required; [Go Runtime](GOOS is the running program's operating system target)
See [AWS Lambda deployment package in Go](https://docs.aws.amazon.com/lambda/latest/dg/golang-package.html) for further instructions on how to package a Go lambda function. 

## Stripe signature verification
Stripe can optionally [sign](https://stripe.com/docs/webhooks/signatures) the webhook events it sends to your endpoints by including a signature in each event’s Stripe-Signature header. This check is added in the [stripeWebhookHandler](https://github.com/cdugga/eventbridge-stripe-go/tree/master/lambda/stripe-webhook-handler) function. It relies on a secure token read from AWS Secrets Manager

```
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
```

## Tools used
* [Go 1.15](https://golang.org/)
* [Python 3.8](https://www.python.org/)
* [AWS CDK 1.62.0](https://github.com/aws/aws-cdk)
* [AWS APIGateway](https://aws.amazon.com/api-gateway/)
* [AWS Lambda](https://aws.amazon.com/lambda/)
* [AWS Eventbridge](https://aws.amazon.com/eventbridge/)
* [AWS Dynamodb](https://aws.amazon.com/dynamodb/)
* [AWS SecretsManager](https://aws.amazon.com/secrets-manager/)
