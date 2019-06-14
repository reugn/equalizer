# equalizer
[![GoDoc](https://godoc.org/github.com/reugn/equalizer?status.svg)](https://godoc.org/github.com/reugn/equalizer)
[![Go Report Card](https://goreportcard.com/badge/github.com/reugn/equalizer)](https://goreportcard.com/report/github.com/reugn/equalizer)

API call rate limiter

## About
API Quota is a very valuable resource. And while using it very carefully how do we back off in case of quota exceeded. If we continue with requests it can cause waste load and account lock in specific cases.  
Distributed quota management service is a solution but what if this is not a case? Equalizer provides the ability to manage quota based on previous responses and slow down or accelerate accordingly. It uses a bitmap of variable size which is quick and efficient.

## Usage
```go
offset := NewRandomOffset(96)
eq := NewEqualizer(96, 16, offset)

//claim quota
haveQuota := eq.Claim()

//nofify equalizer
eq.Notify(true, 10)
```

## Benchmark results
```go
BenchmarkEqualizerShortClaimStep-8      10000000               144 ns/op               0 B/op          0 allocs/op
BenchmarkEqualizerShortClaimRandom-8    10000000               155 ns/op               0 B/op          0 allocs/op
BenchmarkEqualizerShortNotify-8          5000000               266 ns/op               0 B/op          0 allocs/op
BenchmarkEqualizerLongClaimStep-8       10000000               143 ns/op               0 B/op          0 allocs/op
BenchmarkEqualizerLongClaimRandom-8     10000000               146 ns/op               0 B/op          0 allocs/op
BenchmarkEqualizerLongNotify-8             30000             48446 ns/op               0 B/op          0 allocs/op
```

## License
Licensed under the MIT License.
