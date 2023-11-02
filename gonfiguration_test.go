package gonfiguration_test

import (
	"testing"

	"github.com/psyb0t/gonfiguration"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	type SampleConfig struct {
		Key1 string
		Key2 int
	}

	viper.Set("Key1", "test_key1")
	viper.Set("Key2", 42)

	err := gonfiguration.Parse("")
	require.Error(t, err)

	var sampleConfig SampleConfig
	err = gonfiguration.Parse(&sampleConfig)
	require.NoError(t, err)

	require.Equal(t, "test_key1", sampleConfig.Key1)
	require.Equal(t, 42, sampleConfig.Key2)
}
