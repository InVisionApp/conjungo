package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/InVisionApp/conjungo"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(log.InfoLevel)
}

func main() {
	fmt.Println("Simple map merge")
	SimpleMap()

	fmt.Println()
	fmt.Println("Simple struct merge")
	SimpleStruct()

	fmt.Println()
	fmt.Println("Custom merge func")
	CustomMerge()

	fmt.Println()
	fmt.Println("Custom interface merge func")
	CustomInterfaceMerge()

	fmt.Println()
	fmt.Println("Custom struct merge func")
	CustomStructMerge()

	fmt.Println()
	fmt.Println("No overwrite")
	NoOverwrite()

	fmt.Println()
	fmt.Println("From JSON")
	FromJSON()
}

func SimpleMap() {
	targetMap := map[string]interface{}{
		"A": "wrong",
		"B": 1,
		"C": map[string]interface{}{"foo": "unchanged", "bar": "orig"},
		"D": []interface{}{"unchanged", 0},
	}

	sourceMap := map[string]interface{}{
		"A": "correct",
		"B": 2,
		"C": map[string]interface{}{"bar": "newVal", "safe": "added"},
		"D": []interface{}{"added", 1},
	}

	err := conjungo.Merge(&targetMap, sourceMap, nil)
	if err != nil {
		log.Error(err)
	}

	marshalIndentPrint(targetMap)
}

func SimpleStruct() {
	type Foo struct {
		Name    string
		Size    int
		Special bool
		SubMap  map[string]string
	}

	targetStruct := Foo{
		Name:    "target",
		Size:    2,
		Special: false,
		SubMap:  map[string]string{"foo": "unchanged", "bar": "orig"},
	}

	sourceStruct := Foo{
		Name:    "source",
		Size:    4,
		Special: true,
		SubMap:  map[string]string{"bar": "newVal", "safe": "added"},
	}

	err := conjungo.Merge(&targetStruct, sourceStruct, nil)
	if err != nil {
		log.Error(err)
	}

	marshalIndentPrint(targetStruct)
}

func CustomMerge() {
	type foo struct {
		Bar string
	}

	targetMap := map[string]interface{}{
		"A": "wrong",
		"B": 1,
		"C": foo{"target"},
	}

	sourceMap := map[string]interface{}{
		"A": "correct",
		"B": 2,
		"C": foo{"source"},
	}

	opts := conjungo.NewOptions()
	opts.SetTypeMergeFunc(
		reflect.TypeOf(0),
		// merge two 'int' types by adding them together
		func(t, s reflect.Value, o *conjungo.Options) (reflect.Value, error) {
			iT, _ := t.Interface().(int)
			iS, _ := s.Interface().(int)
			return reflect.ValueOf(iT + iS), nil
		},
	)

	opts.SetKindMergeFunc(
		reflect.TypeOf(struct{}{}).Kind(),
		// merge two 'struct' kinds by replacing the target with the source
		// provides a mechanism to set override = true for just structs
		func(t, s reflect.Value, o *conjungo.Options) (reflect.Value, error) {
			return s, nil
		},
	)

	err := conjungo.Merge(&targetMap, sourceMap, opts)
	if err != nil {
		log.Error(err)
	}

	marshalIndentPrint(targetMap)
}

type customInterfaceFoo struct {
	Bar string
	Baz customInterfaceSpecial
	Qux *customInterfaceSpecial
}

type customInterfaceSpecialFace interface {
	HasMethod()
}

type customInterfaceSpecial string

func (s *customInterfaceSpecial) HasMethod() {}

func CustomInterfaceMerge() {
	dog := customInterfaceSpecial("dog")
	target := customInterfaceFoo{
		Bar: "hello",
		Baz: customInterfaceSpecial("beautiful"),
		Qux: &dog,
	}

	world := customInterfaceSpecial("dog")
	source := customInterfaceFoo{
		// Does not overwrite because opts.Overwrite is not set
		Bar: "aloha!",
		// Does not trigger because not a pointer, but would overwrite if opts.Overwrite is set
		Baz: customInterfaceSpecial("ugly"),
		// Triggers custom interface func
		Qux: &world,
	}

	opts := conjungo.NewOptions()
	opts.Overwrite = false
	opts.SetInterfaceMergeFunc(
		reflect.TypeOf((*customInterfaceSpecialFace)(nil)).Elem(),
		// if triggered returns pointer to type customInterfaceSpecial with value "galaxy"
		func(t, s reflect.Value, o *conjungo.Options) (reflect.Value, error) {
			galaxy := customInterfaceSpecial("galaxy")
			return reflect.ValueOf(&galaxy), nil
		},
	)

	err := conjungo.Merge(&target, source, opts)
	if err != nil {
		log.Error(err)
	}

	marshalIndentPrint(target)
}

