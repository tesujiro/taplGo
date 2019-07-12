all : tapl.go ./parser/*.y ./parser/*.go ./ast/*.go ./vm/*.go ./parser/grammar.go
	go build -o tapl .

./parser/grammar.go : ./parser/grammar.go.y ./ast/*.go
	goyacc -o ./parser/grammar.go ./parser/grammar.go.y
	gofmt -s -w ./parser

.PHONY: test
test: ./*_test.go ./parser/grammar.go
	go vet ./...
	go test -v -count=1 . -coverpkg ./...

.PHONY: cover
cover:
	go test -v . -coverpkg ./... -coverprofile=cover.out && go tool cover -html=cover.out -o cover.html

