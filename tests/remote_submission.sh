#!/usr/bin/env bash

# Runs some nodes and hosts boards on those nodes.
# These boards are capable of remote content submission.

source "include/include.sh"

# Run a messenger server (ADDRESS).

DISCOVERY_ADDRESS=[::]:8008
RunMS ${DISCOVERY_ADDRESS}

# Wait for messenger server to start (assuming 5s is enough).

pv2 "SLEEP 5s"
sleep 5

# Run some nodes (HTTP | CXO | RPC | GUI).

RunNode ${DISCOVERY_ADDRESS} 5410 5412 5414 false
RunNode ${DISCOVERY_ADDRESS} 7410 7412 7414 true

# Wait for nodes to start running (assuming 10s is enough).

pv2 "SLEEP 15s"
sleep 15

# Host some boards on the nodes (HTTP | SEED | SUB).

NewBoard 5414 "Board A" "Board generated with 'a'." "a"
NewBoard 5414 "Board B" "Board generated with 'b'." "b"
NewBoard 7414 "Board C" "Board generated with 'c'." "c"

# Subscribe.

NewSubscription 7414 "032ffee44b9554cd3350ee16760688b2fb9d0faae7f3534917ff07e971eb36fd6b"
NewSubscription 7414 "02c9d0d1faca3c852c307b4391af5f353e63a296cded08c1a819f03b7ae768530b"

# Add some threads.

for i in {1..3}
do
    NewThread 7414 "032ffee44b9554cd3350ee16760688b2fb9d0faae7f3534917ff07e971eb36fd6b" \
        "Test Thread ${i}" "A description of thread ${i}." \
        "8705518acec973239f704aa1bdbf7f5300f006682d8f6b435976e49c8b62aab0" 0
done

for i in {1..3}
do
    NewThread 7414 "02c9d0d1faca3c852c307b4391af5f353e63a296cded08c1a819f03b7ae768530b" \
        "Test Thread ${i}" "A description of thread ${i}." \
        "8705518acec973239f704aa1bdbf7f5300f006682d8f6b435976e49c8b62aab0" 0
done

# All done.
sleep 2
pv2 "ALL DONE"

wait