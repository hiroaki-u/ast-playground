package main

var typeZeroValueMap = map[string]string{
	"bool":    "false",
	"byte":    "0",
	"error":   "nil",
	"float32": "0.0",
	"float64": "0.0",
	"int":     "0",
	"int8":    "0",
	"int16":   "0",
	"int32":   "0",
	"int64":   "0",
	"string":  "''",
	"uint":    "0",
	"uint8":   "0",
	"uint16":  "0",
	"uint32":  "0",
	"uint64":  "0",
	"uintptr": "0",
}

// astのIdentがプリミティブ型かどうかを判定する
func isPrimitive(typeName string) bool {
	_, ok := typeZeroValueMap[typeName]
	return ok
}
