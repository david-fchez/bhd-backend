package app

import "C"
import "encoding/json"

type ApiReturnStruct struct {
	ErrorID          int    `json:"errorId"`
	ErrorDescription string `json:"errorDescription"`
	Content          string `json:"content"`
}

// ToJsonString converts ApiReturnStruct to json string
func (ars ApiReturnStruct) ToJsonString() string {
	content, err := json.MarshalIndent(ars, "", "    ")
	if err != nil {
		return "{errorId:\"9999\",errorDescription:\"Failed to marshal ApiReturnStruct to json\"}"
	}
	return string(content)
}

func (ars ApiReturnStruct) IsSuccess() bool {
	return ars.ErrorID == 0
}

func FromJsonString(jsonStr string) ApiReturnStruct {
	var retVal ApiReturnStruct
	json.Unmarshal([]byte(jsonStr), &retVal)
	return retVal
}

func FromJsonCString(str *C.char) ApiReturnStruct {
	return FromJsonString(FromCString(str))
}
