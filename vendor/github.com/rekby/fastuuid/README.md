[![PkgGoDev](https://pkg.go.dev/badge/github.com/rekby/fastuuid)](https://pkg.go.dev/github.com/rekby/fastuuid)
[![Go Report Card](https://goreportcard.com/badge/github.com/rekby/fastuuid)](https://goreportcard.com/report/github.com/rekby/fastuuid)
[![codecov](https://codecov.io/github/rekby/fastuuid/branch/main/graph/badge.svg?token=AIGUP7QYV2)](https://codecov.io/github/rekby/fastuuid)

Fast uuid library, now implemented only UUIDv4 (crypto-random) generators.
It generate UUID to string about 50-100ns with one allocation and good paralellism by cpu.

Command for install:

```bash
go get github.com/rekby/fastuuid
```

Example:

```golang
package main

import "github.com/rekby/fastuuid"

func main(){
	fmt.Println(fastuuid.MustUUIDv4String())
}

```

For run benchmarks:

```
git clone https://github.com/rekby/fastuuid.git
cd fastuuid/benchmarks
go test -test.bench=. -v
```

Result for macbook pro M1:

```
goos: darwin
goarch: arm64
pkg: github.com/rekby/fastuuid/benchmarks
BenchmarkRekbyUUID/one-thread                                        13688407               85.92 ns/op            48 B/op          1 allocs/op
BenchmarkRekbyUUID/multi-thread                                      25657815               51.81 ns/op            47 B/op          0 allocs/op
BenchmarkGoogleUUID4/one-thread                                       2256043               530.5 ns/op            64 B/op          2 allocs/op
BenchmarkGoogleUUID4/multi-thread                                     1897683               634.9 ns/op            63 B/op          1 allocs/op
BenchmarkSatoriUUID4/one-thread                                       2277092               531.1 ns/op            64 B/op          2 allocs/op
BenchmarkSatoriUUID4/multi-thread                                     1950952               624.3 ns/op            62 B/op          1 allocs/op
BenchmarkGofrs/one-thread                                             2265781               530.3 ns/op            64 B/op          2 allocs/op
BenchmarkGofrs/multi-thread                                           1912966               635.9 ns/op            63 B/op          1 allocs/op
BenchmarkRogpeppeUnsecuredBecauseItCounter/one-thread                28806307               40.35 ns/op            48 B/op          1 allocs/op
BenchmarkRogpeppeUnsecuredBecauseItCounter/multi-thread               7198768               168.0 ns/op            47 B/op          0 allocs/op
BenchmarkJakehl/one-thread                                            1956322               611.8 ns/op           144 B/op          5 allocs/op
BenchmarkJakehl/multi-thread                                          1963675               643.4 ns/op           142 B/op          4 allocs/op
BenchmarkRwxrob/one-thread                                            1620858               740.5 ns/op           184 B/op          7 allocs/op
BenchmarkRwxrob/multi-thread                                          1870315               644.6 ns/op           182 B/op          6 allocs/op
```