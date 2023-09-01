# Regexcmp
A tool to compare different *regexp* libraries in Go.

## Usage
* Build docker image first:
```
docker build -t regexcmp .
```

* Run:
```go
# To benchmark all libs with 100MB payload 
docker run -it regexcmp 1000000 

# To benchmark specific lib
docker run -it regexcmp 1000000 rure

# Use help to see additional arguments
docker run -it regexcmp -h
```
* Example output:
```go
Generate data...
Test data size: 100.00MB

Run RURE:
  [bitcoin] count=1000000, mem=16007.26KB, time=2.393203s 
  [non_matching_email] count=0, mem=0.10KB, time=4.978ms
  [non_matching_tel] count=0, mem=0.10KB, time=19.4ms 
  [email] count=1000000, mem=16007.26KB, time=180.339ms 
  [uri] count=1000000, mem=16007.26KB, time=72.79ms 
  [tel] count=1000000, mem=16007.26KB, time=87.417ms 
  [non_matching_bitcoin] count=0, mem=0.10KB, time=4.701ms
  [non_matching_ssn] count=0, mem=0.10KB, time=4.776ms
  [non_matching_uri] count=0, mem=0.10KB, time=3.953ms
  [ssn] count=1000000, mem=16007.26KB, time=63.091ms 
Total. Counted: 5000000, Memory: 80.04MB, Duration: 2.834648s

Run PCRE:
  [non_matching_bitcoin] count=0, mem=2.83KB, time=5.752038s 
  [non_matching_ssn] count=0, mem=0.38KB, time=4.310001s 
  [non_matching_uri] count=0, mem=0.58KB, time=12.255856s 
  [ssn] count=1000000, mem=185728.36KB, time=2.158912s 
  [uri] count=1000000, mem=185728.62KB, time=2.932053s 
  [tel] count=1000000, mem=185730.73KB, time=374.353ms 
  [non_matching_tel] count=0, mem=0.54KB, time=9.50965s 
  [email] count=1000000, mem=185737.35KB, time=50.398534s 
  [bitcoin] count=1000000, mem=185728.50KB, time=2.209771s 
  [non_matching_email] count=0, mem=0.90KB, time=1m46.619988s 
Total. Counted: 5000000, Memory: 928.67MB, Duration: 3m16.521157s

Run RE2:
  [non_matching_uri] count=0, mem=100475.89KB, time=464.34ms 
  [ssn] count=1000000, mem=254206.04KB, time=955.798ms 
  [uri] count=1000000, mem=254201.70KB, time=1.0334s 
  [tel] count=1000000, mem=254210.47KB, time=676.73ms 
  [non_matching_bitcoin] count=0, mem=101309.31KB, time=486.574ms 
  [non_matching_ssn] count=0, mem=100473.73KB, time=484.58ms 
  [email] count=1000000, mem=254214.71KB, time=993.26ms 
  [bitcoin] count=1000000, mem=255037.29KB, time=1.085351s 
  [non_matching_email] count=0, mem=100473.73KB, time=460.276ms 
  [non_matching_tel] count=0, mem=100482.41KB, time=131.845ms 
Total. Counted: 5000000, Memory: 1775.10MB, Duration: 6.772154s

Run HYPER:
  [non_matching_email] count=0, mem=0.68KB, time=4.664ms
  [non_matching_tel] count=0, mem=0.68KB, time=4.608ms 
  [email] count=1000000, mem=313924.02KB, time=503.998ms 
  [bitcoin] count=1000000, mem=313917.78KB, time=2.009167s 
  [tel] count=1000000, mem=313917.78KB, time=1.098001s 
  [non_matching_bitcoin] count=0, mem=0.68KB, time=4.763ms
  [non_matching_ssn] count=0, mem=0.68KB, time=5.217ms
  [non_matching_uri] count=0, mem=0.68KB, time=99.45ms 
  [ssn] count=1000000, mem=313917.78KB, time=176.874ms 
  [uri] count=1000000, mem=313917.78KB, time=695.772ms 
Total. Counted: 5000000, Memory: 1569.61MB, Duration: 4.602514s

Run REGEXP2:
  [ssn] count=1000000, mem=853734.35KB, time=2.875972s 
  [uri] count=1000000, mem=853739.38KB, time=21.677275s 
  [tel] count=1000000, mem=853741.90KB, time=1.518279s 
  [non_matching_bitcoin] count=0, mem=396009.41KB, time=3.874902s 
  [non_matching_ssn] count=0, mem=396006.73KB, time=3.208081s 
  [non_matching_uri] count=0, mem=396011.63KB, time=33.08703s 
  [email] count=1000000, mem=997750.80KB, time=2m6.951302s 
  [bitcoin] count=1000000, mem=901736.76KB, time=3.49451s 
  [non_matching_email] count=0, mem=396023.10KB, time=2m21.045346s 
  [non_matching_tel] count=0, mem=396014.55KB, time=18.878766s 
Total. Counted: 5000000, Memory: 6440.79MB, Duration: 5m56.611463s

Run DEFAULT:
  [ssn] count=1000000, mem=143931.77KB, time=2.667804s 
  [uri] count=1000000, mem=143926.62KB, time=3.576647s 
  [tel] count=1000000, mem=143935.06KB, time=1.525194s 
  [non_matching_bitcoin] count=0, mem=65.23KB, time=3.034315s 
  [non_matching_ssn] count=0, mem=3.75KB, time=2.589903s 
  [non_matching_uri] count=0, mem=6.07KB, time=3.605346s 
  [email] count=1000000, mem=191935.21KB, time=11.309224s 
  [bitcoin] count=1000000, mem=159985.27KB, time=3.130667s 
  [non_matching_email] count=0, mem=22.83KB, time=11.396394s 
  [non_matching_tel] count=0, mem=14.94KB, time=2.437024s 
Total. Counted: 5000000, Memory: 783.83MB, Duration: 45.272517s
```

## Used libraries
* [go-re2](https://github.com/wasilibs/go-re2)
* [regexp2](https://github.com/dlclark/regexp2)
* [go-pcre](https://github.com/GRbit/go-pcre)
* [rure-go](https://github.com/BurntSushi/rure-go)
* [gohs](https://github.com/flier/gohs)
* [go-yara](https://github.com/hillu/go-yara)
