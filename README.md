## Ludicrous MV Server

[![Build Status](https://travis-ci.org/Ludicrous-MV/server.svg?branch=master)](https://travis-ci.org/Ludicrous-MV/server)
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;[![wercker status](https://app.wercker.com/status/037ce014d4ef61782a039dad204b2349/s "wercker status")](https://app.wercker.com/project/bykey/037ce014d4ef61782a039dad204b2349)

## Installation

### Git

    $ git clone https://github.com/Ludicrous-MV/server.git
    $ cd server
    $ go get ./...
    $ go get labix.org/v2/mgo
    $ make
    $ make install

### Go Install

    $ go get github.com/Ludicrous-MV/server
    $ go install github.com/Ludicrous-MV/server

## Assumptions

The server automatically attempts to connect to a MongoDB instance hosted on `localhost` and creates/uses a database called `Ludicrous-MV`, storing uploaded files in a collections called `Files`.

If you would like to change any of this (use a different host, database, etc.) you will need to pull down the source and modify the constants defined in `server.go` before building it.

## Usage

    $ lmv-server -h
    Usage of lmv-server:
      -host="127.0.0.1":
      -pid=false: Save the PID to lmv-server.pid
      -port="5688":

## API Endpoints


| Method | Endpoint | Parameters | Response |
|:------:|:---------|:-----|:---------|
| GET | `/files/` | | __200__<br>Array of JSON objects representing files on the server |
| GET | `/files/<token>` | | __200__<br>JSON object representing the file associated with the token<br>__404__<br>When the token requested doesn't exist on the server|
| POST | `/files/` | __Headers__<br>`Content-Type: application/json`<br>__Body__<br>JSON representation of the file | __200__<br>JSON object containing the token associated with the new file |
