package api

import "testing"

// TEMPORARY — Stage 0.3 gate verification. This file proves CI reports red and
// attributes failures to the right job. It is reverted before merge.

func Test_GateVerification_UnitGateReportsFailure(t *testing.T) {
	got := 1 + 1
	want := 3
	if got != want {
		t.Fatalf("deliberate failure to verify the unit gate: got %d, want %d", got, want)
	}
}

// Deliberately misformatted to verify the gofmt check in the lint job.
func   Test_GateVerification_LintGateReportsFailure( t *testing.T )   {
        t.Log( "misformatted on purpose" )
}
