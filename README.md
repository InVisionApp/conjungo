#go-merge

[![Build Status](https://travis-ci.com/InVisionApp/go-merge.svg?token=KosA43m1X3ikri8JEukQ&branch=master)](https://travis-ci.com/InVisionApp/go-merge)

A merge utility designed for flexibility and customizability.

Currently supports the merging of two `map[string]interface{}` because this is used for JSON merging.
Custom merge functions can be defined for any type.

##Usage
Merge two `map[string]interface{}` together:
```go
targetMap := map[string]interface{}{
	"A": "wrong",
	"B": 1,
	"C": map[string]string{"foo": "unchanged", "bar": "orig"},
}

sourceMap := map[string]interface{}{
	"A": "correct",
	"B": 2,
	"C": map[string]string{"bar": "newVal", "safe": "added"},
}

// use the main merge func
merged, err := merge.Merge(targetMap, sourceMap, merge.NewOptions())
if err != nil {
	log.Error(err)
}
mergedMap, ok := merged.(map[string]interface{})
if !ok {
	log.Error("did not return map")
}

// OR 

// use the map merge wrapper
merged, err := merge.MergeMapStrIface(targetMap, sourceMap, merge.NewOptions())
if err != nil {
	log.Error(err)
}
```

Define a custom merge function:
```go
opts := merge.NewOptions()
opts.MergeFuncs.SetTypeMergeFunc(
	reflect.TypeOf(float64(0)),
	func(t, s interface{}, o *merge.Options) (interface{}, error) {
		iT, _ := t.(float64)
		iS, _ := s.(float64)
		return iT + iS, nil
	},
)
```
See [examples](example/main.go) for more details.
