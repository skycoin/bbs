#!/usr/bin/env bash

# Runs some nodes and hosts boards on those nodes.
# These boards are capable of remote content submission.

source "include/include.sh"

# Run some nodes (HTTP | SUB | CXO | GUI).

RunNode 5410 5411 5412 false
RunNode 7410 7411 7412 true

# Wait for nodes to start running (assuming 10s is enough).

pv2 "SLEEP 10s"
sleep 10

# Login.

NewUser 5410 user1
NewUser 7410 user2

Login 5410 user1
Login 7410 user2

# Host some boards on the nodes (HTTP | SEED | SUB).

NewBoard 5410 a 5411
sleep 1

NewBoard 5410 b 5411
sleep 1

NewBoard 7410 c 7411
sleep 1

# Connect and subscribe.

NewConnection 7410 "[::]:5412"
sleep 1

NewSubscription 7410 "032ffee44b9554cd3350ee16760688b2fb9d0faae7f3534917ff07e971eb36fd6b"
sleep 1

NewSubscription 7410 "02c9d0d1faca3c852c307b4391af5f353e63a296cded08c1a819f03b7ae768530b"
sleep 1

# Add some threads.

for i in {1..9}
do
NewTestThread 5410 "032ffee44b9554cd3350ee16760688b2fb9d0faae7f3534917ff07e971eb36fd6b" ${i} &
done
NewTestThread 5410 "032ffee44b9554cd3350ee16760688b2fb9d0faae7f3534917ff07e971eb36fd6b" 10

for i in {1..9}
do
NewTestThread 5410 "02c9d0d1faca3c852c307b4391af5f353e63a296cded08c1a819f03b7ae768530b" ${i} &
done
NewTestThread 5410 "02c9d0d1faca3c852c307b4391af5f353e63a296cded08c1a819f03b7ae768530b" 10

# All done.
sleep 1
pv2 "ALL DONE"

wait