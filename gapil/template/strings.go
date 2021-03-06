// Copyright (C) 2017 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package template

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"unicode"
)

type stringList []string

// String returns the concatenation of all the string segments with no separator.
func (l stringList) String() string {
	return strings.Join([]string(l), "")
}

// stringify transforms the input parameters into a string list. Arrays and
// slices are flattened into a sequential list of strings.
func stringify(v ...interface{}) stringList {
	out := stringList{}
	for _, v := range v {
		switch v := v.(type) {
		case string:
			out = append(out, v)
		case []string:
			out = append(out, v...)
		case stringList:
			out = append(out, v...)
		default:
			switch reflect.TypeOf(v).Kind() {
			case reflect.Array, reflect.Slice:
				v := reflect.ValueOf(v)
				for i, c := 0, v.Len(); i < c; i++ {
					out = append(out, stringify(v.Index(i).Interface())...)
				}
			default:
				out = append(out, fmt.Sprintf("%v", v))
			}
		}
	}
	// Filter out any empty strings
	count := 0
	for _, s := range out {
		if len(s) > 0 {
			out[count] = s
			count++
		}
	}
	return out[:count]
}

// Strings returns the arguments as a string list.
func (Functions) Strings(v ...interface{}) stringList {
	return stringify(v...)
}

// JoinWith returns the concatenation of all the string segments with the specified separator.
func (Functions) JoinWith(sep string, v ...interface{}) string {
	l := stringify(v...)
	return strings.Join([]string(l), sep)
}

// SplitOn slices each string segement into all substrings separated by sep. The returned stringList
// will not contain any occurances of sep.
func (Functions) SplitOn(sep string, v ...interface{}) stringList {
	l := stringify(v...)
	out := stringList{}
	for _, s := range l {
		for _, v := range strings.Split(s, sep) {
			if len(v) > 0 {
				out = append(out, v)
			}
		}
	}
	return out
}

// SplitUpperCase slices each string segment before and after each upper-case rune.
func (Functions) SplitUpperCase(v ...interface{}) stringList {
	l := stringify(v...)
	out := stringList{}
	for _, s := range l {
		str := ""
		for _, r := range s {
			if unicode.IsUpper(r) {
				if len(str) > 0 {
					out = append(out, str)
					str = ""
				}
				out = append(out, string(r))
			} else {
				str += string(r)
			}
		}
		if len(str) > 0 {
			out = append(out, str)
		}
	}
	return out
}

// SplitPascalCase slices each string segment at each transition from an letter rune to a upper-case
// letter rune.
func (Functions) SplitPascalCase(v ...interface{}) stringList {
	l := stringify(v...)
	out := stringList{}
	for _, str := range l {
		runes := bytes.Runes([]byte(str))
		str := ""
		p := 'x'
		for _, c := range runes {
			if unicode.IsLetter(p) && unicode.IsUpper(c) {
				if len(str) > 0 {
					out = append(out, str)
				}
				str = string(c)
			} else {
				str += string(c)
			}
			p = c
		}
		if len(str) > 0 {
			out = append(out, str)
		}
	}
	return out
}

// Title capitalizes each letter of each string segment.
func (Functions) Title(v ...interface{}) stringList {
	l := stringify(v...)
	out := make(stringList, len(l))
	for i, s := range l {
		first := true
		out[i] = strings.Map(func(r rune) rune {
			if first {
				first = false
				return unicode.ToTitle(r)
			} else {
				return r
			}
		}, s)
	}
	return out
}

// Untitle lower-cases each letter of each string segment.
func (Functions) Untitle(v ...interface{}) stringList {
	l := stringify(v...)
	out := make(stringList, len(l))
	for i, s := range l {
		first := true
		out[i] = strings.Map(func(r rune) rune {
			if first {
				first = false
				return unicode.ToLower(r)
			}
			return r
		}, s)
	}
	return out
}

// Lower lower-cases all letters of each string segment.
func (Functions) Lower(v ...interface{}) stringList {
	l := stringify(v...)
	out := make(stringList, len(l))
	for i, s := range l {
		out[i] = strings.ToLower(s)
	}
	return out
}

// Upper upper-cases all letters of each string segment.
func (Functions) Upper(v ...interface{}) stringList {
	l := stringify(v...)
	out := make(stringList, len(l))
	for i, s := range l {
		out[i] = strings.ToUpper(s)
	}
	return out
}

// Contains returns true if any string segment contains substr.
func (Functions) Contains(substr string, v ...interface{}) bool {
	l := stringify(v...)
	for _, s := range l {
		if strings.Contains(s, substr) {
			return true
		}
	}
	return false
}

// HasPrefix tests whether the string s begins with prefix.
func (Functions) HasPrefix(s string, prefix string) bool {
	return strings.HasPrefix(s, prefix)
}

// Replace any occurance of old with new in the string segments.
func (Functions) Replace(old string, new string, v ...interface{}) stringList {
	l := stringify(v...)
	out := stringList{}
	for _, s := range l {
		s = strings.Replace(s, old, new, -1)
		if len(s) > 0 {
			out = append(out, s)
		}
	}
	return out
}

func (Functions) TrimLeft(cutset string, v ...interface{}) stringList {
	l := stringify(v...)
	out := stringList{}
	for _, s := range l {
		s = strings.TrimLeft(s, cutset)
		if len(s) > 0 {
			out = append(out, s)
		}
	}
	return out
}

func (Functions) TrimRight(cutset string, v ...interface{}) stringList {
	l := stringify(v...)
	out := stringList{}
	for _, s := range l {
		s = strings.TrimRight(s, cutset)
		if len(s) > 0 {
			out = append(out, s)
		}
	}
	return out
}

func (Functions) TrimPrefix(prefix string, v ...interface{}) stringList {
	l := stringify(v...)
	out := stringList{}
	for _, s := range l {
		s = strings.TrimPrefix(s, prefix)
		if len(s) > 0 {
			out = append(out, s)
		}
	}
	return out
}

// FilterOut returns from with all occurances of v removed.
func (Functions) FilterOut(v, from stringList) stringList {
	m := make(map[string]struct{}, len(v))
	for _, s := range v {
		m[s] = struct{}{}
	}

	out := make(stringList, 0, len(from))
	for _, s := range from {
		if _, found := m[s]; !found {
			out = append(out, s)
		}
	}
	return out
}
