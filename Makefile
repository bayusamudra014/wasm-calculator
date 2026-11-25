.PHONY: build serve test clean

build:
	cp "$$(go env GOROOT)/lib/wasm/wasm_exec.js" web/
	GOOS=js GOARCH=wasm go build -o web/main.wasm .

serve: build
	python3 -m http.server 8083 -d web

# Runs the unit tests for calc.go on your machine (no browser needed).
test:
	go test -v .

clean:
	rm -f web/main.wasm web/wasm_exec.js
