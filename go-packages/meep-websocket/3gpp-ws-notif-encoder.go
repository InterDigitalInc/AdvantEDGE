/*
 * Copyright (c) 2022  The AdvantEDGE Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package websocket

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
)

const seqStr string = "3GPP-WS-Notif-Seq: "
const contentTypeHeader string = "Content-Type"
const contentEncodingHeader string = "Content-Encoding"
const contentLenHeader string = "Content-Length"
const eolStr string = "\r\n"

func EncodeRequest(r *http.Request, seq uint32) ([]byte, error) {
	var req []byte
	var reqBody []byte
	headers := make(map[string]string)

	// Get mandatory Content-Type
	headers[contentTypeHeader] = r.Header.Get(contentTypeHeader)
	if headers[contentTypeHeader] == "" {
		return nil, errors.New("Missing header: " + contentTypeHeader)
	}

	// Get optional Content-Encoding
	headers[contentEncodingHeader] = r.Header.Get(contentEncodingHeader)

	// Get mandatory Content-Length
	headers[contentLenHeader] = r.Header.Get(contentLenHeader)
	if headers[contentLenHeader] != "" {
		// Get request body
		body, err := r.GetBody()
		if err != nil {
			log.Error(err.Error())
			return nil, err
		}
		reqBody, err = ioutil.ReadAll(body)
		if err != nil {
			log.Error(err.Error())
			return nil, err
		}
	} else {
		// Set "Content-Length" to 0 if empty
		headers[contentLenHeader] = "0"
	}

	// Build message
	addSequenceNumber(&req, seq)
	addHeaders(&req, headers)
	addPayload(&req, reqBody)

	return req, nil
}

func DecodeRequest(data []byte) (*http.Request, uint32, error) {
	offset := 0

	// Get Sequence number
	seq, err := getSequenceNumber(data, &offset)
	if err != nil {
		log.Error(err.Error())
		return nil, 0, err
	}
	// Get Headers
	headers, err := getHeaders(data, &offset)
	if err != nil {
		log.Error(err.Error())
		return nil, 0, err
	}
	// Get payload
	payload, err := getPayload(data, &offset)
	if err != nil {
		log.Error(err.Error())
		return nil, 0, err
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", "", bytes.NewBuffer(payload))
	if err != nil {
		log.Error(err.Error())
		return nil, 0, err
	}
	// Set mandatory content type
	contentType, found := headers[contentTypeHeader]
	if !found {
		return nil, 0, errors.New("Missing header: " + contentTypeHeader)
	}
	req.Header.Set(contentTypeHeader, contentType)
	// Set optional content encoding
	contentEncoding, found := headers[contentEncodingHeader]
	if found {
		req.Header.Set(contentEncodingHeader, contentEncoding)
	}
	// Set mandatory content length
	contentLen, found := headers[contentLenHeader]
	if !found {
		return nil, 0, errors.New("Missing header: " + contentLenHeader)
	}
	req.Header.Set(contentLenHeader, contentLen)

	return req, seq, nil
}

func EncodeResponse(r *http.Response, seq uint32) ([]byte, error) {
	var resp []byte
	var respBody []byte
	var err error
	headers := make(map[string]string)

	// Get mandatory Content-Type
	headers[contentTypeHeader] = r.Header.Get(contentTypeHeader)
	if headers[contentTypeHeader] == "" {
		return nil, errors.New("Missing header: " + contentTypeHeader)
	}

	// Get optional Content-Encoding
	headers[contentEncodingHeader] = r.Header.Get(contentEncodingHeader)

	// Get mandatory Content-Length
	headers[contentLenHeader] = r.Header.Get(contentLenHeader)
	if headers[contentLenHeader] != "" {
		defer r.Body.Close()
		respBody, err = ioutil.ReadAll(r.Body)
		if err != nil {
			log.Error(err.Error())
			return nil, err
		}
	} else {
		// Set "Content-Length" to 0 if empty
		headers[contentLenHeader] = "0"
	}

	// Build message
	addSequenceNumber(&resp, seq)
	addStatus(&resp, r.StatusCode, r.Status)
	addHeaders(&resp, headers)
	addPayload(&resp, respBody)

	return resp, nil
}

func DecodeResponse(data []byte) (*http.Response, uint32, error) {
	offset := 0

	// Get Sequence number
	seq, err := getSequenceNumber(data, &offset)
	if err != nil {
		log.Error(err.Error())
		return nil, 0, err
	}
	// Get Status
	statusCode, status, err := getStatus(data, &offset)
	if err != nil {
		log.Error(err.Error())
		return nil, 0, err
	}
	// Get Headers
	headers, err := getHeaders(data, &offset)
	if err != nil {
		log.Error(err.Error())
		return nil, 0, err
	}
	// Get payload
	payload, err := getPayload(data, &offset)
	if err != nil {
		log.Error(err.Error())
		return nil, 0, err
	}

	// Create HTTP response
	resp := &http.Response{
		Status:     status,
		StatusCode: statusCode,
		Header:     make(http.Header),
		Body:       ioutil.NopCloser(bytes.NewReader(payload)),
	}
	// Set mandatory content type
	contentType, found := headers[contentTypeHeader]
	if !found {
		return nil, 0, errors.New("Missing header: " + contentTypeHeader)
	}
	resp.Header.Set(contentTypeHeader, contentType)
	// Set optional content encoding
	contentEncoding, found := headers[contentEncodingHeader]
	if found {
		resp.Header.Set(contentEncodingHeader, contentEncoding)
	}
	// Set mandatory content length
	contentLen, found := headers[contentLenHeader]
	if !found {
		return nil, 0, errors.New("Missing header: " + contentLenHeader)
	}
	resp.Header.Set(contentLenHeader, contentLen)
	// Set content length in response
	len, err := strconv.ParseInt(contentLen, 10, 64)
	if err != nil {
		return nil, 0, errors.New("Invalid content length: " + contentLen)
	}
	resp.ContentLength = len

	return resp, seq, nil
}

// Sequence number
func addSequenceNumber(data *[]byte, seq uint32) {
	// Convert sequent number to a 4-byte array
	seqNum := make([]byte, 4)
	binary.BigEndian.PutUint32(seqNum, seq)
	// Add sequence string & num
	*data = append(*data, []byte(seqStr)...)
	*data = append(*data, seqNum...)
	*data = append(*data, []byte(eolStr)...)
}

func getSequenceNumber(data []byte, offset *int) (uint32, error) {
	curIndex := *offset

	// Make sure offset is within data range
	if len(data) < *offset {
		return 0, errors.New("Missing sequence line")
	}

	// Find EOL
	eolOffset := bytes.Index(data[curIndex:], []byte(eolStr))
	if eolOffset == -1 {
		return 0, errors.New("Missing sequence line")
	}
	byteCount := eolOffset + len(eolStr)
	eolIndex := *offset + eolOffset

	// Compare sequence string
	nextIndex := curIndex + len(seqStr)
	if eolIndex < nextIndex {
		return 0, errors.New("Missing sequence string: " + seqStr)
	}
	if !bytes.Equal(data[curIndex:nextIndex], []byte(seqStr)) {
		return 0, errors.New("Invalid sequence string: " + string(data[curIndex:nextIndex]))
	}
	curIndex = nextIndex

	// Get sequence number
	nextIndex = curIndex + 4
	if eolIndex != nextIndex {
		return 0, errors.New("Missing sequence number")
	}
	seq := binary.BigEndian.Uint32(data[curIndex:nextIndex])

	// Update offset
	*offset += byteCount
	return seq, nil
}

// Status
func addStatus(data *[]byte, code int, status string) {
	*data = append(*data, []byte(strconv.Itoa(code)+" "+status+eolStr)...)
}

func getStatus(data []byte, offset *int) (int, string, error) {
	// Make sure offset is within data range
	if len(data) < *offset {
		return 0, "", errors.New("Missing status line")
	}

	// Find EOL
	eolOffset := bytes.Index(data[*offset:], []byte(eolStr))
	if eolOffset == -1 {
		return 0, "", errors.New("Missing status line")
	}
	byteCount := eolOffset + len(eolStr)
	eolIndex := *offset + eolOffset

	// Extract status
	statusOffset := bytes.Index(data[*offset:eolIndex], []byte(" "))
	if statusOffset == -1 {
		return 0, "", errors.New("Invalid status line format")
	}
	statusIndex := *offset + statusOffset
	status := string(data[statusIndex:eolIndex])

	// Extract status code
	statusCodeStr := string(data[*offset:statusIndex])
	statusCode, err := strconv.Atoi(statusCodeStr)
	if err != nil || http.StatusText(statusCode) == "" {
		return 0, "", errors.New("Invalid status code: " + statusCodeStr)
	}

	// Update offset
	*offset += byteCount
	return 204, status, nil
}

// Header
func addHeaders(data *[]byte, headers map[string]string) {
	// Add headers
	for header, val := range headers {
		// Add header, val & CRLF
		if val != "" {
			*data = append(*data, []byte(header+": "+val+eolStr)...)
		}
	}
	// Add CRLF
	*data = append(*data, []byte(eolStr)...)
}

func getHeaders(data []byte, offset *int) (map[string]string, error) {
	headers := make(map[string]string)

	// Make sure offset is within data range
	if len(data) < *offset {
		return nil, errors.New("Missing headers section")
	}

	// Find end of headers section --> 2 x EOL
	eolOffset := bytes.Index(data[*offset:], []byte(eolStr+eolStr))
	if eolOffset == -1 {
		return nil, errors.New("Missing headers section")
	}
	byteCount := eolOffset + len(eolStr+eolStr)
	eolIndex := *offset + eolOffset

	// Extract headers
	headerLines := strings.Split(string(data[*offset:eolIndex]), eolStr)
	for _, headerLine := range headerLines {
		parts := strings.Split(headerLine, ": ")
		if len(parts) != 2 {
			return nil, errors.New("Invalid header line: " + headerLine)
		}
		headers[parts[0]] = parts[1]
	}

	// Update offset
	*offset += byteCount
	return headers, nil
}

// Payload
func addPayload(data *[]byte, payload []byte) {
	if len(payload) > 0 {
		*data = append(*data, payload...)
	}
}

func getPayload(data []byte, offset *int) ([]byte, error) {
	// Make sure offset is within data range
	if len(data) < *offset {
		return nil, errors.New("Missing payload")
	}

	// Return payload
	return data[*offset:], nil
}
