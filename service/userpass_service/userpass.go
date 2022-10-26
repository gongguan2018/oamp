package userpass_service

import "reflect"

func UserAndPass(form interface{}) interface{} {
	sv := reflect.ValueOf(form)
	svs := sv.Slice(0, sv.Len())
	for i := 0; i < svs.Len(); i++ {
		e := svs.Index(i).Interface()
		switch e.(type) {
		case string:
			return e
		}
	}
	return ""
}
