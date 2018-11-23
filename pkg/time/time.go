package time

import (
	"fmt"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

func FromStringTimestamp(rawTimestamp string) (time.Time, error) {
	timestamp, err := strconv.ParseInt(rawTimestamp, 10, 64)
	if err != nil {
		return time.Now(), errors.Wrap(err, fmt.Sprintf("can't parse unix timestamp '%s'", rawTimestamp))
	}

	return time.Unix(timestamp, 0), nil
}
