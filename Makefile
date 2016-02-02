build: clean
	go build

run: build
	./twitter2rss -config twitter2rss.hcl

clean:
	rm -rf twitter2rss
