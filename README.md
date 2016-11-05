[![Build Status](https://travis-ci.org/edu-xored/tyt.svg?branch=master)](https://travis-ci.org/edu-xored/tyt)

# tyt

Track your time

## How to build and run

Inside project root:

* `go get -u github.com/constabulary/gb/...`
* `gb vendor restore`
* `gb build`
* `./bin/tyt`

## How to init database

Inside project root:

* run `./bin/tyt`
* `npm install`
* `cd scripts` 
* `node init.js`

## Watch dev mode

* `go get -u github.com/githubnemo/CompileDaemon`
* `./run.sh`
