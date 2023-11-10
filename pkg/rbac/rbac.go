package rbac

func IntToBinary(n, len int) []int {
	binary := make([]int, len)
	for i := 0; i < len; i++ {
		binary[i] = n % 2
		n = n / 2
	}

	return binary
}
