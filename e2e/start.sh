#!/bin/bash

file=$1
if [ "$file" == "" ]
then
    file=compose.yaml
fi

docker compose version
docker compose -f "$file" down
docker compose -f "$file" build extension
docker compose -f "$file" up --exit-code-from testing --remove-orphans
