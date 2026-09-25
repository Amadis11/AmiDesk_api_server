#!/bin/sh

if [ ! -f /usr/local/bin/amidesk-api-server ]; then
    ln -s /app/amidesk-api-server /usr/local/bin/amidesk-api-server
fi

mkdir -p /app/data

if [ ! -f /app/data/server.yaml ]; then
    cp /app/server.yaml /app/data/server.yaml
fi

cd "$RUSTDESK_API_CONFIG_DIR"

#if [ ! -f /app/server.db ]; then # This is not good if one wants to upgrade instance
/app/amidesk-api-server sync
#fi

if [ ! -f /app/data/.init.lock ] && [ -n "$ADMIN_USER" ] && [ -n "$ADMIN_PASS" ]; then
    /app/amidesk-api-server user add $ADMIN_USER $ADMIN_PASS --admin
    touch /app/data/.init.lock
fi

/app/amidesk-api-server start
