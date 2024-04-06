// MIT License
//
// Copyright (c) 2023 Bret McKee
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package goconfig

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/spf13/pflag"
)

const (
	testBadDelimiter   = "---"
	testBadFileName    = "/this/file/does/not/exist"
	testDataDir        = "testdata"
	testDefaultValue1  = 1
	testDefaultValue2  = 2
	testDefaultValue3  = 3
	testDelimiter      = "."
	testEnv1           = "testenv"
	testFlagsetName    = "TestFlagsetName"
	testInvalidOption  = "-this-is-a-bad-option"
	testKey1           = "value1"
	testKey2           = "value2"
	testKey3           = "value3"
	testNestedTag      = "nested"
	testNestedKey      = "nestedvalue"
	testNoHelpMessage  = ""
	testNonInteger     = "this is not an integer"
	testPrefix         = "TEST_"
	testValue1         = 101
	testValue2         = 102
	testValue3         = 103
	testGoodJSONConfig = "good.json" // Sets value=101 val=102
)

type nameValue struct {
	name  string
	value string
}

type testConfig1 struct {
	NestedVal int `koanf:"nestedvalue"`
}

type testConfig struct {
	Value1 int         `koanf:"value1"`
	Value2 int         `koanf:"value2"`
	Value3 int         `koanf:"value3"`
	Nested testConfig1 `koanf:"nested"`
}

func TestNew(t *testing.T) {
	c, err := New(testPrefix, testDelimiter)
	if err != nil {
		t.Fatalf("New err: got=%v want=nil", err)
	}
	if got, want := c.prefix, testPrefix; got != want {
		t.Errorf("New prefix: got=%q want=%q", got, want)
	}
	if got, want := c.delimiter, testDelimiter; got != want {
		t.Errorf("New delimiter: got=%q want=%q", got, want)
	}
}

func TestNewError(t *testing.T) {
	if _, err := New(testPrefix, testBadDelimiter); err == nil {
		t.Fatalf("New err: got=nil want=non-nil")
	}
}

func Test_updateEnv(t *testing.T) {
	cases := []struct {
		name    string
		value   string
		want    string
		wantErr bool
	}{
		{
			name: "Empty String",
		},
		{
			name:  "No underscores",
			value: testEnv1,
			want:  testEnv1,
		},
		{
			name:  "No underscores uppercase",
			value: strings.ToUpper(testEnv1),
			want:  testEnv1,
		},
		{
			name:  "underscore",
			value: testEnv1 + "_" + testEnv1,
			want:  testEnv1 + "." + testEnv1,
		},
		{
			name:  "underscore uppercase",
			value: strings.ToUpper(testEnv1 + "_" + testEnv1),
			want:  testEnv1 + "." + testEnv1,
		},
		{
			name:  "multiple underscores",
			value: testEnv1 + "_" + testEnv1 + "_",
			want:  testEnv1 + "." + testEnv1 + ".",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c, err := New(testPrefix, testDelimiter)
			if err != nil {
				t.Fatalf("New failed unexpectedly: %v", err)
			}
			if got, want := c.updateEnv(tc.value), tc.want; got != want {
				t.Errorf("updateEnv: got=%q want=%q", got, want)
			}
		})
	}
}

func TestLoadUnchangedForNoInput(t *testing.T) {
	var got, want testConfig

	f := pflag.NewFlagSet(testFlagsetName, pflag.ExitOnError)

	c, err := New(testPrefix, testDelimiter)
	if err != nil {
		t.Fatalf("New failed unexpectedly: %v", err)
	}
	if err := c.Load(f, &got); err != nil {
		t.Fatalf("Load err: got=%v, want=nil", err)
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Set() mismatch (-want +got):\n%s", diff)
	}
}

