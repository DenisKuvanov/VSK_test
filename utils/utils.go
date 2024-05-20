package utils

import (
	"strconv"
	"strings"
)

// GetIdFromUrl extracts the ID from a URL string and returns it as an integer.
//
// Parameters:
// - url: a string representing the URL from which to extract the ID.
//
// Returns:
// - int: the extracted ID as an integer.
// - error: an error if the ID cannot be extracted or converted to an integer.
func GetIdFromUrl(url string) (int, error) {
	idStr := strings.Split(url, "/")[len(strings.Split(url, "/"))-1]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, err
	}
	return id, nil
}
