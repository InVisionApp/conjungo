package main

import (
	"encoding/json"
	"fmt"
	"github.com/InVisionApp/go-merge"
	"github.com/InVisionApp/go-merge/util"
	log "github.com/Sirupsen/logrus"
	"reflect"
)

type Foo struct {
	A string             `json:"a,omitempty"`
	B int64              `json:"b,omitempty"`
	C map[string]string  `json:"c,omitempty"`
	D *map[string]string `json:"d,omitempty"`
	E []string           `json:"e,omitempty"`
	F []int              `json:"f,omitempty"`
}

func init() {
	log.SetLevel(log.InfoLevel)

}
func main() {
	destFoo := Foo{
		A: "wrong",
		B: 1,
		C: map[string]string{"foo": "unchanged", "bar": "orig"},
		D: &map[string]string{"foo": "unchanged", "bar": "orig"},
		E: []string{"old"},
		F: []int{1},
	}

	newFoo := Foo{
		A: "correct",
		B: 2,
		C: map[string]string{"bar": "newVal", "safe": "added"},
		D: &map[string]string{"bar": "newVal", "safe": "added"},
		E: []string{"new"},
	}

	destMap, err := destFoo.toMap()
	if err != nil {
		log.Fatal(err)
	}

	destMap2, err := destFoo.toMap()
	if err != nil {
		log.Fatal(err)
	}

	newFMap, err := newFoo.toMap()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("NEWMAP: %v", newFMap)

	opts := merge.NewOptions()
	opts.SetMergeFunc(
		reflect.TypeOf(float64(0)),
		func(t, s interface{}, o *merge.Options) (interface{}, error) {
			iT, _ := t.(float64)
			iS, _ := s.(float64)
			return iT + iS, nil
		},
	)

	destMapOver, err := merge.Merge(destMap, newFMap, opts)
	if err != nil {
		log.Fatal(err)
	}

	noOverOpt := opts
	noOverOpt.Overwrite = false
	merMap, err := merge.Merge(destMap2, newFMap, noOverOpt)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\n\nOUTPUT:\n")
	fmt.Println("overwrite")
	util.IndentMarshalPrint(destMapOver)

	fmt.Println("\n")

	fmt.Println("no overwrite")
	util.IndentMarshalPrint(merMap)

	fmt.Println("\nEND")

}

func (f *Foo) toMap() (map[string]interface{}, error) {
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
