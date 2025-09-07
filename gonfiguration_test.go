package gonfiguration_test

import (
	"fmt"
	"os"
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
		input         any
		defaultConfig map[string]any
		envConfig     map[string]any
		expected      any
		expectError   bool
	}{
		{
			name:  "simple config vals from defaults",
			input: &SimpleConfig{},
			defaultConfig: map[string]any{
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
			envConfig: map[string]any{
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
			defaultConfig: map[string]any{
				"NAME":   "cuc",
				"AGE":    79,
				"ACTIVE": false,
				"HEIGHT": 10.1,
			},
			envConfig: map[string]any{
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

func TestSetDefaults(t *testing.T) {
	defer gonfiguration.Reset()

	defaults := map[string]any{
		"TEST_KEY1": "value1",
		"TEST_KEY2": 42,
		"TEST_KEY3": true,
	}

	gonfiguration.SetDefaults(defaults)

	allDefaults := gonfiguration.GetDefaults()
	require.Equal(t, "value1", allDefaults["TEST_KEY1"])
	require.Equal(t, 42, allDefaults["TEST_KEY2"])
	require.Equal(t, true, allDefaults["TEST_KEY3"])
}

func TestGetDefaults(t *testing.T) {
	defer gonfiguration.Reset()

	gonfiguration.SetDefault("TEST_KEY", "test_value")

	defaults := gonfiguration.GetDefaults()
	require.Equal(t, "test_value", defaults["TEST_KEY"])
}

func TestGetEnvVars(t *testing.T) {
	defer gonfiguration.Reset()

	t.Setenv("TEST_ENV_VAR", "test_value")
	
	type TestConfig struct {
		TestValue string `env:"TEST_ENV_VAR"`
	}
	
	config := &TestConfig{}
	err := gonfiguration.Parse(config)
	require.NoError(t, err)
	
	envVars := gonfiguration.GetEnvVars()
	require.Equal(t, "test_value", envVars["TEST_ENV_VAR"])
}

func TestParseWithUintValues(t *testing.T) {
	defer gonfiguration.Reset()

	type UintConfig struct {
		UintVal   uint   `env:"UINT_VAL"`
		Uint8Val  uint8  `env:"UINT8_VAL"`
		Uint16Val uint16 `env:"UINT16_VAL"`
		Uint32Val uint32 `env:"UINT32_VAL"`
		Uint64Val uint64 `env:"UINT64_VAL"`
	}

	t.Setenv("UINT_VAL", "42")
	t.Setenv("UINT8_VAL", "255")
	t.Setenv("UINT16_VAL", "65535")
	t.Setenv("UINT32_VAL", "4294967295")
	t.Setenv("UINT64_VAL", "18446744073709551615")

	config := &UintConfig{}
	err := gonfiguration.Parse(config)
	require.NoError(t, err)

	require.Equal(t, uint(42), config.UintVal)
	require.Equal(t, uint8(255), config.Uint8Val)
	require.Equal(t, uint16(65535), config.Uint16Val)
	require.Equal(t, uint32(4294967295), config.Uint32Val)
	require.Equal(t, uint64(18446744073709551615), config.Uint64Val)
}

func TestParseWithFloatValues(t *testing.T) {
	defer gonfiguration.Reset()

	type FloatConfig struct {
		Float32Val float32 `env:"FLOAT32_VAL"`
		Float64Val float64 `env:"FLOAT64_VAL"`
	}

	t.Setenv("FLOAT32_VAL", "3.14")
	t.Setenv("FLOAT64_VAL", "2.718281828459045")

	config := &FloatConfig{}
	err := gonfiguration.Parse(config)
	require.NoError(t, err)

	require.InDelta(t, 3.14, config.Float32Val, 0.001)
	require.InDelta(t, 2.718281828459045, config.Float64Val, 0.000000000000001)
}

func TestParseWithBoolValues(t *testing.T) {
	defer gonfiguration.Reset()

	type BoolConfig struct {
		TrueVal  bool `env:"TRUE_VAL"`
		FalseVal bool `env:"FALSE_VAL"`
	}

	t.Setenv("TRUE_VAL", "true")
	t.Setenv("FALSE_VAL", "false")

	config := &BoolConfig{}
	err := gonfiguration.Parse(config)
	require.NoError(t, err)

	require.True(t, config.TrueVal)
	require.False(t, config.FalseVal)
}

func TestParseErrorCases(t *testing.T) {
	defer gonfiguration.Reset()

	testCases := []struct {
		name    string
		envVars map[string]string
		config  any
		errorContains string
	}{
		{
			name: "invalid uint",
			envVars: map[string]string{
				"UINT_VAL": "invalid",
			},
			config: &struct {
				UintVal uint `env:"UINT_VAL"`
			}{},
			errorContains: "Failed to parse uint",
		},
		{
			name: "invalid float",
			envVars: map[string]string{
				"FLOAT_VAL": "invalid",
			},
			config: &struct {
				FloatVal float64 `env:"FLOAT_VAL"`
			}{},
			errorContains: "Failed to parse float",
		},
		{
			name: "invalid bool",
			envVars: map[string]string{
				"BOOL_VAL": "invalid",
			},
			config: &struct {
				BoolVal bool `env:"BOOL_VAL"`
			}{},
			errorContains: "Failed to parse bool",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			defer gonfiguration.Reset()

			for key, value := range tc.envVars {
				t.Setenv(key, value)
			}

			err := gonfiguration.Parse(tc.config)
			require.Error(t, err)
			require.Contains(t, err.Error(), tc.errorContains)
		})
	}
}

func TestGetAllValuesWithEmptyState(t *testing.T) {
	defer gonfiguration.Reset()

	allValues := gonfiguration.GetAllValues()
	require.Empty(t, allValues)
}

func TestGetAllValuesWithDefaultsAndEnvVars(t *testing.T) {
	defer gonfiguration.Reset()

	gonfiguration.SetDefault("DEFAULT_KEY", "default_value")
	
	t.Setenv("ENV_KEY", "env_value")
	t.Setenv("SHARED_KEY", "env_overrides_default")
	
	gonfiguration.SetDefault("SHARED_KEY", "default_value")

	type TestConfig struct {
		EnvKey    string `env:"ENV_KEY"`
		SharedKey string `env:"SHARED_KEY"`
	}
	
	config := &TestConfig{}
	err := gonfiguration.Parse(config)
	require.NoError(t, err)

	allValues := gonfiguration.GetAllValues()
	
	require.Equal(t, "default_value", allValues["DEFAULT_KEY"])
	require.Equal(t, "env_value", allValues["ENV_KEY"])
	require.Equal(t, "env_overrides_default", allValues["SHARED_KEY"])
}

func TestDefaultValueTypeMismatch(t *testing.T) {
	defer gonfiguration.Reset()

	gonfiguration.SetDefault("MISMATCH_KEY", "string_value")

	type MismatchConfig struct {
		IntVal int `env:"MISMATCH_KEY"`
	}

	config := &MismatchConfig{}
	err := gonfiguration.Parse(config)
	require.Error(t, err)
	require.Contains(t, err.Error(), "Default value type mismatch")
}

func TestInvalidEnvVariable(t *testing.T) {
	defer gonfiguration.Reset()
	
	originalEnv := os.Environ()
	defer func() {
		os.Clearenv()
		for _, env := range originalEnv {
			if key, val, found := splitEnvVar(env); found {
				os.Setenv(key, val)
			}
		}
	}()

	os.Clearenv()
	os.Setenv("VALID_KEY", "valid_value")
	
	type SimpleConfig struct {
		ValidKey string `env:"VALID_KEY"`
	}

	config := &SimpleConfig{}
	err := gonfiguration.Parse(config)
	require.NoError(t, err)
	require.Equal(t, "valid_value", config.ValidKey)
}

func splitEnvVar(env string) (key, val string, found bool) {
	eq := fmt.Sprintf("%s", "=")
	if idx := fmt.Sprintf("%s", env); len(idx) > 0 {
		parts := []string{}
		current := ""
		for i, char := range env {
			if string(char) == eq && len(parts) == 0 {
				parts = append(parts, current)
				current = env[i+1:]
				break
			}
			current += string(char)
		}
		if len(current) > 0 {
			parts = append(parts, current)
		}
		if len(parts) == 2 {
			return parts[0], parts[1], true
		}
	}
	return "", "", false
}
