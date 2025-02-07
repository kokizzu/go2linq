//go:build go1.18

package go2linq

// https://docs.microsoft.com/dotnet/api/system.linq.enumerable.prepend

// Prepend adds a value to the beginning of the sequence.
func Prepend[Source any](source Enumerator[Source], element Source) (Enumerator[Source], error) {
	if source == nil {
		return nil, ErrNilSource
	}
	return Concat(RepeatMust(element, 1), source)
}

// PrependMust is like Prepend but panics in case of error.
func PrependMust[Source any](source Enumerator[Source], element Source) Enumerator[Source] {
	r, err := Prepend[Source](source, element)
	if err != nil {
		panic(err)
	}
	return r
}
