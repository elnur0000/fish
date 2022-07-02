build-wasm:
	tinygo build -o ./bin/main.wasm -target wasm ./cmd/wasm/wasm.go
start:
	go build -o ./bin/server ./cmd/server/main.go 
	./bin/server
dev:
	go build -gcflags="all=-N -l"  -o ./bin/server ./cmd/server/main.go 
	./bin/server
dev-watch:
	nodemon --watch './**/*.go' --signal SIGTERM --exec 'make dev'
build-server:
	go build -race -o ./bin/server ./cmd/server/main.go