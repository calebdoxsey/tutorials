#!/bin/bash
set -eumo pipefail

echo "starting cockroach"
cockroach start \
    --insecure \
    --store=/tmp/cockroach-data \
    --log-dir=/tmp/cockroach-logs \
    --listen-addr=localhost \
    --advertise-addr=localhost &

function kill-on-exit() {
  kill %1 || true
}
trap kill-on-exit EXIT

while ! echo exit | nc localhost 26257; do sleep 1; done

echo "creating SQL tables"
cockroach sql \
  --execute '
CREATE TABLE IF NOT EXISTS book (
  id BIGSERIAL,
  title TEXT,
  url TEXT
);
' \
  --host=localhost:26257 \
  --insecure

fg %1
