// Copyright 2023 Shannon Wynter
//
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.

/*
Package ghostport provides a http.RoundTripper for recording a number of
recent requests/responses along with some other useful information for
the purposes of diagnosting failures in applications that make a lot of
http requests.
*/
package ghostport
