package cov

import (
	"bytes"
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"testing"
)

var validMockData2 = []*Cov{
	&Cov{Hash: "asd1", Cov: 78.1335},
	&Cov{Hash: "asd2", Cov: 78.1335},
	&Cov{Hash: "asd3", Cov: 78.1336},
	&Cov{Hash: "asd4", Cov: 78.1337},
	&Cov{Hash: "asd5", Cov: 78.1338},
	&Cov{Hash: "asd6", Cov: 78.1339},
	&Cov{Hash: "asd7", Cov: 78.1340},
	&Cov{Hash: "asd8", Cov: 78.1341},
	&Cov{Hash: "asd9", Cov: 78.1342},
	&Cov{Hash: "asd10", Cov: 78.1347},
}

var validMockData = []*Cov{
	&Cov{Hash: "hash0", Cov: 78.1234},
	&Cov{Hash: "hash1", Cov: 78.1235},
	&Cov{Hash: "hash2", Cov: 78.1236},
	&Cov{Hash: "hash3", Cov: 78.1237},
	&Cov{Hash: "hash4", Cov: 78.1238},
	&Cov{Hash: "hash5", Cov: 78.1239},
	&Cov{Hash: "hash6", Cov: 78.1240},
	&Cov{Hash: "hash7", Cov: 78.1241},
	&Cov{Hash: "hash8", Cov: 78.1242},
	&Cov{Hash: "hash9", Cov: 78.1243},
}

var validMockDataTooLong = append(validMockData, []*Cov{
	&Cov{Hash: "hash10", Cov: 78.1244},
	&Cov{Hash: "hash11", Cov: 78.1245},
}...)

var validMockDataCollidingHashesTooLong = append(validMockData, []*Cov{
	&Cov{Hash: "hash0", Cov: 78.1244},
	&Cov{Hash: "hash1", Cov: 78.1245},
	&Cov{Hash: "other", Cov: 43.1231},
	&Cov{Hash: "other", Cov: 93.2},
}...)

func TestNewCovs(t *testing.T) {
	t.Parallel()
	cs := NewCovs(100)
	if len(cs.content) != 0 {
		t.Fatal("Covs has payload when shouldn't")
	}
	if cs.maxSize != 100 {
		t.Fatal("Covs has payload when shouldn't")
	}
}

func TestCovsInit(t *testing.T) {
	t.Parallel()
	cs := NewCovs(10)
	if err := cs.Init(validMockData); err != nil {
		t.Log(err)
		if err == hashMatch {
			t.Log("There is a matching hash")
		}
		t.Fatal("Initialization failed")
	}
	contentSame(t, cs, validMockData)
	t.Run("Add method, limited size", testCovsAdd)
	t.Run("Add method, no size", testCovsAddNoSize)
}

func testCovsAdd(t *testing.T) {
	t.Parallel()
	cs := NewCovs(10)
	if err := cs.Init(validMockData); err != nil {
		t.Fatal(err)
	}
	cs.Add(validMockData2)
	contentSame(t, cs, validMockData2)
}

func testCovsAddNoSize(t *testing.T) {
	t.Parallel()
	cs := NewCovs(0)
	if err := cs.Init(validMockData); err != nil {
		t.Fatal(err)
	}
	cs.Add(validMockData2)
	contentSame(t, cs, append(validMockData, validMockData2...))
}

func TestCovsUnmarshal(t *testing.T) {
	t.Parallel()
	// test with empty init
	raw, err := json.Marshal(validMockData)
	if err != nil {
		t.Fatal(err)
	}
	cs := NewCovs(10)
	if err := json.NewDecoder(bytes.NewReader(raw)).Decode(cs); err != nil {
		t.Fatal(err)
	}
	contentSame(t, cs, validMockData)

	// test with existing payload
	raw, err = json.Marshal(validMockData2)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.NewDecoder(bytes.NewReader(raw)).Decode(cs); err != nil {
		t.Fatal(err)
	}
	contentSame(t, cs, validMockData2)

	if err := json.NewDecoder(strings.NewReader(`"{"WRONG": "JSON", "DATA": 3}`)).Decode(cs); errors.Is(err, NonFatalError) {
		t.Log("error.Is(err, NonFatalError:", errors.Is(NonFatalError, err))
		t.Fatalf("Expected error not to be: %v, got: %v", NonFatalError, err)
	}

	if err := json.NewDecoder(strings.NewReader("INVALID DATA")).Decode(cs); errors.Is(err, NonFatalError) {
		t.Log("error.Is(err, NonFatalError:", errors.Is(NonFatalError, err))
		t.Fatalf("Expected error not to be: %v, got: %v", NonFatalError, err)
	}
}

func TestCovsMarshal(t *testing.T) {
	t.Parallel()
	cs := NewCovs(10)
	if err := cs.Init(validMockData); err != nil {
		t.Fatal(err)
	}
	expectedRaw, _ := json.Marshal(validMockData)
	b := new(strings.Builder)
	// NOTE: returns a newline at end of JSON, no idea why. Trim that before comparing
	if err := json.NewEncoder(b).Encode(cs); err != nil || strings.TrimSuffix(b.String(), "\n") != string(expectedRaw) {
		t.Log("Error:", err)
		t.Log("Encode is:", strings.TrimSuffix(b.String(), "\n"))
		t.Log("Should be:", string(expectedRaw))
		t.Log(strings.EqualFold(b.String(), string(expectedRaw)))
		t.Fatal("Failed to encode data using json.Encoder")
	}
}

