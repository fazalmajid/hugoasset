GO=	env GOPATH=`pwd` go

all: hugoasset

DEPS=	src/github.com/lukasbob/srcset \

#	github.com/jaytaylor/html2text \
#	src/github.com/gohugoio/hugo

hugoasset: $(DEPS) hugoasset.go
	$(GO) build  hugoasset.go

src/github.com/jaytaylor/html2text:
	$(GO) get -f -t -u -v --tags fts5 github.com/jaytaylor/html2text

src/github.com/lukasbob/srcset:
	$(GO) get -f -t -u -v --tags fts5 github.com/lukasbob/srcset

src/github.com/gohugoio/hugo:
	$(GO) get -f -t -u -v --tags fts5 github.com/gohugoio/hugo

test:
	$(GO) test

clean:
	-rm -rf src pkg hugoasset *~ core search.db
