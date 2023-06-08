// Copyright 2023 Shannon Wynter
//
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.

package ghostport

import (
	"container/ring"
	"errors"
	"net/http"
	"net/http/httputil"
	"sync"
)

// DefaultRingSize is 3 because that seems like a sensible number of requests
// to hang on to.
const DefaultRingSize = 3

// Transport implements the following interfaces;
// * http.RoundTripper to provide the ghostport.
// * fmt.Stringer so you can just print.
type Transport struct {
	transport http.RoundTripper
	ring      *ring.Ring
	mu        sync.Mutex

	requestCensor  Censor
	responseCensor Censor
}

var (
	errRequestNil  = errors.New("request was nil")
	errResponseNil = errors.New("response was nil")
)

// New *Transport.
func New(opts ...Option) *Transport {
	var transport Transport

	for _, fn := range opts {
		fn(&transport)
	}

	if transport.transport == nil {
		transport.transport = http.DefaultTransport
	}

	if transport.ring == nil {
		transport.ring = ring.New(DefaultRingSize)
	}

	if transport.requestCensor == nil {
		transport.requestCensor = NoopCensor
	}

	if transport.responseCensor == nil {
		transport.responseCensor = NoopCensor
	}

	return &transport
}

// RoundTrip fulfills the http.RoundTripper inteface and will record the time
// the trip starts, the time it ends, the request, and response from the
// underlying http.RoundTripper along with any errors it encounters along
// the way.
func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Start a trip.
	trip := newRoundTrip(t)

	// Record that the trip is done.
	defer trip.done()

	// Record the request, not sure why you'd ever send a nil request but if you do
	// lets make that note and leave it to whatever other transport is involved to do
	// the complaining.
	if req == nil {
		trip.requestErr = errRequestNil
	} else {
		trip.request, trip.requestErr = httputil.DumpRequestOut(req, true)
	}

	// Permit the underlying transport to do the request.
	resp, err := t.transport.RoundTrip(req)

	// Record the error if there is one.
	trip.roundTripErr = err

	// If the underlying transport did error then there's a good chance that resp
	// is nil, make that note, otherwise save it.
	if resp == nil {
		trip.responseErr = errResponseNil
	} else {
		trip.response, trip.responseErr = httputil.DumpResponse(resp, true)
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	// Stash this trip in the ring.
	t.ring.Value = trip
	t.ring = t.ring.Next()

	return resp, err
}

// Trips returns the list of trips stored from the ring in the order of execution.
func (t *Transport) Trips() []*RoundTrip {
	var trips []*RoundTrip

	t.ring.Do(func(v interface{}) {
		if v == nil {
			return
		}

		trip, isa := v.(*RoundTrip)
		if !isa {
			return
		}

		trips = append(trips, trip)
	})

	return trips
}

// String fulfills the fmt.Stringer interface.
func (t *Transport) String() string {
	var out string

	t.ring.Do(func(v interface{}) {
		if v == nil {
			return
		}

		trip, isa := v.(*RoundTrip)
		if !isa {
			return
		}

		out += trip.String()
		out += "\n"
	})

	return out
}
