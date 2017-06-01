package main

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/InVisionApp/conjungo"
	log "github.com/Sirupsen/logrus"
)

func init() {
	log.SetLevel(log.InfoLevel)
}

func main() {
	fmt.Println("Simple merge")
	SimpleMap()

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

	newMap, err := conjungo.MergeMapStrIface(targetMap, sourceMap, nil)
	if err != nil {
		log.Fatal(err)
	}

	marshalIndentPrint(newMap)
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
	opts.MergeFuncs.SetTypeMergeFunc(
		reflect.TypeOf(0),
		// merge two 'int' types by adding them together
		func(t, s reflect.Value, o *conjungo.Options) (reflect.Value, error) {
			iT, _ := t.Interface().(int)
			iS, _ := s.Interface().(int)
			return reflect.ValueOf(iT + iS), nil
		},
	)

	opts.MergeFuncs.SetKindMergeFunc(
		reflect.TypeOf(struct{}{}).Kind(),
		// merge two 'struct' kinds by replacing the target with the source
		// provides a mechanism to set override = true for just structs
		func(t, s reflect.Value, o *conjungo.Options) (reflect.Value, error) {
			return s, nil
		},
	)

	newMap, err := conjungo.MergeMapStrIface(targetMap, sourceMap, opts)
	if err != nil {
		log.Fatal(err)
	}

	marshalIndentPrint(newMap)
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

	opts := conjungo.NewOptions()
	opts.Overwrite = false
	newMap, err := conjungo.MergeMapStrIface(targetMap, sourceMap, opts)
	if err != nil {
		log.Fatal(err)
	}

	marshalIndentPrint(newMap)
}

func FromJSON() {
	type jsonString string

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
	opts.MergeFuncs.SetTypeMergeFunc(
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

			merged, err := conjungo.MergeMapStrIface(targetMap, sourceMap, o)
			if err != nil {
				return reflect.Value{}, err
			}

			return reflect.ValueOf(merged), nil
		},
	)

	resultMap, err := conjungo.Merge(targetJSON, sourceJSON, opts)
	if err != nil {
		log.Fatal(err)
	}

	marshalIndentPrint(resultMap)
}

func marshalIndentPrint(i interface{}) error {
	jBody, err := json.MarshalIndent(i, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(jBody))
	return nil
}
