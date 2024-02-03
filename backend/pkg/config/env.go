package config

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// Takes a config and updates it with values from the environment
// Results in an environment variable of SERVER__PORT updating config.Server.Port
func updateConfigFromEnv(config interface{}) error {
	configValue := reflect.ValueOf(config).Elem()
	configType := reflect.TypeOf(config).Elem()

	return updateStructFromEnv([]string{}, configValue, configType)
}

// Walks the struct and updates the fields with values from the environment
func updateStructFromEnv(parentstructs []string, structValue reflect.Value, structType reflect.Type) error {
	for i := 0; i < structValue.NumField(); i++ {
		field := structValue.Field(i)
		fieldType := structType.Field(i)

		if fieldType.Type.Kind() == reflect.Struct {
			err := updateStructFromEnv(append(parentstructs, fieldType.Name), field, fieldType.Type)
			if err != nil {
				return err
			}
		} else {
			envVarName := getEnvVarName(parentstructs, structType, fieldType)
			envVarValue, exists := os.LookupEnv(envVarName)
			if exists {
				err := setFieldValue(&field, envVarValue)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func setFieldValue(field *reflect.Value, value string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int:
		intValue, err := parseInt(value)
		if err != nil {
			return err
		}
		field.SetInt(int64(intValue))
	default:
		return fmt.Errorf("unsupported field type: %v", field.Kind())
	}
	return nil
}

func parseInt(s string) (int, error) {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("unable to parse integer: %v", err)
	}
	return int(i), nil
}

func getEnvVarName(parentstructs []string, structType reflect.Type, field reflect.StructField) string {
	var envVarName string
	if len(parentstructs) == 0 {
		envVarName = strings.ToUpper(structType.Name() + "__" + field.Name)
	} else {
		envVarName = strings.ToUpper(strings.Join(parentstructs, "__") + "__" + field.Name)
	}
	return envVarName
}
