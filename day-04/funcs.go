package main

import "errors"

func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("Can't divide by 0")
	}

	return a / b, nil

}

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
