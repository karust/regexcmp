# Regexcmp
A tool to compare different *regexp* libraries in Go.

## Usage
* Build docker image first:
```
docker build -t regexcmp .
```

* Run:
```shell
# To benchmark all libs with 100MB payload 
docker run -it regexcmp 1000000 

# To benchmark specific lib
docker run -it regexcmp 1000000 rure

# Use help to see additional arguments
docker run -it regexcmp -h
```

## Used libraries
* [go-re2](https://github.com/wasilibs/go-re2)
* [regexp2](https://github.com/dlclark/regexp2)
* [go-pcre](https://github.com/GRbit/go-pcre)
* [rure-go](https://github.com/BurntSushi/rure-go)
* [gohs](https://github.com/flier/gohs)
* [go-yara](https://github.com/hillu/go-yara)
