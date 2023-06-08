// Copyright 2023 Shannon Wynter
//
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.

package ghostport

import (
	"container/ring"
	"net/http"
)

// Option funcs allow the optional configuration of the Transport.
//
// Generally you don't need to know or understand what this is, just look for
// the With* functions, for example: WithTransport.
type Option func(*Transport)

// WithTransport permits providing your own http.RoundTripper to be the
// underlying transport. If not used ghostport will use http.DefaultTransport.
func WithTransport(transport http.RoundTripper) Option {
	return func(t *Transport) {
		t.transport = transport
	}
}

// WithRingSize lets you specify the size of the ring buffer used to
// store the round trips, default is DefaultRingSize.
func WithRingSize(size int) Option {
	return func(t *Transport) {
		t.ring = ring.New(size)
	}
}

// WithRequestCensor permits configuring a function to censor the output
// of requests before being returned.
func WithRequestCensor(c Censor) Option {
	return func(t *Transport) {
		t.requestCensor = c
	}
}

// WithResponseCensor permits configuring a function to censor the output
// of responses before being returned.
func WithResponseCensor(c Censor) Option {
	return func(t *Transport) {
		t.responseCensor = c
	}
}
