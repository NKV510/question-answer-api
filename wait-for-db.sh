#!/bin/sh
# wait-for-db.sh

set -e

host="$1"
port="$2"
shift 2
cmd="$@"

echo "🔍 Ожидание подключения к PostgreSQL на $host:$port..."

# Ожидаем доступность порта
until nc -z "$host" "$port"; do
  >&2 echo "⏳ PostgreSQL недоступен, ждем..."
  sleep 1
done

# Проверяем, что БД готова принимать подключения
until PGPASSWORD=$DB_PASSWORD psql -h "$host" -p "$port" -U "$DB_USER" -d "$DB_NAME" -c '\q' 2>/dev/null; do
  >&2 echo "⏳ PostgreSQL загружается, ждем..."
  sleep 1
done

>&2 echo "✅ PostgreSQL доступен, запускаем команду: $cmd"
exec $cmd