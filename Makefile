
GO_BINARY_NAME=main

check-mod:
ifndef MOD
	$(error MOD is undefined)
endif

init: check-mod
	go mod init ${MOD}
	echo 'package main\n\nimport "fmt"\n\nfunc main() {\nfmt.Println("hello world")\n}' > main.go
	go run main.go

api:
	echo 'package main\n\nimport (\n\t"io"\n\t"net/http"\n\t"log"\n)\n\nfunc main() {\n\thelloHandler := func(w http.ResponseWriter, req *http.Request) {\n\t\tio.WriteString(w, "Hello, world!") \n\t}\n\n\thttp.HandleFunc("/hello", helloHandler)\n\tlog.Fatal(http.ListenAndServe(":8000", nil))\n}' > main.go

clean:
	find . -type f -name "*.go" | xargs rm -f
	find . -type f -name "*.mod" | xargs rm -f
	find . -type f -name "main" | xargs rm -f

fmt:
	find . -type f -name "*.go" | grep -v -E '^./vendor|^./third_party|^./_examples' | xargs -L1 dirname | sort | uniq | xargs gofmt -l

build: fmt
	go build -o ./tmp/$(GO_BINARY_NAME) main.go

run: fmt
	go run .

test: fmt
	find . -type f -name "*.go" | grep -v -E '^./vendor|^./third_party|^./_examples' | xargs -L1 dirname | sort | uniq | xargs go test -v -race

test-clean:
	go clean -testcache

vet:
	find . -type f -name "*.go" | grep -v -E '^./vendor' | xargs -L1 dirname | sort | uniq | xargs go vet 2>&1