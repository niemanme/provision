#!/usr/bin/env bash

set -e

[[ $GOPATH ]] || export GOPATH="$HOME/go"
fgrep -q "$GOPATH/bin" <<< "$PATH" || export PATH="$PATH:$GOPATH/bin"

[[ -d "$GOPATH/src/github.com/digitalrebar/provision" ]] || go get github.com/digitalrebar/provision

cd "$GOPATH/src/github.com/digitalrebar/provision"
if ! which go &>/dev/null; then
        echo "Must have go installed"
        exit 255
fi

# Work out the GO version we are working with:
GO_VERSION=$(go version | awk '{ print $3 }' | sed 's/go//')
WANTED_VER=(1 8)
if ! [[ "$GO_VERSION" =~ ([0-9]+)\.([0-9]+) ]]; then
    echo "Cannot figure out what version of Go is installed"
    exit 1
elif ! (( ${BASH_REMATCH[1]} > ${WANTED_VER[0]} || ${BASH_REMATCH[2]} >= ${WANTED_VER[1]} )); then
    echo "Go Version needs to be 1.8 or higher: currently $GO_VERSION"
    exit -1
fi

for tool in go-bindata swagger glide; do
    which "$tool" &>/dev/null && continue
    case $tool in
        go-bindata) go get -u github.com/jteeuwen/go-bindata/...;;
        swagger)    go get -u github.com/go-swagger/go-swagger/cmd/swagger;;
        glide)
            go get -v github.com/Masterminds/glide
            (cd "$GOPATH/src/github.com/Masterminds/glide" && git checkout tags/v0.12.3 && go install);;
        *) echo "Don't know how to install $tool"; exit 1;;
    esac
done

glide install
rm -rf client models embedded/assets/swagger.json
go generate server/assets.go

# Update cli docs if needed. - does change date.
go build -o drpcli-docs cmds/drpcli-docs/drpcli-docs.go

. tools/version.sh

echo "Version = $Prepart$MajorV.$MinorV.$PatchV$Extra-$GITHASH"

VERFLAGS="-X github.com/digitalrebar/provision.RS_MAJOR_VERSION=$MajorV \
          -X github.com/digitalrebar/provision.RS_MINOR_VERSION=$MinorV \
          -X github.com/digitalrebar/provision.RS_PATCH_VERSION=$PatchV \
          -X github.com/digitalrebar/provision.RS_EXTRA=$Extra \
          -X github.com/digitalrebar/provision.RS_PREPART=$Prepart \
          -X github.com/digitalrebar/provision.BuildStamp=`date -u '+%Y-%m-%d_%I:%M:%S%p'` \
          -X github.com/digitalrebar/provision.GitHash=$GITHASH"

arches=("amd64" "386" "arm")
oses=("linux" "darwin" "windows")
for arch in "${arches[@]}"; do
    for os in "${oses[@]}"; do
        (
            export GOOS="$os" GOARCH="$arch"
            echo "Building binaries for ${arch} ${os}"
            binpath="bin/$os/$arch"
            mkdir -p "$binpath"
            go build -ldflags "$VERFLAGS" -o "$binpath/dr-provision" cmds/dr-provision/dr-provision.go
            go build -ldflags "$VERFLAGS" -o "$binpath/drpcli" cmds/drpcli/drpcli.go
        )
        done
done
echo "To run tests, run: tools/test.sh"
