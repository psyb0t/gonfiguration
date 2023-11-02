package gonfiguration

import (
	"reflect"

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
	destValue := reflect.ValueOf(destination)
	if destValue.Kind() != reflect.Ptr {
		return ErrTargetNotPointer
	}

	// Dereferencing to get the actual value and checking if it's a struct.
	destValue = destValue.Elem()
	if destValue.Kind() != reflect.Struct {
		return ErrDestinationNotStruct
	}

	if err := processSimpleStructFields(destValue); err != nil {
		return err
	}

	viper.AutomaticEnv()

	if err := viper.Unmarshal(destination); err != nil {
		return errors.Wrap(ErrParsingConfig, err.Error())
	}

	return nil
}

func isSimpleType(kind reflect.Kind) bool {
	switch kind { //nolint:exhaustive
	case reflect.String,
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
		reflect.Bool:
		return true
	default:
		return false
	}
}

func processSimpleStructFields(value reflect.Value) error {
	for i := 0; i < value.NumField(); i++ {
		field := value.Type().Field(i)
		fieldValue := value.Field(i)

		if !isSimpleType(fieldValue.Kind()) {
			return errors.Wrap(ErrFieldIsNotOfSimpleType, field.Name)
		}

		mapstructureTag := field.Tag.Get("mapstructure")
		if mapstructureTag == "" {
			continue
		}

		if viper.Get(mapstructureTag) == nil {
			viper.Set(mapstructureTag, "")
		}
	}

	return nil
}
