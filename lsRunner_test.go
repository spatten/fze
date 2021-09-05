package fze

import (
	"testing"
)

func TestLsWithMatch(t *testing.T) {
	opts := RunnerOptions{
		TestFilter: "ba",
	}
	err := lsRunner([]string{"testFixtures/lsTest"}, opts)
	if err != nil {
		t.Errorf("err not expected: %v", err)
	}
}

func TestLsWithNoMatch(t *testing.T) {
	opts := RunnerOptions{
		TestFilter: "baaaldkjasldkjsa",
	}
	err := lsRunner([]string{"testFixtures/lsTest"}, opts)
	if err == nil {
		t.Errorf("err expected: %v", err)
	}
}
