package exercise

import (
	"fmt"
	"sort"
)

// a test sample for array
func initArray(groupCount int) [][3]int {
	arrayGrouping := make([][3]int, groupCount)

	if groupCount > 2 {
		fmt.Printf("arrayGrouping[0] %p\n", &arrayGrouping[0])
		fmt.Printf("arrayGrouping[0][0] %p\n", &arrayGrouping[0][0])
		fmt.Printf("arrayGrouping[0][1] %p\n", &arrayGrouping[0][1])
		fmt.Printf("arrayGrouping[0][2] %p\n", &arrayGrouping[0][2])
		fmt.Printf("arrayGrouping[1] %p\n", &arrayGrouping[1])
	}
	fmt.Println("declare in array ", arrayGrouping)

	return arrayGrouping
}

// a test sample for slice
func initSlice(groupCount int) [][]int {
	sliceGrouping := make([][]int, groupCount)
	for i := range sliceGrouping {
		sliceGrouping[i] = make([]int, 3)
	}

	if groupCount > 2 {
		fmt.Printf("sliceGrouping[0] %p\n", &sliceGrouping[0])
		fmt.Printf("sliceGrouping[0][0] %p\n", &sliceGrouping[0][0])
		fmt.Printf("sliceGrouping[0][1] %p\n", &sliceGrouping[0][1])
		fmt.Printf("sliceGrouping[0][2] %p\n", &sliceGrouping[0][2])
		fmt.Printf("sliceGrouping[1] %p\n", &sliceGrouping[1])
	}
	fmt.Println("declare in slice ", sliceGrouping)
	return sliceGrouping
}

func swapItem(items [][]int, newItem []int) {
	items[0] = newItem
	var tmp = items[1]
	items[1] = items[2]
	items[2] = tmp
}

func divideArray(nums []int, k int) [][]int {
	var groupSize = 3
	if len(nums) == 0 || len(nums)%groupSize != 0 {
		return [][]int{}
	}

	groupCount := len(nums) / groupSize
	sliceGrouping := make([][]int, groupCount)

	// go has sort library!
	sort.Ints(nums)

	for i := 0; i < groupCount; i++ {
		// directly assign the slice to the array is more efficient than appending
		sliceGrouping[i] = nums[i*groupSize : (i+1)*groupSize]

		if sliceGrouping[i][groupSize-1]-sliceGrouping[i][0] > k {
			return [][]int{}
		} else if sliceGrouping[i][groupSize-2]-sliceGrouping[i][0] > k {
			return [][]int{}
		}
	}
	return sliceGrouping
}

func RunLeetCode() {
	var resource = [...]int{1, 3, 4, 8, 7, 9, 3, 5, 100}
	result := divideArray(resource[:], 2)
	fmt.Println("result ", result)
}
