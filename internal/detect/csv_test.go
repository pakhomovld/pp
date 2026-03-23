package detect

import "testing"

func TestCSVDetector(t *testing.T) {
	d := &CSVDetector{}

	tests := []struct {
		name       string
		input      string
		wantFormat Format
		wantConf   Confidence
	}{
		{
			"basic csv",
			"name,age,city\nAlice,30,NYC\nBob,25,LA",
			CSV, Medium,
		},
		{
			"5+ lines csv",
			"a,b,c\n1,2,3\n4,5,6\n7,8,9\n10,11,12\n13,14,15",
			CSV, High,
		},
		{
			"tsv",
			"name\tage\tcity\nAlice\t30\tNYC\nBob\t25\tLA",
			TSV, Medium,
		},
		{
			"single line",
			"name,age,city",
			CSV, None,
		},
		{
			"json",
			`{"key": "value"}`,
			CSV, None,
		},
		{
			"inconsistent columns",
			"a,b,c\n1,2\n3,4,5,6",
			CSV, None,
		},
		{
			"empty",
			"",
			CSV, None,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := d.Detect([]byte(tt.input))
			if r.Confidence != tt.wantConf {
				t.Errorf("confidence = %v, want %v", r.Confidence, tt.wantConf)
			}
			if r.Confidence > None && r.Format != tt.wantFormat {
				t.Errorf("format = %v, want %v", r.Format, tt.wantFormat)
			}
		})
	}
}
