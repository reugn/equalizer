# equalizer
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

## License
Licensed under the MIT License.
