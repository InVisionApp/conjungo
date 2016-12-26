package main

import (
	"encoding/json"
	"fmt"
	"github.com/InVisionApp/go-merge"
	log "github.com/Sirupsen/logrus"
)

type Foo struct {
	A string `json:"a"`
	B int64  `json:"b"`
	C map[string]string
	D *map[string]string
}

func init() {
	log.SetLevel(log.InfoLevel)

}
func main() {
	destFoo := Foo{
		A: "wrong",
		C: map[string]string{"foo": "unchanged", "bar": "orig"},
		D: &map[string]string{"foo": "unchanged", "bar": "orig"},
	}

	newFoo := Foo{
		A: "correct",
		B: 2,
		C: map[string]string{"bar": "newVal", "safe": "added"},
		D: &map[string]string{"bar": "newVal", "safe": "added"},
	}

	destMap, err := destFoo.toMap()
	if err != nil {
		log.Fatal(err)
	}

	newFMap, err := newFoo.toMap()
	if err != nil {
		log.Fatal(err)
	}

	destMapOver, err := merge.Merge(destMap, newFMap, &merge.Options{Overwrite: true})
	if err != nil {
		log.Fatal(err)
	}

	merMap, err := merge.Merge(destMap, newFMap, &merge.Options{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\n\nOUTPUT:\n")
	fmt.Println("overwrite")
	indentMarshalPrint(destMapOver)

	fmt.Println("\n")

	fmt.Println("no overwrite")
	indentMarshalPrint(merMap)

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

func indentMarshalPrint(i interface{}) error {
	jBody, err := json.MarshalIndent(i, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(jBody))
	return nil
}
