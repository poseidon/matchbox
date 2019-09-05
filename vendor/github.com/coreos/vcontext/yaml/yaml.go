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

package json

import (
	"github.com/coreos/vcontext/tree"

	"gopkg.in/yaml.v3"
)

func UnmarshalToContext(raw []byte) (tree.Node, error) {
	var ast yaml.Node
	if err := yaml.Unmarshal(raw, &ast); err != nil {
		return nil, err
	}
	return fromYamlNode(ast), nil
}

func fromYamlNode(n yaml.Node) tree.Node {
	m := tree.Marker{
		StartP: &tree.Pos{
			Line:   int64(n.Line),
			Column: int64(n.Column),
		},
	}
	switch n.Kind {
	case 0:
		// empty
		return nil
	case yaml.DocumentNode:
		if len(n.Content) == 0 {
			return nil
		}
		return fromYamlNode(*n.Content[0])
	case yaml.MappingNode:
		ret := tree.MapNode{
			Marker:   m,
			Children: make(map[string]tree.Node, len(n.Content)/2),
			Keys:     make(map[string]tree.Leaf, len(n.Content)/2),
		}
		// MappingNodes list keys and values like [k, v, k, v...]
		for i := 0; i < len(n.Content); i += 2 {
			key := *n.Content[i]
			value := *n.Content[i+1]
			ret.Keys[key.Value] = tree.Leaf{
				Marker: tree.Marker{
					StartP: &tree.Pos{
						Line:   int64(key.Line),
						Column: int64(key.Column),
					},
				},
			}
			ret.Children[key.Value] = fromYamlNode(value)
		}
		return ret
	case yaml.SequenceNode:
		ret := tree.SliceNode{
			Marker:   m,
			Children: make([]tree.Node, 0, len(n.Content)),
		}
		for _, child := range n.Content {
			ret.Children = append(ret.Children, fromYamlNode(*child))
		}
		return ret
	default: // scalars and aliases
		return tree.Leaf{
			Marker: m,
		}
	}
}
