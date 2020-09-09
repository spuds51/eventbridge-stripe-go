import json
import pytest

from aws_cdk import core
from eventbridge-stripe-go.eventbridge_stripe_go_stack import EventbridgeStripeGoStack


def get_template():
    app = core.App()
    EventbridgeStripeGoStack(app, "eventbridge-stripe-go")
    return json.dumps(app.synth().get_stack("eventbridge-stripe-go").template)


def test_sqs_queue_created():
    assert("AWS::SQS::Queue" in get_template())


def test_sns_topic_created():
    assert("AWS::SNS::Topic" in get_template())
