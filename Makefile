all: clean tracker

clean:

	rm -f lmv-server

tracker:

	go build lmv-tracker.go

install:

	mv lmv-tracker /usr/local/bin/.

uninstall:

	rm -f /usr/local/bin/lmv-tracker
