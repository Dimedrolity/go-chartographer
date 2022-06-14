default: clean test build run

test:
	go test ./...

build:
	go build -o build/app cmd/main/main.go

build_windows:
	GOOS=windows go build -o build/app_windows.exe cmd/main/main.go

build_linux:
	GOOS=linux go build -o build/app_linux cmd/main/main.go

run:
	./build/app ${ARGS}

clean:
	rm -rf build