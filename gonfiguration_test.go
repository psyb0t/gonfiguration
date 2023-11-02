package gonfiguration_test

import (
	"testing"
	"time"

	"github.com/psyb0t/gonfiguration"
	"github.com/stretchr/testify/require"
)

//nolint:funlen
func TestParse(t *testing.T) {
	type SimpleConfig struct {
		Name   string  `mapstructure:"name"`
		Age    int     `mapstructure:"age"`
		Active bool    `mapstructure:"active"`
		Height float64 `mapstructure:"height"`
	}

	type ComplexConfig struct {
		Name string `mapstructure:"name"`
		Info struct {
			Age    int
			Height float64
		}
	}

	type PtrFieldConfig struct {
		Name *string `mapstructure:"name"`
		Age  *int    `mapstructure:"age"`
	}

	type TimeFieldConfig struct {
		Birthday time.Time `mapstructure:"birthday"`
	}

	nonStructConfig := "notastruct"

	testCases := []struct {
		name          string
		input         interface{}
		defaultConfig map[string]interface{}
		expected      interface{}
		expectError   bool
	}{
		{
			name:  "simple config",
			input: &SimpleConfig{},
			defaultConfig: map[string]interface{}{
				"name":   "cuc",
				"age":    79,
				"active": false,
				"height": 10.1,
			},
			expected: &SimpleConfig{
				Name:   "cuc",
				Age:    79,
				Active: false,
				Height: 10.1,
			},
			expectError: false,
		},
		{
			name:        "complex config with nested struct",
			input:       &ComplexConfig{},
			expectError: true,
		},
		{
			name:        "config with pointer fields",
			input:       &PtrFieldConfig{},
			expectError: true,
		},
		{
			name:        "config with time field",
			input:       &TimeFieldConfig{},
			expectError: true,
		},
		{
			name:        "non-pointer dest",
			input:       SimpleConfig{},
			expectError: true,
		},
		{
			name:        "non-struct dest",
			input:       &nonStructConfig,
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			for key, value := range tc.defaultConfig {
				gonfiguration.SetDefault(gonfiguration.Default{Key: key, Value: value})
			}

			err := gonfiguration.Parse(tc.input)

			if tc.expectError {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.expected, tc.input)
		})
	}
}
