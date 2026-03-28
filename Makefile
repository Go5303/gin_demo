.PHONY: build run clean tidy

APP_NAME := woda_oa
CMD_DIR := ./cmd

build:
	go build -o $(APP_NAME) $(CMD_DIR)/main.go

run:
	go run $(CMD_DIR)/main.go

run-config:
	go run $(CMD_DIR)/main.go -config config/config.yaml

clean:
	rm -f $(APP_NAME)

tidy:
	go mod tidy
