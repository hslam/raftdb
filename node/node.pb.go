// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package node

import (
	"fmt"
	"github.com/hslam/code"
)

// SetCommand represents a set command.
type SetCommand struct {
	Key   string
	Value string
}

// Size returns the size of the buffer required to represent the data when encoded.
func (d *SetCommand) Size() int {
	var size uint64
	size += 11 + uint64(len(d.Key))
	size += 11 + uint64(len(d.Value))
	return int(size)
}

// Marshal returns the encoded bytes.
func (d *SetCommand) Marshal() ([]byte, error) {
	size := d.Size()
	buf := make([]byte, size)
	n, err := d.MarshalTo(buf[:size])
	return buf[:n], err
}

// MarshalTo marshals into buf and returns the number of bytes.
func (d *SetCommand) MarshalTo(buf []byte) (int, error) {
	var size = uint64(d.Size())
	if uint64(cap(buf)) >= size {
		buf = buf[:size]
	} else {
		return 0, fmt.Errorf("proto: buf is too short")
	}
	var offset uint64
	var n uint64
	if len(d.Key) > 0 {
		buf[offset] = 1<<3 | 2
		offset++
		n = code.EncodeString(buf[offset:], d.Key)
		offset += n
	}
	if len(d.Value) > 0 {
		buf[offset] = 2<<3 | 2
		offset++
		n = code.EncodeString(buf[offset:], d.Value)
		offset += n
	}
	return int(offset), nil
}

// Unmarshal unmarshals from data.
func (d *SetCommand) Unmarshal(data []byte) error {
	var length = uint64(len(data))
	var offset uint64
	var n uint64
	var tag uint64
	var fieldNumber int
	var wireType uint8
	for {
		if offset < length {
			tag = uint64(data[offset])
			offset++
		} else {
			break
		}
		fieldNumber = int(tag >> 3)
		wireType = uint8(tag & 0x7)
		switch fieldNumber {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Key", wireType)
			}
			n = code.DecodeString(data[offset:], &d.Key)
			offset += n
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Value", wireType)
			}
			n = code.DecodeString(data[offset:], &d.Value)
			offset += n
		}
	}
	return nil
}

// Request represents an RPC request.
type Request struct {
	Key   string
	Value string
}

// Size returns the size of the buffer required to represent the data when encoded.
func (d *Request) Size() int {
	var size uint64
	size += 11 + uint64(len(d.Key))
	size += 11 + uint64(len(d.Value))
	return int(size)
}

// Marshal returns the encoded bytes.
func (d *Request) Marshal() ([]byte, error) {
	size := d.Size()
	buf := make([]byte, size)
	n, err := d.MarshalTo(buf[:size])
	return buf[:n], err
}

// MarshalTo marshals into buf and returns the number of bytes.
func (d *Request) MarshalTo(buf []byte) (int, error) {
	var size = uint64(d.Size())
	if uint64(cap(buf)) >= size {
		buf = buf[:size]
	} else {
		return 0, fmt.Errorf("proto: buf is too short")
	}
	var offset uint64
	var n uint64
	if len(d.Key) > 0 {
		buf[offset] = 1<<3 | 2
		offset++
		n = code.EncodeString(buf[offset:], d.Key)
		offset += n
	}
	if len(d.Value) > 0 {
		buf[offset] = 2<<3 | 2
		offset++
		n = code.EncodeString(buf[offset:], d.Value)
		offset += n
	}
	return int(offset), nil
}

// Unmarshal unmarshals from data.
func (d *Request) Unmarshal(data []byte) error {
	var length = uint64(len(data))
	var offset uint64
	var n uint64
	var tag uint64
	var fieldNumber int
	var wireType uint8
	for {
		if offset < length {
			tag = uint64(data[offset])
			offset++
		} else {
			break
		}
		fieldNumber = int(tag >> 3)
		wireType = uint8(tag & 0x7)
		switch fieldNumber {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Key", wireType)
			}
			n = code.DecodeString(data[offset:], &d.Key)
			offset += n
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Value", wireType)
			}
			n = code.DecodeString(data[offset:], &d.Value)
			offset += n
		}
	}
	return nil
}

// Response represents an RPC response.
type Response struct {
	Ok     bool
	Result string
	Leader string
}

// Size returns the size of the buffer required to represent the data when encoded.
func (d *Response) Size() int {
	var size uint64
	size += 11
	size += 11 + uint64(len(d.Result))
	size += 11 + uint64(len(d.Leader))
	return int(size)
}

// Marshal returns the encoded bytes.
func (d *Response) Marshal() ([]byte, error) {
	size := d.Size()
	buf := make([]byte, size)
	n, err := d.MarshalTo(buf[:size])
	if err != nil {
		return nil, err
	}
	return buf[:n], nil
}

// MarshalTo marshals into buf and returns the number of bytes.
func (d *Response) MarshalTo(buf []byte) (int, error) {
	var size = uint64(d.Size())
	if uint64(cap(buf)) >= size {
		buf = buf[:size]
	} else {
		return 0, fmt.Errorf("proto: buf is too short")
	}
	var offset uint64
	var n uint64
	if d.Ok {
		buf[offset] = 1<<3 | 0
		offset++
		n = code.EncodeBool(buf[offset:], d.Ok)
		offset += n
	}
	if len(d.Result) > 0 {
		buf[offset] = 2<<3 | 2
		offset++
		n = code.EncodeString(buf[offset:], d.Result)
		offset += n
	}
	if len(d.Leader) > 0 {
		buf[offset] = 3<<3 | 2
		offset++
		n = code.EncodeString(buf[offset:], d.Leader)
		offset += n
	}
	return int(offset), nil
}

// Unmarshal unmarshals from data.
func (d *Response) Unmarshal(data []byte) error {
	var length = uint64(len(data))
	var offset uint64
	var n uint64
	var tag uint64
	var fieldNumber int
	var wireType uint8
	for {
		if offset < length {
			tag = uint64(data[offset])
			offset++
		} else {
			break
		}
		fieldNumber = int(tag >> 3)
		wireType = uint8(tag & 0x7)
		switch fieldNumber {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Ok", wireType)
			}
			n = code.DecodeBool(data[offset:], &d.Ok)
			offset += n
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Result", wireType)
			}
			n = code.DecodeString(data[offset:], &d.Result)
			offset += n
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Leader", wireType)
			}
			n = code.DecodeString(data[offset:], &d.Leader)
			offset += n
		}
	}
	return nil
}
