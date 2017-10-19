#!/usr/bin/env bash

# Runs some nodes and hosts boards on those nodes.
# These boards are capable of remote content submission.

source "include/include.sh"

# Run a messenger server (ADDRESS).

RunMS :8080

# Wait for messenger server to start (assuming 5s is enough).

pv2 "SLEEP 5s"
sleep 5

# Run some nodes (HTTP | CXO | GUI).

RunNodeWithMessengerSkycoinNet 5410 5412 false
RunNodeWithMessengerSkycoinNet 7410 7412 true

# Wait for nodes to start running (assuming 10s is enough).

pv2 "SLEEP 15s"
sleep 15

# Login.

NewUser 5410 user1
Login 5410 user1

for i in {1..10}
do
    NewUser 7410 "user${i}"
done

# Host some boards on the nodes (HTTP | SEED | SUB).

NewBoard 5410 a
NewBoard 5410 b
NewBoard 7410 c

# Connect and subscribe.

#NewConnection 7410 "[::]:5412"
NewSubscription 7410 "032ffee44b9554cd3350ee16760688b2fb9d0faae7f3534917ff07e971eb36fd6b"
NewSubscription 7410 "02c9d0d1faca3c852c307b4391af5f353e63a296cded08c1a819f03b7ae768530b"

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
sleep 2
pv2 "ALL DONE"

wait