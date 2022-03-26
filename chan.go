// Licensed under the MIT license, see LICENSE file for details.

package iterate

// FromChannel returns an iterator producing values from the given channel. The
// iteration stops when the channel is closed. The returned error is always nil.
func FromChannel[T any](ch <-chan T) Iterator[T] {
	return &channelIterator[T]{
		ch: ch,
	}
}

type channelIterator[T any] struct {
	ch    <-chan T
	value T
}

// Next implements Iterator[T].Next.
func (it *channelIterator[T]) Next() bool {
	val, ok := <-it.ch
	if ok {
		it.value = val
		return true
	}
	it.value = *new(T)
	return false
}

// Value implements Iterator[T].Value by returning values from the channel.
func (it *channelIterator[T]) Value() T {
	return it.value
}

// Err implements Iterator[T].Err. The returned error is always nil.
func (it *channelIterator[T]) Err() error {
	return nil
}
