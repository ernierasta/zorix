package config

import (
	"reflect"
	"testing"
)

// TestCommentedConfig tests, that every specified paramater is in commented config. This is mainly documentation test.
func TestCommentedConfig(t *testing.T) {
	c := New("../config.commented.toml")
	u, err := c.Read()
	if err != nil {
		t.Error(err)
	}

	// test if config has some variables which have wrong names or
	// are missing in config package structs
	if len(u) != 0 {
		t.Errorf("there are some variables unread from config file: %v", u)
	}

	// awesome test, checks which variables are uninitialized
	// if config has always non default examples we easily catch
	// forgotten vars.
	DeepFields(*c, t)
}

func DeepFields(iface interface{}, t *testing.T) {
	ifval := reflect.ValueOf(iface)
	iftype := reflect.TypeOf(iface)
	ifkind := ifval.Kind()

	if ifkind == reflect.Struct || ifkind == reflect.Slice {
		for i := 0; i < iftype.NumField(); i++ {
			v := ifval.Field(i)

			switch v.Kind() {
			case reflect.Struct:
				DeepFields(v.Interface(), t)
			case reflect.Slice:
				if v.Len() == 0 {
					t.Error("in struct:", iftype.Name()+",", iftype.Field(i).Name, ":", v, "is empty list!")

				}
				for i := 0; i < v.Len(); i++ {
					DeepFields(v.Index(i).Interface(), t)
				}
			default:
				if isDefault(v.Interface()) {
					t.Error("in struct:", iftype.Name()+",", iftype.Field(i).Name, ":", v, "has default value!")
				}
			}
		}
	} else {
		if isDefault(ifval) {
			t.Error("in struct:", iftype.Name()+",", ":", ifval, "has default value!")

			//fmt.Println(iftype.Name(), ":", ifval)
		}
	}
}

func isDefault(iface interface{}) bool {
	v := reflect.ValueOf(iface)
	return v.Interface() == reflect.Zero(v.Type()).Interface()
}
