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
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/posflag"
	"github.com/knadh/koanf/v2"
	"github.com/spf13/pflag"
)

var (
	BadDelimiterError = errors.New("delimiter must contain exactly 1 character")
)

// FileArgName is the name that is used to specify configuration files.
const FileArgName = "config"

// Config holds the data necessary to process configuration data.
type Config struct {
	prefix    string
	delimiter string
}

// New returns a Config initialized with prefix and delimiter. For information
// about how these values are used see the description of load.
func New(envPrefix, flagDelimiter string) (Config, error) {
	if len(flagDelimiter) != 1 {
		return Config{}, fmt.Errorf("invalid delimiter %q: %w", flagDelimiter, BadDelimiterError)
	}
	return Config{
		prefix:    envPrefix,
		delimiter: flagDelimiter,
	}, nil
}

func (c Config) updateEnv(s string) string {
	return strings.Replace(strings.ToLower(strings.TrimPrefix(s, c.prefix)), "_", c.delimiter, -1)
}

// Load loads values into cfg from environment variables, flags and yaml files.
func (c Config) Load(f *pflag.FlagSet, cfg interface{}) error {
	const unmarshalEverything = ""

	k := koanf.New(c.delimiter)

	// Load the config files provided on the commandline if there is an argument
	// named FileArgName.
	if p := f.Lookup(FileArgName); p != nil {
		ss, err := f.GetStringSlice(FileArgName)
		if err != nil {
			return fmt.Errorf("Load GetStringSlice: %v", err)
		}
		for _, c := range ss {
			if err := k.Load(file.Provider(c), yaml.Parser()); err != nil {
				return fmt.Errorf("Load file %s: %v", c, err)
			}
		}
	}

	if err := k.Load(env.Provider(c.prefix, c.delimiter, c.updateEnv), nil); err != nil {
		return fmt.Errorf("Load env: %v", err)
	}

	if err := k.Load(posflag.Provider(f, ".", k), nil); err != nil {
		log.Fatalf("Load flags: %v", err)
	}

	if err := k.Unmarshal(unmarshalEverything, cfg); err != nil {
		return fmt.Errorf("Load unmarshal: %v", err)
	}

	return nil
}
