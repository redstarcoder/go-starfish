default:
	go fmt
	GOOS=js GOARCH=wasm go build -o starfish.wasm

# If running with standard *><>, this will always fail, currently.
tiny:
	go fmt
	GOOS=js GOARCH=wasm tinygo build -o starfish.wasm
