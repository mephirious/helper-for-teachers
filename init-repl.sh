#!/bin/bash

# Название контейнера и имя хоста
CONTAINER_NAME="final-mongo"
MONGO_HOST="final-mongo:27017"

# Команда инициализации replica set
INIT_CMD="rs.initiate({_id: 'rs0', members: [{ _id: 0, host: '$MONGO_HOST' }]})"
STATUS_CMD="rs.status()"

echo "🧩 Checking if MongoDB replica set is already initialized..."

docker exec -i "$CONTAINER_NAME" mongosh --quiet --eval "$STATUS_CMD" | grep '"ok" : 1' &> /dev/null

if [ $? -eq 0 ]; then
  echo "✅ Replica set already initialized!"
else
  echo "🚀 Initializing replica set..."
  docker exec -i "$CONTAINER_NAME" mongosh --quiet --eval "$INIT_CMD"
  sleep 3
  echo "🔍 Verifying status..."
  docker exec -i "$CONTAINER_NAME" mongosh --quiet --eval "$STATUS_CMD"
fi
