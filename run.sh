#!/bin/bash

docker network create training-network

NETWORK=training-network docker-compose up -d dependencies

NETWORK=training-network docker-compose up -d gateway