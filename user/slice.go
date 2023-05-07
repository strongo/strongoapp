package user

// RemoveInPlace removes an item from a slice in place without new memory allocation
func removeInPlace[T comparable](v T, slice []T) []T {
	shift := 0
	for i, s := range slice {
		if s == v {
			shift++
			continue
		}
		slice[i-shift] = s
	}
	return slice[:len(slice)-shift]
}
