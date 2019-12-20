package timestamp

import (
	"fmt"
	"time"
)

// FromNow Offset is a number that will multiply TimUnit
// and the result will be added to the current time
type FromNow struct {
	Offset   int
	TimeUnit time.Duration
}

func (ts FromNow) String() string {
	t := time.Now().Add(ts.TimeUnit * time.Duration(ts.Offset))
	timeStamp := t.Format("2006-01-02 15:04:05")

	return fmt.Sprintf("%v", timeStamp)
}
