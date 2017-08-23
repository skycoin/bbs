#!/usr/bin/env bash

source "include/include.sh"

# HTTP | SUB | CXO
RunNode 5100 5200 5300

pv2 "SLEEP 10s"
sleep 10

NewBoard 5100 a 5200
sleep 5

wait