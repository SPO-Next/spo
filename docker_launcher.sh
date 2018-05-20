#!/bin/sh
COMMAND="spo --data-dir $DATA_DIR --wallet-dir $WALLET_DIR $@"

adduser -D -u 10000 spo

if [[ \! -d $DATA_DIR ]]; then
    mkdir -p $DATA_DIR
fi
if [[ \! -d $WALLET_DIR ]]; then
    mkdir -p $WALLET_DIR
fi

chown -R spo:spo $( realpath $DATA_DIR )
chown -R spo:spo $( realpath $WALLET_DIR )

su spo -c "$COMMAND"
