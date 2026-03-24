package detect

import "testing"

func TestSQLDetector(t *testing.T) {
	d := &SQLDetector{}

	tests := []struct {
		name     string
		input    string
		wantConf Confidence
	}{
		{"select with where", "SELECT * FROM users WHERE id = 1", High},
		{"select with from", "SELECT id, name FROM users", High},
		{"create table", "CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(255) NOT NULL)", High},
		{"insert into", "INSERT INTO users VALUES (1, 'Alice')", High},
		{"update set where", "UPDATE users SET name = 'Bob' WHERE id = 1", High},
		{"delete from where", "DELETE FROM users WHERE id = 1", High},
		{"with cte", "WITH active AS (SELECT * FROM users WHERE active = true) SELECT * FROM active", High},
		{"select alone", "SELECT 1", Medium},
		{"begin alone", "BEGIN", Medium},
		{"case insensitive", "select * from users where id = 1", High},
		{"leading whitespace", "  SELECT id FROM users", High},
		{"multiline", "SELECT id, name\nFROM users\nWHERE active = true\nORDER BY name", High},
		{"not sql plain text", "The select few went home", None},
		{"json with select", `{"query": "SELECT * FROM users"}`, None},
		{"xml", "<query>SELECT</query>", None},
		{"plain text", "just some regular text here", None},
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
