// Licensed under the MIT license, see LICENSE file for details.

package iterate

// FromChannel returns an iterator producing values from the given channel. The
// iteration stops when the channel is closed. The returned error is always nil.
func FromChannel[T any](ch <-chan T) Iterator[T] {
	return IteratorFunc[T](func(v *T) bool {
		val, ok := <-ch
		*v = val
		return ok
	})
}
