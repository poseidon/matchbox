// Copyright 2016 CoreOS, Inc.
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
// limitations under the License.

package astyaml

import (
	"errors"
	"io"
	"strings"

	yaml "github.com/ajeddeloh/yaml"
	"github.com/coreos/ignition/config/validate/astnode"
)

var (
	ErrNotDocumentNode = errors.New("Can only convert from document node")
	ErrNotMappingNode  = errors.New("Tried to change the key of a node which is not a mapping node")
	ErrKeyNotFound     = errors.New("Key to be replaced not found")
)

type YamlNode struct {
	tag string
	key yaml.Node
	yaml.Node
}

func FromYamlDocumentNode(n yaml.Node) (YamlNode, error) {
	if n.Kind != yaml.DocumentNode {
		return YamlNode{}, ErrNotDocumentNode
	}

	return YamlNode{
		key:  n,
		tag:  "yaml",
		Node: *n.Children[0],
	}, nil
}

func (n YamlNode) ValueLineCol(source io.ReadSeeker) (int, int, string) {
	return n.Line + 1, n.Column + 1, ""
}

func (n YamlNode) KeyLineCol(source io.ReadSeeker) (int, int, string) {
	return n.key.Line + 1, n.key.Column + 1, ""
}

func (n YamlNode) LiteralValue() interface{} {
	return n.Value
}

func (n YamlNode) SliceChild(index int) (astnode.AstNode, bool) {
	if n.Kind != yaml.SequenceNode {
		return nil, false
	}
	if index >= len(n.Children) {
		return nil, false
	}

	return YamlNode{
		key:  yaml.Node{},
		tag:  n.tag,
		Node: *n.Children[index],
	}, true
}

func (n YamlNode) KeyValueMap() (map[string]astnode.AstNode, bool) {
	if n.Kind != yaml.MappingNode {
		return nil, false
	}

	kvmap := map[string]astnode.AstNode{}
	for i := 0; i < len(n.Children); i += 2 {
		key := *n.Children[i]
		if n.tag == "json" {
			key.Value = getIgnKeyName(key.Value)
		}
		value := *n.Children[i+1]
		kvmap[key.Value] = YamlNode{
			key:  key,
			tag:  n.tag,
			Node: value,
		}
	}
	return kvmap, true
}

// ChangeKey replaces the oldkey with a new key/value pair. Useful for patching
// up a tree parsed from yaml but then used for validating an ignition structure
func (n *YamlNode) ChangeKey(oldKeyName, newKeyName string, newValue YamlNode) error {
	if n.Kind != yaml.MappingNode {
		return ErrNotMappingNode
	}
	for i := 0; i < len(n.Children); i += 2 {
		key := n.Children[i]
		if key.Value == oldKeyName {
			//key.Value = newKeyName
			(*n.Children[i]).Value = newKeyName
			*n.Children[i+1] = newValue.Node
			return nil
		}
	}

	return ErrKeyNotFound
}

// getIgnKeyName converts a snake_case (used by clct) to a camelCase (used by
// ignition)
func getIgnKeyName(keyname string) string {
	words := strings.Split(keyname, "_")
	for i, word := range words[1:] {
		words[i+1] = strings.Title(word)
	}
	return strings.Join(words, "")
}

func (n YamlNode) Tag() string {
	return n.tag
}

// ChangeTreeTag changes the value Tag() returns to newTag
func (n *YamlNode) ChangeTreeTag(newTag string) {
	n.tag = newTag
}
