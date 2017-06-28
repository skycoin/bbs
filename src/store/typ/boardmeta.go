package typ

import (
	"encoding/json"
	"github.com/pkg/errors"
	"strings"
)

// BoardMeta contains the meta data of a board.
type BoardMeta struct {
	SubmissionAddresses []string `json:"submission_addresses"`
}

// Marshal marshals BoardMeta to json.
func (bm *BoardMeta) Marshal() ([]byte, error) {
	return json.Marshal(*bm)
}

// Unmarshal unmarshals BoardMeta from json.
func (bm *BoardMeta) Unmarshal(data []byte) error {
	return json.Unmarshal(data, bm)
}

// AddSubmissionAddress adds a submission address.
func (bm *BoardMeta) AddSubmissionAddress(address string) error {
	address = strings.TrimSpace(address)
	for _, old := range bm.SubmissionAddresses {
		if address == old {
			return errors.Errorf("address %s already exists", address)
		}
	}
	bm.SubmissionAddresses = append(bm.SubmissionAddresses, address)
	return nil
}

// RemoveSubmissionAddress removes a submission address.
func (bm *BoardMeta) RemoveSubmissionAddress(address string) {
	address = strings.TrimSpace(address)
	for i, old := range bm.SubmissionAddresses {
		if address == old {
			bm.SubmissionAddresses[i], bm.SubmissionAddresses[0] =
				bm.SubmissionAddresses[0], bm.SubmissionAddresses[i]
			bm.SubmissionAddresses = bm.SubmissionAddresses[1:]
		}
	}
}

// Trim trims spaces in addresses.
func (bm *BoardMeta) Trim() {
	for i := len(bm.SubmissionAddresses) - 1; i >= 0; i-- {
		bm.SubmissionAddresses[i] =
			strings.TrimSpace(bm.SubmissionAddresses[i])
		if bm.SubmissionAddresses[i] == "" {
			bm.SubmissionAddresses[0], bm.SubmissionAddresses[i] =
				bm.SubmissionAddresses[i], bm.SubmissionAddresses[0]
			bm.SubmissionAddresses = bm.SubmissionAddresses[1:]
		}
	}
}
