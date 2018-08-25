// Copyright (c) 2016 Uber Technologies, Inc.
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


// Code is modified by souhup.
// Original file is json_encoder.go

package zapcore

import (
	"encoding/base64"
	"encoding/json"
	"math"
	"sync"
	"time"
	"unicode/utf8"

	"go.uber.org/zap/buffer"
	"go.uber.org/zap/internal/bufferpool"
)

var _simplePool = sync.Pool{New: func() interface{} {
	return &simpleEncoder{}
}}

func getSimpleEncoder() *simpleEncoder {
	return _simplePool.Get().(*simpleEncoder)
}

func putSimpleEncoder(enc *simpleEncoder) {
	if enc.reflectBuf != nil {
		enc.reflectBuf.Free()
	}
	enc.EncoderConfig = nil
	enc.buf = nil
	enc.spaced = false
	enc.openNamespaces = 0
	enc.reflectBuf = nil
	enc.reflectEnc = nil
	_simplePool.Put(enc)
}

type simpleEncoder struct {
	*EncoderConfig
	buf            *buffer.Buffer
	spaced         bool // include spaces after colons and commas
	openNamespaces int

	// for encoding generic values by reflection
	reflectBuf *buffer.Buffer
	reflectEnc *json.Encoder
}

// NewSimpleEncoder creates a fast, low-allocation JSON encoder. The encoder
// appropriately escapes all field keys and values.
//
// Note that the encoder doesn't deduplicate keys, so it's possible to produce
// a message like
//   {"foo":"bar","foo":"baz"}
// This is permitted by the JSON specification, but not encouraged. Many
// libraries will ignore duplicate key-value pairs (typically keeping the last
// pair) when unmarshaling, but users should attempt to avoid adding duplicate
// keys.
func NewSimpleEncoder(cfg EncoderConfig) Encoder {
	return newSimpleEncoder(cfg, false)
}

func newSimpleEncoder(cfg EncoderConfig, spaced bool) *simpleEncoder {
	return &simpleEncoder{
		EncoderConfig: &cfg,
		buf:           bufferpool.Get(),
		spaced:        spaced,
	}
}

func (enc *simpleEncoder) AddArray(key string, arr ArrayMarshaler) error {
	enc.addKey(key)
	return enc.AppendArray(arr)
}

func (enc *simpleEncoder) AddObject(key string, obj ObjectMarshaler) error {
	enc.addKey(key)
	return enc.AppendObject(obj)
}

func (enc *simpleEncoder) AddBinary(key string, val []byte) {
	enc.AddString(key, base64.StdEncoding.EncodeToString(val))
}

func (enc *simpleEncoder) AddByteString(key string, val []byte) {
	enc.addKey(key)
	enc.AppendByteString(val)
}

func (enc *simpleEncoder) AddBool(key string, val bool) {
	enc.addKey(key)
	enc.AppendBool(val)
}

func (enc *simpleEncoder) AddComplex128(key string, val complex128) {
	enc.addKey(key)
	enc.AppendComplex128(val)
}

func (enc *simpleEncoder) AddDuration(key string, val time.Duration) {
	enc.addKey(key)
	enc.AppendDuration(val)
}

func (enc *simpleEncoder) AddFloat64(key string, val float64) {
	enc.addKey(key)
	enc.AppendFloat64(val)
}

func (enc *simpleEncoder) AddInt64(key string, val int64) {
	enc.addKey(key)
	enc.AppendInt64(val)
}

func (enc *simpleEncoder) resetReflectBuf() {
	if enc.reflectBuf == nil {
		enc.reflectBuf = bufferpool.Get()
		enc.reflectEnc = json.NewEncoder(enc.reflectBuf)
	} else {
		enc.reflectBuf.Reset()
	}
}

func (enc *simpleEncoder) AddReflected(key string, obj interface{}) error {
	enc.resetReflectBuf()
	err := enc.reflectEnc.Encode(obj)
	if err != nil {
		return err
	}
	enc.reflectBuf.TrimNewline()
	enc.addKey(key)
	_, err = enc.buf.Write(enc.reflectBuf.Bytes())
	return err
}

func (enc *simpleEncoder) OpenNamespace(key string) {
	enc.addKey(key)
	enc.buf.AppendByte('{')
	enc.openNamespaces++
}

