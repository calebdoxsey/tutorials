#!/usr/bin/env bash
set -eumo pipefail

function run-book() {
  exec go run github.com/cespare/reflex -s -r '\.go$' -- sh -c "cd $1 && go run ."
}

function run-cockroach() {
  echo "starting cockroach"
  cockroach start \
    --insecure \
    --store=/tmp/cockroach-data \
    --log-dir=/tmp/cockroach-logs \
    --listen-addr=localhost \
    --advertise-addr=localhost &

  function kill-on-exit() {
    echo "killing cockroach"
    kill %1
  }
  trap kill-on-exit EXIT SIGINT SIGTERM

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

  wait %1
}

function run-jaeger() {
  exec jaeger-all-in-one
}

function run-redis() {
  redis-server --save "" --appendonly no &

  function kill-on-exit() {
    echo "killing redis"
    kill %1
  }
  trap kill-on-exit EXIT SIGINT SIGTERM

  while ! redis-cli PING; do sleep 1; done

  echo "creating consumer group"
  redis-cli XGROUP CREATE jobs workers $ MKSTREAM

  wait %1
}

case "$1" in
book-*)
  run-book "$1"
  ;;
"cockroach")
  run-cockroach
  ;;
"jaeger")
  run-jaeger
  ;;
"redis")
  run-redis
  ;;
*)
  echo "unknown command: $1"
  exit 1
  ;;
esac
