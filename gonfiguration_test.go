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
		Name       string  `env:"NAME"`
		Age        int     `env:"AGE"`
		Active     bool    `env:"ACTIVE"`
		Height     float64 `env:"HEIGHT"`
		IsDead     bool    `env:"IS_DEAD"`
		OtherField string  `mapstructure:"OTHER_FIELD"`
	}

	type ComplexConfigInfo struct {
		Age    int
		Height float64
	}

	type ComplexConfig struct {
		Name string            `env:"NAME"`
		Info ComplexConfigInfo `env:"CUQUE"`
	}

	type PtrFieldConfig struct {
		Name *string `env:"NAME"`
		Age  *int    `env:"AGE"`
	}

	type TimeFieldConfig struct {
		Birthday time.Time `env:"BIRTHDAY"`
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
				"NAME":        "cuc",
				"AGE":         79,
				"ACTIVE":      true,
				"HEIGHT":      10.1,
				"IS_DEAD":     true,
				"OTHER_FIELD": "caras",
			},
			expected: &SimpleConfig{
				Name:   "cuc",
				Age:    79,
				Active: true,
				Height: 10.1,
				IsDead: true,
			},
			expectError: false,
		},
		{
			name:  "simple config vals from env",
			input: &SimpleConfig{},
			envConfig: map[string]interface{}{
				"NAME":        "cuc",
				"AGE":         79,
				"ACTIVE":      true,
				"HEIGHT":      10.1,
				"OTHER_FIELD": "gibelio",
			},
			expected: &SimpleConfig{
				Name:   "cuc",
				Age:    79,
				Active: true,
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
			defer gonfiguration.Reset()

			for key, value := range tc.defaultConfig {
				gonfiguration.SetDefault(key, value)
			}

			for key, value := range tc.envConfig {
				t.Logf("setting env var: %s=%v", key, value)
				t.Setenv(key, fmt.Sprintf("%v", value))
			}

			for k, v := range gonfiguration.GetAllValues() {
				t.Logf("%s=%v", k, v)
			}

			err := gonfiguration.Parse(tc.input)

			if tc.expectError {
				require.NotNil(t, err)

				return
			}

			require.Nil(t, err)
			require.Equal(t, tc.expected, tc.input)
		})
	}
}
