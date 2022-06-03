package validator

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type Validator interface {
	Validate(input interface{}) error
}

func Validate(input interface{}) (err error) {
	value := reflect.ValueOf(input)
	kind := value.Type()

	for i := 0; i < value.NumField(); i++ {
		name := kind.Field(i).Name
		field := value.Field(i)

		if field.Kind() == reflect.Ptr && field.IsNil() {
			continue
		}

		switch data := field.Interface().(type) {
		case string:
			err = validateString(kind.Field(i), data, name)
		case *string:
			err = validateString(kind.Field(i), *data, name)
		case *float32:
			err = validateFloat32(kind.Field(i), *data, name)
		case float32:
			err = validateFloat32(kind.Field(i), data, name)
		default:
			err = fmt.Errorf("unsupported datatype")
		}

		if err != nil {
			return err
		}
	}

	return
}

func validateString(field reflect.StructField, text, name string) error {
	re := field.Tag.Get("re")
	if re == "" {
		return nil
	}

	regex, err := regexp.Compile(re)
	if err != nil {
		return err
	}

	if !regex.MatchString(text) {
		return fmt.Errorf("field %s not match regex %s", name, re)
	}

	return nil
}

func validateFloat32(field reflect.StructField, number float32, name string) error {
	scope := field.Tag.Get("range")
	ranges := strings.Split(scope, ",")

	if len(ranges) != 2 {
		return fmt.Errorf("field %s has more than 2 values in range tag", name)
	}

	min, err := strconv.ParseFloat(ranges[0], 32)
	if err == nil && number < float32(min) {
		return fmt.Errorf("field %s is less than %v", name, min)
	}

	max, err := strconv.ParseFloat(ranges[1], 32)
	if err == nil && number > float32(max) {
		return fmt.Errorf("field %s is more than %v", name, max)
	}

	return nil
}
