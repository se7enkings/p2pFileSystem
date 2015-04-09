
PROJECT="$(basename `pwd`)"

checkErr () {
    EXIT_STATUS=$1
    if [ "$EXIT_STATUS" != "0" ]; then
        echo "Build failed"
        exit $EXIT_STATUS
    fi
}

# Prepare files
cd data
if [ "$1" = "debug" ]; then
    coffee -b -c -m js/*.coffee
    $GOPATH/bin/go-bindata -debug -pkg="data" -o data.go -ignore=\.\*\\.go  ./...
else
    coffee -b -c js/*.coffee
    $GOPATH/bin/go-bindata -pkg="data" -o data.go -ignore=\.\*\\.go\|\.\*\\.coffee  ./...
fi
cd ..

if [ "$1" = "release" ]; then
    oss="windows darwin linux"
    archs="amd64 386"
    for os in $oss
    do
        for arch in $archs
        do
            export GOOS=$os
            export GOARCH=$arch
            filename=${PROJECT}_${os}_${arch}
            if [ "$os" = "windows" ]; then
                filename=$filename.exe
            fi
            go build -o $filename
            checkErr $?
        done
    done
fi

unset GOARCH
unset GOOS
go build
checkErr $?
 
echo "Build succeeded"
exit 0

