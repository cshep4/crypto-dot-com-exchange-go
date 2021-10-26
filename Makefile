LINTER_VERSION=v1.42.1

get-linter:
	command -v golangci-lint || curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b ${GOPATH}/bin ${LINTER_VERSION}

lint: get-linter
	golangci-lint run --timeout=5m


gen:
	go install github.com/golang/mock/mockgen
	go generate -mod vendor ./...

clean:
	find internal -iname '*.gen.go' -exec rm {} \;

regenerate: clean gen

vendor:
	go mod vendor

test: lint
	go test -v -mod vendor -race -coverprofile=race.out ./...

cover-ci:
	go tool cover -func=race.out

cover:
	go test -race ./... -vet all -coverprofile=coverage.out
	go tool cover -html=coverage.out


