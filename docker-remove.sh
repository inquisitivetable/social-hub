#!/bin/bash

# Stop containers and removes containers, networks, volumes, and images created by docker compose up
docker compose down --rmi local

# Delete build cache objects

docker system prune -f
