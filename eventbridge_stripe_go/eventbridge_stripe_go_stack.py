from aws_cdk import (
    aws_iam as iam,
    aws_apigateway as _apigw,
    aws_lambda as _lambda,
    aws_events as events,
    aws_events_targets as targets,
    aws_dynamodb as ddb,
    aws_secretsmanager as sm,
    core
)

class EventbridgeStripeGoStack(core.Stack):

    def __init__(self, scope: core.Construct, id: str, **kwargs) -> None:
        super().__init__(scope, id, **kwargs)

        table = ddb.Table(
            self, 'StripeSampleCustomers',
            partition_key={'name': 'CustomerID', 'type': ddb.AttributeType.STRING}
        )

        bus = events.EventBus(self, 'stripeAppEventBus', event_bus_name='stripeAppEventBus')

        lambda_role_for_go = iam.Role(self,
                                   "Role",role_name='stripeAppRole',
                                   assumed_by=iam.ServicePrincipal("lambda.amazonaws.com"),
                                   managed_policies=[iam.ManagedPolicy.from_aws_managed_policy_name("service-role/AWSLambdaBasicExecutionRole"),
                                                     iam.ManagedPolicy.from_aws_managed_policy_name("AmazonEventBridgeFullAccess"),
                                                     iam.ManagedPolicy.from_aws_managed_policy_name("SecretsManagerReadWrite")]
                                   )

        customer_created_handler = _lambda.Function(self, "createStripeCustomerHandler",
                                                  runtime=_lambda.Runtime.GO_1_X,
                                                  code=_lambda.Code.asset('lambda/stripe-create-customer'),
                                                  handler='createCustomerHandler',
                                                  timeout=core.Duration.seconds(8),
                                                  role=lambda_role_for_go,
                                                  environment={
                                                      'CUSTOMER_TABLE_NAME': table.table_name,
                                                  }
                                                  )
        table.grant_read_write_data(customer_created_handler)

        go_lambda = _lambda.Function(self, "stripeWebhookEventHandler",
                                     runtime=_lambda.Runtime.GO_1_X,
                                     code=_lambda.Code.asset('lambda/stripe-webhook-handler'),
                                     handler='stripeWebhookHandler',
                                     timeout=core.Duration.seconds(8),
                                     role=lambda_role_for_go,
                                     )

        _apigw.LambdaRestApi(self, "stripeWebhookAPI", handler = go_lambda)

        customer_created_handler.add_permission("createStripeCustomerHandlerPermission",
                                              principal=iam.ServicePrincipal("events.amazonaws.com"),
                                              action='lambda:InvokeFunction',
                                              source_arn=go_lambda.function_arn
                                              )

        go_lambda.add_permission("stripeWebhookHandlerPermission",
                                 principal=iam.ServicePrincipal("lambda.amazonaws.com"),
                                 action='lambda:InvokeFunction',
                                 source_arn=customer_created_handler.function_arn
                                 )

        event = events.Rule(self, 'stripeWebhookEventRule',
                            rule_name='stripeWebhookEventRule',
                            enabled=True,
                            event_bus=bus,
                            description='all success events are caught here and logged centrally',
                            event_pattern=events.EventPattern(
                                detail = {"stripeEvent": ["customer.subscription.created"]},
                                source = ["stripeWebHookHandler.lambda"]
                            ))

        event.add_target(targets.LambdaFunction(customer_created_handler))
