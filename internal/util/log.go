package util

import (
	"encoding/json"
	"fmt"
)

func LogStruct(s any) {
	marshal, err := json.Marshal(s)
	if err != nil {
		fmt.Printf("Error marshaling struct: %s\n", err)
		return
	}

	fmt.Printf("Struct is: %s\n", string(marshal))
}
