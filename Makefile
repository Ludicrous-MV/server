all: clean dependencies tracker

dependencies:

	go get github.com/dchest/uniuri
	go get github.com/gin-gonic/gin
	go get github.com/jinzhu/gorm
	go get github.com/mattn/go-sqlite3

clean:

	rm -f lmv-server

tracker:

	go build lmv-tracker.go

install:

	mv lmv-tracker /usr/local/bin/.

uninstall:

	rm -f /usr/local/bin/lmv-tracker
