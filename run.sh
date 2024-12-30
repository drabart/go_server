#!/bin/bash

# Load secrets
set -a # automatically export all variables
source .env
set +a

templ generate

go run main.go
