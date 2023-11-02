package gonfiguration

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestSetDefault(t *testing.T) {
	def := Default{
		Key:   "test_key",
		Value: "test_value",
	}

	SetDefault(def)
	require.Equal(t, "test_value", viper.Get("test_key"))
}

func TestSetDefaults(t *testing.T) {
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
