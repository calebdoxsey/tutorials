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
  --execute "
CREATE TABLE IF NOT EXISTS book (
  id BIGSERIAL,
  url TEXT NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS book_job_status (
  book_id BIGINT NOT NULL,
  job_type TEXT NOT NULL,
  status TEXT NOT NULL,

  UNIQUE(book_id, job_type)
);

CREATE TABLE IF NOT EXISTS book_download (
  book_id BIGINT NOT NULL UNIQUE,
  path TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS book_stat (
  book_id BIGINT NOT NULL UNIQUE,
  number_of_words INT NOT NULL,
  longest_word TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS book_review (
  book_id BIGINT NOT NULL,
  username TEXT NOT NULL,
  review TEXT NOT NULL
);
" \
  --host=localhost:26257 \
  --insecure

fg %1
