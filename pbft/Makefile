LIBFILE=pbft

ifeq ($(mode),debug)
	LDFLAGS="-X 'main.BUILD_TIME=`date`' -X 'main.GO_VERSION=`go version`' -X main.GIT_HASH=`git rev-parse HEAD`"
else
	LDFLAGS="-s -w -X 'main.BUILD_TIME=`date`' -X 'main.GO_VERSION=`go version`' -X main.GIT_HASH=`git rev-parse HEAD`"
endif

$(LIBFILE):*.go
	go build -ldflags $(LDFLAGS) -o $(LIBFILE)
clean:
	rm -f $(LIBFILE)
