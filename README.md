<h1 align="center">equalizer</h1>
<p align="center">
  <a href="https://travis-ci.org/reugn/equalizer"><img src="https://travis-ci.org/reugn/equalizer.svg?branch=master"></a>
  <a href="https://pkg.go.dev/github.com/reugn/equalizer"><img src="https://pkg.go.dev/badge/github.com/reugn/equalizer"></a>
  <a href="https://goreportcard.com/report/github.com/reugn/equalizer"><img src="https://goreportcard.com/badge/github.com/reugn/equalizer"></a>
  <a href="https://codecov.io/gh/reugn/equalizer"><img src="https://codecov.io/gh/reugn/equalizer/branch/master/graph/badge.svg"></a>
</p>
<p align="center">A rate limiters package for Go.</p>

Pick one of the rate limiters to throttle requests and control quota.
* [Equalizer](#equalizer)
* [Slider](#slider)
* [TokenBucket](#tokenbucket)

## Equalizer
`Equalizer` is a rate limiter that manages quota based on previous requests' statuses and slows down or accelerates accordingly.

### Usage
```go
offset := equalizer.NewRandomOffset(96)

// an Equalizer with the bitmap size of 96 with 16 reserved
// positive bits and the random offset manager
eq := equalizer.NewEqualizer(96, 16, offset)

// non-blocking quota request
haveQuota := eq.Ask()

// update with ten previous successful requests
eq.Notify(true, 10)
```

### Benchmarks
```sh
BenchmarkEqualizerShortAskStep-16       30607452                37.5 ns/op             0 B/op          0 allocs/op
BenchmarkEqualizerShortAskRandom-16     31896340                34.5 ns/op             0 B/op          0 allocs/op
BenchmarkEqualizerShortNotify-16        12715494                81.9 ns/op             0 B/op          0 allocs/op
BenchmarkEqualizerLongAskStep-16        34627239                35.4 ns/op             0 B/op          0 allocs/op
BenchmarkEqualizerLongAskRandom-16      32399748                34.0 ns/op             0 B/op          0 allocs/op
BenchmarkEqualizerLongNotify-16            59935               20343 ns/op             0 B/op          0 allocs/op
```

## Slider
`Slider` rate limiter is based on a sliding window with a specified quota capacity.
Implements the `Limiter` interface.

### Usage
```go
// a Slider with one second window size, 100 millis sliding interval
// and the capacity of 32
slider := equalizer.NewSlider(time.Second, time.Millisecond*100, 32)

// non-blocking quota request
haveQuota := slider.Ask()

// blocking call
slider.Take()
```

### Benchmarks
```sh
BenchmarkSliderShortWindow-16           123488035                9.67 ns/op            0 B/op          0 allocs/op
BenchmarkSliderLongerWindow-16          128023276                9.76 ns/op            0 B/op          0 allocs/op
```

## TokenBucket
`TokenBucket` rate limiter is based on the token bucket algorithm with a refill interval.
Implements the `Limiter` interface.

### Usage
```go
// a TokenBucket with the capacity of 32 and 100 millis refill interval
tokenBucket := equalizer.NewTokenBucket(32, time.Millisecond*100)

// non-blocking quota request
haveQuota := tokenBucket.Ask()

// blocking call
tokenBucket.Take()
```

### Benchmarks
```sh
BenchmarkTokenBucketDenseRefill-16      212631714                5.64 ns/op            0 B/op          0 allocs/op
BenchmarkTokenBucketSparseRefill-16     211491368                5.63 ns/op            0 B/op          0 allocs/op
```

## License
Licensed under the MIT License.
