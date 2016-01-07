package api

// RequirementSet is a map of key:value equality requirements which
// match against any Labels which are supersets.
type RequirementSet map[string]string

// Matches returns true if the given labels satisfy all the requirements,
// false otherwise.
func (r RequirementSet) Matches(labels Labels) bool {
	for k, v := range r {
		if labels.Get(k) != v {
			return false
		}
	}
	return true
}

// Labels present key to value mappings, independent of their storage.
type Labels interface {
	// Get returns the value for the given label.
	Get(label string) string
}

// LabelSet is a map of key:value labels.
type LabelSet map[string]string

// Get returns the value for the given label.
func (ls LabelSet) Get(label string) string {
	return ls[label]
}