func TestLoadViaEnv(t *testing.T) {
	cases := []struct {
		name    string
		nvs     []nameValue
		want    testConfig
		wantErr bool
	}{

		{
			name: "no values",
		},
		{
			name: "values from env lowercase",
			nvs: []nameValue{
				{testPrefix + testKey1, fmt.Sprintf("%v", testValue1)},
				{testPrefix + testNestedTag + "_" + testNestedKey, fmt.Sprintf("%v", testValue2)},
			},
			want: testConfig{
				Value1: testValue1,
				Nested: testConfig1{
					NestedVal: testValue2,
				},
			},
		},
		{
			name: "values from env uppercase",
			nvs: []nameValue{
				{strings.ToUpper(testPrefix + testKey1), fmt.Sprintf("%v", testValue1)},
				{strings.ToUpper(testPrefix + testNestedTag + "_" + testNestedKey), fmt.Sprintf("%v", testValue2)},
			},
			want: testConfig{
				Value1: testValue1,
				Nested: testConfig1{
					NestedVal: testValue2,
				},
			},
		},
		{
			name: "bad values",
			nvs: []nameValue{
				{testPrefix + testKey1, testNonInteger},
				{testPrefix + testNestedTag + "_" + testNestedKey, testNonInteger},
			},
			wantErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			for _, nv := range tc.nvs {
				if err := os.Setenv(nv.name, nv.value); err != nil {
					t.Fatalf("os.Setenv failed unexpectedl: %v", err)
				}
				defer func(k string) {
					if err := os.Unsetenv(k); err != nil {
						t.Fatalf("os.Unsetenv failed unexpectedly: %v", err)
					}
				}(nv.name)
			}

			f := pflag.NewFlagSet(testFlagsetName, pflag.ExitOnError)

			c, err := New(testPrefix, testDelimiter)
			if err != nil {
				t.Fatalf("New failed unexpectedly: %v", err)
			}

			var cfg testConfig
			err = c.Load(f, &cfg)

			if tc.wantErr {
				if err == nil {
					t.Errorf("Load err: got=nil want=<non-nil>")
				}
				return
			}
			if err != nil {
				t.Fatalf("Load err: got=%v want=nil", err)
			}
			if diff := cmp.Diff(tc.want, cfg); diff != "" {
				t.Errorf("Load cfg mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestLoadViaFlag(t *testing.T) {
	cases := []struct {
		name         string
		args         []string
		want         testConfig
		wantParseErr bool
		wantLoadErr  bool
	}{

		{
			name: "invalid option",
			args: []string{
				fmt.Sprintf("%s=%d", testInvalidOption, testValue1),
			},
			wantParseErr: true,
		},
		{
			name: "loads default value",
			want: testConfig{
				Value1: testDefaultValue1,
				Nested: testConfig1{
					NestedVal: testDefaultValue2,
				},
			},
		},
		{
			name: "good values",
			args: []string{
				fmt.Sprintf("--%s=%d", testKey1, testValue1),
				fmt.Sprintf("--%s.%s=%d", testNestedTag, testNestedKey, testValue2),
			},
			want: testConfig{
				Value1: testValue1,
				Nested: testConfig1{
					NestedVal: testValue2,
				},
			},
		},
		{
			name:        "bad value",
			wantLoadErr: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f := pflag.NewFlagSet(testFlagsetName, pflag.ContinueOnError)
			if tc.wantLoadErr {
				f.String(testKey1, testNonInteger, testNoHelpMessage)
			} else {
				f.Int(testKey1, testDefaultValue1, testNoHelpMessage)
			}
			f.Int(testNestedTag+"."+testNestedKey, testDefaultValue2, testNoHelpMessage)
			// TODO: f.SetOutput()

			err := f.Parse(tc.args)
			if tc.wantParseErr {
				if err == nil {
					t.Fatalf("f.Parse got=nil want=non-nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("f.Parse got=%v want=nil", err)
			}

			c, err := New(testPrefix, testDelimiter)
			if err != nil {
				t.Fatalf("New failed unexpectedly: %v", err)
			}

			var cfg testConfig
			err = c.Load(f, &cfg)
			if tc.wantLoadErr {
				if err == nil {
					t.Errorf("Load err: got=nil want=<non-nil>")
				}
				return
			}
			if err != nil {
				t.Fatalf("Load err: got=%v want=nil", err)
			}
			if diff := cmp.Diff(tc.want, cfg); diff != "" {
				t.Errorf("Load cfg mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestLoadViaConfigFailsForBadType(t *testing.T) {
	f := pflag.NewFlagSet(testFlagsetName, pflag.ContinueOnError)
	f.Int(FileArgName, testDefaultValue1, testNoHelpMessage)

	c, err := New(testPrefix, testDelimiter)
	if err != nil {
		t.Fatalf("New failed unexpectedly: %v", err)
	}

	var cfg testConfig
	if err := c.Load(f, &cfg); err == nil {
		t.Fatalf("Load: got=nil want=non-nil")
	}
}

func TestLoadViaConfigFailsForMissingFile(t *testing.T) {
	f := pflag.NewFlagSet(testFlagsetName, pflag.ContinueOnError)
	f.StringSlice(FileArgName, nil, testNoHelpMessage)

	args := []string{
		fmt.Sprintf("--%s=%s", FileArgName, testBadFileName),
	}

	if err := f.Parse(args); err != nil {
		t.Fatalf("f.Parse failed unexpectedly: %v", err)
	}

	c, err := New(testPrefix, testDelimiter)
	if err != nil {
		t.Fatalf("New failed unexpectedly: %v", err)
	}

	var cfg testConfig
	if err := c.Load(f, &cfg); err == nil {
		t.Fatalf("Load: got=nil want=non-nil")
	}
}

func testFileName(file string) string {
	return path.Join(testDataDir, file)
}

func TestLoadViaConfig(t *testing.T) {
	cases := []struct {
		name        string
		file        string
		want        testConfig
		wantLoadErr bool
	}{
		{
			name: "empty file keeps defaults",
			file: testFileName("empty.json"),
			want: testConfig{
				Value1: testDefaultValue1,
			},
		},
		{
			name: "good values overwrite defaults",
			file: testFileName(testGoodJSONConfig),
			want: testConfig{
				Value1: testValue1,
				Nested: testConfig1{
					NestedVal: testValue2,
				},
			},
		},
		{
			name:        "bad values",
			file:        testFileName("bad.json"),
			wantLoadErr: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f := pflag.NewFlagSet(testFlagsetName, pflag.ContinueOnError)
			f.Int(testKey1, testDefaultValue1, testNoHelpMessage)
			f.StringSlice(FileArgName, nil, testNoHelpMessage)

			args := []string{
				fmt.Sprintf("--%s=%s", FileArgName, tc.file),
			}
			if err := f.Parse(args); err != nil {
				t.Fatalf("f.Parse failed unexpectedly: %v", err)
			}

			c, err := New(testPrefix, testDelimiter)
			if err != nil {
				t.Fatalf("New failed unexpectedly: %v", err)
			}

			var cfg testConfig
			err = c.Load(f, &cfg)
			if tc.wantLoadErr {
				if err == nil {
					t.Errorf("Load err: got=nil want=<non-nil>")
				}
				return
			}
			if err != nil {
				t.Fatalf("Load err: got=%v want=nil", err)
			}
			if diff := cmp.Diff(tc.want, cfg); diff != "" {
				t.Errorf("Load cfg mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestEnvIsAfterFile(t *testing.T) {
	k := strings.ToUpper(testPrefix + testKey1)
	if err := os.Setenv(k, strconv.Itoa(testValue2)); err != nil {
		t.Fatalf("os.Setenv failed unexpetedly: %v", err)
	}
	defer func(k string) {
		if err := os.Unsetenv(k); err != nil {
			t.Fatalf("os.Unsetenv failed unexpectedly: %v", err)
		}
	}(k)

	f := pflag.NewFlagSet(testFlagsetName, pflag.ContinueOnError)
	f.StringSlice(FileArgName, nil, testNoHelpMessage)

	args := []string{
		fmt.Sprintf("--%s=%s", FileArgName, testFileName(testGoodJSONConfig)),
	}

	if err := f.Parse(args); err != nil {
		t.Fatalf("f.Parse failed unexpectedly: %v", err)
	}

	c, err := New(testPrefix, testDelimiter)
	if err != nil {
		t.Fatalf("New failed unexpectedly: %v", err)
	}

	var cfg testConfig
	if err := c.Load(f, &cfg); err != nil {
		t.Fatalf("c.Load: got=%v want=nil", err)
	}

	if got, want := cfg.Value1, testValue2; got != want {
		t.Errorf("Value: got=%d want=%d", got, want)
	}
}

func TestFlagIsAfterEnv(t *testing.T) {
	k := strings.ToUpper(testPrefix + testKey1)
	if err := os.Setenv(k, strconv.Itoa(testValue1)); err != nil {
		t.Fatalf("os.Setenv failed unexpectedly: %v", err)
	}
	defer func(k string) {
		if err := os.Unsetenv(k); err != nil {
			t.Fatalf("os.Unsetenv failed unexpectedly: %v", err)
		}
	}(k)

	f := pflag.NewFlagSet(testFlagsetName, pflag.ContinueOnError)
	f.Int(testKey1, testValue2, testNoHelpMessage)
	f.StringSlice(FileArgName, nil, testNoHelpMessage)

	args := []string{
		fmt.Sprintf("--%s=%d", testKey1, testValue3),
	}

	if err := f.Parse(args); err != nil {
		t.Fatalf("f.Parse failed unexpectedly: %v", err)
	}

	c, err := New(testPrefix, testDelimiter)
	if err != nil {
		t.Fatalf("New failed unexpectedly: %v", err)
	}

	var cfg testConfig
	if err := c.Load(f, &cfg); err != nil {
		t.Fatalf("c.Load: got=%v want=nil", err)
	}

	if got, want := cfg.Value1, testValue3; got != want {
		t.Errorf("Value: got=%d want=%d", got, want)
	}
}

func TestFlagIsAfterFile(t *testing.T) {
	f := pflag.NewFlagSet(testFlagsetName, pflag.ContinueOnError)
	f.Int(testKey1, testDefaultValue1, testNoHelpMessage)
	f.StringSlice(FileArgName, nil, testNoHelpMessage)

	args := []string{
		fmt.Sprintf("--%s=%s", FileArgName, testFileName(testGoodJSONConfig)),
		fmt.Sprintf("--%s=%d", testKey1, testValue3),
	}

	if err := f.Parse(args); err != nil {
		t.Fatalf("f.Parse failed unexpectedly: %v", err)
	}

	c, err := New(testPrefix, testDelimiter)
	if err != nil {
		t.Fatalf("New failed unexpectedly: %v", err)
	}

	var cfg testConfig
	if err := c.Load(f, &cfg); err != nil {
		t.Fatalf("c.Load: got=%v want=nil", err)
	}

	if got, want := cfg.Value1, testValue3; got != want {
		t.Errorf("Value: got=%d want=%d", got, want)
	}
}
