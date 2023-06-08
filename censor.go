// Copyright 2023 Shannon Wynter
//
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.

package ghostport

// Censor funcs will let you filter input bytes to remove sensitive data before it's returned.
type Censor func(in []byte) []byte

// NoopCensor will simply do nothing but waste a couple of cpu cycles.
func NoopCensor(in []byte) []byte {
	return in
}
