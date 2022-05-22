/*
Copyright 2016 The Kubernetes Authors.

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

package preflight

import (
	"bytes"
	"fmt"
	"net"
	"net/http"
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/utils/exec"
)

type preflightCheckTest struct {
	msg string
}

func (pfct preflightCheckTest) Name() string {
	return "preflightCheckTest"
}

func (pfct preflightCheckTest) Check() (warning, errorList []error) {
	if pfct.msg == "warning" {
		return []error{errors.New("warning")}, nil
	}
	if pfct.msg != "" {
		return nil, []error{errors.New("fake error")}
	}
	return
}

func TestFileExistingCheck(t *testing.T) {
	f, err := os.CreateTemp("", "file-exist-check")
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}
	defer os.Remove(f.Name())
	var tests = []struct {
		name          string
		check         FileExistingCheck
		expectedError bool
	}{
		{
			name: "File does not exist, so it's not available",
			check: FileExistingCheck{
				Path: "/does/not/exist",
			},
			expectedError: true,
		},
		{
			name: "File exists and is available",
			check: FileExistingCheck{
				Path: f.Name(),
			},
			expectedError: false,
		},
	}
	for _, rt := range tests {
		output := rt.check.Check()
		if (output != nil) != rt.expectedError {
			t.Errorf(
				"Failed FileExistingCheck:%v\n\texpectedError: %t\n\t  actual: %t",
				rt.name,
				rt.expectedError,
				(output != nil),
			)
		}
	}
}

func TestFileAvailableCheck(t *testing.T) {
	f, err := os.CreateTemp("", "file-avail-check")
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}
	defer os.Remove(f.Name())
	var tests = []struct {
		name          string
		check         FileAvailableCheck
		expectedError bool
	}{
		{
			name: "The file does not exist",
			check: FileAvailableCheck{
				Path: "/does/not/exist",
			},
			expectedError: false,
		},
		{
			name: "The file exists",
			check: FileAvailableCheck{
				Path: f.Name(),
			},
			expectedError: true,
		},
	}
	for _, rt := range tests {
		output := rt.check.Check()
		if (output != nil) != rt.expectedError {
			t.Errorf(
				"Failed FileAvailableCheck:%v\n\texpectedError: %t\n\t  actual: %t",
				rt.name,
				rt.expectedError,
				(output != nil),
			)
		}
	}
}

func TestFileContentCheck(t *testing.T) {
	f, err := os.CreateTemp("", "file-content-check")
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}
	defer os.Remove(f.Name())
	var tests = []struct {
		name          string
		check         FileContentCheck
		expectedError bool
	}{
		{
			name: "File exists and has matching content",
			check: FileContentCheck{
				Path:    f.Name(),
				Content: []byte("Test FileContentCheck"),
			},
			expectedError: false,
		},
		{
			name: "File exists, content is nil",
			check: FileContentCheck{
				Path:    f.Name(),
				Content: nil,
			},
			expectedError: false,
		},
		{
			name: "File exists but has unexpected content",
			check: FileContentCheck{
				Path:    f.Name(),
				Content: []byte("foo"),
			},
			expectedError: true,
		},
		{
			name: "File does not exist, content is not nil",
			check: FileContentCheck{
				Path:    "/does/not/exist",
				Content: []byte("foo"),
			},
			expectedError: true,
		},
		{
			name: "File dose not exist, content is nil",
			check: FileContentCheck{
				Path:    "/does/not/exist",
				Content: nil,
			},
			expectedError: true,
		},
	}
	if _, err = f.WriteString("Test FileContentCheck"); err != nil {
		t.Fatalf("Failed to write to file: %v", err)
	}
	for _, rt := range tests {
		output := rt.check.Check()
		if (output != nil) != rt.expectedError {
			t.Errorf(
				"Failed FileContentCheck:%v\n\texpectedError: %t\n\t  actual: %t",
				rt.name,
				rt.expectedError,
				(output != nil),
			)
		}
	}
}

func TestDirAvailableCheck(t *testing.T) {
	fileDir, err := os.MkdirTemp("", "dir-avail-check")
	if err != nil {
		t.Fatalf("failed creating directory: %v", err)
	}
	defer os.RemoveAll(fileDir)
	var tests = []struct {
		name          string
		check         DirAvailableCheck
		expectedError bool
	}{
		{
			name: "Directory exists and is empty",
			check: DirAvailableCheck{
				Path: fileDir,
			},
			expectedError: false,
		},
		{
			name: "Directory exists and has something",
			check: DirAvailableCheck{
				Path: os.TempDir(), // a directory was created previously in this test
			},
			expectedError: true,
		},
		{
			name: "Directory does not exist",
			check: DirAvailableCheck{
				Path: "/does/not/exist",
			},
			expectedError: false,
		},
	}
	for _, rt := range tests {
		output := rt.check.Check()
		if (output != nil) != rt.expectedError {
			t.Errorf(
				"Failed DirAvailableCheck:%v\n\texpectedError: %t\n\t  actual: %t",
				rt.name,
				rt.expectedError,
				(output != nil),
			)
		}
	}
}

func TestPortOpenCheck(t *testing.T) {
	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("could not listen on local network: %v", err)
	}
	defer ln.Close()
	var tests = []struct {
		name          string
		check         PortOpenCheck
		expectedError bool
	}{
		{
			name:          "Port is available",
			check:         PortOpenCheck{port: 0},
			expectedError: false,
		},
		{
			name:          "Port is not available",
			check:         PortOpenCheck{port: ln.Addr().(*net.TCPAddr).Port},
			expectedError: true,
		},
	}
	for _, rt := range tests {
		output := rt.check.Check()
		if (output != nil) != rt.expectedError {
			t.Errorf(
				"Failed PortOpenCheck:%v\n\texpectedError: %t\n\t  actual: %t",
				rt.name,
				rt.expectedError,
				(output != nil),
			)
		}
	}
}

func TestRunChecks(t *testing.T) {
	var tokenTest = []struct {
		p        []Checker
		expected bool
		output   string
	}{
		{[]Checker{}, true, ""},
		// {[]Checker{preflightCheckTest{"warning"}}, true, "\t[WARNING preflightCheckTest]: warning\n"}, // should just print warning
		// {[]Checker{preflightCheckTest{"error"}}, false, ""},
		// {[]Checker{preflightCheckTest{"test"}}, false, ""},
		{[]Checker{DirAvailableCheck{Path: "/does/not/exist"}}, true, ""},
		{[]Checker{DirAvailableCheck{Path: "/"}}, false, ""},
		{[]Checker{FileAvailableCheck{Path: "/does/not/exist"}}, true, ""},
		{[]Checker{FileContentCheck{Path: "/does/not/exist"}}, false, ""},
		{[]Checker{FileContentCheck{Path: "/"}}, true, ""},
		{[]Checker{FileContentCheck{Path: "/", Content: []byte("does not exist")}}, false, ""},
		{[]Checker{InPathCheck{executable: "foobarbaz", exec: exec.New()}}, true, "\t[WARNING FileExisting-foobarbaz]: foobarbaz not found in system path\n"},
		{[]Checker{InPathCheck{executable: "foobarbaz", mandatory: true, exec: exec.New()}}, false, ""},
		{[]Checker{InPathCheck{executable: "foobar", mandatory: false, exec: exec.New(), suggestion: "install foobar"}}, true, "\t[WARNING FileExisting-foobar]: foobar not found in system path\nSuggestion: install foobar\n"},
	}
	for _, rt := range tokenTest {
		buf := new(bytes.Buffer)
		actual := RunChecks(rt.p, buf, sets.NewString())
		if (actual == nil) != rt.expected {
			t.Errorf(
				"failed RunChecks:\n\texpected: %t\n\t  actual: %t",
				rt.expected,
				(actual == nil),
			)
		}
		if buf.String() != rt.output {
			t.Errorf(
				"failed RunChecks:\n\texpected: %s\n\t  actual: %s",
				rt.output,
				buf.String(),
			)
		}
	}
}

func TestHTTPProxyCIDRCheck(t *testing.T) {
	var tests = []struct {
		check          HTTPProxyCIDRCheck
		expectWarnings bool
	}{
		{
			check: HTTPProxyCIDRCheck{
				Proto: "https",
				CIDR:  "127.0.0.0/8",
			}, // Loopback addresses never should produce proxy warnings
			expectWarnings: false,
		},
		{
			check: HTTPProxyCIDRCheck{
				Proto: "https",
				CIDR:  "10.96.0.0/12",
			}, // Expected to be accessed directly, we set NO_PROXY to 10.0.0.0/8
			expectWarnings: false,
		},
		{
			check: HTTPProxyCIDRCheck{
				Proto: "https",
				CIDR:  "192.168.0.0/16",
			}, // Expected to go via proxy as this range is not listed in NO_PROXY
			expectWarnings: true,
		},
		{
			check: HTTPProxyCIDRCheck{
				Proto: "https",
				CIDR:  "2001:db8::/56",
			}, // Expected to be accessed directly, part of 2001:db8::/48 in NO_PROXY
			expectWarnings: false,
		},
		{
			check: HTTPProxyCIDRCheck{
				Proto: "https",
				CIDR:  "2001:db8:1::/56",
			}, // Expected to go via proxy, range is not in 2001:db8::/48
			expectWarnings: true,
		},
	}

	// Save current content of *_proxy and *_PROXY variables.
	savedEnv := resetProxyEnv(t)
	defer restoreEnv(savedEnv)

	for _, rt := range tests {
		warning := rt.check.Check()
		if (warning != nil) != rt.expectWarnings {
			t.Errorf(
				"failed HTTPProxyCIDRCheck:\n\texpected: %t\n\t  actual: %t (CIDR:%s). Warnings: %v",
				rt.expectWarnings,
				(warning != nil),
				rt.check.CIDR,
				warning,
			)
		}
	}
}

func TestHTTPProxyCheck(t *testing.T) {
	var tests = []struct {
		name           string
		check          HTTPProxyCheck
		expectWarnings bool
	}{
		{
			name: "Loopback address",
			check: HTTPProxyCheck{
				Proto: "https",
				Host:  "127.0.0.1",
			}, // Loopback addresses never should produce proxy warnings
			expectWarnings: false,
		},
		{
			name: "IPv4 direct access",
			check: HTTPProxyCheck{
				Proto: "https",
				Host:  "10.96.0.1",
			}, // Expected to be accessed directly, we set NO_PROXY to 10.0.0.0/8
			expectWarnings: false,
		},
		{
			name: "IPv4 via proxy",
			check: HTTPProxyCheck{
				Proto: "https",
				Host:  "192.168.0.1",
			}, // Expected to go via proxy as this range is not listed in NO_PROXY
			expectWarnings: true,
		},
		{
			name: "IPv6 direct access",
			check: HTTPProxyCheck{
				Proto: "https",
				Host:  "[2001:db8::1:15]",
			}, // Expected to be accessed directly, part of 2001:db8::/48 in NO_PROXY
			expectWarnings: false,
		},
		{
			name: "IPv6 via proxy",
			check: HTTPProxyCheck{
				Proto: "https",
				Host:  "[2001:db8:1::1:15]",
			}, // Expected to go via proxy, range is not in 2001:db8::/48
			expectWarnings: true,
		},
		{
			name: "IPv6 direct access, no brackets",
			check: HTTPProxyCheck{
				Proto: "https",
				Host:  "2001:db8::1:15",
			}, // Expected to be accessed directly, part of 2001:db8::/48 in NO_PROXY
			expectWarnings: false,
		},
		{
			name: "IPv6 via proxy, no brackets",
			check: HTTPProxyCheck{
				Proto: "https",
				Host:  "2001:db8:1::1:15",
			}, // Expected to go via proxy, range is not in 2001:db8::/48
			expectWarnings: true,
		},
	}

	// Save current content of *_proxy and *_PROXY variables.
	savedEnv := resetProxyEnv(t)
	defer restoreEnv(savedEnv)

	for _, rt := range tests {
		warning := rt.check.Check()
		if (warning != nil) != rt.expectWarnings {
			t.Errorf(
				"%s failed HTTPProxyCheck:\n\texpected: %t\n\t  actual: %t (Host:%s). Warnings: %v",
				rt.name,
				rt.expectWarnings,
				(warning != nil),
				rt.check.Host,
				warning,
			)
		}
	}
}

// resetProxyEnv is helper function that unsets all *_proxy variables
// and return previously set values as map. This can be used to restore
// original state of the environment.
func resetProxyEnv(t *testing.T) map[string]string {
	savedEnv := make(map[string]string)
	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		if strings.HasSuffix(strings.ToLower(pair[0]), "_proxy") {
			savedEnv[pair[0]] = pair[1]
			os.Unsetenv(pair[0])
		}
	}
	t.Log("Saved environment: ", savedEnv)

	os.Setenv("HTTP_PROXY", "http://proxy.example.com:3128")
	os.Setenv("HTTPS_PROXY", "https://proxy.example.com:3128")
	os.Setenv("NO_PROXY", "example.com,10.0.0.0/8,2001:db8::/48")
	// Check if we can reliably execute tests:
	// ProxyFromEnvironment caches the *_proxy environment variables and
	// if ProxyFromEnvironment already executed before our test with empty
	// HTTP_PROXY it will make these tests return false positive failures
	req, err := http.NewRequest("GET", "http://host.fake.tld/", nil)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	proxy, err := http.ProxyFromEnvironment(req)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if proxy == nil {
		t.Skip("test skipped as ProxyFromEnvironment already initialized in environment without defined HTTP proxy")
	}
	t.Log("http.ProxyFromEnvironment is usable, continue executing test")
	return savedEnv
}

// restoreEnv is helper function to restores values
// of environment variables from saved state in the map
func restoreEnv(e map[string]string) {
	for k, v := range e {
		os.Setenv(k, v)
	}
}

func TestSetHasItemOrAll(t *testing.T) {
	var tests = []struct {
		ignoreSet      sets.String
		testString     string
		expectedResult bool
	}{
		{sets.NewString(), "foo", false},
		{sets.NewString("all"), "foo", true},
		{sets.NewString("all", "bar"), "foo", true},
		{sets.NewString("bar"), "foo", false},
		{sets.NewString("baz", "foo", "bar"), "foo", true},
		{sets.NewString("baz", "bar", "foo"), "Foo", true},
	}

	for _, rt := range tests {
		t.Run(rt.testString, func(t *testing.T) {
			result := setHasItemOrAll(rt.ignoreSet, rt.testString)
			if result != rt.expectedResult {
				t.Errorf(
					"setHasItemOrAll: expected: %v actual: %v (arguments: %q, %q)",
					rt.expectedResult, result,
					rt.ignoreSet,
					rt.testString,
				)
			}
		})
	}
}

func TestNumCPUCheck(t *testing.T) {
	var tests = []struct {
		numCPU      int
		numErrors   int
		numWarnings int
	}{
		{0, 0, 0},
		{999999999, 1, 0},
	}

	for _, rt := range tests {
		t.Run(fmt.Sprintf("number of CPUs: %d", rt.numCPU), func(t *testing.T) {
			errors := NumCPUCheck{NumCPU: rt.numCPU}.Check()

			if errors != nil {
				t.Errorf("expected %d warning(s) but err: %q", rt.numErrors, errors)
			}
		})
	}
}

func TestMemCheck(t *testing.T) {
	// skip this test, if OS in not Linux, since it will ONLY pass on Linux.
	if runtime.GOOS != "linux" {
		t.Skip("unsupported OS for memory check test ")
	}

	var tests = []struct {
		minimum        uint64
		expectedErrors int
	}{
		{0, 0},
		{9999999999999999, 1},
	}

	for _, rt := range tests {
		t.Run(fmt.Sprintf("MemoryCheck{%d}", rt.minimum), func(t *testing.T) {
			errors := MemCheck{Mem: rt.minimum}.Check()
			if errors != nil {
				t.Errorf("expected 0 warnings but got %q", errors)
			}
		})
	}
}
