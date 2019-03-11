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
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/url"
	"path"

	"github.com/coreos/container-linux-config-transpiler/config/astyaml"
	"github.com/coreos/container-linux-config-transpiler/internal/util"

	ignTypes "github.com/coreos/ignition/config/v2_2/types"
	"github.com/coreos/ignition/config/validate/astnode"
	"github.com/coreos/ignition/config/validate/report"
	"github.com/vincent-petithory/dataurl"
)

var (
	DefaultFileMode = 0644
	DefaultDirMode  = 0755

	WarningUnsetFileMode = fmt.Errorf("mode unspecified for file, defaulting to %#o", DefaultFileMode)
	WarningUnsetDirMode  = fmt.Errorf("mode unspecified for directory, defaulting to %#o", DefaultDirMode)

	ErrTooManyFileSources = errors.New("only one of the following can be set: local, inline, remote.url")
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
	User       *FileUser    `yaml:"user"`
	Group      *FileGroup   `yaml:"group"`
	Mode       *int         `yaml:"mode"`
	Contents   FileContents `yaml:"contents"`
	Overwrite  *bool        `yaml:"overwrite"`
	Append     bool         `yaml:"append"`
}

type FileContents struct {
	Remote Remote `yaml:"remote"`
	Inline string `yaml:"inline"`
	Local  string `yaml:"local"`
}

type Remote struct {
	Url          string       `yaml:"url"`
	Compression  string       `yaml:"compression"`
	Verification Verification `yaml:"verification"`
}

type Directory struct {
	Filesystem string     `yaml:"filesystem"`
	Path       string     `yaml:"path"`
	User       *FileUser  `yaml:"user"`
	Group      *FileGroup `yaml:"group"`
	Mode       *int       `yaml:"mode"`
	Overwrite  *bool      `yaml:"overwrite"`
}

type Link struct {
	Filesystem string     `yaml:"filesystem"`
	Path       string     `yaml:"path"`
	User       *FileUser  `yaml:"user"`
	Group      *FileGroup `yaml:"group"`
	Hard       bool       `yaml:"hard"`
	Target     string     `yaml:"target"`
	Overwrite  *bool      `yaml:"overwrite"`
}

func (f File) ValidateMode() report.Report {
	if f.Mode == nil {
		return report.ReportFromError(WarningUnsetFileMode, report.EntryWarning)
	}
	return report.Report{}
}

func (d Directory) ValidateMode() report.Report {
	if d.Mode == nil {
		return report.ReportFromError(WarningUnsetDirMode, report.EntryWarning)
	}
	return report.Report{}
}

func (fc FileContents) Validate() report.Report {
	i := 0
	if fc.Remote.Url != "" {
		i++
	}
	if fc.Inline != "" {
		i++
	}
	if fc.Local != "" {
		i++
	}
	if i > 1 {
		return report.ReportFromError(ErrTooManyFileSources, report.EntryError)
	}
	return report.Report{}
}

func init() {
	register(func(in Config, ast astnode.AstNode, out ignTypes.Config, platform string) (ignTypes.Config, report.Report, astnode.AstNode) {
		r := report.Report{}
		files_node, _ := getNodeChildPath(ast, "storage", "files")
		for i, file := range in.Storage.Files {
			if file.Mode == nil {
				file.Mode = util.IntToPtr(DefaultFileMode)
			}
			file_node, _ := getNodeChild(files_node, i)
			newFile := ignTypes.File{
				Node: ignTypes.Node{
					Filesystem: file.Filesystem,
					Path:       file.Path,
					Overwrite:  file.Overwrite,
				},
				FileEmbedded1: ignTypes.FileEmbedded1{
					Mode:   file.Mode,
					Append: file.Append,
				},
			}
			if file.User != nil {
				newFile.User = &ignTypes.NodeUser{
					ID:   file.User.Id,
					Name: file.User.Name,
				}
			}
			if file.Group != nil {
				newFile.Group = &ignTypes.NodeGroup{
					ID:   file.Group.Id,
					Name: file.Group.Name,
				}
			}

			if file.Contents.Inline != "" {
				newFile.Contents = ignTypes.FileContents{
					Source: (&url.URL{
						Scheme: "data",
						Opaque: "," + dataurl.EscapeString(file.Contents.Inline),
					}).String(),
				}
			}

			if file.Contents.Local != "" {
				// The provided local file path is relative to the value of the
				// --files-dir flag.
				filesDir := flag.Lookup("files-dir")
				if filesDir == nil || filesDir.Value.String() == "" {
					err := errors.New("local files require setting the --files-dir flag to the directory that contains the file")
					flagReport := report.ReportFromError(err, report.EntryError)
					if n, err := getNodeChildPath(file_node, "contents", "local"); err == nil {
						line, col, _ := n.ValueLineCol(nil)
						flagReport.AddPosition(line, col, "")
					}
					r.Merge(flagReport)
					continue
				}
				localPath := path.Join(filesDir.Value.String(), file.Contents.Local)
				contents, err := ioutil.ReadFile(localPath)
				if err != nil {
					// If the file could not be read, record error and continue.
					convertReport := report.ReportFromError(err, report.EntryError)
					if n, err := getNodeChildPath(file_node, "contents", "local"); err == nil {
						line, col, _ := n.ValueLineCol(nil)
						convertReport.AddPosition(line, col, "")
					}
					r.Merge(convertReport)
					continue
				}

				// Include the contents of the local file as if it were provided inline.
				newFile.Contents = ignTypes.FileContents{
					Source: (&url.URL{
						Scheme: "data",
						Opaque: "," + dataurl.Escape(contents),
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
					r.Merge(convertReport)
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
			if dir.Mode == nil {
				dir.Mode = util.IntToPtr(DefaultDirMode)
			}
			newDir := ignTypes.Directory{
				Node: ignTypes.Node{
					Filesystem: dir.Filesystem,
					Path:       dir.Path,
					Overwrite:  dir.Overwrite,
				},
				DirectoryEmbedded1: ignTypes.DirectoryEmbedded1{
					Mode: dir.Mode,
				},
			}
			if dir.User != nil {
				newDir.User = &ignTypes.NodeUser{
					ID:   dir.User.Id,
					Name: dir.User.Name,
				}
			}
			if dir.Group != nil {
				newDir.Group = &ignTypes.NodeGroup{
					ID:   dir.Group.Id,
					Name: dir.Group.Name,
				}
			}
			out.Storage.Directories = append(out.Storage.Directories, newDir)
		}
		for _, link := range in.Storage.Links {
			newLink := ignTypes.Link{
				Node: ignTypes.Node{
					Filesystem: link.Filesystem,
					Path:       link.Path,
					Overwrite:  link.Overwrite,
				},
				LinkEmbedded1: ignTypes.LinkEmbedded1{
					Hard:   link.Hard,
					Target: link.Target,
				},
			}
			if link.User != nil {
				newLink.User = &ignTypes.NodeUser{
					ID:   link.User.Id,
					Name: link.User.Name,
				}
			}
			if link.Group != nil {
				newLink.Group = &ignTypes.NodeGroup{
					ID:   link.Group.Id,
					Name: link.Group.Name,
				}
			}
			out.Storage.Links = append(out.Storage.Links, newLink)
		}
		return out, r, ast
	})
}
