package wordtools

import (
	"testing"
)

func TestSubstitutionCyper(t *testing.T) {
	input := "aäeiywuøorlnmbpvfgkdtzsžšh"
	expected := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	output := SubstitutionCypher(input)

	if output != expected {
		t.Errorf("Output was incorrect, got: %s, want: %s.", output, expected)
	}
}
