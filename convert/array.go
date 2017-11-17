package convert

import "fmt"

const (
	typeString  = "string"
	typeInteger = "integer"
	typeNumber  = "number"

	formatInt32  = "int32"
	formatInt64  = "int64"
	formatFloat  = "float"
	formatDouble = "double"
)

// Array converts array of strings according to type and format of array items type
// described in OAS 2.0.
// https://github.com/OAI/OpenAPI-Specification/blob/master/versions/2.0.md#parameterObject
func Array(vals []string, itemsType, itemsFormat string) (value interface{}, err error) {
	switch itemsType {
	case typeString:
		switch itemsFormat {
		case "":
			return stringArray(vals)
		default:
			// For formats that are currently unsupported.
			return stringArray(vals)
		}
	case typeInteger:
		switch itemsFormat {
		case formatInt32:
			return int32Array(vals)
		case formatInt64:
			return int64Array(vals)
		default:
			// For formats that are currently unsupported.
			return int64Array(vals)
		}
	case typeNumber:
		switch itemsFormat {
		case formatFloat:
			return floatArray(vals)
		case formatDouble:
			return doubleArray(vals)
		default:
			// For formats that are currently unsupported.
			return doubleArray(vals)
		}
	default:
		return nil, fmt.Errorf("unsupported (not implemented yet?) items type %s for type array", itemsType)
	}
}

func stringArray(vals []string) (value interface{}, err error) {
	ps := make([]string, len(vals))
	for i, v := range vals {
		p, err := Primitive(v, typeString, "")
		if err != nil {
			// This should actually never happen.
			return nil, err
		}
		ps[i] = p.(string)
	}
	return ps, nil
}

func int32Array(vals []string) (value interface{}, err error) {
	ps := make([]int32, len(vals))
	for i, v := range vals {
		p, err := Primitive(v, typeInteger, formatInt32)
		if err != nil {
			return nil, err
		}
		ps[i] = p.(int32)
	}
	return ps, nil
}

func int64Array(vals []string) (value interface{}, err error) {
	ps := make([]int64, len(vals))
	for i, v := range vals {
		p, err := Primitive(v, typeInteger, formatInt64)
		if err != nil {
			return nil, err
		}
		ps[i] = p.(int64)
	}
	return ps, nil
}

func floatArray(vals []string) (value interface{}, err error) {
	ps := make([]float32, len(vals))
	for i, v := range vals {
		p, err := Primitive(v, typeNumber, formatFloat)
		if err != nil {
			return nil, err
		}
		ps[i] = p.(float32)
	}
	return ps, nil
}

func doubleArray(vals []string) (value interface{}, err error) {
	ps := make([]float64, len(vals))
	for i, v := range vals {
		p, err := Primitive(v, typeNumber, formatDouble)
		if err != nil {
			return nil, err
		}
		ps[i] = p.(float64)
	}
	return ps, nil
}
