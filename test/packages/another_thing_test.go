package thing_test

import (
	"testing"

	"github.com/object88/cprofile/test/packages"
)

func Test_Another_Thing(t *testing.T) {
	th := &thing.Thing{A: 1, B: 2}

	if th.B != 2 {
		t.Fatalf("th.B=%d", th.B)
	}
}
