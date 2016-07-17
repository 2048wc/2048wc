![Build Status](https://travis-ci.org/2048wc/2048wc.svg?branch=master)

# 2048wc

2048wc is an implementation of a business idea for a website, where people can challenge their facebook, twitter or gmail (google+) friends in 2048 and win (or loose) bragging rights or tiny amount of money from each other.

After cloning this repository execute `source setup2048` to:

1. Check if you have go environment installed. [You have to install it yourself](https://golang.org/doc/install).
2. Set up a pre-commit hook, which will bar you from pushing broken or unformatted code to git.
3. Set up GOPATH environment variable.

You should execute the script every time before you start working on the project.

You have to execute this script in the root of the repository (directory where the `setup2048` file lives).

## Module Specification

All modules depend on API2048, which contains a single file with all interfaces implemented by this project.

### boardLib

boardLib is a library implementing game logic of 2048. It supports serialisation of a game board into and from json. Because of this serialisation capability, boardLib can be combined with a database and a dynamic web server to build a gaming platform.

boardLib implements API2048.Move and API2048.MoveCreator interfaces.

boardLib does not have any dependencies.

boardLib is a library. It can't be run.

### simpleCLI

Simple command line interface, which talks directly to boardLib without any database or web server. It doesn't depend on ncurses either, and is therefore a bit crude (no screen refereshing, you have to choose direction and press enter).

simpleCLI does not implement any interfaces.

simpleCLI depends on boardLib.

simpleCLI can be run with `go build ./... && bin/simpleCLI`

### mockDB

mockDB is an in-memory pure-Go implementation of the database part of 2048wc. As the data does not get persisted to disk, this can't be used in production, but is useful for testing, debugging, presentations and any other situation when running a database is not preferred.

mockDB implements API2048.QueryCallback, API2048.Query and API2048.QueryBuilder interfaces.

mockDB does not have any dependencies.

mockDB is and can't be run directly.

### userAuthLib

TODO

### dynamicWebServer

TODO
