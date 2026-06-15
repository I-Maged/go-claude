package main

func minMax(nums []int) (min, max int) {
	min, max = nums[0], nums[0]
	for _, v := range nums {
		if min < v {
			min = v
		}
		if max > v {
			max = v
		}
	}
	return
}
