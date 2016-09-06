.PHONY: all build deploy clean

build: clean
	hugo

deploy: build
	go run cmd_deploy.go

clean:
	/bin/rm -rf public sbinet.github.io
