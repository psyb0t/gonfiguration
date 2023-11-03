package gonfiguration

import (
	"reflect"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestSetDefault(t *testing.T) {
	viper.Reset()
	defer viper.Reset()

	def := Default{
		Key:   "test_key",
		Value: "test_value",
	}

	SetDefault(def)
	require.Equal(t, "test_value", viper.Get("test_key"))
}

func TestSetDefaults(t *testing.T) {
	viper.Reset()
	defer viper.Reset()

	defaults := []Default{
		{
			Key:   "key1",
			Value: "value1",
		},
		{
			Key:   "key2",
			Value: "value2",
		},
	}

	SetDefaults(defaults...)

	require.Equal(t, "value1", viper.Get("key1"))
	require.Equal(t, "value2", viper.Get("key2"))
}

func TestIsSimpleType(t *testing.T) {
	viper.Reset()
	defer viper.Reset()

	testCases := []struct {
		name     string
		kind     reflect.Kind
		expected bool
	}{
		{"string", reflect.String, true},
		{"int", reflect.Int, true},
		{"int8", reflect.Int8, true},
		{"int16", reflect.Int16, true},
		{"int32", reflect.Int32, true},
		{"int64", reflect.Int64, true},
		{"uint", reflect.Uint, true},
		{"uint8", reflect.Uint8, true},
		{"uint16", reflect.Uint16, true},
		{"uint32", reflect.Uint32, true},
		{"uint64", reflect.Uint64, true},
		{"float32", reflect.Float32, true},
		{"float64", reflect.Float64, true},
		{"bool", reflect.Bool, true},
		{"struct", reflect.Struct, false},
		{"interface", reflect.Interface, false},
		{"array", reflect.Array, false},
		{"chan", reflect.Chan, false},
		{"func", reflect.Func, false},
		{"map", reflect.Map, false},
		{"ptr", reflect.Ptr, false},
		{"slice", reflect.Slice, false},
		{"unsafePointer", reflect.UnsafePointer, false},
		{"complex64", reflect.Complex64, false},
		{"complex128", reflect.Complex128, false},
		{"invalid", reflect.Invalid, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := isSimpleType(tc.kind)
			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestSetNonNilOnNonDefaultedFields(t *testing.T) {
	viper.Reset()
	defer viper.Reset()

	type config struct {
		Name    string  `mapstructure:"NAME"`
		Address string  `mapstructure:"ADDRESS"`
		Email   string  `mapstructure:"EMAIL"`
		Height  float64 `mapstructure:"HEIGHT"`
		Age     int     `mapstructure:"AGE"`
		Active  bool    `mapstructure:"ACTIVE"`
	}

	defaultName := Default{
		Key:   "NAME",
		Value: "default_name",
	}

	defaultAddress := Default{
		Key:   "ADDRESS",
		Value: "default_address",
	}

	SetDefaults(
		defaultName,
		defaultAddress,
	)

	reflectValue, err := getDestinationStructValue(&config{})
	require.Nil(t, err)

	err = setNonNilOnNonDefaultedFields(reflectValue)
	require.Nil(t, err)

	actualName := viper.Get(defaultName.Key)
	require.Equal(t, defaultName.Value, actualName)

	actualAddress := viper.Get(defaultAddress.Key)
	require.Equal(t, defaultAddress.Value, actualAddress)

	expectedEmail := ""
	require.Equal(t, expectedEmail, viper.Get("EMAIL"))

	expectedHeight := ""
	require.Equal(t, expectedHeight, viper.Get("HEIGHT"))

	expectedAge := ""
	require.Equal(t, expectedAge, viper.Get("AGE"))

	expectedActive := ""
	require.Equal(t, expectedActive, viper.Get("ACTIVE"))
}
