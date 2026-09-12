package main

import (
	"slices"
	"testing"
)

func TestCIWorkflowRunsTestsOnAllSupportedDesktopPlatforms(t *testing.T) {
	job := ciTestJob(t)
	got := make(map[string]int)
	for _, row := range job.Strategy.Matrix.Include {
		got[row["os"]]++
	}

	want := map[string]int{
		"ubuntu-latest":  1,
		"macos-latest":   1,
		"windows-latest": 3,
	}
	if len(got) != len(want) {
		t.Fatalf("test matrix operating systems = %v, want %v", got, want)
	}
	for osName, count := range want {
		if got[osName] != count {
			t.Errorf("test matrix rows for %q = %d, want %d", osName, got[osName], count)
		}
	}
}

func TestCIWorkflowUsesRaceTestsOnUnixRunners(t *testing.T) {
	job := ciTestJob(t)
	commands := workflowCommandsMatching(job.Steps, func(step wfStep) bool {
		return exactRunnerOSCondition(step.If, "!=", "Windows")
	})

	var raceTests []workflowCommand
	for _, command := range commands {
		if command.name == "go" && slices.Equal(command.args, []string{"test", "-race", "./..."}) {
			raceTests = append(raceTests, command)
		}
	}
	if len(raceTests) != 1 {
		t.Fatalf("Unix-only go test -race ./... commands = %d, want 1; normalized commands: %#v", len(raceTests), commands)
	}
}

// TestCIWorkflowMakesForkSafetyAncestryAvailable reproduces the Linux/macOS
// PR failure where the root ancestry regression ran in checkout's default
// depth-1 clone: verify-fork-safety correctly refused because PR 11's commit
// object was absent. The whole test job needs the real graph so that positive
// ancestry means retained history rather than an environment-dependent skip.
func TestCIWorkflowMakesForkSafetyAncestryAvailable(t *testing.T) {
	job := ciTestJob(t)
	var checkouts []wfStep
	for _, step := range job.Steps {
		if step.Uses == "actions/checkout@v6" {
			checkouts = append(checkouts, step)
		}
	}
	if len(checkouts) != 1 {
		t.Fatalf("test job checkout steps = %d, want 1", len(checkouts))
	}
	if got := checkouts[0].With["fetch-depth"]; got != "0" {
		t.Fatalf("test job checkout fetch-depth = %q, want 0 so fork ancestry is available", got)
	}
}
