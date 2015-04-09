
cd data

# Prepare files
if [ "$1" = "debug" ]; then
    coffee -b -c -m js/*.coffee
    $GOPATH/bin/go-bindata -debug -pkg="data" -o data.go -ignore=\.\*\\.go  ./...
else
    coffee -b -c js/*.coffee
    $GOPATH/bin/go-bindata -pkg="data" -o data.go -ignore=\.\*\\.go\|\.\*\\.coffee  ./...
fi

cd ..

go build
 
EXIT_STATUS=$?
 
if [ "$EXIT_STATUS" = "0" ]; then
  echo "Build succeeded"
else
  echo "Build failed"
fi
 
exit $EXIT_STATUS

