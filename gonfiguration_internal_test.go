package gonfiguration

import (
	"reflect"
	"testing"
	"time"

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

func TestUnsupportedFieldType(t *testing.T) {
	type UnsupportedStruct struct {
		MapField map[string]string `env:"MAP_FIELD"`
	}

	envVars := map[string]string{
		"MAP_FIELD": "test",
	}

	dst := UnsupportedStruct{}
	err := parseDstFields(reflect.ValueOf(&dst).Elem(), envVars)
	
	require.Error(t, err)
	require.Contains(t, err.Error(), "Type not supported")
}

func TestSetEnvVarValueUnsupportedType(t *testing.T) {
	envVal := "test_value"
	
	mapValue := reflect.ValueOf(make(map[string]string))
	err := setEnvVarValue(mapValue, envVal)
	
	require.Error(t, err)
	require.Contains(t, err.Error(), "Unsupported field type")
}

func TestIsSupportedType(t *testing.T) {
	supportedTypes := []reflect.Value{
		reflect.ValueOf(""),
		reflect.ValueOf(int(0)),
		reflect.ValueOf(int8(0)),
		reflect.ValueOf(int16(0)),
		reflect.ValueOf(int32(0)),
		reflect.ValueOf(int64(0)),
		reflect.ValueOf(uint(0)),
		reflect.ValueOf(uint8(0)),
		reflect.ValueOf(uint16(0)),
		reflect.ValueOf(uint32(0)),
		reflect.ValueOf(uint64(0)),
		reflect.ValueOf(float32(0)),
		reflect.ValueOf(float64(0)),
		reflect.ValueOf(false),
		reflect.ValueOf(time.Duration(0)), // Add time.Duration test
		reflect.ValueOf([]string{}),       // Add []string test
	}

	for _, val := range supportedTypes {
		require.True(t, isSupportedType(val), "Expected %v to be supported", val.Type())
	}

	unsupportedTypes := []reflect.Value{
		reflect.ValueOf(make(map[string]string)),
		reflect.ValueOf([1]string{}),
		reflect.ValueOf(struct{}{}),
		reflect.ValueOf(&struct{}{}),
		reflect.ValueOf(make(chan int)),
		reflect.ValueOf(func() {}),
	}

	for _, val := range unsupportedTypes {
		require.False(t, isSupportedType(val), "Expected %v to be unsupported", val.Type())
	}
}
