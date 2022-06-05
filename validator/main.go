package validator

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Validator interface {
	Validate(input interface{}) error
}

func Validate(input interface{}) (err error) {
	value := reflect.ValueOf(input)
	kind := value.Type()

	for i := 0; i < value.NumField(); i++ {
		name := kind.Field(i).Name
		tag := kind.Field(i).Tag
		field := value.Field(i)

		if field.Kind() == reflect.Ptr && field.IsNil() {
			continue
		}

		switch data := field.Interface().(type) {
		case string:
			err = validateString(tag.Get("re"), data, name)
		case *string:
			err = validateString(tag.Get("re"), *data, name)
		case float32:
			err = validateFloat32(tag.Get("range"), data, name)
		case *float32:
			err = validateFloat32(tag.Get("range"), *data, name)
		case time.Time:
			err = validateTime(tag.Get("period"), data, name)
		case *time.Time:
			err = validateTime(tag.Get("period"), *data, name)
		default:
			err = fmt.Errorf("unsupported datatype")
		}

		if err != nil {
			return err
		}
	}

	return
}

func validateString(tag, text, name string) error {
	re, err := regexp.Compile(tag)
	if err != nil {
		return err
	}

	if !re.MatchString(text) {
		return fmt.Errorf("field %s not match regex %s", name, tag)
	}

	return nil
}

func validateFloat32(tag string, number float32, name string) error {
	ranges := strings.Split(tag, ",")

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

func validateTime(tag string, date time.Time, name string) error {
	periods := strings.Split(tag, ",")
	if len(periods) != 2 {
		return fmt.Errorf("field %s has more than 2 values in range tag", name)
	}

	re, err := regexp.Compile(`([+-])(\d+)([a-z]+)`)
	if err != nil {
		return err
	}

	matches := re.FindAllStringSubmatch(periods[0], -1)
	before := time.Now()
	for _, match := range matches {
		if err := updateDate(&before, match); err != nil {
			return err
		}
	}

	if len(matches) != 0 && !before.Before(date) {
		return fmt.Errorf("field %s is before than %v", name, before)
	}

	matches = re.FindAllStringSubmatch(periods[1], -1)
	after := time.Now()
	for _, match := range matches {
		if err := updateDate(&after, match); err != nil {
			return err
		}
	}

	if len(matches) != 0 && !after.After(date) {
		return fmt.Errorf("field %s is after than %v", name, after)
	}

	return nil
}

func updateDate(date *time.Time, match []string) error {
	count, err := strconv.Atoi(match[2])
	if err != nil {
		return err
	}

	if match[1] == "-" {
		count *= -1
	}

	switch match[3] {
	case "years":
		*date = date.AddDate(count, 0, 0)
	case "months":
		*date = date.AddDate(0, count, 0)
	case "days":
		*date = date.AddDate(0, 0, count)
	case "hours":
		*date = date.Add(time.Duration(count) * time.Hour)
	case "minutes":
		*date = date.Add(time.Duration(count) * time.Minute)
	case "seconds":
		*date = date.Add(time.Duration(count) * time.Second)
	case "nanoseconds":
		*date = date.Add(time.Duration(count) * time.Nanosecond)
	default:
		return fmt.Errorf("not supported unit")
	}

	return nil
}
