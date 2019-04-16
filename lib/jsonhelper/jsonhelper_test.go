package jsonhelper_test

import (
	"testing"
	"transport/lib/jsonhelper"

	"github.com/stretchr/testify/assert"
)

func TestExtractNested(t *testing.T) {
	jsonString := `{"data": {"list": [{"id": "ID1"},{"id": "ID2"}]}}`
	// The path representing the property we want to extract
	pathToNestedProperty := "data.list.#.id"

	expected := []string{"ID1", "ID2"}
	actual := jsonhelper.ExtractNested(jsonString, pathToNestedProperty)

	assert.Equal(t, expected, actual)
}

func TestExtractNestedEmptyList(t *testing.T) {
	jsonString := `{"data": {"list": []}}`
	// The path representing the property we want to extract
	pathToNestedProperty := "data.list.#.id"

	// We expect a `nil` slice - no slice should have been allocated
	var expected []string
	actual := jsonhelper.ExtractNested(jsonString, pathToNestedProperty)

	assert.Equal(t, expected, actual)
}

func TestExtractNestedInvalidPath(t *testing.T) {
	jsonString := `{"data": {"list": []}}`
	// This path doesn't exist
	invalidPath := "data.invalidPath.#.id"

	// We expect a `nil` slice - no slice should have been allocated
	var expected []string
	actual := jsonhelper.ExtractNested(jsonString, invalidPath)

	assert.Equal(t, expected, actual)
}

func TestExtractNonArrayPath(t *testing.T) {
	// The array in the JSON has been replaced with an object
	jsonString := `{"data": {"list": { "item1": [{"id": "ID1"}], "item2": {"id": "ID2"} } }}`
	invalidPath := "data.list.#.id"

	// We expect a `nil` slice - no slice should have been allocated
	var expected []string
	actual := jsonhelper.ExtractNested(jsonString, invalidPath)

	assert.Equal(t, expected, actual)
}
