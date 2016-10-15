package storagepb

import "encoding/json"

// ParseMetadata parses bytes into a Metadata.
func ParseMetadata(data []byte) (*Metadata, error) {
	richMetadata := new(RichMetadata)
	err := json.Unmarshal(data, richMetadata)
	if err != nil {
		return nil, err
	}
	metadata, err := richMetadata.ToMetadata()
	if err != nil {
		return nil, err
	}
	return metadata, err
}

// ToRichMetadata converts a Metadata into a RichMetadata suitable for writing and
// user manipulation.
func (m *Metadata) ToRichMetadata() (*RichMetadata, error) {
	metadata := make(map[string]interface{})
	if m.Metadata != nil {
		err := json.Unmarshal(m.Metadata, &metadata)
		if err != nil {
			return nil, err
		}
	}
	return &RichMetadata{
		Id:       m.Id,
		Name:     m.Name,
		Metadata: metadata,
	}, nil
}

// RichMetadata is a user provided Metadata definition.
type RichMetadata struct {
	// machine readable Id
	Id string `json:"id,omitempty"`
	// Human readable name
	Name string `json:"name,omitempty"`
	// Metadata
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// ToMetadata converts a user provided RichMetadata into a Metadata which can be
// serialized as a protocol buffer.
func (rg *RichMetadata) ToMetadata() (*Metadata, error) {
	var metadata []byte
	if rg.Metadata != nil {
		var err error
		metadata, err = json.Marshal(rg.Metadata)
		if err != nil {
			return nil, err
		}
	}
	return &Metadata{
		Id:       rg.Id,
		Name:     rg.Name,
		Metadata: metadata,
	}, nil
}
