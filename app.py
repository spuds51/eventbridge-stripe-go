#!/usr/bin/env python3

from aws_cdk import core

from eventbridge_stripe_go.eventbridge_stripe_go_stack import EventbridgeStripeGoStack


app = core.App()
EventbridgeStripeGoStack(app, "eventbridge-stripe-go", env={'region': 'us-west-2'})

app.synth()
