package gonfiguration_test

import (
	"fmt"
	"os"
	"sync"
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

	type DurationFieldConfig struct {
		Timeout    time.Duration `env:"TIMEOUT"`
		RetryDelay time.Duration `env:"RETRY_DELAY"`
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
			name:  "duration config with defaults and env vars",
			input: &DurationFieldConfig{},
			defaultConfig: map[string]any{
				"TIMEOUT":     30 * time.Second,
				"RETRY_DELAY": time.Minute,
			},
			envConfig: map[string]any{
				"TIMEOUT":     "45s",
				"RETRY_DELAY": "2m30s",
			},
			expected: &DurationFieldConfig{
				Timeout:    45 * time.Second,
				RetryDelay: 2*time.Minute + 30*time.Second,
			},
			expectError: false,
		},
		{
			name:  "duration config with only env vars",
			input: &DurationFieldConfig{},
			envConfig: map[string]any{
				"TIMEOUT":     "1h",
				"RETRY_DELAY": "500ms",
			},
			expected: &DurationFieldConfig{
				Timeout:    time.Hour,
				RetryDelay: 500 * time.Millisecond,
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
		name          string
		envVars       map[string]string
		config        any
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

func TestTimeDurationDefaults(t *testing.T) {
	defer gonfiguration.Reset()

	type TimingConfig struct {
		Timeout    time.Duration `env:"TIMEOUT"`
		RetryDelay time.Duration `env:"RETRY_DELAY"`
	}

	// Set time.Duration defaults
	gonfiguration.SetDefaults(map[string]any{
		"TIMEOUT":     30 * time.Second,
		"RETRY_DELAY": 5 * time.Minute,
	})

	cfg := TimingConfig{}
	require.NoError(t, gonfiguration.Parse(&cfg))

	// Should use defaults
	require.Equal(t, 30*time.Second, cfg.Timeout)
	require.Equal(t, 5*time.Minute, cfg.RetryDelay)

	// Override with env vars (string format)
	t.Setenv("TIMEOUT", "1m30s")
	t.Setenv("RETRY_DELAY", "2h")

	cfg2 := TimingConfig{}
	require.NoError(t, gonfiguration.Parse(&cfg2))

	// Should use env var values
	require.Equal(t, 1*time.Minute+30*time.Second, cfg2.Timeout)
	require.Equal(t, 2*time.Hour, cfg2.RetryDelay)
}

func TestStringSliceDefaults(t *testing.T) {
	defer gonfiguration.Reset()

	type SliceConfig struct {
		Tags     []string `env:"TAGS"`
		Features []string `env:"FEATURES"`
	}

	// Set []string defaults
	gonfiguration.SetDefaults(map[string]any{
		"TAGS":     []string{"default", "production"},
		"FEATURES": []string{"auth", "logging"},
	})

	cfg := SliceConfig{}
	require.NoError(t, gonfiguration.Parse(&cfg))

	// Should use defaults
	require.Equal(t, []string{"default", "production"}, cfg.Tags)
	require.Equal(t, []string{"auth", "logging"}, cfg.Features)

	// Override with env vars (comma-separated format)
	t.Setenv("TAGS", "test,staging,prod")
	t.Setenv("FEATURES", "auth, metrics , caching")

	cfg2 := SliceConfig{}
	require.NoError(t, gonfiguration.Parse(&cfg2))

	// Should use env var values (with whitespace trimmed)
	require.Equal(t, []string{"test", "staging", "prod"}, cfg2.Tags)
	require.Equal(t, []string{"auth", "metrics", "caching"}, cfg2.Features)
}

func TestParseStringSlice(t *testing.T) {
	type StringSliceConfig struct {
		Tags     []string `env:"TAGS"`
		Features []string `env:"FEATURES"`
		Empty    []string `env:"EMPTY"`
	}

	gonfiguration.Reset()

	// Set environment variables
	t.Setenv("TAGS", "go,rust,python")
	t.Setenv("FEATURES", "feature1, feature2 , feature3")
	t.Setenv("EMPTY", "")

	cfg := StringSliceConfig{}
	require.NoError(t, gonfiguration.Parse(&cfg))

	// Test comma-separated values
	require.Equal(t, []string{"go", "rust", "python"}, cfg.Tags)

	// Test comma-separated values with spaces (should be trimmed)
	require.Equal(t, []string{"feature1", "feature2", "feature3"}, cfg.Features)

	// Test empty string (should result in empty slice)
	require.Equal(t, []string{}, cfg.Empty)
}

func TestConcurrentAccess(t *testing.T) {
	defer gonfiguration.Reset()

	type ConcurrentConfig struct {
		Value1 string `env:"VALUE1"`
		Value2 int    `env:"VALUE2"`
		Value3 bool   `env:"VALUE3"`
	}

	// Set some initial defaults
	gonfiguration.SetDefaults(map[string]any{
		"VALUE1": "default1",
		"VALUE2": 42,
		"VALUE3": false,
	})

	const numGoroutines = 50
	const numOperations = 100

	// Test concurrent Parse operations
	t.Run("ConcurrentParse", func(t *testing.T) {
		t.Setenv("VALUE1", "test")
		t.Setenv("VALUE2", "123")
		t.Setenv("VALUE3", "true")

		var wg sync.WaitGroup
		results := make([]ConcurrentConfig, numGoroutines)

		for i := range numGoroutines {
			wg.Add(1)
			go func(index int) {
				defer wg.Done()
				cfg := ConcurrentConfig{}
				require.NoError(t, gonfiguration.Parse(&cfg))
				results[index] = cfg
			}(i)
		}

		wg.Wait()

		// All results should be identical
		expected := ConcurrentConfig{
			Value1: "test",
			Value2: 123,
			Value3: true,
		}
		for i, result := range results {
			require.Equal(t, expected, result, "goroutine %d result mismatch", i)
		}
	})

	// Test concurrent SetDefault operations
	t.Run("ConcurrentSetDefaults", func(t *testing.T) {
		var wg sync.WaitGroup

		for i := range numGoroutines {
			wg.Add(1)
			go func(index int) {
				defer wg.Done()
				for j := range numOperations {
					key := fmt.Sprintf("CONCURRENT_KEY_%d_%d", index, j)
					value := fmt.Sprintf("value_%d_%d", index, j)
					gonfiguration.SetDefault(key, value)
				}
			}(i)
		}

		wg.Wait()

		// Verify all defaults were set correctly
		defaults := gonfiguration.GetDefaults()
		for i := range numGoroutines {
			for j := range numOperations {
				key := fmt.Sprintf("CONCURRENT_KEY_%d_%d", i, j)
				expectedValue := fmt.Sprintf("value_%d_%d", i, j)
				actualValue, exists := defaults[key]
				require.True(t, exists, "key %s should exist", key)
				require.Equal(t, expectedValue, actualValue, "value mismatch for key %s", key)
			}
		}
	})

	// Test concurrent mixed operations
	t.Run("ConcurrentMixedOperations", func(t *testing.T) {
		t.Setenv("MIXED_VALUE", "env_value")

		type MixedConfig struct {
			MixedValue string `env:"MIXED_VALUE"`
		}

		var wg sync.WaitGroup

		// Concurrent Parse operations
		for range numGoroutines / 2 {
			wg.Add(1)
			go func() {
				defer wg.Done()
				cfg := MixedConfig{}
				require.NoError(t, gonfiguration.Parse(&cfg))
				require.Equal(t, "env_value", cfg.MixedValue)
			}()
		}

		// Concurrent GetAllValues operations
		for range numGoroutines / 4 {
			wg.Add(1)
			go func() {
				defer wg.Done()
				values := gonfiguration.GetAllValues()
				require.NotNil(t, values)
			}()
		}

		// Concurrent GetDefaults operations
		for range numGoroutines / 4 {
			wg.Add(1)
			go func() {
				defer wg.Done()
				defaults := gonfiguration.GetDefaults()
				require.NotNil(t, defaults)
			}()
		}

		wg.Wait()
	})
}
