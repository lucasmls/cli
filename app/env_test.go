// The MIT License
//
// Copyright (c) 2022 Temporal Technologies Inc.  All rights reserved.
//
// Copyright (c) 2020 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package app_test

import (
	"log"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

const (
	testEnvName = "tctl-test-env"
)

func (s *cliAppSuite) TestSetEnvValue() {
	defer setupConfig(s.app)()

	err := s.app.Run([]string{"", "env", "set", testEnvName + ".address", "0.0.0.0:00000"})
	s.NoError(err)

	config := readConfig()
	s.Contains(config, "tctl-test-env:")
	s.Contains(config, "address: 0.0.0.0:00000")
	s.Contains(config, "namespace: tctl-test-namespace")
}

func (s *cliAppSuite) TestDeleteEnvProperty() {
	defer setupConfig(s.app)()

	err := s.app.Run([]string{"", "env", "set", testEnvName + ".address", "1.2.3.4:5678"})
	s.NoError(err)

	err = s.app.Run([]string{"", "env", "delete", testEnvName + ".address"})
	s.NoError(err)

	config := readConfig()
	s.Contains(config, "tctl-test-env:")
	s.Contains(config, "namespace: tctl-test-namespace")
	s.NotContains(config, "address: 1.2.3.4:5678")
}

func (s *cliAppSuite) TestDeleteEnv() {
	defer setupConfig(s.app)()

	err := s.app.Run([]string{"", "env", "set", testEnvName + ".address", "1.2.3.4:5678"})
	s.NoError(err)

	err = s.app.Run([]string{"", "env", "delete", testEnvName})
	s.NoError(err)

	config := readConfig()
	s.NotContains(config, "tctl-test-env:")
	s.NotContains(config, "namespace: tctl-test-namespace")
	s.NotContains(config, "address: 1.2.3.4:5678")
}

func setupConfig(app *cli.App) func() {
	err := app.Run([]string{"", "env", "set", testEnvName + ".namespace", "tctl-test-namespace"})
	if err != nil {
		log.Fatal(err)
	}

	return func() {
		err := app.Run([]string{"", "env", "delete", testEnvName})
		if err != nil {
			log.Printf("unable to unset test env: %s", err)
		}
	}
}

func readConfig() string {
	path := getConfigPath()
	content, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return string(content)
}

func getConfigPath() string {
	dpath, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	path := filepath.Join(dpath, ".config", "temporalio", "temporal.yaml")

	return path
}
