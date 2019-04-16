package jsonhelper

import "github.com/tidwall/gjson"

func ExtractNested(json string, path string) []string {
	nestedProperties := gjson.Get(json, path).Array()
	return resultArrayToStringArray(nestedProperties)
}

func resultArrayToStringArray(resultArray []gjson.Result) []string {
	if len(resultArray) == 0 {
		return nil
	}
	stringArray := make([]string, len(resultArray))
	for i, result := range resultArray {
		stringArray[i] = result.String()
	}
	return stringArray
}
