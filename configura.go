package configura

/*
Original from
 https://github.com/agonzalezro/configura/blob/master/configura.go
under the Apache 2.0 license.

Modifications from original:
 - removed prefix
 - renamed flag from `configura` to `env`
 - cleaned up documentation
 - changed empty default behavior for numbers & bools: interpreted as 0 or false
*/

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func mismatchError(n string, i interface{}, t reflect.Kind) error {
	return fmt.Errorf("%s=%v must be %s", n, i, t)
}

func getStructInfo(v reflect.StructField) (fieldName, envVar, defVal string) {
	fieldName = v.Name
	tags := strings.Split(v.Tag.Get("env"), ",")
	envVar = tags[0]
	if len(tags) > 1 {
		defVal = tags[1]
	}
	return
}

// LoadEnv will go through all the fields defined
// in the struct and load their values from system
// environment variables.
//
// The variable name can set using struct tags.
// A default value can also be optionally set.
//
// For example:
//
// var config = struct {
// 	LogPrefix   string `env:"LOG_PREFIX"`
// 	Port        int    `env:"PORT,8888"`
// 	Development bool   `env:"DEVELOPMENT"`
// }{}
//
// err = LoadEnv(&config)
func LoadEnv(c interface{}) error {
	t := reflect.TypeOf(c)
	te := t.Elem()
	v := reflect.ValueOf(c)
	ve := v.Elem()

	if te.Kind() != reflect.Struct {
		return errors.New("the config must be a struct")
	}

	for i := 0; i < te.NumField(); i++ {
		sf := te.Field(i)
		fieldName, envVar, defVal := getStructInfo(sf)
		field := ve.FieldByName(fieldName)

		if envVar == "" {
			envVar = strings.ToUpper(fieldName)
		}
		env := os.Getenv(envVar)

		if env == "" {
			env = defVal
		}

		kind := field.Kind()

		switch kind {
		case reflect.Int, reflect.Int64, reflect.Float32, reflect.Float64:
			if env == "" {
				env = "0"
			}
		case reflect.Bool:
			if env == "" {
				env = "false"
			}
		}

		switch kind {
		case reflect.String:
			field.SetString(env)
		case reflect.Int:
			n, err := strconv.Atoi(env)
			if err != nil {
				return mismatchError(fieldName, n, kind)
			}
			field.SetInt(int64(n))
		case reflect.Float32, reflect.Float64:
			bitSize := 32
			if kind == reflect.Float64 {
				bitSize = 64
			}
			n, err := strconv.ParseFloat(env, bitSize)
			if err != nil {
				return mismatchError(fieldName, n, kind)
			}
			field.SetFloat(n)
		case reflect.Bool:
			b, err := strconv.ParseBool(env)
			if err != nil {
				return mismatchError(fieldName, b, kind)
			}
			field.SetBool(b)
		case reflect.Int64: // time.Duration
			t, err := time.ParseDuration(env)
			if err != nil {
				return mismatchError(fieldName, t, kind)
			}
			field.Set(reflect.ValueOf(t))
		default:
			return fmt.Errorf("%s is not parsable", kind)
		}
	}

	return nil
}
