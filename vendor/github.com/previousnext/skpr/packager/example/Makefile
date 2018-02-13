#!/usr/bin/make -f

PROJECT=github.com/previousnext/skpr/packager/example

build:
	gox -os='linux darwin' -arch='amd64' -output='bin/packager_{{.OS}}_{{.Arch}}' -ldflags='-extldflags "-static"' $(PROJECT)
