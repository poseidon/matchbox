// Copyright 2015 CoreOS, Inc.
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

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/coreos/fuze/config"
)

func stderr(f string, a ...interface{}) {
	out := fmt.Sprintf(f, a...)
	fmt.Fprintln(os.Stderr, strings.TrimSuffix(out, "\n"))
}

func main() {
	flags := struct {
		help    bool
		pretty  bool
		inFile  string
		outFile string
	}{}

	flag.BoolVar(&flags.help, "help", false, "print help and exit")
	flag.BoolVar(&flags.pretty, "pretty", false, "indent the output file")
	flag.StringVar(&flags.inFile, "in-file", "/dev/stdin", "input file (YAML)")
	flag.StringVar(&flags.outFile, "out-file", "/dev/stdout", "output file (JSON)")

	flag.Parse()

	if flags.help {
		flag.Usage()
		return
	}

	dataIn, err := ioutil.ReadFile(flags.inFile)
	if err != nil {
		stderr("Failed to read: %v", err)
		os.Exit(1)
	}

	fuzeCfg, report := config.Parse(dataIn)
	stderr(report.String())
	if report.IsFatal() {
		stderr("Failed to parse fuze config")
		os.Exit(1)
	}

	cfg, report := config.ConvertAs2_0_0(fuzeCfg)
	stderr(report.String())
	if report.IsFatal() {
		stderr("Generated Ignition config was invalid.")
		os.Exit(1)
	}

	var dataOut []byte
	if flags.pretty {
		dataOut, err = json.MarshalIndent(&cfg, "", "  ")
		dataOut = append(dataOut, '\n')
	} else {
		dataOut, err = json.Marshal(&cfg)
	}
	if err != nil {
		stderr("Failed to marshal output: %v", err)
		os.Exit(1)
	}

	if err := ioutil.WriteFile(flags.outFile, dataOut, 0640); err != nil {
		stderr("Failed to write: %v", err)
		os.Exit(1)
	}
}
