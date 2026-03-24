package detect

import "testing"

func TestHCLDetector(t *testing.T) {
	d := &HCLDetector{}

	tests := []struct {
		name     string
		input    string
		wantConf Confidence
	}{
		{
			"terraform resource blocks",
			`resource "aws_instance" "web" {
  ami           = "ami-123"
  instance_type = "t2.micro"
}

resource "aws_s3_bucket" "data" {
  bucket = "my-bucket"
}`,
			High,
		},
		{
			"variable and output",
			`variable "region" {
  default = "us-east-1"
}

output "instance_ip" {
  value = aws_instance.web.public_ip
}`,
			High,
		},
		{
			"single block header",
			`resource "aws_instance" "web" {
  ami = "ami-123"
}`,
			Medium,
		},
		{
			"assignments with braces no toml sections",
			`name = "myapp"
region = "us-east-1"
count  = 3
tags = {
  env = "prod"
}`,
			Medium,
		},
		{
			"toml-like sections prefer toml",
			"[server]\nhost = \"localhost\"\nport = 8080",
			None,
		},
		{"json", `{"key": "value"}`, None},
		{"xml", "<root/>", None},
		{"plain text", "just some text here", None},
		{"empty", "", None},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := d.Detect([]byte(tt.input))
			if r.Confidence != tt.wantConf {
				t.Errorf("confidence for %q = %v, want %v", tt.name, r.Confidence, tt.wantConf)
			}
		})
	}
}
