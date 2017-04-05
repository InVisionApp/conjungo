# go-merge

[![Build Status](https://travis-ci.com/InVisionApp/go-merge.svg?token=KosA43m1X3ikri8JEukQ&branch=master)](https://travis-ci.com/InVisionApp/go-merge)

A merge utility designed for flexibility and customizability.
The library has a simple interface that uses a set of default merge functions that will fit most basic use 
cases. From there, specific customizations can be made to merge things in any particular way that is needed.

Merge any two things of the same type, including maps, slices, and structs.
By default, the target value will be overwritten by the source. If the overwrite option is turned off, only 
new values in source that do not already exist in target will be added.  
If you would like to change the way two items of a particular type get merged, custom merge functions 
can be defined for any type or kind (see below).  

## Usage
Merge two structs together:
```go
type Foo struct {
	Name    string
	Size    int
	Special bool
}

targetStruct := Foo{
	Name:    "target",
	Size:    2,
	Special: false,
}

sourceStruct := Foo{
	Name:    "source",
	Size:    4,
	Special: true,
}

merged, err := merge.Merge(targetStruct, sourceStruct, nil)
if err != nil {
	log.Error(err)
}
```

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

// use the map merge wrapper
merged, err := merge.MergeMapStrIface(targetMap, sourceMap, nil)
if err != nil {
	log.Error(err)
}

// OR 

// use the main merge func
merged, err := merge.Merge(targetMap, sourceMap, nil)
if err != nil {
	log.Error(err)
}
mergedMap, _ := merged.(map[string]interface{})
```

Define a custom merge function for a type:
```go
opts := merge.NewOptions()
opts.MergeFuncs.SetTypeMergeFunc(
	reflect.TypeOf(0),
	// merge two 'int' types by adding them together
	func(t, s interface{}, o *merge.Options) (interface{}, error) {
		iT, _ := t.(int)
		iS, _ := s.(int)
		return iT + iS, nil
	},
)

merged, err := merge.Merge(1, 2, opts)
if err != nil {
	log.Error(err)
}
// merged == 3
```

or for a kind:
```go
opts := merge.NewOptions()
opts.MergeFuncs.SetKindMergeFunc(
	reflect.TypeOf(struct{}{}).Kind(),
	// merge two 'struct' kinds by replacing the target with the source
	func(t, s interface{}, o *merge.Options) (interface{}, error) {
		return s, nil
	},
)
```

See [examples](example/main.go) for more details.
