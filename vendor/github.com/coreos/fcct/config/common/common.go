// Copyright 2019 Red Hat, Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.)

package common

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/coreos/vcontext/tree"
	vyaml "github.com/coreos/vcontext/yaml"
	"gopkg.in/yaml.v3"
)

type TranslateOptions struct {
	Pretty bool
	Strict bool
}

type Common struct {
	Version string `yaml:"version"`
	Variant string `yaml:"variant"`
}

// Misc helpers

// Unmarshal unmarshals the data to "to" and also returns a context tree for the source. If strict
// is set it errors out on unused keys.
func Unmarshal(data []byte, to interface{}, strict bool) (tree.Node, error) {
	dec := yaml.NewDecoder(bytes.NewReader(data))
	dec.KnownFields(strict)
	if err := dec.Decode(to); err != nil {
		return nil, err
	}
	return vyaml.UnmarshalToContext(data)
}

// Marshal is a wrapper for marshaling to json with or without pretty-printing the output
func Marshal(from interface{}, pretty bool) ([]byte, error) {
	if pretty {
		return json.MarshalIndent(from, "", "  ")
	}
	return json.Marshal(from)
}

// camel takes a snake_case string and converting it to camelCase
func camel(in string) string {
	words := strings.Split(in, "_")
	for i, word := range words[1:] {
		words[i+1] = strings.Title(word)
	}
	return strings.Join(words, "")
}

// ToCamelCase converts the keys in a context tree from snake_case to camelCase
func ToCamelCase(t tree.Node) tree.Node {
	switch n := t.(type) {
	case tree.MapNode:
		m := tree.MapNode{
			Children: make(map[string]tree.Node, len(n.Children)),
			Keys:     make(map[string]tree.Leaf, len(n.Keys)),
			Marker:   n.Marker,
		}
		for k, v := range n.Children {
			m.Children[camel(k)] = ToCamelCase(v)
		}
		for k, v := range n.Keys {
			m.Keys[camel(k)] = v
		}
		return m
	case tree.SliceNode:
		s := tree.SliceNode{
			Children: make([]tree.Node, 0, len(n.Children)),
			Marker:   n.Marker,
		}
		for _, v := range n.Children {
			s.Children = append(s.Children, ToCamelCase(v))
		}
		return s
	default: // leaf
		return t
	}
}
