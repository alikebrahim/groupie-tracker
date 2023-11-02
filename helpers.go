package main

import (
	"net/http"
	"strconv"
	"strings"
)

func GetID(r *http.Request) int {
	path := r.URL.Path
	parts := strings.Split(path, "/")
	idStr := parts[len(parts)-1]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0
	}
	return id
}

func FormatText(input string) string {
	// Convert words to title case
	input = strings.ToTitle(input)

	// Replace underscores with spaces
	input = strings.ReplaceAll(input, "_", " ")

	// Replace hyphens with commas
	input = strings.ReplaceAll(input, "-", ", ")

	return input
}
