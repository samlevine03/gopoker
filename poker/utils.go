package poker

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

// Combinations generates all combinations of size k from the input array (of size n)
func Combinations(arr []uint32, size int) [][]uint32 {
	n := len(arr)
	result := make([][]uint32, 0, binomial(n, size))
	indices := make([]int, size)
	for i := range indices {
		indices[i] = i
	}

	for {
		combination := make([]uint32, size)
		for i, index := range indices {
			combination[i] = arr[index]
		}
		result = append(result, combination)

		// Generate next combination
		i := size - 1
		for i >= 0 && indices[i] == n-size+i {
			i--
		}
		if i < 0 {
			break
		}
		indices[i]++
		for j := i + 1; j < size; j++ {
			indices[j] = indices[j-1] + 1
		}
	}

	return result
}
