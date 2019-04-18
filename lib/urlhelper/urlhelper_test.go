package urlhelper_test

import (
	"testing"
	"transport/lib/urlhelper"

	"github.com/stretchr/testify/assert"
)

var inputs = []map[string]string{
	{"key1": "val1", "key2": "val2"},
	{"key1": "val1"},
	{"key_%1": "val%_1"},
	{},
}

var expectedOutputs = []string{
	"?key1=val1&key2=val2",
	"?key1=val1",
	"?key_%1=val%_1",
	"",
}

func TestBuildQueryString(t *testing.T) {
	for i, params := range inputs {
		assert.Equal(t, expectedOutputs[i], urlhelper.BuildQueryString(params))
	}
}
