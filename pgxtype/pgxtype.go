package pgxtype

import (
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func Text(val string) pgtype.Text {
	return pgtype.Text{String: val, Valid: val != ""}
}

func Numeric(val float64, exp int32) pgtype.Numeric {
	if val == 0 {
		return pgtype.Numeric{Valid: false}
	}
	return pgtype.Numeric{Int: big.NewInt(int64(val * 100)), Exp: exp, Valid: true}
}

func Float4(val float32) pgtype.Float4 {
	return pgtype.Float4{Float32: val, Valid: true}
}

func Int4(val int32) pgtype.Int4 {
	return pgtype.Int4{Int32: val, Valid: true}
}

func Int8(val int64) pgtype.Int8 {
	return pgtype.Int8{Int64: val, Valid: true}
}

func Int2(val int16) pgtype.Int2 {
	return pgtype.Int2{Int16: val, Valid: true}
}

func Bool(val bool) pgtype.Bool {
	return pgtype.Bool{Bool: val, Valid: true}
}

func Time(val time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{Time: val, Valid: !val.IsZero()}
}

func SliceOfDates(val []time.Time) []pgtype.Date {
	out := make([]pgtype.Date, len(val))
	for i, v := range val {
		out[i] = Date(v)
	}

	return out
}

func Date(val time.Time) pgtype.Date {
	return pgtype.Date{Time: val, Valid: !val.IsZero()}
}

func AnyToText(val interface{}) pgtype.Text {
	if val == nil {
		return pgtype.Text{String: "", Valid: false}
	}
	s, ok := val.(string)
	if !ok || s == "" {
		return pgtype.Text{String: "", Valid: false}
	}
	return pgtype.Text{String: s, Valid: true}
}

func AnyToTime(val interface{}) pgtype.Timestamptz {
	if val == nil {
		return pgtype.Timestamptz{Valid: false}
	}

	t, ok := val.(time.Time)
	if !ok || t.IsZero() {
		return pgtype.Timestamptz{Valid: false}
	}

	return pgtype.Timestamptz{Time: t, Valid: true}
}

func IgnoreSqlErrNoRows(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return nil
	}

	return fmt.Errorf("query: %w", err)
}
