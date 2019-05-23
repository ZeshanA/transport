package math

// SliceMeanRounded takes a slice of integers and returns
// their mean as a truncated int. Empty slices return 0.
func SliceMeanRounded(nums []int) int {
	if len(nums) == 0 {
		return 0
	}
	total := 0
	for _, num := range nums {
		total += num
	}
	return int(total / len(nums))
}
