![Build Status](https://travis-ci.org/2048wc/2048wc.svg?branch=master)

# 2048wc

2048wc is an implementation of a business idea for a website, where people can challenge their facebook, twitter or gmail (google+) friends in 2048 and win (or loose) bragging rights or tiny amount of money from each other.

After cloning this repository execute `./setup` to:

1. Check if you have go environment installed (you have to install it yourself).
2. Set up a pre-commit hook, which will bar you from pushing broken or unformatted code to git.

# boardLib

boardLib is a library implementing game logic of 2048. It supports serialisation of a game board into and from json. Because of this serialisation capability, boardLib can be combined with a database and a dynamic web server to build a gaming platform.

# simpleCLI

Simple command line interface, which talks directly to boardLib without any database or web server. It doesn't depend on ncurses either, and is therefore a bit crude (no screen refereshing, you have to choose direction and press enter).

# userAuthLib

TODO

# dynamicWebServer

TODO