func CustomStructMerge() {
	type Foo struct {
		Name string
		Size int
	}

	target := Foo{
		Name: "bar",
		Size: 25,
	}

	source := Foo{
		Name: "baz",
		Size: 35,
	}

	opts := conjungo.NewOptions()
	opts.SetTypeMergeFunc(
		reflect.TypeOf(Foo{}),
		// merge two 'int' types by adding them together
		func(t, s reflect.Value, o *conjungo.Options) (reflect.Value, error) {
			tFoo := t.Interface().(Foo)
			sFoo := s.Interface().(Foo)

			// names are merged by concatenating them
			tFoo.Name = tFoo.Name + "." + sFoo.Name
			// sizes are merged by averaging them
			tFoo.Size = (tFoo.Size + sFoo.Size) / 2

			return reflect.ValueOf(tFoo), nil
		},
	)

	err := conjungo.Merge(&target, source, opts)
	if err != nil {
		log.Error(err)
	}

	marshalIndentPrint(target)
}

func NoOverwrite() {
	targetMap := map[string]interface{}{
		"A": "not overwritten",
		"B": 1,
		"C": map[string]string{"foo": "unchanged", "bar": "orig"},
	}

	sourceMap := map[string]interface{}{
		"A": "overwritten",
		"B": 2,
		"C": map[string]string{"bar": "newVal", "safe": "added"},
	}

	opts := conjungo.NewOptions()
	opts.Overwrite = false
	err := conjungo.Merge(&targetMap, sourceMap, opts)
	if err != nil {
		log.Error(err)
	}

	marshalIndentPrint(targetMap)
}

type jsonString string

func FromJSON() {

	var targetJSON jsonString = `
	{
	  "a": "wrong",
	  "b": 1,
	  "c": {"bar": "orig", "foo": "unchanged"},
	  "d": ["old"],
	  "e": [1]
	}`

	var sourceJSON jsonString = `
	{
	  "a": "correct",
	  "b": 2,
	  "c": {"bar": "newVal", "safe": "added"},
	  "d": ["new"]
	}`

	opts := conjungo.NewOptions()
	opts.SetTypeMergeFunc(
		reflect.TypeOf(jsonString("")),
		// merge two json strings by unmarshalling them to maps
		func(t, s reflect.Value, o *conjungo.Options) (reflect.Value, error) {
			targetStr, _ := t.Interface().(jsonString)
			sourceStr, _ := s.Interface().(jsonString)

			targetMap := map[string]interface{}{}
			if err := json.Unmarshal([]byte(targetStr), &targetMap); err != nil {
				return reflect.Value{}, err
			}

			sourceMap := map[string]interface{}{}
			if err := json.Unmarshal([]byte(sourceStr), &sourceMap); err != nil {
				return reflect.Value{}, err
			}

			err := conjungo.Merge(&targetMap, sourceMap, o)
			if err != nil {
				return reflect.Value{}, err
			}

			mergedJSON, err := json.Marshal(targetMap)
			if err != nil {
				return reflect.Value{}, err
			}

			return reflect.ValueOf(jsonString(mergedJSON)), nil
		},
	)

	err := conjungo.Merge(&targetJSON, sourceJSON, opts)
	if err != nil {
		log.Error(err)
	}

	fmt.Println(targetJSON)
}

func marshalIndentPrint(i interface{}) error {
	jBody, err := json.MarshalIndent(i, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(jBody))
	return nil
}

// pretty print
func (s jsonString) String() string {
	out := bytes.Buffer{}
	if err := json.Indent(&out, []byte(string(s)), "", "  "); err != nil {
		log.Fatal(err)
	}

	return out.String()
}
