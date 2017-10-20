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

# seed "b" generates (for Board B):
B_PK=02c9d0d1faca3c852c307b4391af5f353e63a296cded08c1a819f03b7ae768530b
B_SK=2ad6269d4c9b0a7b995d3173f51a9be84361d978c95eaebf8197047992eb7bcf

# ADDRESS
RunMS :8080
sleep 5

# HTTP | CXO | RPC | GUI
RunNode 7410 7412 7414 false
sleep 10

# RPC | Name | Body | Seed
NewBoard 7414 "Board A" "Board with seed 'a'." "a"
NewBoard 7414 "Board B" "Board with seed 'b'." "b"

# Create threads in Board A.
for i in {1..2}
do
    # RPC | BPK | NAME | BODY | USK
    NewThread 7414 ${A_PK} "Thread ${i}" "This is a thread of index ${i}." ${USER_SK}
done

# Export (RPC | BPK | LOC) : Board A -> file.
ExportBoard 7414 ${A_PK} a.json

# Import (RPC | BSK | LOC) : file -> Board B.
ImportBoard 7414 ${B_SK} a.json

# Finish.
sleep 1
pv2 "ALL DONE"
wait