package templa

import "time"

// Date returns the current date in YYYY-MM-DD format
func Date() string {
	return time.Now().Format("2006-01-02")
}
