package gonfiguration

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Default struct {
	Key   string
	Value interface{}
}

func SetDefault(def Default) {
	viper.SetDefault(def.Key, def.Value)
}

func SetDefaults(defaults ...Default) {
	for _, def := range defaults {
		SetDefault(def)
	}
}

func Parse(destination interface{}) error {
	viper.AutomaticEnv()

	if err := viper.Unmarshal(destination); err != nil {
		return errors.Wrap(ErrParsingConfig, err.Error())
	}

	return nil
}
