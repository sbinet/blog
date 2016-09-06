.PHONY: all build deploy clean

build: clean
	hugo

deploy: build
	go run deploy.go

clean:
	/bin/rm -rf public _blog
