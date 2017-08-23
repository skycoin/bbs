#!/usr/bin/env bash

source "include/include.sh"

RunNode 5100 5200 5300

pv2 "SLEEP 10s"
sleep 10

Login 5100 "User5100"
sleep 1

Logout 5100

wait