func (enc *simpleEncoder) AddString(key, val string) {
	enc.addKey(key)
	enc.AppendString(val)
}

func (enc *simpleEncoder) AddTime(key string, val time.Time) {
	enc.addKey(key)
	enc.AppendTime(val)
}

func (enc *simpleEncoder) AddUint64(key string, val uint64) {
	enc.addKey(key)
	enc.AppendUint64(val)
}

func (enc *simpleEncoder) AppendArray(arr ArrayMarshaler) error {
	enc.addElementSeparator()
	enc.buf.AppendByte('[')
	err := arr.MarshalLogArray(enc)
	enc.buf.AppendByte(']')
	return err
}

func (enc *simpleEncoder) AppendObject(obj ObjectMarshaler) error {
	enc.addElementSeparator()
	enc.buf.AppendByte('{')
	err := obj.MarshalLogObject(enc)
	enc.buf.AppendByte('}')
	return err
}

func (enc *simpleEncoder) AppendBool(val bool) {
	enc.addElementSeparator()
	enc.buf.AppendBool(val)
}

func (enc *simpleEncoder) AppendByteString(val []byte) {
	enc.addElementSeparator()
	enc.buf.AppendByte('"')
	enc.safeAddByteString(val)
	enc.buf.AppendByte('"')
}

func (enc *simpleEncoder) AppendComplex128(val complex128) {
	enc.addElementSeparator()
	// Cast to a platform-independent, fixed-size type.
	r, i := float64(real(val)), float64(imag(val))
	enc.buf.AppendByte('"')
	// Because we're always in a quoted string, we can use strconv without
	// special-casing NaN and +/-Inf.
	enc.buf.AppendFloat(r, 64)
	enc.buf.AppendByte('+')
	enc.buf.AppendFloat(i, 64)
	enc.buf.AppendByte('i')
	enc.buf.AppendByte('"')
}

func (enc *simpleEncoder) AppendDuration(val time.Duration) {
	cur := enc.buf.Len()
	enc.EncodeDuration(val, enc)
	if cur == enc.buf.Len() {
		// User-supplied EncodeDuration is a no-op. Fall back to nanoseconds to keep
		// JSON valid.
		enc.AppendInt64(int64(val))
	}
}

func (enc *simpleEncoder) AppendInt64(val int64) {
	enc.addElementSeparator()
	enc.buf.AppendInt(val)
}

func (enc *simpleEncoder) AppendReflected(val interface{}) error {
	enc.resetReflectBuf()
	err := enc.reflectEnc.Encode(val)
	if err != nil {
		return err
	}
	enc.reflectBuf.TrimNewline()
	enc.addElementSeparator()
	_, err = enc.buf.Write(enc.reflectBuf.Bytes())
	return err
}

func (enc *simpleEncoder) AppendString(val string) {
	enc.addElementSeparator()
	//enc.buf.AppendByte('"')
	enc.safeAddString(val)
	//enc.buf.AppendByte('"')
}

func (enc *simpleEncoder) AppendTime(val time.Time) {
	cur := enc.buf.Len()
	enc.EncodeTime(val, enc)
	if cur == enc.buf.Len() {
		// User-supplied EncodeTime is a no-op. Fall back to nanos since epoch to keep
		// output JSON valid.
		enc.AppendInt64(val.UnixNano())
	}
}

func (enc *simpleEncoder) AppendUint64(val uint64) {
	enc.addElementSeparator()
	enc.buf.AppendUint(val)
}

