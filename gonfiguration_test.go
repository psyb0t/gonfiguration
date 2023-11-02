package gonfiguration_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/psyb0t/gonfiguration"
	"github.com/stretchr/testify/require"
)

//nolint:funlen
func TestParse(t *testing.T) {
	type SimpleConfig struct {
		Name   string  `mapstructure:"NAME"`
		Age    int     `mapstructure:"AGE"`
		Active bool    `mapstructure:"ACTIVE"`
		Height float64 `mapstructure:"HEIGHT"`
	}

	type ComplexConfig struct {
		Name string `mapstructure:"NAME"`
		Info struct {
			Age    int
			Height float64
		}
	}

	type PtrFieldConfig struct {
		Name *string `mapstructure:"NAME"`
		Age  *int    `mapstructure:"AGE"`
	}

	type TimeFieldConfig struct {
		Birthday time.Time `mapstructure:"BIRTHDAY"`
	}

	nonStructConfig := "notastruct"

	testCases := []struct {
		name          string
		input         interface{}
		defaultConfig map[string]interface{}
		envConfig     map[string]interface{}
		expected      interface{}
		expectError   bool
	}{
		{
			name:  "simple config vals from defaults",
			input: &SimpleConfig{},
			defaultConfig: map[string]interface{}{
				"NAME":   "cuc",
				"AGE":    79,
				"ACTIVE": false,
				"HEIGHT": 10.1,
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
			name:  "simple config vals from env",
			input: &SimpleConfig{},
			envConfig: map[string]interface{}{
				"NAME":   "cuc",
				"AGE":    79,
				"ACTIVE": false,
				"HEIGHT": 10.1,
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
			name:  "simple config with defaults and vals from env",
			input: &SimpleConfig{},
			defaultConfig: map[string]interface{}{
				"NAME":   "cuc",
				"AGE":    79,
				"ACTIVE": false,
				"HEIGHT": 10.1,
			},
			envConfig: map[string]interface{}{
				"NAME":   "cucu",
				"AGE":    22,
				"ACTIVE": true,
				"HEIGHT": 5.4,
			},
			expected: &SimpleConfig{
				Name:   "cucu",
				Age:    22,
				Active: true,
				Height: 5.4,
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

			for key, value := range tc.envConfig {
				t.Setenv(key, fmt.Sprintf("%v", value))
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
