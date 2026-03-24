package detect

import (
	"bytes"
	"regexp"
)

var (
	sqlStatementRe = regexp.MustCompile(`(?i)^\s*(SELECT|INSERT\s+INTO|UPDATE|DELETE\s+FROM|CREATE\s+(TABLE|INDEX|VIEW|DATABASE|SCHEMA)|ALTER\s+TABLE|DROP\s+(TABLE|INDEX|VIEW)|WITH\s+\w+\s+AS|EXPLAIN|BEGIN|COMMIT|GRANT|REVOKE)\b`)
	sqlClauseRe    = regexp.MustCompile(`(?i)\b(FROM|WHERE|JOIN|LEFT\s+JOIN|RIGHT\s+JOIN|INNER\s+JOIN|GROUP\s+BY|ORDER\s+BY|HAVING|LIMIT|OFFSET|SET|VALUES|INTO|ON\s|AND|OR|NOT\s+NULL|PRIMARY\s+KEY|FOREIGN\s+KEY|REFERENCES|DEFAULT|UNION|INTERSECT|EXCEPT)\b`)
)

// SQLDetector detects SQL queries and statements.
type SQLDetector struct{}

func (d *SQLDetector) Detect(sample []byte) Result {
	trimmed := bytes.TrimSpace(sample)
	if len(trimmed) == 0 {
		return Result{Format: SQL, Confidence: None}
	}

	// Skip JSON, JSON arrays, and XML.
	if trimmed[0] == '{' || trimmed[0] == '[' || trimmed[0] == '<' {
		return Result{Format: SQL, Confidence: None}
	}

	hasStarter := sqlStatementRe.Match(trimmed)
	if !hasStarter {
		return Result{Format: SQL, Confidence: None}
	}

	clauseMatches := sqlClauseRe.FindAll(trimmed, -1)
	if len(clauseMatches) >= 1 {
		return Result{Format: SQL, Confidence: High}
	}

	return Result{Format: SQL, Confidence: Medium}
}
