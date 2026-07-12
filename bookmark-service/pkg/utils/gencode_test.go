package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenCode_Generate(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string 

		expectedLen int
		expectedError error
	}{
		{
			name: "success",
			expectedLen: 12,
			expectedError: nil,
		},
		{
			name: "success with custom length",
			expectedLen: 10000,
			expectedError: nil,
		},
		{
			name: "success with custom length",
			expectedLen: 1,
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			codeGen := NewGenCode()

			code, err := codeGen.Generate(tc.expectedLen)

			assert.ErrorIs(t, err, tc.expectedError)
			assert.Equal(t, tc.expectedLen, len(code))
		})
	}
}