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
    http://127.0.0.1:$a_bbsgui/api/get_threads \
    | ydump

sleep 1

echo "> ADDING A FEW POSTS FROM BBS NODE 'B' ..."

curl \
    -X POST \
    -F "board=032ffee44b9554cd3350ee16760688b2fb9d0faae7f3534917ff07e971eb36fd6b" \
    -F "thread=8d26f218cb37fadb931fb081808037c6241d3f3b5958d1175642264e4757d1f6" \
    -F "title=Post 1" \
    -F "body=This is post 1 added from B." \
    http://127.0.0.1:$b_bbsgui/api/new_post \
    | ydump

sleep 1

curl \
    -X POST \
    -F "board=032ffee44b9554cd3350ee16760688b2fb9d0faae7f3534917ff07e971eb36fd6b" \
    -F "thread=8d26f218cb37fadb931fb081808037c6241d3f3b5958d1175642264e4757d1f6" \
    -F "title=Post 2" \
    -F "body=This is post 2 added from B." \
    http://127.0.0.1:$b_bbsgui/api/new_post \
    | ydump

sleep 1

curl \
    -X POST \
    -F "board=032ffee44b9554cd3350ee16760688b2fb9d0faae7f3534917ff07e971eb36fd6b" \
    -F "thread=8d26f218cb37fadb931fb081808037c6241d3f3b5958d1175642264e4757d1f6" \
    -F "title=Post 3" \
    -F "body=This is post 3 added from B." \
    http://127.0.0.1:$b_bbsgui/api/new_post \
    | ydump

sleep 1

echo "> OBTAIN THREADPAGE FROM BBS NODE 'A' ..."

curl \
    -X POST \
    -F "board=032ffee44b9554cd3350ee16760688b2fb9d0faae7f3534917ff07e971eb36fd6b" \
    -F "thread=8d26f218cb37fadb931fb081808037c6241d3f3b5958d1175642264e4757d1f6" \
    http://127.0.0.1:$a_bbsgui/api/get_threadpage \
    | ydump

sleep 1

echo "> ADDING A BOARD TO BBS NODE 'B' ..."

curl \
    -X POST \
    -F "seed=b" \
    -F "name=Board B" \
    -F "description=Board on BBS Node B with seed 'b'." \
    http://127.0.0.1:$b_bbsgui/api/new_board \
    | ydump

sleep 1

echo "[ TESTING IMPORT THREAD ]"

echo "> IMPORT THREAD FROM 'A' to 'B' ..."

curl \
    -X POST \
    -F "from_board=032ffee44b9554cd3350ee16760688b2fb9d0faae7f3534917ff07e971eb36fd6b" \
    -F "to_board=02c9d0d1faca3c852c307b4391af5f353e63a296cded08c1a819f03b7ae768530b" \
    -F "thread=8d26f218cb37fadb931fb081808037c6241d3f3b5958d1175642264e4757d1f6" \
    http://127.0.0.1:$b_bbsgui/api/import_thread \
    | ydump

sleep 1

