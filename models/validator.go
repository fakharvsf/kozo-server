package models

import (
	"fmt"
	"net/mail"
	"reflect"
)

type ValidatorClosure func() string

func NoValidator() string {
	return ""
}

func NullValidator(key string, value interface{}) string {
	valueType := reflect.TypeOf(value).Name()

	var error string

	if valueType == "string" {
		if value.(string) == "" {
			error = fmt.Sprintf("%s can't be null.", key)
		}
	}

	if valueType == "int" {
		// Do nothing. As we allow the value to be 0.
		error = ""
	}

	if valueType == "uint" {
		if value.(uint) == 0 {
			error = fmt.Sprintf("%s can't be null.", key)
		}
	}
	
	return error
}

func IDValidator(key string, value uint) string {
	var error string
	if value == 0 {
		error = fmt.Sprintf("%s can't be null.", key)
	}
	return error
}

func EmailValidator(key string, value string) string {
	var error string
	if value == "" {
		error = fmt.Sprintf("%s can't be null.", key)
	} else {
		_, errorParseAddress := mail.ParseAddress(value)
		if errorParseAddress != nil {
			error = fmt.Sprintf("%s is invalid.", key)
		}
	}
	return error
}

func LengthValidator(key string, value string, length int) string {
	var error string
	if value == "" {
		error = fmt.Sprintf("%s can't be null.", key)
	} else if len([]byte(value)) < length {
		error = fmt.Sprintf("%s can't be less than %d characters.", key, length)
	}
	return error
}

var PersonalTaskStatusesEnums = map[string]string{
	"created": "created",
	"assigned": "assigned",
	"completed": "completed",
}

func PersonalTaskStatusValidator(key string, value string) string {
	var error string
	if value == "" {
		error = fmt.Sprintf("%s can't be null.", key)
	} else {
		_, exists := PersonalTaskStatusesEnums[value]
		if !exists {
			error = fmt.Sprintf("%s is not valid for %s.", value, key)
		}
	}
	return error
}

var FriendRequestsStatusesEnums = map[string]string{
	"sent": "sent",
	"accepted": "accepted",
	"rejected": "rejected",
}

func Validator(validatorClosures []ValidatorClosure) []string {
	errors := []string{}

	for i := 0; i < len(validatorClosures); i++ {
		f := validatorClosures[i]
		error := f()
		if error != "" {
			errors = append(errors, error)
		}
	}

	return errors
}