package format

import (
	"bytes"
	"strings"
	"testing"
)

func TestNDJSONFormatter_NoColor(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			"two objects",
			"{\"b\":2,\"a\":1}\n{\"d\":4,\"c\":3}\n",
			"{\n  \"a\": 1,\n  \"b\": 2\n}\n{\n  \"c\": 3,\n  \"d\": 4\n}\n",
		},
		{
			"invalid line passes through",
			"{\"a\":1}\nnot json\n{\"b\":2}\n",
			"{\n  \"a\": 1\n}\nnot json\n{\n  \"b\": 2\n}\n",
		},
		{
			"empty lines preserved",
			"{\"a\":1}\n\n{\"b\":2}\n",
			"{\n  \"a\": 1\n}\n\n{\n  \"b\": 2\n}\n",
		},
	}

	f := &NDJSONFormatter{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := f.Format(&buf, strings.NewReader(tt.input), nil)
			if err != nil {
				t.Fatal(err)
			}
			got := buf.String()
			if got != tt.want {
				t.Errorf("got:\n%s\nwant:\n%s", got, tt.want)
			}
		})
	}
}
