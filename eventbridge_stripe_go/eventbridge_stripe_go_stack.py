from aws_cdk import (
    aws_iam as iam,
    aws_apigateway as _apigw,
    aws_lambda as _lambda,
    aws_events as events,
    aws_events_targets as targets,
    aws_dynamodb as ddb,
    core
)

class EventbridgeStripeGoStack(core.Stack):

    def __init__(self, scope: core.Construct, id: str, **kwargs) -> None:
        super().__init__(scope, id, **kwargs)

        table = ddb.Table(
            self, 'StripeCustomers',
            partition_key={'name': 'customerID', 'type': ddb.AttributeType.STRING}
        )

        bus = events.EventBus(self, 'appEventBus', event_bus_name='appEventBus')

        lambdaRoleForGo = iam.Role(self,
                                   "Role",role_name='hgcomRole',
                                   assumed_by=iam.ServicePrincipal("lambda.amazonaws.com"),
                                   managed_policies=[iam.ManagedPolicy.from_aws_managed_policy_name("service-role/AWSLambdaBasicExecutionRole"),
                                                     iam.ManagedPolicy.from_aws_managed_policy_name("AmazonEventBridgeFullAccess")]
                                   )

        customerCreatedHandler = _lambda.Function(self, "createCustomer",
                                                  runtime=_lambda.Runtime.PYTHON_3_8,
                                                  handler="success.handler",
                                                  code=_lambda.Code.asset('lambda'),
                                                  timeout=core.Duration.seconds(8),
                                                  role=lambdaRoleForGo,
                                                  environment={
                                                      'HG_TABLE_NAME': table.table_name,
                                                  }
                                                  )
        table.grant_read_write_data(customerCreatedHandler)

        go_lambda = _lambda.Function(self, "stripeWebhookHandler1",
                                     runtime=_lambda.Runtime.GO_1_X,
                                     code=_lambda.Code.asset('lambdago/stripeWebhookMod'),
                                     handler='stripeWebhookHandler',
                                     timeout=core.Duration.seconds(8),
                                     role=lambdaRoleForGo
                                     )

        _apigw.LambdaRestApi(self, "stripeWebhook", handler = go_lambda)

        customerCreatedHandler.add_permission("succesLambdaPolicy",
                                              principal=iam.ServicePrincipal("events.amazonaws.com"),
                                              action='lambda:InvokeFunction',
                                              source_arn=go_lambda.function_arn
                                              )

        go_lambda.add_permission("succesLambdaPolicy",
                                 principal=iam.ServicePrincipal("lambda.amazonaws.com"),
                                 action='lambda:InvokeFunction',
                                 source_arn=customerCreatedHandler.function_arn
                                 )

        # eventObj = {"stripeEvent": ["customer.created"]}

        event = events.Rule(self, 'successWebHookRule',
                            rule_name='successWebHookRule',
                            enabled=True,
                            event_bus=bus,
                            description='all success events are caught here and logged centrally',
                            event_pattern=events.EventPattern(
                                detail = {"stripeEvent": ["customer.created"]},
                                source = ["stripeHandler.lambda"]
                            ))

        event.add_target(targets.LambdaFunction(customerCreatedHandler))
