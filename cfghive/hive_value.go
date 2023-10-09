package cfghive

import (
	"errors"
	"fmt"
)

const (
	HiveTypeBool = iota
	HiveTypeByte
	HiveTypeInt64
	HiveTypeUint64
	HiveTypeFloat64
	HiveTypeInt
	HiveTypeUint
	HiveTypeFloat32
	HiveTypeString
	HiveTypeBytes
	HiveTypeSub
)

var HiveTypeMap map[int]string = map[int]string{
	0:  "bool",
	1:  "byte",
	2:  "int64",
	3:  "uint64",
	4:  "float64",
	5:  "int",
	6:  "uint",
	7:  "float32",
	8:  "string",
	9:  "bytes",
	10: "sub",
}

type HiveValue struct {
	value      interface{}
	vlen       uint64
	storedType byte
}

func NewHiveValue(v interface{}) (HiveValue, error) {
	hv := HiveValue{nil, 0, 0}
	switch v.(type) {
	case bool:
		hv.storedType = HiveTypeBool
		hv.value = v.(bool)
	case byte:
		hv.storedType = HiveTypeByte
		hv.value = v.(byte)
	case int64:
		hv.storedType = HiveTypeInt64
		hv.value = v.(int64)
	case uint64:
		hv.storedType = HiveTypeUint64
		hv.value = v.(uint64)
	case float64:
		hv.storedType = HiveTypeFloat64
		hv.value = v.(float64)
	case int:
		hv.storedType = HiveTypeInt
		hv.value = v.(int)
	case uint:
		hv.storedType = HiveTypeUint
		hv.value = v.(uint)
	case float32:
		hv.storedType = HiveTypeFloat32
		hv.value = v.(float32)
	case string:
		hv.storedType = HiveTypeString
		hv.vlen = uint64(len(v.(string)))
		hv.value = v.(string)
	case []byte:
		hv.storedType = HiveTypeBytes
		hv.vlen = uint64(len(v.([]byte)))
		hv.value = v.([]byte)
	case map[string]HiveValue:
		hv.storedType = HiveTypeSub
		hv.vlen = uint64(len(v.(map[string]HiveValue)))
		hv.value = v.(map[string]HiveValue)
	case map[string]interface{}:
		hv.storedType = HiveTypeSub
		hv.vlen = uint64(len(v.(map[string]interface{})))
		value, err := GenericMapToSubMap(v.(map[string]interface{}))
		if err != nil {
			return hv, err
		}
		hv.value = value
	default:
		return hv, fmt.Errorf("invalid type %T", v)
	}
	return hv, nil
}

func (v *HiveValue) Type() byte {
	return v.storedType
}

func (v *HiveValue) Len() uint64 {
	return v.vlen
}

func (v *HiveValue) IsStoredType(t byte) bool {
	return v.storedType == byte(t)
}

func (v *HiveValue) Value() interface{} {
	return v.value
}

func (v *HiveValue) Bool() (bool, error) {
	if v.storedType != HiveTypeBool {
		return false, errors.New("stored type is not bool")
	}
	return v.value.(bool), nil
}

func (v *HiveValue) Byte() (byte, error) {
	if v.storedType != HiveTypeByte {
		return 0, errors.New("stored type is not byte")
	}
	return v.value.(byte), nil
}

func (v *HiveValue) Int64() (int64, error) {
	if v.storedType != HiveTypeInt64 {
		return 0, errors.New("stored type is not int64")
	}
	return v.value.(int64), nil
}

func (v *HiveValue) Uint64() (uint64, error) {
	if v.storedType != HiveTypeUint64 {
		return 0, errors.New("stored type is not uint64")
	}
	return v.value.(uint64), nil
}

func (v *HiveValue) Float64() (float64, error) {
	if v.storedType != HiveTypeFloat64 {
		return 0, errors.New("stored type is not float64")
	}
	return v.value.(float64), nil
}

func (v *HiveValue) Int() (int, error) {
	if v.storedType != HiveTypeInt {
		return 0, errors.New("stored type is not int")
	}
	return v.value.(int), nil
}

func (v *HiveValue) Uint() (uint, error) {
	if v.storedType != HiveTypeUint {
		return 0, errors.New("stored type is not uint")
	}
	return v.value.(uint), nil
}

func (v *HiveValue) Float32() (float32, error) {
	if v.storedType != HiveTypeFloat32 {
		return 0, errors.New("stored type is not float32")
	}
	return v.value.(float32), nil
}

func (v *HiveValue) String() (string, error) {
	if v.storedType != HiveTypeString {
		return "", errors.New("stored type is not string")
	}
	return v.value.(string), nil
}

func (v *HiveValue) Bytes() ([]byte, error) {
	if v.storedType != HiveTypeBytes {
		return nil, errors.New("stored type is not bytes")
	}
	return v.value.([]byte), nil
}

func (v *HiveValue) Sub() (map[string]HiveValue, error) {
	if v.storedType != HiveTypeSub {
		return nil, errors.New("stored type is not sub")
	}
	return v.value.(map[string]HiveValue), nil
}

func (v *HiveValue) TypeString() string {
	return HiveTypeMap[int(v.storedType)]
}

func GenericMapToSubMap(v map[string]interface{}) (map[string]HiveValue, error) {
	sub := make(map[string]HiveValue)
	for k, v := range v {
		hv, err := NewHiveValue(v)
		if err != nil {
			return nil, err
		}
		sub[k] = hv
	}
	return sub, nil
}

func HiveMapToGeneric(hive map[string]HiveValue) map[string]interface{} {
	generic := make(map[string]interface{})
	for k, v := range hive {
		switch v.Type() {
		case HiveTypeBool:
			generic[k], _ = v.Bool()
		case HiveTypeByte:
			generic[k], _ = v.Byte()
		case HiveTypeInt64:
			generic[k], _ = v.Int64()
		case HiveTypeUint64:
			generic[k], _ = v.Uint64()
		case HiveTypeFloat64:
			generic[k], _ = v.Float64()
		case HiveTypeInt:
			generic[k], _ = v.Int()
		case HiveTypeUint:
			generic[k], _ = v.Uint()
		case HiveTypeFloat32:
			generic[k], _ = v.Float32()
		case HiveTypeString:
			generic[k], _ = v.String()
		case HiveTypeBytes:
			generic[k], _ = v.Bytes()
		case HiveTypeSub:
			generic[k] = HiveMapToGeneric(v.value.(map[string]HiveValue))
		}
	}
	return generic
}
