/*
Copyright 2016 The Kubernetes Authors All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package command

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"time"

	"golang.org/x/sync/syncmap"

	"github.com/golang/glog"
	"github.com/pkg/errors"

	"k8s.io/minikube/pkg/minikube/assets"
)

// FakeCommandRunner mocks command output without running the Commands
//
// It implements the CommandRunner interface and is used for testing.
type FakeCommandRunner struct {
	cmdMap  syncmap.Map
	fileMap syncmap.Map
}

// NewFakeCommandRunner returns a new FakeCommandRunner
//
// The expected output of commands should be set with SetCommandToOutput
func NewFakeCommandRunner() *FakeCommandRunner {
	return &FakeCommandRunner{}
}

// RunCmd implements the Command Runner interface to run a exec.Cmd object
func (f *FakeCommandRunner) RunCmd(cmd *exec.Cmd) (*RunResult, error) {
	rr := &RunResult{Args: cmd.Args}
	glog.Infof("(FakeCommandRunner) Run:  %v", rr.Command())

	start := time.Now()

	out, ok := f.cmdMap.Load(strings.Join(rr.Args, " "))
	if !ok {
		return rr, fmt.Errorf("unregistered command: `%s`\nexpected: %v\n", rr.Command(), f.commands())
	}

	var buf bytes.Buffer
	outStr := ""
	if out != nil {
		outStr = out.(string)
	}
	_, err := buf.WriteString(outStr)
	if err != nil {
		return rr, errors.Wrap(err, "Writing outStr to FakeCommandRunner's buffer")
	}
	rr.Stdout = buf
	rr.Stderr = buf

	elapsed := time.Since(start)

	// Reduce log spam
	if elapsed > (1 * time.Second) {
		glog.Infof("(FakeCommandRunner) Done: %v: (%s)", rr.Command(), elapsed)
	}
	return rr, nil
}

// Copy adds the filename, file contents key value pair to the stored map.
func (f *FakeCommandRunner) Copy(file assets.CopyableFile) error {
	var b bytes.Buffer
	_, err := io.Copy(&b, file)
	if err != nil {
		return errors.Wrapf(err, "error reading file: %+v", file)
	}
	f.fileMap.Store(file.GetAssetName(), b.String())
	return nil
}

// Remove removes the filename, file contents key value pair from the stored map
func (f *FakeCommandRunner) Remove(file assets.CopyableFile) error {
	f.fileMap.Delete(file.GetAssetName())
	return nil
}

// SetFileToContents stores the file to contents map for the FakeCommandRunner
func (f *FakeCommandRunner) SetFileToContents(fileToContents map[string]string) {
	for k, v := range fileToContents {
		f.fileMap.Store(k, v)
	}
}

// SetCommandToOutput stores the file to contents map for the FakeCommandRunner
func (f *FakeCommandRunner) SetCommandToOutput(cmdToOutput map[string]string) {
	for k, v := range cmdToOutput {
		glog.Infof("fake command %q -> %q", k, v)
		f.cmdMap.Store(k, v)
	}
}

// GetFileToContents stores the file to contents map for the FakeCommandRunner
func (f *FakeCommandRunner) GetFileToContents(filename string) (string, error) {
	contents, ok := f.fileMap.Load(filename)
	if !ok {
		return "", fmt.Errorf("unavailable file: %s", filename)
	}
	return contents.(string), nil
}

func (f *FakeCommandRunner) commands() []string {
	cmds := []string{}
	f.cmdMap.Range(func(k, v interface{}) bool {
		cmds = append(cmds, fmt.Sprintf("`%s`", k))
		return true
	})
	return cmds
}

// DumpMaps prints out the list of stored commands and stored filenames.
func (f *FakeCommandRunner) DumpMaps(w io.Writer) {
	fmt.Fprintln(w, "Commands:")
	f.cmdMap.Range(func(k, v interface{}) bool {
		fmt.Fprintf(w, "%s:%s", k, v)
		return true
	})
	fmt.Fprintln(w, "Filenames: ")
	f.fileMap.Range(func(k, v interface{}) bool {
		fmt.Fprint(w, k)
		return true
	})
}
