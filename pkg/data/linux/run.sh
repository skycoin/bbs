#!/usr/bin/env bash

WHERE="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

./bbsnode \
    -web-gui-dir="${WHERE}/static/dist" \
    -enforced-messenger-addresses="messenger.skycoin.net:8080" \
    -enforced-subscriptions="03588a2c8085e37ece47aec50e1e856e70f893f7f802cb4f92d52c81c4c3212742" \
    -open-browser=true