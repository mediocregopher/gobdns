all:
	(cd http/console && make all)
	go-bindata -o=http/bindata.go -prefix=http/console/build -pkg=http http/console/build
	go build

clean:
	rm http/bindata.go
