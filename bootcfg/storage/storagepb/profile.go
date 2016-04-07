package storagepb

import (
	"encoding/json"
	"errors"
)

var (
	ErrIdRequired = errors.New("Id is required")
)

// ParseProfile parses bytes into a Profile.
func ParseProfile(data []byte) (*Profile, error) {
	profile := new(Profile)
	err := json.Unmarshal(data, profile)
	return profile, err
}

// AssertValid validates a Profile. Returns nil if there are no validation
// errors.
func (p *Profile) AssertValid() error {
	// Id is required
	if p.Id == "" {
		return ErrIdRequired
	}
	return nil
}
