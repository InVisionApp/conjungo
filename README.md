#go-merge

A merge utility designed for flexibility and customizability.

Currently supports the merging of two `map[string]interface{}` because this is used for JSON merging.
Custom merge functions can be defined for any type.

##Usage
Basic use:
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

mergedMap, err := merge.Merge(targetMap, sourceMap, merge.NewOptions())
if err != nil {
	log.Fatal(err)
}
```

Define a custom merge function:
```go
opts := merge.NewOptions()
opts.SetMergeFunc(
	reflect.TypeOf(float64(0)),
	func(t, s interface{}, o *merge.Options) (interface{}, error) {
		iT, _ := t.(float64)
		iS, _ := s.(float64)
		return iT + iS, nil
	},
)
```
See [examples](example/main.go) for more details.
