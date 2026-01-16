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

	t.Run("all fields", func(t *testing.T) {
		envVars := map[string]string{
			"STRING_FIELD": "test",
			"INT_FIELD":    "42",
		}

		dst := EnvTestStruct{}
		err := parseDstFields(reflect.ValueOf(&dst).Elem(), envVars)

		require.NoError(t, err)
		require.Equal(t, "test", dst.StringField)
		require.Equal(t, 42, dst.IntField)
	})

	t.Run("wrong int format", func(t *testing.T) {
		envVars := map[string]string{
			"INT_FIELD": "not_an_int",
		}

		dst := EnvTestStruct{}
		err := parseDstFields(reflect.ValueOf(&dst).Elem(), envVars)

		require.Error(t, err)
	})
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
	require.ErrorIs(t, err, ErrUnsupportedFieldType)
}

func TestSetEnvVarValueUnsupportedType(t *testing.T) {
	envVal := "test_value"

	mapValue := reflect.ValueOf(make(map[string]string))
	err := setEnvVarValue(mapValue, envVal)

	require.Error(t, err)
	require.ErrorIs(t, err, ErrUnsupportedFieldType)
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
		reflect.ValueOf(time.Duration(0)),
		reflect.ValueOf([]string{}),
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
