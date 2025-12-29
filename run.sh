#!/bin/bash
# Принудительно используем порт из main.go
exec go build -o tmp/main cmd/web/main.go && ./tmp/main
exec go run cmd/web/main.go