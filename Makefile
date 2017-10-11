.PHONY: all build deploy clean

build: clean
	hugo

deploy: build
	go run cmd_deploy.go

clean:
	/bin/rm -rf _build

serve:
	hugo serve -w -p 8080

