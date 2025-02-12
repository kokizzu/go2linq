//go:build go1.18

package go2linq

// https://docs.microsoft.com/dotnet/api/system.linq.enumerable.append

// Append appends a value to the end of the sequence.
func Append[Source any](source Enumerator[Source], element Source) (Enumerator[Source], error) {
	if source == nil {
		return nil, ErrNilSource
	}
	return Concat(source, RepeatMust(element, 1))
}

// AppendMust is like Append but panics in case of error.
func AppendMust[Source any](source Enumerator[Source], element Source) Enumerator[Source] {
	r, err := Append[Source](source, element)
	if err != nil {
		panic(err)
	}
	return r
}
