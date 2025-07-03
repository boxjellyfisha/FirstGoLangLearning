package exercise

import (
	"reflect"
	"testing"
)

func TestDivideArray(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		k        int
		expected [][]int
	}{
		{
			name:     "正常情況 - 可以分組",
			nums:     []int{1, 3, 4, 8, 7, 9, 3, 5, 100},
			k:        2,
			expected: [][]int{},
		},
		{
			name:     "空陣列",
			nums:     []int{},
			k:        2,
			expected: [][]int{},
		},
		{
			name:     "陣列長度不是3的倍數",
			nums:     []int{1, 2, 3, 4},
			k:        2,
			expected: [][]int{},
		},
		{
			name:     "可以成功分組的情況",
			nums:     []int{1, 2, 3, 4, 5, 6},
			k:        2,
			expected: [][]int{{1, 2, 3}, {4, 5, 6}},
		},
		{
			name:     "相鄰元素差值超過k",
			nums:     []int{1, 10, 20, 30, 40, 50},
			k:        5,
			expected: [][]int{},
		},
		{
			name:     "單一元素組",
			nums:     []int{1, 2, 3},
			k:        1,
			expected: [][]int{{1, 2, 3}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := divideArray(tt.nums, tt.k)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("divideArray() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestInitArray(t *testing.T) {
	tests := []struct {
		name        string
		groupCount  int
		expectedLen int
	}{
		{
			name:        "初始化4組",
			groupCount:  4,
			expectedLen: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := initArray(tt.groupCount)
			if len(result) != tt.expectedLen {
				t.Errorf("initArray() 長度 = %v, want %v", len(result), tt.expectedLen)
			}
			// 檢查每個元素是否為 [3]int 類型
			for i, arr := range result {
				if len(arr) != 3 {
					t.Errorf("initArray()[%d] 長度 = %v, want 3", i, len(arr))
				}
			}
		})
	}
}

func TestInitSlice(t *testing.T) {
	tests := []struct {
		name        string
		groupCount  int
		expectedLen int
	}{
		{
			name:        "初始化3組",
			groupCount:  3,
			expectedLen: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := initSlice(tt.groupCount)
			if len(result) != tt.expectedLen {
				t.Errorf("initSlice() 長度 = %v, want %v", len(result), tt.expectedLen)
			}
			// 檢查每個子切片是否為長度3的切片
			for i, slice := range result {
				if len(slice) != 3 {
					t.Errorf("initSlice()[%d] 長度 = %v, want 3", i, len(slice))
				}
			}
		})
	}
}

func TestSwapItem(t *testing.T) {
	tests := []struct {
		name        string
		resource [][]int
		expectedFirstItem []int
		originalFirstItem []int
	}{
		{
			name:        "modify the first item and swap second and third item",
			originalFirstItem: []int{1, 2, 3},
			resource:  [][]int{ {1, 2, 3}, {4, 5, 6}, {7, 8, 9}},
			expectedFirstItem: []int{11, 12, 13},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var originalFirstItem = tt.originalFirstItem
			swapItem(tt.resource, tt.expectedFirstItem)	

			if reflect.DeepEqual(tt.resource[0], originalFirstItem) {
				t.Errorf("first item should be different from original first item = %v, want %v", tt.resource[0], originalFirstItem)
			}
			if !reflect.DeepEqual(tt.resource[0], tt.expectedFirstItem) {
				t.Errorf("after swap first item should be the same as expected first item = %v, want %v", tt.resource[0], tt.expectedFirstItem)
			}
			
			t.Logf("swapItem() = %v", tt.resource)
		})
	}
}

func BenchmarkDivideArray(b *testing.B) {
	nums := make([]int, 300000)
	for i := range nums {
		nums[i] = i
	}
	for i := 0; i < b.N; i++ {
		divideArray(nums, 3)
	}
}
