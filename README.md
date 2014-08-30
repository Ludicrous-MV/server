## Ludicrous MV Server

[![Build Status](https://travis-ci.org/Ludicrous-MV/server.svg?branch=master)](https://travis-ci.org/Ludicrous-MV/server)
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;[![wercker status](https://app.wercker.com/status/037ce014d4ef61782a039dad204b2349/s "wercker status")](https://app.wercker.com/project/bykey/037ce014d4ef61782a039dad204b2349)

## Installation

    $ git clone https://github.com/citruspi/Ludicrous-MV-Tracker.git
    $ cd Ludicrous-MV-Tracker
    $ make
    $ make install

## Usage

    $ ./lmv-tracker --help
    Usage of ./lmv-tracker:
      -pid=false: Save the PID to lmv-server.pid

## API Endpoints


| Method | Endpoint | Parameters | Response |
|:------:|:---------|:-----|:---------|
| GET | `/files/` | | __200__<br>Array of JSON objects representing files on the server |
| GET | `/files/<token>` | | __200__<br>JSON object representing the file associated with the token<br>__404__<br>When the token requested doesn't exist on the server|
| POST | `/files/` | __Headers__<br>`Content-Type: application/json`<br>__Body__<br>JSON representation of the file | __200__<br>JSON object containing the token associated with the new file |
