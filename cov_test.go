package cov

import (
	"errors"
	"strings"
	"testing"
)

const (
	hashA = "SOMEHASHCODE"
	hashB = "OTHERHASH"
	covA  = 72.6532058094226
	covB  = 74.6532058094226
)

func TestCovError(t *testing.T) {
	t.Parallel()
	t.Logf("%d, %d", hashMatch, diffCov)
	if !hashMatch.contains(hashMatch) {
		t.Error("hashCov type should contain itself")
	}
	if !diffCov.contains(diffCov) {
		t.Error("diffCov type should contain itself")
	}

	if hashMatch.contains(diffCov) {
		t.Error("hashMatch claims to contain diffCov")
	}
	if diffCov.contains(hashMatch) {
		t.Error("diffCov claims to contain hashMatch")
	}

	if match := diffCov | hashMatch; diffCov.contains(match) {
		t.Logf("%d", match)
		t.Log("Error string:", match)
		t.Error("hashMatch with diffCov should contain diffCov")
	}
	if match := hashMatch | diffCov; diffCov.contains(match) {
		t.Logf("%d", match)
		t.Log("Error string:", match)
		t.Error("hashMatch with diffCov should contain diffCov")
	}
}

func TestCovErrorWrapsNonFatalError(t *testing.T) {
	t.Parallel()
	for _, err := range []covError{
		hashMatch,
		diffCov,
		noError,
		hashMatch | diffCov,
		diffCov | hashMatch,
	} {
		if !errors.Is(err, NonFatalError) {
			t.Logf("Cov err '%v' is not '%v'", err, NonFatalError)
		}
	}
}

func TestCov(t *testing.T) {
	t.Parallel()
	csA := &Cov{Hash: hashA, Cov: covA}
	csB := &Cov{Hash: hashA, Cov: covB}
	t.Logf("Object a: %v, Object b: %v", *csA, *csB)
	// DON'T MODIFY COVS IN TESTS
	testMatch(t, csA, csB)
	testMatch(t, csB, csA)
}

func testMatch(t *testing.T, csA, csB *Cov) {
	t.Logf("diffCov: %d", diffCov)
	t.Logf("hashMatch: %d", hashMatch)
	t.Logf("Both: %d", hashMatch|diffCov)
	if err := csA.Matches(csB); err == 0 {
		t.Fatal("Returned nil for no matches when there should've been one")
	}
	if match := csA.Matches(csB); hashMatch.contains(match) {
		t.Logf("%d", match)
		t.Log(match)
		t.Fatal("Didn't detect same hash")
	}
	if match := csA.Matches(csB); diffCov.contains(match) {
		t.Logf("%d", match)
		t.Log(match)
		t.Fatal("Didn't detect diff coverage")
	}
	if match := csA.Matches(csB); !strings.Contains(match.Error(), "There is a matching coverage for hash") {
		t.Log(match)
		t.Fatal("Error message doesn't have mention of missing hash match")
	}
	if match := csA.Matches(csB); !strings.Contains(match.Error(), " and coverage is different than old") {
		t.Logf("%d", match)
		t.Log(match)
		t.Fatal("Error message doesn't have mention of different hash")
	}
}
