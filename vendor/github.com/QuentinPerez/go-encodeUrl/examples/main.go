package main

import (
	"fmt"
	"reflect"

	"github.com/QuentinPerez/go-encodeUrl"
	"github.com/Sirupsen/logrus"
)

type ID struct {
	Name        string `url:"name,ifStringIsNotEmpty"`
	DisplayName string `url:"display-name,ifStringIsNotEmpty"`
}

func nameEqualBadCode(interface{}) (string, bool, error) {
	return "0xBadC0de", true, nil
}

func overloadFieldNameAndTag(fieldName string, fieldTag reflect.StructTag) (fieldTagOverloaded reflect.StructTag) {
	fieldTagOverloaded = fieldTag
	if fieldName == "Name" {
		fieldTagOverloaded = `url:"overloaded,nameEqualBadCode"`
	}
	return
}

// changTag overloads the tag field
func changeNameAndTag() {
	encurl.AddEncodeFunc(nameEqualBadCode)
	values, errs := encurl.Translate(&ID{"qperez", "Quentin Perez"}, overloadFieldNameAndTag)
	if errs != nil {
		logrus.Fatal(errs)
	}
	fmt.Printf("https://example.com/?%v\n", values.Encode())
}

func main() {
	values, errs := encurl.Translate(&ID{"qperez", "Quentin Perez"})
	if errs != nil {
		logrus.Fatal(errs)
	}
	fmt.Printf("https://example.com/?%v\n", values.Encode())
	changeNameAndTag()
}
