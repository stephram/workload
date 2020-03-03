package ulid

import (
	crand "crypto/rand"
	"fmt"
	"math/rand"
	"time"

	"github.com/oklog/ulid"
)

// New returns a new ULID
func New() string {
	return ulid.MustNew(ulid.Now(), crand.Reader).String()
}

func NewTime() string {
	t := time.Unix(1000000, 0)
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	return fmt.Sprint(ulid.MustNew(ulid.Timestamp(t), entropy))
}
