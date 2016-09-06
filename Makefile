.PHONY: all build deploy

build:
	hugo

deploy: build
	/bin/rm -rf _blog
	git clone git@github.com:sbinet/sbinet.github.io _blog
	/bin/cp -fr public/* _blog/.
	(cd _blog && git add -A . && git commit -m "update" && git push origin master)