func (enc *simpleEncoder) AddComplex64(k string, v complex64) { enc.AddComplex128(k, complex128(v)) }
func (enc *simpleEncoder) AddFloat32(k string, v float32)     { enc.AddFloat64(k, float64(v)) }
func (enc *simpleEncoder) AddInt(k string, v int)             { enc.AddInt64(k, int64(v)) }
func (enc *simpleEncoder) AddInt32(k string, v int32)         { enc.AddInt64(k, int64(v)) }
func (enc *simpleEncoder) AddInt16(k string, v int16)         { enc.AddInt64(k, int64(v)) }
func (enc *simpleEncoder) AddInt8(k string, v int8)           { enc.AddInt64(k, int64(v)) }
func (enc *simpleEncoder) AddUint(k string, v uint)           { enc.AddUint64(k, uint64(v)) }
func (enc *simpleEncoder) AddUint32(k string, v uint32)       { enc.AddUint64(k, uint64(v)) }
func (enc *simpleEncoder) AddUint16(k string, v uint16)       { enc.AddUint64(k, uint64(v)) }
func (enc *simpleEncoder) AddUint8(k string, v uint8)         { enc.AddUint64(k, uint64(v)) }
func (enc *simpleEncoder) AddUintptr(k string, v uintptr)     { enc.AddUint64(k, uint64(v)) }
func (enc *simpleEncoder) AppendComplex64(v complex64)        { enc.AppendComplex128(complex128(v)) }
func (enc *simpleEncoder) AppendFloat64(v float64)            { enc.appendFloat(v, 64) }
func (enc *simpleEncoder) AppendFloat32(v float32)            { enc.appendFloat(float64(v), 32) }
func (enc *simpleEncoder) AppendInt(v int)                    { enc.AppendInt64(int64(v)) }
func (enc *simpleEncoder) AppendInt32(v int32)                { enc.AppendInt64(int64(v)) }
func (enc *simpleEncoder) AppendInt16(v int16)                { enc.AppendInt64(int64(v)) }
func (enc *simpleEncoder) AppendInt8(v int8)                  { enc.AppendInt64(int64(v)) }
func (enc *simpleEncoder) AppendUint(v uint)                  { enc.AppendUint64(uint64(v)) }
func (enc *simpleEncoder) AppendUint32(v uint32)              { enc.AppendUint64(uint64(v)) }
func (enc *simpleEncoder) AppendUint16(v uint16)              { enc.AppendUint64(uint64(v)) }
func (enc *simpleEncoder) AppendUint8(v uint8)                { enc.AppendUint64(uint64(v)) }
func (enc *simpleEncoder) AppendUintptr(v uintptr)            { enc.AppendUint64(uint64(v)) }

func (enc *simpleEncoder) Clone() Encoder {
	clone := enc.clone()
	clone.buf.Write(enc.buf.Bytes())
	return clone
}

func (enc *simpleEncoder) clone() *simpleEncoder {
	clone := getSimpleEncoder()
	clone.EncoderConfig = enc.EncoderConfig
	clone.spaced = enc.spaced
	clone.openNamespaces = enc.openNamespaces
	clone.buf = bufferpool.Get()
	return clone
}

func (enc *simpleEncoder) EncodeEntry(ent Entry, fields []Field) (*buffer.Buffer, error) {
	final := enc.clone()
	//final.buf.AppendByte('{')

	if final.TimeKey != "" {
		final.AddTime(final.TimeKey, ent.Time)
	}
	if final.LevelKey != "" {
		final.addKey(final.LevelKey)
		cur := final.buf.Len()
		final.EncodeLevel(ent.Level, final)
		if cur == final.buf.Len() {
			// User-supplied EncodeLevel was a no-op. Fall back to strings to keep
			// output JSON valid.
			final.AppendString(ent.Level.String())
		}
	}
	if ent.LoggerName != "" && final.NameKey != "" {
		final.addKey(final.NameKey)
		cur := final.buf.Len()
		nameEncoder := final.EncodeName

		// if no name encoder provided, fall back to FullNameEncoder for backwards
		// compatibility
		if nameEncoder == nil {
			nameEncoder = FullNameEncoder
		}

		nameEncoder(ent.LoggerName, final)
		if cur == final.buf.Len() {
			// User-supplied EncodeName was a no-op. Fall back to strings to
			// keep output JSON valid.
			final.AppendString(ent.LoggerName)
		}
	}
	if enc.buf.Len() > 0 {
		final.addElementSeparator()
		final.buf.Write(enc.buf.Bytes())
	}
	addFields(final, fields)
	//for i := range fields {
	//	if i != 0 {
	//		//final.buf.AppendByte('\t')
	//	}
	//	fields[i].AddTo(final)
	//}
	if ent.Caller.Defined && final.CallerKey != "" {
		final.addKey(final.CallerKey)
		cur := final.buf.Len()
		final.EncodeCaller(ent.Caller, final)
		if cur == final.buf.Len() {
			// User-supplied EncodeCaller was a no-op. Fall back to strings to
			// keep output JSON valid.
			final.AppendString(ent.Caller.String())
		}
	}
	final.closeOpenNamespaces()
	if final.MessageKey != "" {
		final.addKey(enc.MessageKey)
		final.AppendString(ent.Message)
	}
	if ent.Stack != "" && final.StacktraceKey != "" {
		final.AddString(final.StacktraceKey, ent.Stack)
	}
	//final.buf.AppendByte('}')
	if final.LineEnding != "" {
		final.buf.AppendString(final.LineEnding)
	} else {
		final.buf.AppendString(DefaultLineEnding)
	}

	ret := final.buf
	putSimpleEncoder(final)
	return ret, nil
}

