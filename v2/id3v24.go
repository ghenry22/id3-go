// Copyright 2013 Michael Yang. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package v2

import (
	"bytes"
	"github.com/mikkyang/id3-go/encodedbytes"
	"io"
)

var (
	// Common frame IDs
	V24CommonFrame = map[string]FrameType{
		"Title":    V23FrameTypeMap["TIT2"],
		"Artist":   V23FrameTypeMap["TPE1"],
		"Album":    V23FrameTypeMap["TALB"],
		"Year":     V23FrameTypeMap["TDRC"],
		"Genre":    V23FrameTypeMap["TCON"],
		"Comments": V23FrameTypeMap["COMM"],
	}

	// V23DeprecatedTypeMap contains deprecated frame IDs from ID3v2.2
	V24DeprecatedTypeMap = V23DeprecatedTypeMap

	// V23FrameTypeMap specifies the frame IDs and constructors allowed in ID3v2.3
	V24FrameTypeMap = V23FrameTypeMap
)

func ParseV24Frame(reader io.Reader) Framer {
	data := make([]byte, FrameHeaderSize)
	if n, err := io.ReadFull(reader, data); n < FrameHeaderSize || err != nil {
		return nil
	}

	id := string(bytes.Trim(data[:4], "\x00"))
	t, ok := V24FrameTypeMap[id]
	if !ok {
		t = FrameType{id: id, description: "Unknown frame", constructor: ParseDataFrame}
	}

	size, err := encodedbytes.SynchInt(data[4:8])
	if err != nil {
		return nil
	}

	if id == "" && size == 0 {
		return nil
	}

	h := FrameHead{
		FrameType:   t,
		statusFlags: data[8],
		formatFlags: data[9],
		size:        size,
	}

	frameData := make([]byte, size)
	if n, err := io.ReadFull(reader, frameData); n < int(size) || err != nil {
		return nil
	}

	return t.constructor(h, frameData)
}

func V24Bytes(f Framer) []byte {
	headBytes := make([]byte, 0, FrameHeaderSize)

	headBytes = append(headBytes, f.Id()...)
	headBytes = append(headBytes, encodedbytes.SynchBytes(uint32(f.Size()))...)
	headBytes = append(headBytes, f.StatusFlags(), f.FormatFlags())

	return append(headBytes, f.Bytes()...)
}
