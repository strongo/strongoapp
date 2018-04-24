package slices

func MergeInt64s(source []int64, slices ...[]int64) (result []int64, changed bool) {
	result = source[:]
	for _, slice := range slices {
		for _, v := range slice {
			for _, r := range result {
				if r == v {
					goto found
				}
			}
			result = append(result, v)
			changed = true
		found:
		}
	}
	return
}

func MergeStrings(source []string, slices ...[]string) (result []string, changed bool) {
	result = source[:]
	for _, slice := range slices {
		for _, v := range slice {
			for _, r := range result {
				if r == v {
					goto found
				}
			}
			result = append(result, v)
			changed = true
		found:
		}
	}
	return
}
