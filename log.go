// Copyright (c) 2017 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.
//
// Copyright (c) 2009 The Go Authors. All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//    * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//    * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//    * Neither the name of Google Inc. nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package btclog

import (
	"io"
	"os"
	"runtime"
	"strings"
	"time"
)

// Disabled is a Logger that will never output anything.
var Disabled Logger

// defaultFlags specifies changes to the default logger behavior.  It is set
// during package init and configured using the LOGFLAGS environment variable.
// New logger backends can override these default flags using WithFlags.
var defaultFlags uint32

// Flags to modify Backend's behavior.
const (
	// Llongfile modifies the logger output to include full path and line number
	// of the logging callsite, e.g. /a/b/c/main.go:123.
	Llongfile uint32 = 1 << iota

	// Lshortfile modifies the logger output to include filename and line number
	// of the logging callsite, e.g. main.go:123.  Overrides Llongfile.
	Lshortfile
)

// From stdlib log package.
// Cheap integer to fixed-width decimal ASCII.  Give a negative width to avoid
// zero-padding.
func itoa(buf *buffer, i int, wid int) {
	// Assemble decimal in reverse order.
	var b [20]byte
	bp := len(b) - 1
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	// i < 10
	b[bp] = byte('0' + i)
	buf.writeBytes(b[bp:])
}

// writeTimestamp writes the date in the format 'YYYY-MM-DD hh:mm:ss.sss' to the
// buffer.
func writeTimestamp(buf *buffer, t time.Time) {
	year, month, day := t.Date()
	hour, min, sec := t.Clock()
	ms := t.Nanosecond() / 1e6

	itoa(buf, year, 4)
	buf.writeByte('-')
	itoa(buf, int(month), 2)
	buf.writeByte('-')
	itoa(buf, day, 2)
	buf.writeByte(' ')
	itoa(buf, hour, 2)
	buf.writeByte(':')
	itoa(buf, min, 2)
	buf.writeByte(':')
	itoa(buf, sec, 2)
	buf.writeByte('.')
	itoa(buf, ms, 3)
	buf.writeByte(' ')
}

// callsite returns the file name and line number of the callsite to the
// subsystem logger.
func callsite(flag uint32, skipDepth int) (string, int) {
	_, file, line, ok := runtime.Caller(skipDepth)
	if !ok {
		return "???", 0
	}
	if flag&Lshortfile != 0 {
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if os.IsPathSeparator(file[i]) {
				short = file[i+1:]
				break
			}
		}
		file = short
	}
	return file, line
}

func init() {
	// Initialise the Disabled logger.
	Disabled = NewSLogger(NewDefaultHandler(io.Discard))
	Disabled.SetLevel(LevelOff)

	// Read logger flags from the LOGFLAGS environment variable. Multiple
	// flags can be set at once, separated by commas.
	for _, f := range strings.Split(os.Getenv("LOGFLAGS"), ",") {
		switch f {
		case "longfile":
			defaultFlags |= Llongfile
		case "shortfile":
			defaultFlags |= Lshortfile
		}
	}
}
