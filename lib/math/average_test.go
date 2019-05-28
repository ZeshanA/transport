package math

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSliceMeanRounded(t *testing.T) {
	inputSlices := [][]int{
		{1, 2, 3, 4, 5},
		{6, 7},
		{8, 9, 10},
		{11},
		{},
	}
	expectedAverages := []int{3, 6, 9, 11, 0}
	for i, slice := range inputSlices {
		calculatedAverage := SliceMeanRounded(slice)
		assert.Equal(t, expectedAverages[i], calculatedAverage)
	}
}
