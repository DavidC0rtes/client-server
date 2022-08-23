package utils

func TruncateFilename(filename string, limit int) []byte {
	temp := []byte(filename)
	if len(filename) >= limit {
		return temp[:limit]
	}
	return temp
}
