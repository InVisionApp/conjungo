package main

import (
	"encoding/json"
	"fmt"
	"github.com/InVisionApp/go-merge"
	"github.com/InVisionApp/go-merge/util"
	log "github.com/Sirupsen/logrus"
	"reflect"
)

func init() {
	log.SetLevel(log.InfoLevel)
}

func main() {
	fmt.Println("Simple merge")
	Simple()

	fmt.Println()
	fmt.Println("Custom merge func")
	CustomMerge()

	fmt.Println()
	fmt.Println("No overwrite")
	NoOverwrite()

	fmt.Println()
	fmt.Println("From JSON")
	FromJSON()
}

func Simple() {
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

	newMap, err := merge.Merge(targetMap, sourceMap, merge.NewOptions())
	if err != nil {
		log.Fatal(err)
	}

	util.MarshalIndentPrint(newMap)
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

	opts.MergeFuncs.SetKindMergeFunc(
		reflect.TypeOf(struct{}{}).Kind(),
		// merge two 'struct' kinds by replacing the target with the source
		func(t, s interface{}, o *merge.Options) (interface{}, error) {
			return s, nil
		},
	)

	newMap, err := merge.Merge(targetMap, sourceMap, opts)
	if err != nil {
		log.Fatal(err)
	}

	util.MarshalIndentPrint(newMap)
}

func NoOverwrite() {
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

	opts := merge.NewOptions()
	opts.Overwrite = false
	newMap, err := merge.Merge(targetMap, sourceMap, opts)
	if err != nil {
		log.Fatal(err)
	}

	util.MarshalIndentPrint(newMap)
}

type Foo struct {
	A string             `json:"a,omitempty"`
	B int64              `json:"b,omitempty"`
	C map[string]string  `json:"c,omitempty"`
	D *map[string]string `json:"d,omitempty"`
	E []string           `json:"e,omitempty"`
	F []int              `json:"f,omitempty"`
}

func FromJSON() {
	targetFoo := Foo{
		A: "wrong",
		B: 1,
		C: map[string]string{"foo": "unchanged", "bar": "orig"},
		D: &map[string]string{"foo": "unchanged", "bar": "orig"},
		E: []string{"old"},
		F: []int{1},
	}

	sourceFoo := Foo{
		A: "correct",
		B: 2,
		C: map[string]string{"bar": "newVal", "safe": "added"},
		D: &map[string]string{"bar": "newVal", "safe": "added"},
		E: []string{"new"},
	}

	targetMap, err := targetFoo.toMapViaJSON()
	if err != nil {
		log.Fatal(err)
	}

	sourceMap, err := sourceFoo.toMapViaJSON()
	if err != nil {
		log.Fatal(err)
	}

	resultMap, err := merge.Merge(targetMap, sourceMap, merge.NewOptions())
	if err != nil {
		log.Fatal(err)
	}

	util.MarshalIndentPrint(resultMap)
}

func (f *Foo) toMapViaJSON() (map[string]interface{}, error) {
	jsonBody, err := json.Marshal(f)
	if err != nil {
		return nil, err
	}

	resultMap := &map[string]interface{}{}
	if err := json.Unmarshal(jsonBody, resultMap); err != nil {
		return nil, err
	}

	return *resultMap, nil
}
