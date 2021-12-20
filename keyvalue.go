// Licensed under the MIT license, see LICENSE file for details.

package iterate

// KeyValue represents a key value pair.
type KeyValue[K comparable, V any] struct {
	Key   K
	Value V
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
	keys   Iterator[K]
	values Iterator[V]
}

// Next implements Iterator[T].Next by producing key/value pairs.
func (it *zipper[K, V]) Next(v *KeyValue[K, V]) bool {
	var key K
	var value V
	if it.keys.Next(&key) && it.values.Next(&value) {
		*v = KeyValue[K, V]{
			Key:   key,
			Value: value,
		}
		return true
	}
	*v = KeyValue[K, V]{}
	return false
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
	u *unzipper[K, V]
}

// Next implements Iterator[T].Next by producing keys from the key/value
// iterator stored in the unzipper.
func (it *keyIterator[K, V]) Next(v *K) bool {
	// TODO(frankban): make this thread safe.
	if len(it.u.keys) != 0 {
		var key K
		key, it.u.keys = it.u.keys[0], it.u.keys[1:]
		*v = key
		return true
	}
	var kv KeyValue[K, V]
	if it.u.kvs.Next(&kv) {
		*v = kv.Key
		it.u.values = append(it.u.values, kv.Value)
		return true
	}
	*v = *new(K)
	return false
}

// Err implements Iterator[T].Err by propagating the error from the key/value
// source iterator.
func (it *keyIterator[K, V]) Err() error {
	return it.u.kvs.Err()
}

type valueIterator[K comparable, V any] struct {
	u *unzipper[K, V]
}

// Next implements Iterator[T].Next by producing values from the key/value
// iterator stored in the unzipper.
func (it *valueIterator[K, V]) Next(v *V) bool {
	// TODO(frankban): make this thread safe.
	if len(it.u.values) != 0 {
		var value V
		value, it.u.values = it.u.values[0], it.u.values[1:]
		*v = value
		return true
	}
	var kv KeyValue[K, V]
	if it.u.kvs.Next(&kv) {
		*v = kv.Value
		it.u.keys = append(it.u.keys, kv.Key)
		return true
	}
	*v = *new(V)
	return false
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
	var kv KeyValue[K, V]
	for it.Next(&kv) {
		m[kv.Key] = kv.Value
	}
	return m, it.Err()
}
