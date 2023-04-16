#!/usr/bin/env sh

export QUICK_CONFIGFILE=config-dev-local.toml
export EXCLUDED_FOLDERS="/vendor/|/docs"

echo [GO-CHECKS] Run golangci-lint
golangci-lint run -v

echo [GO-CHECKS] Run go fmt
go fmt $(go list ./... | grep -vE "${EXCLUDED_FOLDERS}")

echo [GO-CHECKS] Run go vet
go vet $(go list ./... | grep -vE "${EXCLUDED_FOLDERS}")

echo [GO-CHECKS] Run go test
go test -race $(go list ./... | grep -vE "${EXCLUDED_FOLDERS}")
