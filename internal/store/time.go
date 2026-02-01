package store

import "time"

// sqliteTimeLayout matches SQLite datetime('now') output.
const sqliteTimeLayout = "2006-01-02 15:04:05"

func parseTime(s string) time.Time {
	t, _ := time.Parse(sqliteTimeLayout, s)
	if t.IsZero() {
		t, _ = time.Parse(time.RFC3339, s)
	}
	return t
}
