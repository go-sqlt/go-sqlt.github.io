package repository

import (
	"context"
	"math/big"
	"net/url"
	"time"

	"github.com/go-sqlt/sqlt"
	"github.com/google/uuid"
)

type aBool bool

type Dest struct {
	// Scan can be used with types that implements the sql.Scanner interface.
	ID uuid.UUID
	// ScanBytes scans values into a byte slice.
	Bytes []byte
	// ScanTime scans values into time.Time.
	Time time.Time
	// A pointer indicates a nullable value.
	Date *time.Time
	// ScanString can be used with any kind of string values.
	String string
	// ScanInt can be used with any kind of int values.
	Int int32
	// ScanUint can be used with any kind of uint values.
	Uint uint64
	// ScanFloat can be used with any kind of float values.
	Float float32
	// ScanBool can be used with any kind of bool values.
	Bool aBool
	// ScanJSON can be used to unmarshal json values into any type.
	Map map[string]string
	// ScanBinary can be used with types that implements encoding.BinaryUnmarshaler.
	URL url.URL
	// ScanText can be used with types that implements encoding.TextUnmarshaler.
	Big big.Int
	// ScanStringSlice uses a string for scanning and splits its values using a separator.
	StringSlice []string
	// ScanIntSlice uses a string for scanning and splits its values using a separator.
	// Conversion with strconv.ParseInt(str, 10, 64).
	IntSlice []int64
	// ScanUintSlice uses a string for scanning and splits its values using a separator.
	// Conversion with strconv.ParseUint(str, 10, 64).
	UintSlice []uint64
	// ScanFloatSlice uses a string for scanning and splits its values using a separator.
	// Conversion with strconv.ParseFloat(str, 64).
	FloatSlice []float64
	// ScanBoolSlice uses a string for scanning and splits its values using a separator.
	// Conversion with strconv.ParseBool(str).
	BoolSlice []bool
}

var Query = sqlt.Custom[any](
	func(ctx context.Context, db sqlt.DB, expr sqlt.Expression[Dest]) (string, error) {
		return expr.SQL, nil
	},
	sqlt.Parse(`
		SELECT
			id			  					{{ Scan "ID" }}
			, raw							{{ ScanBytes "Bytes" }}
			, date							{{ ScanTime "Time" }}
			, '2020-01-01'					{{ ScanStringTime "Date" "DateOnly" "UTC" }}
			{{/* Raw preserves whitespace; sqlt collapses consecutive spaces by default. */}}
			, {{ Raw "'  i need  space'" }}	{{ ScanString "String" }}
			, 100							{{ ScanInt "Int" }}
			, 100							{{ ScanUint "Uint" }}
			, 1.23							{{ ScanFloat "Float" }}
			, 1								{{ ScanBool "Bool" }}
			, '{"some": "data"}'			{{ ScanJSON "Map" }}
			, 'one,two,three'				{{ ScanStringSlice "StringSlice" "," }}
			, '-10,20,30'					{{ ScanIntSlice "IntSlice" "," }}
			, '10,20,30'					{{ ScanUintSlice "UintSlice" "," }}
			, '1.1,2.2,3.3'					{{ ScanFloatSlice "FloatSlice" "," }}
			, '1,0,t,f'						{{ ScanBoolSlice "BoolSlice" "," }}
		FROM example;
	`),
)

// SELECT id , raw , date , '2020-01-01' , '  i need  space' , 100 , 100 , 1.23 , 1 , '{"some": "data"}' , 'one,two,three' , '-10,20,30' , '10,20,30' , '1.1,2.2,3.3' , '1,0,t,f' FROM example;
