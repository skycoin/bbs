#!/usr/bin/env bash

PORT_A=6490
PORT_A_CXO=8998
PORT_B=6480

echo "[ RUNNING COMMANDS ]"

echo "> ADD FILLED BOARD TO NODE 'A' ..."

curl \
    -X POST \
    -F "seed=a" \
    -F "threads=5" \
    -F "min_posts=1" \
    -F "max_posts=3" \
    -sS http://127.0.0.1:$PORT_A/api/tests/new_filled_board \
    | jq

sleep 1

echo "> ADD EMPTY BOARD TO NODE 'B' ..."

curl \
    -X POST \
    -F "seed=b" \
    -F "name=Board B" \
    -F "description=Board to test sync." \
    -sS http://127.0.0.1:$PORT_B/api/new_board \
    | jq

sleep 1

echo "> IMPORT THREAD FROM 'A' TO 'B' ..."

echo "   - (subscribing to board 032ffee44b9554cd3350ee16760688b2fb9d0faae7f3534917ff07e971eb36fd6b)"

curl \
    -X POST \
    -F "address=127.0.0.1:${PORT_A_CXO}" \
    -F "board=032ffee44b9554cd3350ee16760688b2fb9d0faae7f3534917ff07e971eb36fd6b" \
    -sS http://127.0.0.1:$PORT_B/api/subscribe \
    | jq

echo "   - (waiting 10 seconds after subscription)"

sleep 10

curl \
    -X POST \
    -F "from_board=032ffee44b9554cd3350ee16760688b2fb9d0faae7f3534917ff07e971eb36fd6b" \
    -F "to_board=02c9d0d1faca3c852c307b4391af5f353e63a296cded08c1a819f03b7ae768530b" \
    -F "thread=c0ae8d23fcf299393ee6df2d507d93c0d14487cd36d9b813fd02297d411cd865" \
	-sS http://127.0.0.1:$PORT_B/api/import_thread \
    | jq

sleep 1

echo "> ADD SOME POSTS ..."

echo "   - (from Node B)"

curl \
    -X POST \
    -F "board=032ffee44b9554cd3350ee16760688b2fb9d0faae7f3534917ff07e971eb36fd6b" \
    -F "thread=c0ae8d23fcf299393ee6df2d507d93c0d14487cd36d9b813fd02297d411cd865" \
    -F "title=Added Post 1" \
    -F "body=This is a manually added post from Node B." \
	-sS http://127.0.0.1:$PORT_B/api/new_post \
    | jq

sleep 1

echo "   - (from Node A)"

curl \
    -X POST \
    -F "board=032ffee44b9554cd3350ee16760688b2fb9d0faae7f3534917ff07e971eb36fd6b" \
    -F "thread=c0ae8d23fcf299393ee6df2d507d93c0d14487cd36d9b813fd02297d411cd865" \
    -F "title=Added Post 2" \
    -F "body=This is a manually added post from Node A." \
	-sS http://127.0.0.1:$PORT_A/api/new_post \
    | jq

sleep 1

echo "> WAIT A WHILE FOR SYNC ..."

sleep 10

echo "> SHOW IMPORTED THREADPAGE (FROM B) ..."

curl \
    -X POST \
    -F "board=02c9d0d1faca3c852c307b4391af5f353e63a296cded08c1a819f03b7ae768530b" \
    -F "thread=c0ae8d23fcf299393ee6df2d507d93c0d14487cd36d9b813fd02297d411cd865" \
	-sS http://127.0.0.1:$PORT_B/api/get_threadpage \
    | jq

sleep 1

echo "> SHOW IMPORTED THREADPAGE (FROM A) ..."

curl \
    -X POST \
    -F "board=032ffee44b9554cd3350ee16760688b2fb9d0faae7f3534917ff07e971eb36fd6b" \
    -F "thread=c0ae8d23fcf299393ee6df2d507d93c0d14487cd36d9b813fd02297d411cd865" \
	-sS http://127.0.0.1:$PORT_A/api/get_threadpage \
    | jq

sleep 1

echo "> VOTE FOR A THREAD LOCALLY, THEN VIA RPC ..."

curl \
    -X POST \
    -F "board=032ffee44b9554cd3350ee16760688b2fb9d0faae7f3534917ff07e971eb36fd6b" \
    -F "thread=c0ae8d23fcf299393ee6df2d507d93c0d14487cd36d9b813fd02297d411cd865" \
    -F "mode=1" \
    -sS http://127.0.0.1:$PORT_A/api/add_thread_vote \
    | jq

sleep 1

curl \
    -X POST \
    -F "board=032ffee44b9554cd3350ee16760688b2fb9d0faae7f3534917ff07e971eb36fd6b" \
    -F "thread=c0ae8d23fcf299393ee6df2d507d93c0d14487cd36d9b813fd02297d411cd865" \
    -F "mode=1" \
    -sS http://127.0.0.1:$PORT_B/api/add_thread_vote \
    | jq

sleep 1

echo "> SHOW VOTES FOR THE THREAD ..."

curl \
    -X POST \
    -F "board=032ffee44b9554cd3350ee16760688b2fb9d0faae7f3534917ff07e971eb36fd6b" \
    -F "thread=c0ae8d23fcf299393ee6df2d507d93c0d14487cd36d9b813fd02297d411cd865" \
    -sS http://127.0.0.1:$PORT_A/api/get_thread_votes \
    | jq

curl \
    -X POST \
    -F "board=032ffee44b9554cd3350ee16760688b2fb9d0faae7f3534917ff07e971eb36fd6b" \
    -F "thread=c0ae8d23fcf299393ee6df2d507d93c0d14487cd36d9b813fd02297d411cd865" \
    -sS http://127.0.0.1:$PORT_B/api/get_thread_votes \
    | jq