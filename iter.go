// Licensed under the MIT license, see LICENSE file for details.

package iterate

// Iterator is implemented by types producing values of type T. Implementations
// are typically used in for loops, for instance:
//
//     var v SomeType
//     for iterator.Next(&v) {
//         // Do something with v.
//     }
//     if err := iterator.Err(); err != nil {
//         // Handle error.
//     }
//
// Depending on the implementation, producing values might lead to errors. For
// this reason it is important to always check Err() after iterating.
type Iterator[T any] interface {
	// Next produces the next iterator value and assign it to the variable
	// pointed by v. When the iterator is done, no values are assigned and false
	// is returned. Further calls to Next should just return false with no other
	// side effects. When iterating produces an error, false is returned, and
	// Err() returns the error.
	Next(v *T) bool
	// Err returns the first error occurred while iterating.
	Err() error
}

// Ideas:
// - DropWhile
// - NewReader
// - FromChan
