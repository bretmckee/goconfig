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

// Package cfgloader provides a wrapper for the koanf/v2 configuration package
// that make it easier to use for a preferred use case, configuring it to load
// values (in order) from yaml files, flags and environment variables.
//
// Values are loaded in order from:
// - configuration files in the order given on the command line
// - environment variables
// - flags
//
// If multiple sources contain different values for the same configruation
// field, the last one found is used.
//
// In order for a structure field to be loaded, it must be exported (e.g start
// with a capital letter), and contain a koanf field tag:
//
//	type Config struct {
//	  Value  string         `koanf:"value"`
//	}
//
// Values are loaded from all the sources based on the tag, and would be loaded
// from any of these sources if the were supplied:
//
// $ export VALUE="from the environment"
// or
// $ ./prog -value="from a flag"
// or
// $ cat config.yaml
// value: from a config file
//
// There are two strings that affect the loading of the configuration. There is
// a prefix which can be used to avoid collisons with environment variables.
//
// If the prefix is set to "CONFIG_" when calling New, then the above example
// becomes:
//
// $ export CONFIG_VALUE="from the environment"
//
// Delimiter is used to separate nested configuration structures in flags:
//
//	type Nested struct {
//	  Val int `koanf:"val"`
//	}
//
//	type Config struct {
//	  Value  string `koanf:"value"`
//	  Nested Nested `koanf:"nested"`
//	}
//
// With the delimiter set to ".", to set val, use:
// $ ./prog --nested.val=7
//
// Environment variables always use "_" as their delimter, so without a prefix
// this would be:
//
// $ export NESTED_VAL="from the environment"
package goconfig
