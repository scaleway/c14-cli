package encurl

import (
	"errors"
	"strconv"
)

func ifStringIsNotEmpty(obj interface{}) (string, bool, error) {
	if val, ok := obj.(string); ok {
		if val != "" {
			return val, true, nil
		}
		return "", false, nil
	}
	return "", false, errors.New("this field should be a string")
}

func ifBoolIsFalse(obj interface{}) (string, bool, error) {
	if val, ok := obj.(bool); ok {
		if !val {
			return "False", true, nil
		}
		return "", false, nil
	}
	return "", false, errors.New("this field should be a boolean")
}

func ifBoolIsTrue(obj interface{}) (string, bool, error) {
	if val, ok := obj.(bool); ok {
		if val {
			return "True", true, nil
		}
		return "", false, nil
	}
	return "", false, errors.New("this field should be a boolean")
}

func itoa(obj interface{}) (string, bool, error) {
	if val, ok := obj.(int); ok {
		return strconv.Itoa(val), true, nil
	}
	return "", false, errors.New("this field should be an int")
}

func itoaIfNotNil(obj interface{}) (string, bool, error) {
	if val, ok := obj.(*int); ok {
		if val != nil {
			return strconv.Itoa(*val), true, nil
		}
		return "", false, nil
	}
	return "", false, errors.New("this field should be an *int")
}
