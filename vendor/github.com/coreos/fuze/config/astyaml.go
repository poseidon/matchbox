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

package config

import (
	"errors"
	"io"

	yaml "github.com/ajeddeloh/yaml"
	"github.com/coreos/ignition/config/validate"
)

var (
	ErrNotDocumentNode = errors.New("Can only convert from document node")
)

type YamlNode struct {
	key yaml.Node
	yaml.Node
}

func FromYamlDocumentNode(n yaml.Node) (YamlNode, error) {
	if n.Kind != yaml.DocumentNode {
		return YamlNode{}, ErrNotDocumentNode
	}

	return YamlNode{
		key:  n,
		Node: *n.Children[0],
	}, nil
}

func (n YamlNode) ValueLineCol(source io.ReadSeeker) (int, int, string) {
	return n.Line, n.Column, ""
}

func (n YamlNode) KeyLineCol(source io.ReadSeeker) (int, int, string) {
	return n.key.Line, n.key.Column, ""
}

func (n YamlNode) LiteralValue() interface{} {
	return n.Value
}

func (n YamlNode) SliceChild(index int) (validate.AstNode, bool) {
	if n.Kind != yaml.SequenceNode {
		return nil, false
	}
	if index >= len(n.Children) {
		return nil, false
	}

	return YamlNode{
		key:  yaml.Node{},
		Node: *n.Children[index],
	}, true
}

func (n YamlNode) KeyValueMap() (map[string]validate.AstNode, bool) {
	if n.Kind != yaml.MappingNode {
		return nil, false
	}

	kvmap := map[string]validate.AstNode{}
	for i := 0; i < len(n.Children); i += 2 {
		key := *n.Children[i]
		value := *n.Children[i+1]
		kvmap[key.Value] = YamlNode{
			key:  key,
			Node: value,
		}
	}
	return kvmap, true
}

func (n YamlNode) Tag() string {
	return "yaml"
}
