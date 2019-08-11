# equalizer
[![Build Status](https://travis-ci.org/reugn/equalizer.svg?branch=master)](https://travis-ci.org/reugn/equalizer)
[![GoDoc](https://godoc.org/github.com/reugn/equalizer?status.svg)](https://godoc.org/github.com/reugn/equalizer)
[![Go Report Card](https://goreportcard.com/badge/github.com/reugn/equalizer)](https://goreportcard.com/report/github.com/reugn/equalizer)
[![codecov](https://codecov.io/gh/reugn/equalizer/branch/master/graph/badge.svg)](https://codecov.io/gh/reugn/equalizer)

Rate limiters collection

## About
API Quota is a very valuable resource. And while using it very carefully how do we back off in case of quota exceeded. If we continue with requests it can cause waste load and account lock in specific cases.  
Distributed quota management service is a solution but what if this is not a case? Chose one of the following limiters and your preferred backoff algorithm if needed.

## Equalizer
Equalizer provides the ability to manage quota based on previous responses and slow down or accelerate accordingly. It uses a bitmap of variable size which is quick and efficient.

### Usage
```go
offset := NewRandomOffset(96)
eq := NewEqualizer(96, 16, offset)

//claim quota
haveQuota := eq.Claim()

//nofify equalizer
eq.Notify(true, 10)
```

### Benchmark results
```go
BenchmarkEqualizerShortClaimStep-8      10000000               144 ns/op               0 B/op          0 allocs/op
BenchmarkEqualizerShortClaimRandom-8    10000000               155 ns/op               0 B/op          0 allocs/op
BenchmarkEqualizerShortNotify-8          5000000               266 ns/op               0 B/op          0 allocs/op
BenchmarkEqualizerLongClaimStep-8       10000000               143 ns/op               0 B/op          0 allocs/op
BenchmarkEqualizerLongClaimRandom-8     10000000               146 ns/op               0 B/op          0 allocs/op
BenchmarkEqualizerLongNotify-8             30000             48446 ns/op               0 B/op          0 allocs/op
```

## Slider
Sliding window based rate limiter.

### Usage
```go
//slider with 1 second window size, 100 millis sliding interval and capacity of 32
slider := NewSlider(time.Second, time.Millisecond*100, 32)

//claim quota
haveQuota := slider.Claim()
```

## Token bucket
Custom implementation for token bucket rate limiter with refill interval.

### Usage
```go
//bucket with capacity of 32 and 100 millis refill interval
tokenBucket := NewTokenBucket(32, time.Millisecond*100)

//claim quota
haveQuota := tokenBucket.Claim()
```

## License
Licensed under the MIT License.
