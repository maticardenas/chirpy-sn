package chirptext

import "testing"

func TestReplaceChirpInput(t *testing.T) {
	// slice with test cases
	testCases := []struct {
		input    string
		expected string
	}{
		{
			input:    "Hello kerfuffle",
			expected: "Hello ****",
		},
		{
			input:    "Hello kerfuffle!",
			expected: "Hello kerfuffle!",
		},
		{
			input:    "Hello sharbert",
			expected: "Hello ****",
		},
		{
			input:    "Hello kerfuffle, sharbert and fornax",
			expected: "Hello kerfuffle, **** and ****",
		},
	}

	for _, tc := range testCases {
		result := replaceChirpInput(tc.input)
		if result != tc.expected {
			t.Errorf("Expected %q, but got %q", tc.expected, result)
		}
	}
}
