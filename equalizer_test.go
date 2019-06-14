package equalizer

import (
	"testing"
)

func TestEqualizer(t *testing.T) {
	offset := NewStepOffset(96, 15)
	eq := NewEqualizer(96, 16, offset)

	assertEqual(t, eq.tape.Text(2), "111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111")
	assertEqual(t, eq.mask.Text(2), "111111111111111100000000000000000000000000000000000000000000000000000000000000000000000000000000")

	eq.Notify(false, 50)
	assertEqual(t, eq.tape.Text(2), "111111111111111111111111111111111111111111111100000000000000000000000000000000000000000000000000")

	assertEqual(t, eq.Claim(), false)
	assertEqual(t, eq.Claim(), false)
	assertEqual(t, eq.Claim(), false)
	assertEqual(t, eq.Claim(), true)

	eq.Notify(true, 10)
	assertEqual(t, eq.tape.Text(2), "111111111111111111111111111111111111000000000000000000000000000000000000000000000000001111111111")

	eq.Notify(false, 1)
	assertEqual(t, eq.tape.Text(2), "111111111111111111111111111111111110000000000000000000000000000000000000000000000000011111111110")
}

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Fatalf("%s != %s", a, b)
	}
}
