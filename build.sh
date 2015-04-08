# Save the pwd before we run anything
PRE_PWD=`pwd`
 
# Determine the build script's actual directory, following symlinks
SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ] ; do SOURCE="$(readlink "$SOURCE")"; done
BUILD_DIR="$(cd -P "$(dirname "$SOURCE")" && pwd)"
 
# Derive the project name from the directory
PROJECT="$(basename $BUILD_DIR)"

cd $BUILD_DIR
cd data

# Prepare files
if [ "$1" = "debug" ]; then
    coffee -b -c -m js/*.coffee
    $GOPATH/bin/go-bindata -debug -pkg="data" -o data.go -ignore=\.\*\\.go  ./...
else
    coffee -b -c js/*.coffee
    $GOPATH/bin/go-bindata -pkg="data" -o data.go -ignore=\.\*\\.go\|\.\*\\.coffee  ./...
fi

cd $BUILD_DIR

# Build the project
mkdir -p bin
go build
 
EXIT_STATUS=$?
 
if [ $EXIT_STATUS == 0 ]; then
  echo "Build succeeded"
else
  echo "Build failed"
fi
 
# Change back to where we were
cd $PRE_PWD
 
exit $EXIT_STATUS



