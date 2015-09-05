# 2048wc
2048wc is an implementation of a business idea for a website, where people can challenge their facebook, twitter or gmail (google+) friends in 2048 and win (or loose) bragging rights or tiny amount of money from each other.

To setup your developer environment after cloning the source run "source activate". This bash script is derived from virtualenv's activate, but has been modified to work with go. Go installation is not required -- repository ships its own go. To use it, type "setup" after "source activate". 3rd party library setup is not required, these are shipped as well and are included in the master branch.

To generate documentation for the project run createDoc. To run unit tests go to the test directory and run "go test". To run some code, you can cd to one of the website's packages and run a main program, which lives inside and shows some basic funcionality of a given package.



There is no webserver to run (yet).
