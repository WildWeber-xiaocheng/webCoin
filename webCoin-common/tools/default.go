package tools

import (
	"errors"
	"reflect"
)

// 通过反射给字段生成默认值
func Default(data any) error {
	typeOf := reflect.TypeOf(data)
	valueOf := reflect.ValueOf(data)
	if typeOf.Kind() != reflect.Pointer {
		return errors.New("must be pointer")
	}
	ele := typeOf.Elem()
	valueEle := valueOf.Elem()
	for i := 0; i < ele.NumField(); i++ {
		field := ele.Field(i)
		value := valueEle.Field(i)
		//field.Tag.Get("default") 这里可以获取字段设置的tag
		kind := field.Type.Kind()
		if kind == reflect.Int {
			//这里根据设置的tag进行值的设置
			value.Set(defaultInt())
		}
		if kind == reflect.Int32 {
			value.Set(defaultInt32())
		}
		if kind == reflect.Int64 {
			value.Set(defaultInt64())
		}
		if kind == reflect.String {
			value.Set(defaultString())
		}
		if kind == reflect.Float64 {
			value.Set(defaultFloat64())
		}
		if kind == reflect.Float32 {
			value.Set(defaultFloat32())
		}
	}
	return nil
}

func defaultString() reflect.Value {
	var i = ""
	return reflect.ValueOf(i)
}

func defaultInt() reflect.Value {
	var i int = -1
	return reflect.ValueOf(i)
}

func defaultInt32() reflect.Value {
	var i int32 = -1
	return reflect.ValueOf(i)
}
func defaultInt64() reflect.Value {
	var i int64 = -1
	return reflect.ValueOf(i)
}

func defaultFloat64() reflect.Value {
	var i float64 = -1
	return reflect.ValueOf(i)
}
func defaultFloat32() reflect.Value {
	var i float32 = -1
	return reflect.ValueOf(i)
}
