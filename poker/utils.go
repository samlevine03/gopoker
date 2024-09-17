package poker

var Primes = [IntRanks]int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41}

// O(min(k, n-k)) implementation of binomial coefficient
func binomial(n, k int) int {
	if k > n {
		return 0
	}
	if k == 0 || k == n {
		return 1
	}
	k = min(k, n-k) // Take advantage of symmetry
	c := 1
	for i := 0; i < k; i++ {
		c = c * (n - i) / (i + 1)
	}
	return c
}

func Combinations(arr []uint32, size int) [][]uint32 {
	n := len(arr)
	result := make([][]uint32, 0, binomial(n, size))
	indices := make([]int, size)
	combination := make([]uint32, size)
	for i := 0; i < size; i++ {
		indices[i] = i
	}
	for {
		for i, index := range indices {
			combination[i] = arr[index]
		}

		// result = append(result, combination) // This is a bug
		// Since slices are references to the underlying array, we need to copy the slice
		// Otherwise, the slice will be overwritten in the next iteration

		combCopy := make([]uint32, size)
		copy(combCopy, combination)
		result = append(result, combCopy)

		i := size - 1
		for i >= 0 && indices[i] == len(arr)-size+i {
			i--
		}
		if i < 0 {
			return result
		}

		indices[i]++
		for j := i + 1; j < size; j++ {
			indices[j] = indices[j-1] + 1
		}
	}
}
