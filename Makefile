default: build

build:
	go build -o pkt ./cmd/pkt 

test:
	go test ./...

run: build
	./pkt


clean:
	rm *.bolt
