// Test for parsing data from structs based on their json tags.
package main

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

var jsonRe = regexp.MustCompile(`json:"?([^"]+)"?`)

// UserData contains additional user data.
type UserData struct {
	Admin bool `json:"admin"`
}

// User represents a user.
type User struct {
	Name    string    `json:"name"`
	Data    UserData  `json:"data"`
	PtrData *UserData `json:"ptr_data"`
}

func main() {
	user := &User{"Nero", UserData{false}, &UserData{true}}

	fmt.Println(jsonPluck(user, "name"))
	fmt.Println(jsonPluck(user, "data.admin"))
	fmt.Println(jsonPluck(user, "ptr_data.admin"))
	fmt.Println(jsonPluck(user, "name.fake"))
}

func jsonPluck(s interface{}, path string) (interface{}, error) {
	tags := strings.Split(path, ".")
	for _, tag := range tags {
		s = findFieldByJSONTag(s, tag)
		if s == nil {
			return nil, errors.New("json path returned nil")
		}
	}

	return s, nil
}

func findFieldByJSONTag(s interface{}, tag string) interface{} {
	vs := reflect.ValueOf(s)
	switch vs.Kind() {
	case reflect.Ptr, reflect.Interface:
		vs = vs.Elem()
	}

	if vs.Kind() != reflect.Struct {
		return nil // Not a struct, end search
	}

	typeOfS := vs.Type()
	for i := 0; i < vs.NumField(); i++ {
		if hasJSONTag(typeOfS.Field(i), tag) {
			return vs.Field(i).Interface()
		}
	}

	return nil
}

func hasJSONTag(f reflect.StructField, tag string) bool {
	match := jsonRe.FindStringSubmatch(string(f.Tag))
	if len(match) != 2 {
		return false
	}

	return match[1] == tag
}
