# gobdns

**THIS IS STILL IN DEVELOPMENT AND MAY CHANGE AT ANY MOMENT. IT HAS ALSO NOT
BEEN TESTED AT ALL**

A simple, dynamic dns server written in go. It can direct all requests for an
address, or addresses matching `*.<address>` or `*-<address>` to a particular
ip.

For example, if you set `brian.turtles.com` to point to `127.0.0.5`, then all
these domains will get directed at that address:

* `brian.turtles.com`
* `foo.brian.turtles.com`
* `bar-brian.turtles.com`
* `foo-bar.brian.turtles.com`

and so on. This is useful for setting up development environments for multiple
employees, each one maybe having multiple virtual hosts.

## Features

* TCP and UDP support
* Easy to setup
* Prefix-wildcard matching of requests (most specific wins)
* Simple REST api for retrieving and modifying entries
* Disk persistence
* Basic master/slave replication
* Can forward requests for unknown domains to a different nameserver

## Building

First, clone the repo.

To build the web console stuff you'll need elm 0.13 installed and on your path.
Run the following to install the elm dependencies:

    cd http/console
    elm-get install

Then to get the go dependencies:

    # Installs the go-bindata binary
    go get github.com/jteeuwen/go-bindata/...

    # Retrieves dependencies for this project
    go get ./...

After that:

    # Build the gobdns binary
    make clean all

## Usage

### Running

Run `./gobdns --example` to output the default configuration. This can be piped
to a file, modified and used with the `--config` flag.

### Adding entries

Once running, you can use the REST api to add and remove entries (TODO: hoping
for an actual web console soon).

Assuming the REST interface is on `localhost:8080`:

See all existing entries:

    curl -i localhost:8080/api/domains/all

Set `foo.turtles.com` to point to whatever ip the server sees the request coming
from:

    curl -i -XPOST localhost:8080/api/domains/foo.turtles.com

Set `bar.turtles.com` to point to `127.0.0.5`:

    curl -i -XPOST -d'127.0.0.5' localhost:8080/api/domains/bar.turtles.com

Delete the `foo.turtles.com` entry:

    curl -i -XDELETE localhost:8080/api/domains/foo.turtles.com

## Persistence

If `backup-file` is set (it is by default) then every second a snapshot of the
current data set will be written to disk. On startup (again, if `backup-file` is
set) this file will be read in if it exists and used as the initial set of
mappings.

## Replication

`master-addr` can be set to point to the REST interface of another running
gobdns instance, and every 5 seconds will pull the full list of entries from
that instance and overwrite the current list.
