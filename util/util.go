package util

import (
	"encoding/json"
	"fmt"
)

func MarshalIndentPrint(i interface{}) error {
	jBody, err := json.MarshalIndent(i, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(jBody))
	return nil
}
