#!/usr/bin/env bash

# (1) Exports: board "a" -> file.
# (2) Imports: file -> board "b".

source "include/include.sh"

# seed "user" generates (for user):
USER_PK=0254020da01e33cbaf2ff01e7cf28de4bb6cea43b357153fad3a50a0e7dd728718
USER_SK=8705518acec973239f704aa1bdbf7f5300f006682d8f6b435976e49c8b62aab0

# seed "a" generates (for Board A):
A_PK=032ffee44b9554cd3350ee16760688b2fb9d0faae7f3534917ff07e971eb36fd6b
A_SK=b4f56cab07ea360c16c22ac241738e923b232138b69089fe0134f81a432ffaff

# ADDRESS
RunMS :8080
sleep 5

# HTTP | CXO | RPC | GUI
RunNode 7410 7412 7414 false
sleep 10

# RPC | Name | Body | Seed
NewBoard 7414 "Board A" "Board with seed 'a'." "a"

# Create threads in Board A.
for i in {1..2}
do
    # RPC | BPK | NAME | BODY | USK | TS
    NewThread 7414 ${A_PK} "Thread ${i}" "This is a thread of index ${i}." ${USER_SK} ${i}
done

# Thread hashes:
T_HASH_1=2cb0aa27b98e4cdf043c5f2d9d5e2b2307e90cb216e0e36e847a7ea0cf6e603d
T_HASH_2=d24712541878aae276edf2db957605e177c24ef843cafc8f39d87ae2ab1485d5

# Create posts in threads.
for i in {1..2}
do
    # RPC | BPK | T_HASH | NAME | BODY | CSK | TS
    NewPost 7414 ${A_PK} ${T_HASH_1} "Post ${i}" "This is a post of index ${i}." ${USER_SK} ${i}
done
for i in {1..2}
do
    # RPC | BPK | T_HASH | NAME | BODY | CSK | TS
    NewPost 7414 ${A_PK} ${T_HASH_2} "Post ${i}" "This is a post of index ${i}." ${USER_SK} ${i}
done

# Export (RPC | BPK | LOC) : Board A -> file.
ExportBoard 7414 ${A_PK} a.json

# Import (RPC | BSK | LOC) : file -> Board B.
ImportBoard 7414 a.json

# Finish.
sleep 1
pv2 "ALL DONE"
wait

rm a.json*