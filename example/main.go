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

	newMap, err := merge.MergeMapStrIface(targetMap, sourceMap, merge.NewOptions())
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

	newMap, err := merge.MergeMapStrIface(targetMap, sourceMap, opts)
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
	newMap, err := merge.MergeMapStrIface(targetMap, sourceMap, opts)
	if err != nil {
		log.Fatal(err)
	}

	util.MarshalIndentPrint(newMap)
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

	opts := merge.NewOptions()
	opts.MergeFuncs.SetTypeMergeFunc(
		reflect.TypeOf(jsonString("")),
		// merge two json strings by unmarshalling them to maps
		func(t, s interface{}, o *merge.Options) (interface{}, error) {
			targetStr, _ := t.(jsonString)
			sourceStr, _ := s.(jsonString)

			targetMap := map[string]interface{}{}
			if err := json.Unmarshal([]byte(targetStr), &targetMap); err != nil {
				return nil, err
			}

			sourceMap := map[string]interface{}{}
			if err := json.Unmarshal([]byte(sourceStr), &sourceMap); err != nil {
				return nil, err
			}

			return merge.MergeMapStrIface(targetMap, sourceMap, o)
		},
	)

	resultMap, err := merge.Merge(targetJSON, sourceJSON, opts)
	if err != nil {
		log.Fatal(err)
	}

	util.MarshalIndentPrint(resultMap)
}
