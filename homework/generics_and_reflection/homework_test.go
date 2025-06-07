package main

import (
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type Person struct {
	Name    string `properties:"name"`
	Address string `properties:"address,omitempty"`
	Age     int    `properties:"age"`
	Married bool   `properties:"married"`
}

func Serialize(person Person) string {
	var sb strings.Builder
	v := reflect.ValueOf(person)
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		tag := f.Tag.Get("properties")
		if tag == "" {
			continue
		}
		tagParts := strings.Split(tag, ",")
		propKey := tagParts[0]
		omitEmpty := len(tagParts) > 1 && tagParts[1] == "omitempty"
		val := v.Field(i)

		empty := false
		switch val.Kind() {
		case reflect.String:
			empty = val.Len() == 0
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			empty = val.Int() == 0
		case reflect.Bool:
			empty = !val.Bool()
		default:
			continue
		}
		if omitEmpty && empty {
			continue
		}
		if sb.Len() > 0 {
			sb.WriteByte('\n')
		}
		sb.WriteString(propKey)
		sb.WriteByte('=')
		switch val.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			sb.WriteString(strconv.FormatInt(val.Int(), 10))
		case reflect.Bool:
			sb.WriteString(strconv.FormatBool(val.Bool()))
		case reflect.String:
			sb.WriteString(val.String())
		}
	}
	return sb.String()
}

func TestSerialization(t *testing.T) {
	tests := map[string]struct {
		person Person
		result string
	}{
		"test case with empty fields": {
			result: "name=\nage=0\nmarried=false",
		},
		"test case with fields": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
			},
			result: "name=John Doe\nage=30\nmarried=true",
		},
		"test case with omitempty field": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
				Address: "Paris",
			},
			result: "name=John Doe\naddress=Paris\nage=30\nmarried=true",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := Serialize(test.person)
			assert.Equal(t, test.result, result)
		})
	}
}
