default: build

build:
	go build -o pkt ./cmd/pkt 

test:
	go test ./... -v

run: build
	./pkt


clean:
	rm *.bolt
