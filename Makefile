.PHONY: run/server
run/server:
	go run server/main.go

.PHONY: run/client
run/client:
	go run client/main.go ./
