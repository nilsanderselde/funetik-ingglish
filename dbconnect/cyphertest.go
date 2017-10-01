package dbconnect

import "testing"

func TestCypher(t *testing.T) {
	input := "aäeiywuøolrmnbpvfgkdtzsžšh"
	correctOut := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	output := SubstitutionCypher(input)
	if output != correctOut {
		t.Errorf("Output was incorrect, got: %s, want: %s.", input, correctOut)
	}
}
