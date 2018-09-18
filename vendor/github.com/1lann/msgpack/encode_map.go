package msgpack

import (
	"reflect"
	"sort"

	"github.com/1lann/msgpack/codes"
)

func encodeMapValue(e *Encoder, v reflect.Value) error {
	if v.IsNil() {
		return e.EncodeNil()
	}

	if err := e.EncodeMapLen(v.Len()); err != nil {
		return err
	}

	compressed := e.keysToCompressed != nil

	var useKey reflect.Value
	for _, key := range v.MapKeys() {
		useKey = key
		if compressed {
			k, ok := key.Interface().(string)
			if ok {
				keyValue, err := e.keysToCompressed(k)
				if err != nil {
					return err
				}
				useKey = reflect.ValueOf(keyValue)
			}
		}
		if err := e.EncodeValue(useKey); err != nil {
			return err
		}
		if err := e.EncodeValue(v.MapIndex(key)); err != nil {
			return err
		}
	}

	return nil
}

func encodeMapStringStringValue(e *Encoder, v reflect.Value) error {
	if v.IsNil() {
		return e.EncodeNil()
	}

	if err := e.EncodeMapLen(v.Len()); err != nil {
		return err
	}

	m := v.Convert(mapStringStringType).Interface().(map[string]string)
	// if e.sortMapKeys {
	// 	return e.encodeSortedMapStringString(m)
	// }

	compressed := e.keysToCompressed != nil

	var err error
	for mk, mv := range m {
		if compressed {
			mk, err = e.keysToCompressed(mk)
			if err != nil {
				return err
			}
		}

		if err := e.EncodeString(mk); err != nil {
			return err
		}

		if err := e.EncodeString(mv); err != nil {
			return err
		}
	}

	return nil
}

func encodeMapStringInterfaceValue(e *Encoder, v reflect.Value) error {
	if v.IsNil() {
		return e.EncodeNil()
	}

	if err := e.EncodeMapLen(v.Len()); err != nil {
		return err
	}

	m := v.Convert(mapStringInterfaceType).Interface().(map[string]interface{})
	// if e.sortMapKeys {
	// 	return e.encodeSortedMapStringInterface(m)
	// }

	compressed := e.keysToCompressed != nil

	var err error
	for mk, mv := range m {
		if compressed {
			mk, err = e.keysToCompressed(mk)
			if err != nil {
				return err
			}
		}

		if err := e.EncodeString(mk); err != nil {
			return err
		}
		if err := e.Encode(mv); err != nil {
			return err
		}
	}

	return nil
}

func (e *Encoder) encodeSortedMapStringString(m map[string]string) error {
	keys := make([]string, 0, len(m))
	for k, _ := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		err := e.EncodeString(k)
		if err != nil {
			return err
		}
		if err = e.EncodeString(m[k]); err != nil {
			return err
		}
	}

	return nil
}

func (e *Encoder) encodeSortedMapStringInterface(m map[string]interface{}) error {
	keys := make([]string, 0, len(m))
	for k, _ := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		err := e.EncodeString(k)
		if err != nil {
			return err
		}
		if err = e.Encode(m[k]); err != nil {
			return err
		}
	}

	return nil
}

func (e *Encoder) EncodeMapLen(l int) error {
	if l < 16 {
		return e.w.WriteByte(codes.FixedMapLow | byte(l))
	}
	if l < 65536 {
		return e.write2(codes.Map16, uint64(l))
	}
	return e.write4(codes.Map32, uint32(l))
}

func encodeStructValue(e *Encoder, strct reflect.Value) error {
	structFields := structs.Fields(strct.Type())
	if e.structAsArray || structFields.asArray {
		return encodeStructValueAsArray(e, strct, structFields.List)
	}
	fields := structFields.OmitEmpty(strct)

	if err := e.EncodeMapLen(len(fields)); err != nil {
		return err
	}

	compressed := e.keysToCompressed != nil

	var err error
	for _, f := range fields {
		name := f.name
		if compressed {
			name, err = e.keysToCompressed(f.name)
			if err != nil {
				return err
			}
		}

		if err := e.EncodeString(name); err != nil {
			return err
		}
		if err := f.EncodeValue(e, strct); err != nil {
			return err
		}
	}

	return nil
}

func encodeStructValueAsArray(e *Encoder, strct reflect.Value, fields []*field) error {
	if err := e.EncodeArrayLen(len(fields)); err != nil {
		return err
	}
	for _, f := range fields {
		if err := f.EncodeValue(e, strct); err != nil {
			return err
		}
	}
	return nil
}
