package functions

import (
	"errors"
	"reflect"
)

func CompareData(a, b interface{}) error {
	if !reflect.DeepEqual(a, b) {
		return errors.New("ข้อมูลไม่ตรงกัน")
	}
	return nil
}