func (enc *simpleEncoder) truncate() {
	enc.buf.Reset()
}

func (enc *simpleEncoder) closeOpenNamespaces() {
	for i := 0; i < enc.openNamespaces; i++ {
		enc.buf.AppendByte('}')
	}
}

func (enc *simpleEncoder) addKey(key string) {
	//enc.addElementSeparator()
	//enc.buf.AppendByte('"')
	//enc.safeAddString(key)
	//enc.buf.AppendByte('"')
	//enc.buf.AppendByte(':')
	//if enc.spaced {
	//	enc.buf.AppendByte(' ')
	//}
}

func (enc *simpleEncoder) addElementSeparator() {
	last := enc.buf.Len() - 1
	if last < 0 {
		return
	}
	switch enc.buf.Bytes()[last] {
	case '{', '[', ':', ',', ' ':
		return
	default:
		//enc.buf.AppendByte(' ')
		enc.buf.AppendByte('\t')
		//if enc.spaced {
		//	enc.buf.AppendByte(' ')
		//}
	}
}

func (enc *simpleEncoder) appendFloat(val float64, bitSize int) {
	enc.addElementSeparator()
	switch {
	case math.IsNaN(val):
		enc.buf.AppendString(`"NaN"`)
	case math.IsInf(val, 1):
		enc.buf.AppendString(`"+Inf"`)
	case math.IsInf(val, -1):
		enc.buf.AppendString(`"-Inf"`)
	default:
		enc.buf.AppendFloat(val, bitSize)
	}
}

// safeAddString JSON-escapes a string and appends it to the internal buffer.
// Unlike the standard library's encoder, it doesn't attempt to protect the
// user from browser vulnerabilities or JSONP-related problems.
func (enc *simpleEncoder) safeAddString(s string) {
	for i := 0; i < len(s); {
		if enc.tryAddRuneSelf(s[i]) {
			i++
			continue
		}
		r, size := utf8.DecodeRuneInString(s[i:])
		if enc.tryAddRuneError(r, size) {
			i++
			continue
		}
		enc.buf.AppendString(s[i : i+size])
		i += size
	}
}

// safeAddByteString is no-alloc equivalent of safeAddString(string(s)) for s []byte.
func (enc *simpleEncoder) safeAddByteString(s []byte) {
	for i := 0; i < len(s); {
		if enc.tryAddRuneSelf(s[i]) {
			i++
			continue
		}
		r, size := utf8.DecodeRune(s[i:])
		if enc.tryAddRuneError(r, size) {
			i++
			continue
		}
		enc.buf.Write(s[i : i+size])
		i += size
	}
}

// tryAddRuneSelf appends b if it is valid UTF-8 character represented in a single byte.
func (enc *simpleEncoder) tryAddRuneSelf(b byte) bool {
	if b >= utf8.RuneSelf {
		return false
	}
	if 0x20 <= b && b != '\\' && b != '"' {
		enc.buf.AppendByte(b)
		return true
	}
	switch b {
	case '\\', '"':
		enc.buf.AppendByte('\\')
		enc.buf.AppendByte(b)
	case '\n':
		enc.buf.AppendByte('\\')
		enc.buf.AppendByte('n')
	case '\r':
		enc.buf.AppendByte('\\')
		enc.buf.AppendByte('r')
	case '\t':
		enc.buf.AppendByte('\\')
		enc.buf.AppendByte('t')
	default:
		// Encode bytes < 0x20, except for the escape sequences above.
		enc.buf.AppendString(`\u00`)
		enc.buf.AppendByte(_hex[b>>4])
		enc.buf.AppendByte(_hex[b&0xF])
	}
	return true
}

func (enc *simpleEncoder) tryAddRuneError(r rune, size int) bool {
	if r == utf8.RuneError && size == 1 {
		enc.buf.AppendString(`\ufffd`)
		return true
	}
	return false
}
