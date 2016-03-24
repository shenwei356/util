package stringutil

import "github.com/shenwei356/util/byteutil"

// Split splits a byte slice by giveen letters
func Split(slice string, letters string) []string {
	result := byteutil.Split([]byte(slice), []byte(letters))
	result2 := []string{}
	for _, s := range result {
		result2 = append(result2, string(s))
	}
	return result2
}
