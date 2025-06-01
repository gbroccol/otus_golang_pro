package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var (
	errNotStruct        = errors.New("input is not a struct")
	errUnsupportedType  = errors.New("unsupported type")
	errInvalidValidator = errors.New("invalid validator")
	errValidationLen    = errors.New("length mismatch")
	errValidationRegexp = errors.New("regex mismatch")
	errValidationIn     = errors.New("value not in allowed set")
	errValidationMin    = errors.New("value less than minimum")
	errValidationMax    = errors.New("value greater than maximum")
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	errs := make([]string, 0, len(v))
	for _, e := range v {
		errs = append(errs, fmt.Sprintf("field %s: %v", e.Field, e.Err))
	}
	return strings.Join(errs, "; ")
}

func Validate(v interface{}) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Struct {
		return errNotStruct
	}

	var validationErrors ValidationErrors

	typ := val.Type()
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		value := val.Field(i)
		tag := field.Tag.Get("validate")
		if tag == "" {
			continue
		}

		validators := strings.Split(tag, "|")
		errs := validateField(value, field, validators)
		if errs != nil {
			validationErrors = append(validationErrors, errs...)
		}
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}
	return nil
}

func validateField(value reflect.Value, field reflect.StructField, validators []string) ValidationErrors {
	var errs ValidationErrors
	switch value.Kind() {
	case reflect.String:
		err := validateString(field.Name, value.String(), validators)
		if err != nil {
			errs = append(errs, err...)
		}
	case reflect.Int:
		err := validateInt(field.Name, int(value.Int()), validators)
		if err != nil {
			errs = append(errs, err...)
		}
	case reflect.Slice:
		switch value.Type().Elem().Kind() {
		case reflect.String:
			for i := 0; i < value.Len(); i++ {
				err := validateString(field.Name, value.Index(i).String(), validators)
				if err != nil {
					errs = append(errs, err...)
				}
			}
		case reflect.Int:
			for i := 0; i < value.Len(); i++ {
				err := validateInt(field.Name, int(value.Index(i).Int()), validators)
				if err != nil {
					errs = append(errs, err...)
				}
			}
		default:
			errs = append(errs, ValidationError{Field: field.Name, Err: errUnsupportedType})
		}
	default:
		errs = append(errs, ValidationError{Field: field.Name, Err: errUnsupportedType})
	}
	return errs
}

func validateString(field, value string, validators []string) ValidationErrors {
	var errs ValidationErrors
	for _, validator := range validators {
		parts := strings.SplitN(validator, ":", 2)
		if len(parts) != 2 {
			errs = append(errs, ValidationError{Field: field, Err: errInvalidValidator})
			continue
		}
		switch parts[0] {
		case "len":
			length, err := strconv.Atoi(parts[1])
			if err != nil || len(value) != length {
				errs = append(errs, ValidationError{Field: field, Err: errValidationLen})
			}
		case "regexp":
			r, err := regexp.Compile(parts[1])
			if err != nil || !r.MatchString(value) {
				errs = append(errs, ValidationError{Field: field, Err: errValidationRegexp})
			}
		case "in":
			allowed := strings.Split(parts[1], ",")
			found := false
			for _, item := range allowed {
				if item == value {
					found = true
					break
				}
			}
			if !found {
				errs = append(errs, ValidationError{Field: field, Err: errValidationIn})
			}
		default:
			errs = append(errs, ValidationError{Field: field, Err: errInvalidValidator})
		}
	}
	return errs
}

func validateInt(field string, value int, validators []string) ValidationErrors {
	var errs ValidationErrors
	for _, validator := range validators {
		parts := strings.SplitN(validator, ":", 2)
		if len(parts) != 2 {
			errs = append(errs, ValidationError{Field: field, Err: errInvalidValidator})
			continue
		}
		switch parts[0] {
		case "min":
			minExpected, err := strconv.Atoi(parts[1])
			if err != nil || value < minExpected {
				errs = append(errs, ValidationError{Field: field, Err: errValidationMin})
			}
		case "max":
			maxExpected, err := strconv.Atoi(parts[1])
			if err != nil || value > maxExpected {
				errs = append(errs, ValidationError{Field: field, Err: errValidationMax})
			}
		case "in":
			allowed := strings.Split(parts[1], ",")
			found := false
			for _, item := range allowed {
				val, err := strconv.Atoi(item)
				if err == nil && val == value {
					found = true
					break
				}
			}
			if !found {
				errs = append(errs, ValidationError{Field: field, Err: errValidationIn})
			}
		default:
			errs = append(errs, ValidationError{Field: field, Err: errInvalidValidator})
		}
	}
	return errs
}
