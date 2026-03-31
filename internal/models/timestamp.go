package models

import (
	"fmt"
	"io"
	"time"

	"github.com/99designs/gqlgen/graphql"
)

type Timestamp time.Time

const timeLayout = time.RFC3339

func MarshalTimestamp(t Timestamp) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		io.WriteString(w, fmt.Sprintf("\"%s\"", time.Time(t).Format(timeLayout)))
	})
}

func UnmarshalTimestamp(v interface{}) (Timestamp, error) {
	str, ok := v.(string)
	if !ok {
		return Timestamp(time.Time{}), fmt.Errorf("timestamp must be string")
	}

	t, err := time.Parse(timeLayout, str)
	return Timestamp(t), err
}
