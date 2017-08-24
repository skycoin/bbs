#!/usr/bin/env bash

# Runs some nodes and hosts boards on those nodes.
# These boards are capable of remote content submission.

source "include/include.sh"

# Run some nodes (HTTP | SUB | CXO).

RunNode 5100 5200 5300
RunNode 6100 6200 6300

# Wait for nodes to start running (assuming 10s is enough).

pv2 "SLEEP 10s"
sleep 10

# Host some boards on the nodes (HTTP | SEED | SUB).

NewBoard 5100 a 5200
sleep 1

NewBoard 5100 b 5200
sleep 1

NewBoard 6100 c 6200
sleep 1

# Give instructions on how to connect/subscribe.

pv2 "INSTRUCTIONS FOR CONNECTION AND SUBSCRIPTION"
pv " - Connect to '[::]:5300'"
pv " - Subscribe to '032ffee44b9554cd3350ee16760688b2fb9d0faae7f3534917ff07e971eb36fd6b'"
pv " - Subscribe to '02c9d0d1faca3c852c307b4391af5f353e63a296cded08c1a819f03b7ae768530b'"
pv " - Connect to '[::]:6300'"
pv " - Subscribe to '035a630a621aa3483f87cb288438982d7ba8524302ed6f293f667e6d8c9fa369a7'"

wait