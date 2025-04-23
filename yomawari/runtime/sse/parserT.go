package sse

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"math"
	"strconv"
	"sync/atomic"
	"time"
)

const CR byte = '\r'
const LF byte = '\n'
const CRLF string = "\r\n"

var utf8Bom []byte = []byte{0xEF, 0xBB, 0xBF}

const TimeSpanMaxMilliseconds = int64(922337203685477)

type SseParser[T any] struct {
	timeSpan_maxvalue_milliseconds int64
	default_array_pool_rent_size   int
	stream                         io.Reader
	itemParser                     SseItemParser[T]
	used                           int32
	lineBuffer                     []byte
	lineOffset                     int
	lineLength                     int
	newlineIndex                   int
	lastSearchedForNewline         int
	eof                            bool
	dataBuffer                     []byte
	dataLength                     int
	dataAppended                   bool
	eventType                      string
	eventId                        *string
	nextReconnectionInterval       time.Duration
	ReconnectionInterval           time.Duration
	LastEventId                    string
}

func CreateSseParser(stream io.Reader) *SseParser[string] {
	var itemParser SseItemParser[string] = func(_ string, data []byte) string {
		return base64.URLEncoding.EncodeToString(data)
	}
	return NewSseParser(stream, itemParser)
}

func NewSseParser[T any](stream io.Reader, itemParser SseItemParser[T]) *SseParser[T] {
	s := &SseParser[T]{
		timeSpan_maxvalue_milliseconds: math.MaxInt,
		default_array_pool_rent_size:   1024,
		stream:                         stream,
		itemParser:                     itemParser,
		used:                           0,
		lineBuffer:                     []byte{},
		lineOffset:                     0,
		lineLength:                     0,
		newlineIndex:                   0,
		lastSearchedForNewline:         0,
		dataLength:                     0,
		eof:                            false,
		dataBuffer:                     []byte{},
		dataAppended:                   false,
		eventType:                      "message",
		eventId:                        nil,
		nextReconnectionInterval:       0,
		LastEventId:                    "",
		ReconnectionInterval:           0,
	}
	return s
}

func (s *SseParser[T]) skipBomIfPresent() {
	if len(s.lineBuffer) >= 3 && bytes.Equal(s.lineBuffer[:3], utf8Bom) {
		s.lineOffset += 3
		s.lineLength -= 3
	}
}

func (s *SseParser[T]) shiftIfNecessary() {
	if s.lineOffset+s.lineLength == len(s.lineBuffer) {
		if s.lineOffset != 0 {
			if s.lastSearchedForNewline >= 0 {
				s.lastSearchedForNewline -= s.lineOffset
			}
			s.lineOffset = 0
		}
	}
}

func (s *SseParser[T]) throwIfNotFirstEnumeration() error {
	if atomic.SwapInt32(&s.used, 1) != 0 {
		return fmt.Errorf("the enumerable may be enumerated only once")
	}
	return nil
}

func (r *SseParser[T]) fillLineBuffer(ctx context.Context) (int, error) {
	r.shiftIfNecessary()
	offset := r.lineOffset + r.lineLength
	if offset >= len(r.lineBuffer) {
		return 0, fmt.Errorf("buffer overflow")
	}

	n, err := r.stream.Read(r.lineBuffer[offset:])
	if n > 0 {
		r.lineLength += n
	} else if err == nil {
		r.eof = true
	}

	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	default:
	}

	if err != nil && err != io.EOF {
		return n, err
	}
	return n, nil
}

func (s *SseParser[T]) getNewLineLength() int {
	if string(s.lineBuffer[s.newlineIndex:s.lineLength+s.lineOffset]) == CRLF {
		return 2
	}

	return 1
}

