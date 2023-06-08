// Copyright 2023 Shannon Wynter
//
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.

package ghostport

import (
	"fmt"
	"time"
)

// RoundTrip stores the request, response, errors associated with dumping them,
// any error returned by the underlying transport and the start end/time.
type RoundTrip struct {
	start time.Time
	end   time.Time

	request  []byte
	response []byte

	requestErr  error
	responseErr error

	roundTripErr error

	transport *Transport
}

func newRoundTrip(t *Transport) *RoundTrip {
	trip := RoundTrip{
		transport: t,
		start:     time.Now(),
	}

	return &trip
}

func (r *RoundTrip) done() {
	r.end = time.Now()
}

// String dumps a request into a relatively easy to grock string.
func (r RoundTrip) String() string {
	return fmt.Sprintf(`*** Request ***
%[1]q

err: %[3]v

*** Response ***
%[2]q

err: %[4]v

*** RoundTrip ***s

start: %[6]v
end: %[7]v
duration: %[8]v

err: %[5]v

***`, r.Request(), r.Response(), r.requestErr, r.responseErr, r.roundTripErr, r.start, r.end, r.Duration())
}

// Request bytes for request.
func (r RoundTrip) Request() []byte {
	if r.request == nil {
		return nil
	}

	return r.transport.requestCensor(r.request)
}

// Response bytes for response.
func (r RoundTrip) Response() []byte {
	if r.response == nil {
		return nil
	}

	return r.transport.responseCensor(r.response)
}

// RequestErr will be nil unless there was an error dumping the request.
func (r RoundTrip) RequestErr() error {
	return r.requestErr
}

// ResponseErr will be nil unless there was an error dumping the response.
func (r RoundTrip) ResponseErr() error {
	return r.responseErr
}

// RoundTripErr will be nil unless there was an error returned by the underlying transport.
func (r RoundTrip) RoundTripErr() error {
	return r.roundTripErr
}

// Start time of the request.
func (r RoundTrip) Start() time.Time {
	return r.start
}

// End time of the request.
func (r RoundTrip) End() time.Time {
	return r.end
}

// Duration time for the request.
func (r RoundTrip) Duration() time.Duration {
	return r.end.Sub(r.start)
}
