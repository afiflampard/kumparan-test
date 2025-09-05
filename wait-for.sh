#!/bin/bash
set -e

# wait-for.sh usage:
# ./wait-for.sh

echo "Waiting for Postgres..."
until pg_isready -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER"; do
  echo "Postgres is unavailable - sleeping"
  sleep 2
done
echo "Postgres is up!"

echo "Waiting for Elasticsearch..."
until curl -s "$ELASTIC_URL/_cluster/health" | grep '"status":"green"'; do
  echo "Elasticsearch is unavailable - sleeping"
  sleep 5
done
echo "Elasticsearch is up!"

exec "$@"
