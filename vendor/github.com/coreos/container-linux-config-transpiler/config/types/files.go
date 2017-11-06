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

package types

import (
	"net/url"

	"github.com/coreos/container-linux-config-transpiler/config/astyaml"

	ignTypes "github.com/coreos/ignition/config/v2_1/types"
	"github.com/coreos/ignition/config/validate/astnode"
	"github.com/coreos/ignition/config/validate/report"
	"github.com/vincent-petithory/dataurl"
)

type FileUser struct {
	Id   *int   `yaml:"id"`
	Name string `yaml:"name"`
}

type FileGroup struct {
	Id   *int   `yaml:"id"`
	Name string `yaml:"name"`
}

type File struct {
	Filesystem string       `yaml:"filesystem"`
	Path       string       `yaml:"path"`
	User       FileUser     `yaml:"user"`
	Group      FileGroup    `yaml:"group"`
	Mode       int          `yaml:"mode"`
	Contents   FileContents `yaml:"contents"`
}

type FileContents struct {
	Remote Remote `yaml:"remote"`
	Inline string `yaml:"inline"`
}

type Remote struct {
	Url          string       `yaml:"url"`
	Compression  string       `yaml:"compression"`
	Verification Verification `yaml:"verification"`
}

type Directory struct {
	Filesystem string    `yaml:"filesystem"`
	Path       string    `yaml:"path"`
	User       FileUser  `yaml:"user"`
	Group      FileGroup `yaml:"group"`
	Mode       int       `yaml:"mode"`
}

type Link struct {
	Filesystem string    `yaml:"filesystem"`
	Path       string    `yaml:"path"`
	User       FileUser  `yaml:"user"`
	Group      FileGroup `yaml:"group"`
	Hard       bool      `yaml:"hard"`
	Target     string    `yaml:"target"`
}

func init() {
	register2_0(func(in Config, ast astnode.AstNode, out ignTypes.Config, platform string) (ignTypes.Config, report.Report, astnode.AstNode) {
		r := report.Report{}
		files_node, _ := getNodeChildPath(ast, "storage", "files")
		for i, file := range in.Storage.Files {
			file_node, _ := getNodeChild(files_node, i)
			newFile := ignTypes.File{
				Node: ignTypes.Node{
					Filesystem: file.Filesystem,
					Path:       file.Path,
					User: ignTypes.NodeUser{
						ID:   file.User.Id,
						Name: file.User.Name,
					},
					Group: ignTypes.NodeGroup{
						ID:   file.Group.Id,
						Name: file.Group.Name,
					},
				},
				FileEmbedded1: ignTypes.FileEmbedded1{
					Mode: file.Mode,
				},
			}

			if file.Contents.Inline != "" {
				newFile.Contents = ignTypes.FileContents{
					Source: (&url.URL{
						Scheme: "data",
						Opaque: "," + dataurl.EscapeString(file.Contents.Inline),
					}).String(),
				}
			}

			if file.Contents.Remote.Url != "" {
				source, err := url.Parse(file.Contents.Remote.Url)
				if err != nil {
					// if invalid, record error and continue
					convertReport := report.ReportFromError(err, report.EntryError)
					if n, err := getNodeChildPath(file_node, "contents", "remote", "url"); err == nil {
						line, col, _ := n.ValueLineCol(nil)
						convertReport.AddPosition(line, col, "")
					}
					continue
				}

				// patch the yaml tree to look like the ignition tree by making contents
				// the remote section and changing the name from url -> source
				asYamlNode, ok := file_node.(astyaml.YamlNode)
				if ok {
					newContents, _ := getNodeChildPath(file_node, "contents", "remote")
					newContentsAsYaml := newContents.(astyaml.YamlNode)
					asYamlNode.ChangeKey("contents", "contents", newContentsAsYaml)

					url, _ := getNodeChild(newContents.(astyaml.YamlNode), "url")
					newContentsAsYaml.ChangeKey("url", "source", url.(astyaml.YamlNode))
				}

				newFile.Contents = ignTypes.FileContents{Source: source.String()}

			}

			if newFile.Contents == (ignTypes.FileContents{}) {
				newFile.Contents = ignTypes.FileContents{
					Source: "data:,",
				}
			}

			newFile.Contents.Compression = file.Contents.Remote.Compression
			newFile.Contents.Verification = convertVerification(file.Contents.Remote.Verification)

			out.Storage.Files = append(out.Storage.Files, newFile)
		}
		for _, dir := range in.Storage.Directories {
			out.Storage.Directories = append(out.Storage.Directories, ignTypes.Directory{
				Node: ignTypes.Node{
					Filesystem: dir.Filesystem,
					Path:       dir.Path,
					User: ignTypes.NodeUser{
						ID:   dir.User.Id,
						Name: dir.User.Name,
					},
					Group: ignTypes.NodeGroup{
						ID:   dir.Group.Id,
						Name: dir.Group.Name,
					},
				},
				DirectoryEmbedded1: ignTypes.DirectoryEmbedded1{
					Mode: dir.Mode,
				},
			})
		}
		for _, link := range in.Storage.Links {
			out.Storage.Links = append(out.Storage.Links, ignTypes.Link{
				Node: ignTypes.Node{
					Filesystem: link.Filesystem,
					Path:       link.Path,
					User: ignTypes.NodeUser{
						ID:   link.User.Id,
						Name: link.User.Name,
					},
					Group: ignTypes.NodeGroup{
						ID:   link.Group.Id,
						Name: link.Group.Name,
					},
				},
				LinkEmbedded1: ignTypes.LinkEmbedded1{
					Hard:   link.Hard,
					Target: link.Target,
				},
			})
		}
		return out, r, ast
	})
}
