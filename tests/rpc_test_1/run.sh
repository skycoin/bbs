#!/usr/bin/env bash

a_bbsgui=6490
b_bbsgui=6480

# Run commands.
echo "[ RUNNING COMMANDS ]"
echo "> ADDING A BOARD ON BBS NODE 'A' ..."
curl \
    -X POST \
    -F "seed=a" \
    -F "name=Board A" \
    -F "description=Board on BBS Node A with seed 'a'." \
    http://127.0.0.1:$a_bbsgui/api/new_board \
    | ydump
sleep 1
echo "> ADDING A THREAD TO THE BOARD FROM BBS NODE 'B' ..."
curl \
    -X POST \
    -F "board=032ffee44b9554cd3350ee16760688b2fb9d0faae7f3534917ff07e971eb36fd6b" \
    http://127.0.0.1:$b_bbsgui/api/subscribe \
    | ydump
sleep 1
echo ''
curl \
    -X POST \
    -F "board=032ffee44b9554cd3350ee16760688b2fb9d0faae7f3534917ff07e971eb36fd6b" \
    -F "name=Thread Added From Remote" \
    -F "description=Yeah, you know it!" \
    http://127.0.0.1:$b_bbsgui/api/new_thread \
    | ydump
sleep 1
echo "> LISTING THREADS FROM BBS NODE 'A' ..."
curl \
    -X POST \
    -F "board=032ffee44b9554cd3350ee16760688b2fb9d0faae7f3534917ff07e971eb36fd6b" \
    http://127.0.0.1:$b_bbsgui/api/get_threads \
    | ydump
sleep 1

# Cleanup.

wait
echo "[ CLEANING UP ]"
rm cli cxod main *.bak *.json