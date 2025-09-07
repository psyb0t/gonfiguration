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
	supportedTypes := []reflect.Kind{
		reflect.String,
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Float32,
		reflect.Float64,
		reflect.Bool,
	}

	for _, kind := range supportedTypes {
		require.True(t, isSupportedType(kind), "Expected %v to be supported", kind)
	}

	unsupportedTypes := []reflect.Kind{
		reflect.Map,
		reflect.Slice,
		reflect.Array,
		reflect.Struct,
		reflect.Ptr,
		reflect.Interface,
		reflect.Chan,
		reflect.Func,
	}

	for _, kind := range unsupportedTypes {
		require.False(t, isSupportedType(kind), "Expected %v to be unsupported", kind)
	}
}
