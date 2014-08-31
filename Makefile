all: clean dependencies tracker

dependencies:

	go get github.com/tsuru/config
	go get github.com/dchest/uniuri
	go get github.com/gin-gonic/gin
	go get github.com/jinzhu/gorm
	go get github.com/mattn/go-sqlite3
	go get github.com/go-sql-driver/mysql

clean:

	rm -f lmv-tracker
	rm -f lmv-tracker.db
	rm -f lmv-tracker.pid

tracker:

	go build lmv-tracker.go

install:

	mv lmv-tracker /usr/local/bin/.

uninstall:

	rm -f /usr/local/bin/lmv-tracker
