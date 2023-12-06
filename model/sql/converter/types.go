package converter

import (
	"fmt"
	"strings"

	"github.com/zeromicro/ddl-parser/parser"
)

var unsignedTypeMap = map[string]string{
	"int":   "uint",
	"int8":  "uint8",
	"int16": "uint16",
	"int32": "uint32",
	"int64": "uint64",
}

var commonMysqlDataTypeMapInt = map[int]string{
	// For consistency, all integer types are converted to int64
	// number
	parser.Bit:       "byte",
	parser.TinyInt:   "int8",
	parser.SmallInt:  "int16",
	parser.MediumInt: "int32",
	parser.Int:       "int",
	parser.MiddleInt: "int32",
	parser.Int1:      "int",
	parser.Int2:      "int",
	parser.Int3:      "int",
	parser.Int4:      "int",
	parser.Int8:      "int8",
	parser.Integer:   "int",
	parser.BigInt:    "int64",
	parser.Float:     "float64",
	parser.Float4:    "float64",
	parser.Float8:    "float64",
	parser.Double:    "float64",
	parser.Decimal:   "decimal.Decimal",
	parser.Dec:       "decimal.Decimal",
	parser.Fixed:     "float64",
	parser.Numeric:   "float64",
	parser.Real:      "float64",
	// date&time
	parser.Date:      "time.Time",
	parser.DateTime:  "time.Time",
	parser.Timestamp: "time.Time",
	parser.Time:      "string",
	parser.Year:      "int64",
	// string
	parser.Char:            "string",
	parser.VarChar:         "string",
	parser.NVarChar:        "string",
	parser.NChar:           "string",
	parser.Character:       "string",
	parser.LongVarChar:     "string",
	parser.LineString:      "string",
	parser.MultiLineString: "string",
	parser.Binary:          "[]byte",
	parser.VarBinary:       "[]byte",
	parser.TinyText:        "string",
	parser.Text:            "string",
	parser.MediumText:      "string",
	parser.LongText:        "string",
	parser.Enum:            "string",
	parser.Set:             "string",
	parser.Json:            "string",
	parser.Blob:            "[]byte",
	parser.LongBlob:        "[]byte",
	parser.MediumBlob:      "[]byte",
	parser.TinyBlob:        "[]byte",
	// bool
	parser.Bool:    "bool",
	parser.Boolean: "bool",
}

var commonMysqlDataTypeMapString = map[string]string{
	// For consistency, all integer types are converted to int64
	// bool
	"bool":    "bool",
	"_bool":   "pq.BoolArray",
	"boolean": "bool",
	// number
	"tinyint":   "int8",
	"smallint":  "int16",
	"mediumint": "int32",
	"int":       "int",
	"int1":      "int",
	"int2":      "int",
	"_int2":     "pq.Int64Array",
	"int3":      "int",
	"int4":      "int",
	"_int4":     "pq.Int64Array",
	"int8":      "int8",
	"_int8":     "pq.Int64Array",
	"integer":   "int",
	"_integer":  "pq.Int64Array",
	"bigint":    "int64",
	"float":     "float64",
	"float4":    "float64",
	"_float4":   "pq.Float64Array",
	"float8":    "float64",
	"_float8":   "pq.Float64Array",
	"double":    "float64",
	"decimal":   "decimal.Decimal",
	"dec":       "decimal.Decimal",
	"fixed":     "float64",
	"real":      "float64",
	"bit":       "byte",
	// date & time
	"date":      "time.Time",
	"datetime":  "time.Time",
	"timestamp": "time.Time",
	"time":      "string",
	"year":      "int64",
	// string
	"linestring":      "string",
	"multilinestring": "string",
	"nvarchar":        "string",
	"nchar":           "string",
	"char":            "string",
	"bpchar":          "string",
	"_char":           "pq.StringArray",
	"character":       "string",
	"varchar":         "string",
	"_varchar":        "pq.StringArray",
	"binary":          "[]byte",
	"bytea":           "[]byte",
	"longvarbinary":   "[]byte",
	"varbinary":       "[]byte",
	"tinytext":        "string",
	"text":            "string",
	"_text":           "pq.StringArray",
	"mediumtext":      "string",
	"longtext":        "string",
	"enum":            "string",
	"set":             "string",
	"json":            "string",
	"jsonb":           "string",
	"blob":            "[]byte",
	"longblob":        "[]byte",
	"mediumblob":      "[]byte",
	"tinyblob":        "[]byte",
	"ltree":           "[]byte",
}

// ConvertDataType converts mysql column type into golang type
func ConvertDataType(dataBaseType int, isDefaultNull, unsigned, strict bool) (string, error) {
	tp, ok := commonMysqlDataTypeMapInt[dataBaseType]
	if !ok {
		return "", fmt.Errorf("unsupported database type: %v", dataBaseType)
	}

	return mayConvertNullType(tp, isDefaultNull, unsigned, strict), nil
}

// ConvertStringDataType converts mysql column type into golang type
func ConvertStringDataType(dataBaseType string, isDefaultNull, unsigned, strict bool) (
	goType string, isPQArray bool, err error) {
	tp, ok := commonMysqlDataTypeMapString[strings.ToLower(dataBaseType)]
	if !ok {
		return "", false, fmt.Errorf("unsupported database type: %s", dataBaseType)
	}

	if strings.HasPrefix(dataBaseType, "_") {
		return tp, true, nil
	}

	return mayConvertNullType(tp, isDefaultNull, unsigned, strict), false, nil
}

func mayConvertNullType(goDataType string, isDefaultNull, unsigned, strict bool) string {
	if !isDefaultNull {
		if unsigned && strict {
			ret, ok := unsignedTypeMap[goDataType]
			if ok {
				return ret
			}
		}
		return goDataType
	}

	switch goDataType {
	case "int64":
		return "sql.NullInt64"
	case "int32":
		return "sql.NullInt32"
	case "float64":
		return "sql.NullFloat64"
	case "bool":
		return "sql.NullBool"
	case "string":
		return "sql.NullString"
	case "time.Time":
		return "sql.NullTime"
	case "decimal":
		return "*decimal.Decimal"
	default:
		if unsigned {
			ret, ok := unsignedTypeMap[goDataType]
			if ok {
				return ret
			}
		}
		return goDataType
	}
}
