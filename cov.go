package cov

import (
	"encoding/json"
	"fmt"
)

// Creates a new Covs. Don't parallelize
// maxSize = 0 disables size checking/trimming
func NewCovs(maxSize int) *Covs {
	return &Covs{
		content: make([]*Cov, 0, maxSize),
		maxSize: maxSize,
	}
}

// Implements json Decoder and Encoder interfaces.
// Calling it multiple times will add the values, while respecting maxSize
type Covs struct {
	content []*Cov
	maxSize int
}

func (cs *Covs) Init(payload []*Cov) error { return cs.clean(payload) }

func (cs *Covs) Add(payload []*Cov) error { return cs.clean(append(cs.content, payload...)) }

// Appends a new item to Content. There is no locking so make sure the program is not parallized
func (cs *Covs) Append(c *Cov) (err error) {
	err = nil
	for i := len(cs.content) - 1; i >= 0; i-- {
		if cErr := cs.content[i].Matches(c); cErr != noError { // assumes first hit is only possible hit
			err = cErr
			cs.content = append(cs.content[:i], cs.content[i+1:]...)
		}
	}
	cs.content = append(cs.content, c)
	cs.reduce() // we do not run clean since valid data is assumed
	return
}

// allows using directly with json decoder
func (cs *Covs) UnmarshalJSON(payload []byte) error {
	newContent := make([]*Cov, 0, cs.maxSize)
	if err := json.Unmarshal(payload, &newContent); err != nil {
		return err
	}
	if len(cs.content) != 0 {
		return cs.clean(append(cs.content, newContent...))
	}
	return cs.clean(newContent)
}

func (cs *Covs) MarshalJSON() ([]byte, error) { return json.Marshal(cs.content) }

func (cs *Covs) clean(payload []*Cov) (err error) {
	// remove duplicates
	// keeps last value in array (FIFO)
	keys := make(map[string]bool, len(payload))
	newContent := make([]*Cov, 0, len(payload))
	for i := len(payload) - 1; i >= 0; i-- {
		cov := payload[i]
		if _, value := keys[cov.Hash]; value {
			err = hashMatch
			continue
		}
		keys[cov.Hash] = true
		newContent = append(newContent, cov)
	}
	//reverse
	for i := len(newContent)/2 - 1; i >= 0; i-- {
		opp := len(newContent) - 1 - i
		newContent[i], newContent[opp] = newContent[opp], newContent[i]
	}
	cs.content = newContent
	// reduce to maxSize
	cs.reduce()
	return
}

func (cs *Covs) reduce() {
	if cs.maxSize != 0 {
		// reduce to maxSize
		for len(cs.content) > cs.maxSize {
			cs.content = cs.content[1:]
		}
	}
}

type Cov struct {
	Hash string  `json:"hash"`
	Cov  float64 `json:"cov"`
}

func (c *Cov) Matches(tC *Cov) covError {
	if c.Hash == tC.Hash {
		var result covError = hashMatch
		if c.Cov != tC.Cov {
			result |= diffCov
		}
		return result
	}
	return noError
}

type covError uint8

func (e covError) Unwrap() error { return NonFatalError }

func (e covError) Error() string {
	err := "Unknown error"
	if e.contains(hashMatch) {
		err = "There is a matching coverage for hash"
		if e.contains(diffCov) {
			err += " and coverage is different than old"
		}
	}
	return fmt.Sprintf("%v: %v", NonFatalError, err)
}

func (e covError) contains(t covError) bool { return t&e == t }

const (
	noError   covError = 0
	hashMatch covError = 1 << iota
	diffCov
)

var NonFatalError nonFatalError = 99

type nonFatalError uint8

func (e nonFatalError) Error() string { return "Non-Fatal error" }
