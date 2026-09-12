//go:build !windows

package main

import (
	"os/exec"
	"strings"
	"testing"
)

func TestForkSafetyAncestryBuildGate(t *testing.T) {
	cmd := exec.Command("make", "verify-fork-safety")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("required ancestry rejected: %v\n%s", err, out)
	}

	cmd = exec.Command("make", "verify-fork-safety", "FORK_SAFETY_ANCESTORS=1111111111111111111111111111111111111111")
	out, err := cmd.CombinedOutput()
	if err == nil || !strings.Contains(string(out), "required fork/upstream ancestor") {
		t.Fatalf("missing ancestor was not rejected: err=%v\n%s", err, out)
	}
}
