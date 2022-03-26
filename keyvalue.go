// Licensed under the MIT license, see LICENSE file for details.

package iterate

// KeyValue represents a key value pair.
type KeyValue[K comparable, V any] struct {
	Key   K
	Value V
}

// Split return the key and the value.
func (kv KeyValue[K, V]) Split() (K, V) {
	return kv.Key, kv.Value
}

// Zip returns an iterator of key/value pairs produced by the given key/value
// iterators. The shorter of the two iterators is used. An error is returned if
// any of the two iterators returns an error.
func Zip[K comparable, V any](keys Iterator[K], values Iterator[V]) Iterator[KeyValue[K, V]] {
	return &zipper[K, V]{
		keys:   keys,
		values: values,
	}
}

type zipper[K comparable, V any] struct {
	keys    Iterator[K]
	values  Iterator[V]
	stopped bool
}

// Next implements Iterator[T].Next.
func (it *zipper[K, V]) Next() bool {
	if it.keys.Next() && it.values.Next() {
		return true
	}
	it.stopped = true
	return false
}

// Value implements Iterator[T].Value by returning key/value pairs.
func (it *zipper[K, V]) Value() KeyValue[K, V] {
	if it.stopped {
		return KeyValue[K, V]{}
	}
	return KeyValue[K, V]{
		Key:   it.keys.Value(),
		Value: it.values.Value(),
	}
}

// Err implements Iterator[T].Err by propagating the error from the keys and
// values iterators.
func (it *zipper[K, V]) Err() error {
	if err := it.keys.Err(); err != nil {
		return err
	}
	return it.values.Err()
}

// Unzip returns a key iterator and a value iterator with pairs produced by the
// given key/value iterator.
func Unzip[K comparable, V any](kvs Iterator[KeyValue[K, V]]) (Iterator[K], Iterator[V]) {
	u := unzipper[K, V]{
		kvs: kvs,
	}
	return u.iterators()
}

type unzipper[K comparable, V any] struct {
	kvs    Iterator[KeyValue[K, V]]
	keys   []K
	values []V
}

func (u *unzipper[K, V]) iterators() (Iterator[K], Iterator[V]) {
	return &keyIterator[K, V]{
			u: u,
		}, &valueIterator[K, V]{
			u: u,
		}
}

type keyIterator[K comparable, V any] struct {
	u   *unzipper[K, V]
	key K
}

// Next implements Iterator[T].Next.
func (it *keyIterator[K, V]) Next() bool {
	// TODO(frankban): make this thread safe.
	if len(it.u.keys) != 0 {
		it.key, it.u.keys = it.u.keys[0], it.u.keys[1:]
		return true
	}
	if it.u.kvs.Next() {
		kv := it.u.kvs.Value()
		it.key = kv.Key
		it.u.values = append(it.u.values, kv.Value)
		return true
	}
	it.key = *new(K)
	return false
}

// Value implements Iterator[T].Value by returning keys from the key/value
// iterator stored in the unzipper.
func (it *keyIterator[K, V]) Value() K {
	return it.key
}

// Err implements Iterator[T].Err by propagating the error from the key/value
// source iterator.
func (it *keyIterator[K, V]) Err() error {
	return it.u.kvs.Err()
}

type valueIterator[K comparable, V any] struct {
	u     *unzipper[K, V]
	value V
}

// Next implements Iterator[T].Next.
func (it *valueIterator[K, V]) Next() bool {
	// TODO(frankban): make this thread safe.
	if len(it.u.values) != 0 {
		it.value, it.u.values = it.u.values[0], it.u.values[1:]
		return true
	}
	if it.u.kvs.Next() {
		kv := it.u.kvs.Value()
		it.value = kv.Value
		it.u.keys = append(it.u.keys, kv.Key)
		return true
	}
	it.value = *new(V)
	return false
}

// Value implements Iterator[T].Value by returning values from the key/value
// iterator stored in the unzipper.
func (it *valueIterator[K, V]) Value() V {
	return it.value
}

// Err implements Iterator[T].Err by propagating the error from the key/value
// source iterator.
func (it *valueIterator[K, V]) Err() error {
	return it.u.kvs.Err()
}

// ToMap returns a map with the values produced by the given key/value iterator.
// An error is returned if the iterator returns an error, in which case the
// returned map includes the key/value pairs already consumed.
func ToMap[K comparable, V any](it Iterator[KeyValue[K, V]]) (map[K]V, error) {
	m := make(map[K]V)
	for it.Next() {
		kv := it.Value()
		m[kv.Key] = kv.Value
	}
	return m, it.Err()
}
