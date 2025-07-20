#!/bin/bash

modules=("auth" "user" "vehicle" "tracking" "notification")
base_path="internal"

mkdir -p cmd config initiator pkg/{token,utils,entities}

for module in "${modules[@]}"; do
  mkdir -p $base_path/$module/domain/{entity,service,repository}
  mkdir -p $base_path/$module/application/{service,command,query}
  mkdir -p $base_path/$module/port/{inbound,outbound}
  mkdir -p $base_path/$module/adapter/inbound/http
  mkdir -p $base_path/$module/adapter/outbound/mongo
  mkdir -p $base_path/$module/tests
done

# middleware as special shared component
mkdir -p internal/middleware

echo "✅ Fleet folder structure created successfully."
