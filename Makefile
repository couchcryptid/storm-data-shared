.PHONY: test test-cover lint fmt vuln clean

test:
	go test ./... -v -race -count=1

test-cover:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out

lint:
	golangci-lint run ./...

fmt:
	gofmt -w .
	goimports -w .

vuln:
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...

clean:
	rm -f coverage.out
