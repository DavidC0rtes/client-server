package utils

func FindAndDelete(haystack []string, needle string) []string {
	new := make([]string, len(haystack))
	index := 0
	for _, v := range haystack {
		if v != needle {
			new = append(new, v)
			index++
		}
	}
	return new[:index]
}