func (s *SseParser[T]) processLine(advance *int) (*SseItem[T], bool) {
	line := s.lineBuffer[s.lineOffset:s.newlineIndex]

	if len(line) == 0 {
		*advance = s.getNewLineLength()
		if s.dataAppended {
			data := s.itemParser(s.eventType, s.dataBuffer[:s.dataLength])
			sseItem := NewSseItem[T](data, s.eventType)
			sseItem.EventId = s.eventId
			if s.nextReconnectionInterval != 0 {
				sseItem.ReconnectionInterval = s.nextReconnectionInterval
			}

			s.eventType = ""
			s.eventId = nil
			s.nextReconnectionInterval = 0
			s.dataLength = 0
			s.dataAppended = false
			return sseItem, true
		}
		return &SseItem[T]{}, false
	}

	colonPos := bytes.IndexByte(line, ':')
	var fieldName []byte
	var fieldValue []byte
	if colonPos >= 0 {
		fieldName = line[:colonPos]
		fieldValue = line[colonPos+1:]
		if len(fieldValue) > 0 && fieldValue[0] == ' ' {
			fieldValue = fieldValue[1:]
		}
	} else {
		fieldName = line
		fieldValue = []byte{}
	}

	if bytes.Equal(fieldName, []byte("data")) {
		if !s.dataAppended {
			newlineLength := s.getNewLineLength()
			start := s.newlineIndex + newlineLength
			end := start + s.lineLength - len(line) - newlineLength
			remainder := s.lineBuffer[start:end]

			if len(remainder) > 0 && (remainder[0] == LF || (remainder[0] == CR && len(remainder) > 1)) {
				*advance = len(line) + newlineLength
				if bytes.HasPrefix(remainder, []byte(CRLF)) {
					*advance += 2
				} else {
					*advance += 1
				}

				data := s.itemParser(s.eventType, fieldValue)
				sseItem := NewSseItem[T](data, s.eventType)
				sseItem.EventId = s.eventId
				if s.nextReconnectionInterval != 0 {
					sseItem.ReconnectionInterval = s.nextReconnectionInterval
				}

				s.eventType = ""
				s.eventId = nil
				s.nextReconnectionInterval = 0
				return sseItem, true
			}
		}

		if s.dataAppended {
			s.dataLength += 1
			s.dataBuffer[s.dataLength] = LF
		}

		copy(s.dataBuffer[s.dataLength:], fieldValue)
		s.dataLength += len(fieldValue)
		s.dataAppended = true
	} else if bytes.Equal(fieldName, []byte("event")) {
		s.eventType = base64.URLEncoding.EncodeToString(fieldValue)
	} else if bytes.Equal(fieldName, []byte("id")) {
		if bytes.IndexByte(fieldValue, byte('0')) < 0 {
			s.LastEventId = base64.URLEncoding.EncodeToString(fieldValue)
			s.LastEventId = *s.eventId
		}
		s.eventType = base64.URLEncoding.EncodeToString(fieldValue)
	} else if bytes.Equal(fieldName, []byte("retry")) {
		s.tryParseReconnectionInterval(fieldValue)
	}
	*advance = len(line) + s.getNewLineLength()
	return &SseItem[T]{}, true
}

func (s *SseParser[T]) tryParseReconnectionInterval(fieldValue []byte) {
	milliseconds, err := strconv.ParseInt(string(fieldValue), 10, 64)
	if err != nil {
		return
	}

	if milliseconds < 0 || milliseconds > TimeSpanMaxMilliseconds {
		return
	}

	var timeSpan time.Duration
	if milliseconds == TimeSpanMaxMilliseconds {
		timeSpan = time.Duration(1<<63 - 1)
	} else {
		timeSpan = time.Duration(milliseconds) * time.Millisecond
	}

	s.nextReconnectionInterval = timeSpan
	s.ReconnectionInterval = timeSpan
}

func (s *SseParser[T]) getNextSearchOffsetAndLength(searchOffset *int, searchLength *int) {
	if s.lastSearchedForNewline > s.lineOffset {
		*searchOffset = s.lastSearchedForNewline
		*searchLength = s.lineLength - (s.lastSearchedForNewline - s.lineOffset)
	} else {
		*searchOffset = s.lineOffset
		*searchLength = s.lineLength
	}
}

func indexOfAny(data []byte, targets ...byte) int {
	return bytes.IndexAny(data, string(targets))
}

func (sse *SseParser[T]) EnumerateStream(ctx context.Context) (<-chan SseItem[T], <-chan error) {
	outChan := make(chan SseItem[T])
	errChan := make(chan error, 1)

	if err := sse.throwIfNotFirstEnumeration(); err != nil {
		defer close(outChan)
		defer close(errChan)
		errChan <- err
		return outChan, errChan
	}

	go func() {
		defer close(outChan)
		defer close(errChan)
		sse.lineBuffer = []byte{}

		for {
			if _, err := sse.fillLineBuffer(ctx); err != nil {
				errChan <- err
				return
			}

			sse.skipBomIfPresent()

			for {
				searchOffset, searchLength := 0, 0
				sse.getNextSearchOffsetAndLength(&searchOffset, &searchLength)
				sse.newlineIndex = indexOfAny(sse.lineBuffer[searchOffset:searchOffset+searchLength], CR, LF)

				if sse.newlineIndex >= 0 {
					sse.lastSearchedForNewline = -1
					sse.newlineIndex += searchOffset

					// Check newline conditions
					if sse.lineBuffer[sse.newlineIndex] == LF || sse.newlineIndex-sse.lineOffset+1 < sse.lineLength || sse.eof {
						advance := 0
						if item, processed := sse.processLine(&advance); processed {
							select {
							case outChan <- *item:
							case <-ctx.Done():
								errChan <- ctx.Err()
								return
							}

							sse.lineOffset += advance
							sse.lineLength -= advance
							continue
						}
					}
				} else {
					sse.lastSearchedForNewline = searchOffset + searchLength
				}

				if sse.eof {
					return
				}

				if _, err := sse.fillLineBuffer(ctx); err != nil {
					errChan <- err
					return
				}
			}
		}
	}()

	return outChan, errChan
}

func (sse *SseParser[T]) Enumerate(ctx context.Context) ([]SseItem[T], error) {
	itemCh, errCh := sse.EnumerateStream(ctx)

	var results []SseItem[T]
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case err, ok := <-errCh:
			if ok {
				return nil, err
			}
			errCh = nil
		case item, ok := <-itemCh:
			if ok {
				results = append(results, item)
			} else {
				itemCh = nil
			}
		}

		if itemCh == nil && errCh == nil {
			break
		}
	}

	return results, nil
}
