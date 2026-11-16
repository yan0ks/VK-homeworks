package session3

func SumRangeNaive(nums []int, left, right int) int {
	s := 0
	for i := left; i <= right; i++ {
		s += nums[i]
	}
	return s
}

type NumArray struct {
	prefix []int
}

func Constructor(nums []int) NumArray {
	prefix := make([]int, len(nums)+1)
	for index, number := range nums {
		prefix[index+1] = prefix[index] + number
	}
	return NumArray{prefix: prefix}
}

func (na *NumArray) SumRange(left, right int) int {
	return na.prefix[right+1] - na.prefix[left]
}
