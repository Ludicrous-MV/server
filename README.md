## Ludicrous MV Server

[![wercker status](https://app.wercker.com/status/f86323ec0e58822770ce55241591999c/s/master "wercker status")](https://app.wercker.com/project/bykey/f86323ec0e58822770ce55241591999c)

## Installation

    $ git clone https://github.com/citruspi/Ludicrous-MV-Tracker.git
    $ cd Ludicrous-MV-Tracker
    $ make
    $ make install

## Configuration

There must be a configuration file in the same directory as the binary. It should be named `lmv-tracker.yml` and should be formatted as such:

```
tokens:
    pool: ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789
    length: 10
system:
    pid: False
web:
    address: :8080
database:
    type: sqlite3
    source: lmv-tracker.db
```

A sample configuration file is included.

## Usage

    $ ./lmv-tracker
