package thing

import "testing"

func Test_Thing(t *testing.T) {
	th := &Thing{1, 2}

	if th.A != 1 {
		t.Fatalf("th.A=%d", th.A)
	}
}
