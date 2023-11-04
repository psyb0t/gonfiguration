package gonfiguration

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseDstFields(t *testing.T) {
	type EnvTestStruct struct {
		StringField string `env:"STRING_FIELD"`
		IntField    int    `env:"INT_FIELD"`
	}

	tests := []struct {
		name          string
		envVars       map[string]string
		expected      EnvTestStruct
		expectError   bool
		errorContains string
	}{
		{
			name: "all fields",
			envVars: map[string]string{
				"STRING_FIELD": "test",
				"INT_FIELD":    "42",
			},
			expected: EnvTestStruct{
				StringField: "test",
				IntField:    42,
			},
			expectError: false,
		},
		{
			name: "wrong int format",
			envVars: map[string]string{
				"INT_FIELD": "not_an_int",
			},
			expectError:   true,
			errorContains: "Failed to parse int",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dst := EnvTestStruct{}
			err := parseDstFields(reflect.ValueOf(&dst).Elem(), tt.envVars)

			if tt.expectError {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.errorContains)

				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.expected.StringField, dst.StringField)
			require.Equal(t, tt.expected.IntField, dst.IntField)
		})
	}
}
