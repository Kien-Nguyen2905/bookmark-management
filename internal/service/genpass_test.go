package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGeneratePassword(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		length      int
		expectedErr error
	}{
		{name: "normal case - valid length", length: 10},
		{name: "normal case - min length", length: 6},
		{name: "normal case - max length", length: 64},
		{name: "invalid length - too short", length: 4, expectedErr: ErrInvalidLength},
		{name: "invalid length - too long", length: 70, expectedErr: ErrInvalidLength},
		{name: "invalid length - zero", length: 0, expectedErr: ErrInvalidLength},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			svc := NewGenPass()

			password, err := svc.GeneratePassword(tc.length)

			if tc.expectedErr != nil {
				assert.ErrorIs(t, err, tc.expectedErr)
				return
			}

			assert.NoError(t, err)
			assert.Len(t, password, tc.length)
			assert.Regexp(t, `^[a-zA-Z0-9]+$`, password)
		})
	}
}
