package util

import (
	"encoding/json"
	"fmt"
)

func CopyMap(m map[string]interface{}) map[string]interface{} {
	newMap := map[string]interface{}{}
	for k, v := range m {
		newMap[k] = v
	}

	return newMap
}

func MarshalIndentPrint(i interface{}) error {
	jBody, err := json.MarshalIndent(i, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(jBody))
	return nil
}
