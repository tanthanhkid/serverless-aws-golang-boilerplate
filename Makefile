.PHONY: build clean deploy

build:
	cd getFolder && env GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o ../bin/getBin getExample.go && cd ..
	cd getFolder && env GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o ../bin/getQueryBin getQueryExample.go && cd ..
	cd postFolder && env GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o ../bin/postBin ./postExample.go && cd ..
	cd postFolderApi2 && env GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o ../bin/postBinApi2 ./postExample.go && cd ..
	cd postFolderApi3 && env GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o ../bin/postBinApi3 ./postExample.go && cd ..

clean:
	rm -rf ./bin ./vendor Gopkg.lock
 
deploy: clean build
	sls deploy --verbose