func TestCovsClean(t *testing.T) {
	t.Parallel()
	for _, s := range []*scenario{
		&scenario{
			Source:    validMockData,
			Target:    validMockData,
			TargetErr: nil,
		},
		&scenario{
			Source: validMockDataTooLong,
			Target: []*Cov{
				&Cov{Hash: "hash2", Cov: 78.1236},
				&Cov{Hash: "hash3", Cov: 78.1237},
				&Cov{Hash: "hash4", Cov: 78.1238},
				&Cov{Hash: "hash5", Cov: 78.1239},
				&Cov{Hash: "hash6", Cov: 78.1240},
				&Cov{Hash: "hash7", Cov: 78.1241},
				&Cov{Hash: "hash8", Cov: 78.1242},
				&Cov{Hash: "hash9", Cov: 78.1243},
				&Cov{Hash: "hash10", Cov: 78.1244},
				&Cov{Hash: "hash11", Cov: 78.1245},
			},
			TargetErr: nil,
		},
		&scenario{
			Source: validMockDataCollidingHashesTooLong,
			Target: []*Cov{
				&Cov{Hash: "hash3", Cov: 78.1237},
				&Cov{Hash: "hash4", Cov: 78.1238},
				&Cov{Hash: "hash5", Cov: 78.1239},
				&Cov{Hash: "hash6", Cov: 78.1240},
				&Cov{Hash: "hash7", Cov: 78.1241},
				&Cov{Hash: "hash8", Cov: 78.1242},
				&Cov{Hash: "hash9", Cov: 78.1243},
				&Cov{Hash: "hash0", Cov: 78.1244},
				&Cov{Hash: "hash1", Cov: 78.1245},
				&Cov{Hash: "other", Cov: 93.2},
			},
			TargetErr: hashMatch,
		},
	} {
		t.Run("Clean with scenario", s.Matches)
	}
}

type scenario struct {
	Source    []*Cov
	Target    []*Cov
	TargetErr error
}

func (s *scenario) Matches(t *testing.T) {
	t.Parallel()
	cs := NewCovs(10)
	err := cs.clean(s.Source)
	contentSame(t, cs, s.Target)
	checkErr(t, err, s.TargetErr)
}

func checkErr(t *testing.T, this error, is error) {
	if is == nil {
		if this != nil {
			t.Errorf("Invalid error, got: %v, expected: %v", this, is)
		}
	} else {
		if !errors.Is(this, is) {
			t.Errorf("Invalid error, got: %v, expected: %v", this, is)
		}
	}
}

func contentSame(t *testing.T, this *Covs, is []*Cov) {
	if !reflect.DeepEqual(this.content, is) || len(this.content) != len(is) {
		t.Logf("Content length: %d, Mock data length: %d", len(this.content), len(is))
		t.Log("Source:")
		for _, field := range this.content {
			t.Log(field)
		}
		t.Log("Target:")
		for _, field := range is {
			t.Log(field)
		}
		t.Fatal("Source and target do not match")
	}
	if &validMockData == &this.content {
		t.Log(&validMockData, &this.content)
		t.Fatal("Data should not point to the same object as input")
	}
}

type appendScenario struct {
	Target    []*Cov
	NewCov    *Cov
	TargetErr error
}

func (s *appendScenario) Matches(t *testing.T) {
	t.Parallel()
	cs := NewCovs(10)
	if err := cs.Init(validMockData); err != nil {
		t.Fatal(err)
	}
	t.Log("New cov:", s.NewCov)
	err := cs.Append(s.NewCov)
	contentSame(t, cs, s.Target)
	checkErr(t, err, s.TargetErr)
}

func TestCovsAppend(t *testing.T) {
	t.Parallel()
	for _, s := range []*appendScenario{
		&appendScenario{
			NewCov:    &Cov{Hash: "other", Cov: 93.2},
			Target:    append(validMockData[1:], &Cov{Hash: "other", Cov: 93.2}),
			TargetErr: nil,
		},
		&appendScenario{
			NewCov: &Cov{Hash: "hash1", Cov: 93.2},
			Target: []*Cov{
				&Cov{Hash: "hash0", Cov: 78.1234},
				&Cov{Hash: "hash2", Cov: 78.1236},
				&Cov{Hash: "hash3", Cov: 78.1237},
				&Cov{Hash: "hash4", Cov: 78.1238},
				&Cov{Hash: "hash5", Cov: 78.1239},
				&Cov{Hash: "hash6", Cov: 78.1240},
				&Cov{Hash: "hash7", Cov: 78.1241},
				&Cov{Hash: "hash8", Cov: 78.1242},
				&Cov{Hash: "hash9", Cov: 78.1243},
				&Cov{Hash: "hash1", Cov: 93.2},
			},
			TargetErr: hashMatch | diffCov,
		},
		&appendScenario{
			NewCov: &Cov{Hash: "hash1", Cov: 78.1235},
			Target: []*Cov{
				&Cov{Hash: "hash0", Cov: 78.1234},
				&Cov{Hash: "hash2", Cov: 78.1236},
				&Cov{Hash: "hash3", Cov: 78.1237},
				&Cov{Hash: "hash4", Cov: 78.1238},
				&Cov{Hash: "hash5", Cov: 78.1239},
				&Cov{Hash: "hash6", Cov: 78.1240},
				&Cov{Hash: "hash7", Cov: 78.1241},
				&Cov{Hash: "hash8", Cov: 78.1242},
				&Cov{Hash: "hash9", Cov: 78.1243},
				&Cov{Hash: "hash1", Cov: 78.1235},
			},
			TargetErr: hashMatch,
		},
	} {
		t.Run("Clean with scenario", s.Matches)
	}
}
