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

package stringslice

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

const (
	testPiece1 = "p1"
	testPiece2 = "p2"
	testPiece3 = "p3"
)

func TestString(t *testing.T) {
	cases := []struct {
		name   string
		pieces []string
		want   string
	}{
		{
			name: "empty",
		},
		{
			name: "single",
			pieces: []string{
				testPiece1,
			},
			want: testPiece1,
		},
		{
			name: "multiple",
			pieces: []string{
				testPiece1,
				testPiece2,
				testPiece3,
			},
			want: testPiece1 + "," + testPiece2 + "," + testPiece3,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ss := StringSlice(tc.pieces)

			if got, want := ss.String(), tc.want; got != want {
				t.Errorf("String: got=%q want=%q", got, want)
			}
		})
	}
}

func TestSet(t *testing.T) {
	cases := []struct {
		name    string
		value   string
		want    []string
		wantErr bool
	}{
		{
			name:  "empty",
			value: "",
			want:  []string{""},
		},
		{
			name:  "single",
			value: testPiece1,
			want: []string{
				testPiece1,
			},
		},
		{
			name:  "multiple",
			value: testPiece1 + "," + testPiece2 + "," + testPiece3,
			want: []string{
				testPiece1,
				testPiece2,
				testPiece3,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ss := StringSlice{}

			err := ss.Set(tc.value)
			if tc.wantErr {
				if err == nil {
					t.Errorf("Set err: got=nil want=<non-nil>")
				}
				return
			}

			if err != nil {
				t.Fatalf("Set err: got=%v want=nil", err)
			}
			if diff := cmp.Diff(tc.want, []string(ss)); diff != "" {
				t.Errorf("Set() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestGet(t *testing.T) {
	cases := []struct {
		name string
		ss   StringSlice
		want []string
	}{
		{
			name: "empty",
		},
		{
			name: "single",
			ss: StringSlice{
				testPiece1,
			},
			want: []string{
				testPiece1,
			},
		},
		{
			name: "multiple",
			ss: StringSlice{
				testPiece1,
				testPiece2,
				testPiece3,
			},
			want: []string{
				testPiece1,
				testPiece2,
				testPiece3,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.ss.Get()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("Get() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
