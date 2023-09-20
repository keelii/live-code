BINARY_NAME=bin/live-code

all: build

build:
	go build -o $(BINARY_NAME)

build-js:
	esbuild ../deno-toolkit/src/client/src/editor/CodeMirrorEditor.ts \
      --bundle --minify \
      --outfile=static/codemirror-editor.js \
      --global-name=LiveCode

clean:
	go clean
	rm -f $(BINARY_NAME)