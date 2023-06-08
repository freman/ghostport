# Ghostport - Ghost in the Transport

Package ghostport provides a http.RoundTripper for recording a number of
recent requests/responses along with some other useful information for
the purposes of diagnosting failures in applications that make a lot of
http requests.

## Constants

DefaultRingSize is 3 because that seems like a sensible number of requests
to hang on to.

```golang
const DefaultRingSize = 3
```

## Functions

### func [NoopCensor](/censor.go#L12)

`func NoopCensor(in []byte) []byte`

NoopCensor will simply do nothing but waste a couple of cpu cycles.

## Types

### type [Censor](/censor.go#L9)

`type Censor func(in []byte) []byte`

Censor funcs will let you filter input bytes to remove sensitive data before it's returned.

### type [Option](/options.go#L17)

`type Option func(*Transport)`

Option funcs allow the optional configuration of the Transport.

Generally you don't need to know or understand what this is, just look for
the With* functions, for example: WithTransport.

#### func [WithRequestCensor](/options.go#L37)

`func WithRequestCensor(c Censor) Option`

WithRequestCensor permits configuring a function to censor the output
of requests before being returned.

#### func [WithResponseCensor](/options.go#L45)

`func WithResponseCensor(c Censor) Option`

WithResponseCensor permits configuring a function to censor the output
of responses before being returned.

#### func [WithRingSize](/options.go#L29)

`func WithRingSize(size int) Option`

WithRingSize lets you specify the size of the ring buffer used to
store the round trips, default is DefaultRingSize.

#### func [WithTransport](/options.go#L21)

`func WithTransport(transport http.RoundTripper) Option`

WithTransport permits providing your own http.RoundTripper to be the
underlying transport. If not used ghostport will use http.DefaultTransport.

### type [RoundTrip](/roundtrip.go#L15)

`type RoundTrip struct { ... }`

RoundTrip stores the request, response, errors associated with dumping them,
any error returned by the underlying transport and the start end/time.

#### func (RoundTrip) [Duration](/roundtrip.go#L110)

`func (r RoundTrip) Duration() time.Duration`

Duration time for the request.

#### func (RoundTrip) [End](/roundtrip.go#L105)

`func (r RoundTrip) End() time.Time`

End time of the request.

#### func (RoundTrip) [Request](/roundtrip.go#L67)

`func (r RoundTrip) Request() []byte`

Request bytes for request.

#### func (RoundTrip) [RequestErr](/roundtrip.go#L85)

`func (r RoundTrip) RequestErr() error`

RequestErr will be nil unless there was an error dumping the request.

#### func (RoundTrip) [Response](/roundtrip.go#L76)

`func (r RoundTrip) Response() []byte`

Response bytes for response.

#### func (RoundTrip) [ResponseErr](/roundtrip.go#L90)

`func (r RoundTrip) ResponseErr() error`

ResponseErr will be nil unless there was an error dumping the response.

#### func (RoundTrip) [RoundTripErr](/roundtrip.go#L95)

`func (r RoundTrip) RoundTripErr() error`

RoundTripErr will be nil unless there was an error returned by the underlying transport.

#### func (RoundTrip) [Start](/roundtrip.go#L100)

`func (r RoundTrip) Start() time.Time`

Start time of the request.

#### func (RoundTrip) [String](/roundtrip.go#L44)

`func (r RoundTrip) String() string`

String dumps a request into a relatively easy to grock string.

### type [Transport](/transport.go#L23)

`type Transport struct { ... }`

Transport implements the following interfaces;
* http.RoundTripper to provide the ghostport.
* fmt.Stringer so you can just print.

#### func [New](/transport.go#L38)

`func New(opts ...Option) *Transport`

New *Transport.

#### func (*Transport) [RoundTrip](/transport.go#L68)

`func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error)`

RoundTrip fulfills the http.RoundTripper inteface and will record the time
the trip starts, the time it ends, the request, and response from the
underlying http.RoundTripper along with any errors it encounters along
the way.

#### func (*Transport) [String](/transport.go#L129)

`func (t *Transport) String() string`

String fulfills the fmt.Stringer interface.

#### func (*Transport) [Trips](/transport.go#L109)

`func (t *Transport) Trips() []*RoundTrip`

Trips returns the list of trips stored from the ring in the order of execution.