default: build

build:
	go build -o pkt ./cmd/pkt 

test:
	go test -v -cover

clean:
	rm *.bolt
