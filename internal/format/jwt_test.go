package format

import (
	"bytes"
	"strings"
	"testing"
)

func TestJWTFormatter_NoColor(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"

	f := &JWTFormatter{}
	var buf bytes.Buffer
	err := f.Format(&buf, strings.NewReader(token), nil)
	if err != nil {
		t.Fatal(err)
	}

	out := buf.String()
	if !strings.Contains(out, "Header:") {
		t.Error("output should contain 'Header:'")
	}
	if !strings.Contains(out, `"alg"`) {
		t.Error("output should contain decoded 'alg' field")
	}
	if !strings.Contains(out, "Payload:") {
		t.Error("output should contain 'Payload:'")
	}
	if !strings.Contains(out, "John Doe") {
		t.Error("output should contain decoded 'name' value")
	}
	if !strings.Contains(out, "Signature:") {
		t.Error("output should contain 'Signature:'")
	}
}
