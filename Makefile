.PHONY: build clean deploy

build: 
	cd postInsert && env GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o ../bin/postInsert ./postInsert.go && cd ..
	cd postUpdate && env GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o ../bin/postUpdate ./postUpdate.go && cd ..
	cd postDelete && env GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o ../bin/postDelete ./postDelete.go && cd ..
	cd postGetUsers && env GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o ../bin/postGetUsers ./postGetUsers.go && cd ..

clean:
	rm -rf ./bin ./vendor Gopkg.lock
 
deploy: clean build
	sls deploy --verbose
