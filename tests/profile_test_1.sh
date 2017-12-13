#!/usr/bin/env bash

source "include/include.sh"

# Run a messenger server (ADDRESS).

RunMS [::1]:8008

# User 1, seed: 1.
PK1=02f46d2461e2c3aba0585efb5b2ddb8acb34f38a56865f8a2a3f10272e6de257c1
SK1=12348e8a15fcce27de6c187a5ecace09af622d495474cb3280e5e614f8b789b5

# User 2, seed: 2.
PK2=0284da18e80d5ec08cf54ed9c86bbbce6bbd2838b8c700a373a5886e4de44ce895
SK2=9b43b74f9737e15e36921b418ec5c31ebcc92025133240c9061b2846a88f2e0c

# User 3, seed: 3.
PK3=03f5bcfadd87e625bf62900a7d1ed673ce74034dbfc5d5c624cedd4612a8dc6d1c
SK3=9b40df8b560259c41af5ca5049d3fcd010e925511325fe0c11e362a9d0cada60

# User 4, seed: 4.
PK4=032cbd84230bd56decf62366597abcc26656f32d4c739dc43d2fdb64ff7ce34a75
SK4=9acd790243d539acbf3c22cff302f621e757efb60b1ecd2a7089beb449128c7c

# User 5, seed: 5.
PK5=028f0ad89fa9c9f491f269e41a9902b7f43861526fb4cc91d5d58aec69c07e5ed0
SK5=a14d8c6838fba1b4ad5579dd8d8ba4f8df8a2552bece66547c7ca5d572fec523

# Board, seed: a.
B_SEED="a"
BPK=032ffee44b9554cd3350ee16760688b2fb9d0faae7f3534917ff07e971eb36fd6b
BSK=b4f56cab07ea360c16c22ac241738e923b232138b69089fe0134f81a432ffaff

# HTTP | CXO | RPC | GUI
RunNode [::1]:8008 7410 7412 7414 false
sleep 10

# RPC | Name | Body | Seed
NewBoard 7414 "Board" "This is a board." ${B_SEED}

# RPC | BPK | NAME | BODY | USK | TS
NewThread 7414 ${BPK} "Thread 1" "Created by user 1." ${SK1} 0
NewThread 7414 ${BPK} "Thread 2" "Created by user 2." ${SK2} 0
NewThread 7414 ${BPK} "Thread 3" "Created by user 3." ${SK3} 0
NewThread 7414 ${BPK} "Thread 4" "Created by user 4." ${SK4} 0
NewThread 7414 ${BPK} "Thread 5" "Created by user 5." ${SK5} 0

# RPC | BPK | UPK | VALUE | CSK | TS
VoteUser 7414 ${BPK} ${PK1} +1 "trust" ${SK1} 0
VoteUser 7414 ${BPK} ${PK1} +1 "trust" ${SK2} 0
VoteUser 7414 ${BPK} ${PK1} +1 "trust" ${SK3} 0
VoteUser 7414 ${BPK} ${PK1} +1 "trust" ${SK4} 0
VoteUser 7414 ${BPK} ${PK1} +1 "trust" ${SK5} 0

sleep 1
pv2 "ALL DONE"
wait
