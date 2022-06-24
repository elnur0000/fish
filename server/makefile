build-wasm:
	tinygo build -o ./bin/main.wasm -target wasm ./cmd/wasm/wasm.go
start:
	go build -gcflags="all=-N -l"  -o ./bin/server ./cmd/server/main.go 
	./bin/server
build-server:
	go build -race -o ./bin/server ./cmd/server/main.go