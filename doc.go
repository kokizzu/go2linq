//go:build go1.18

// Package go2linq implements .NET's LINQ to Objects
// (https://docs.microsoft.com/dotnet/csharp/programming-guide/concepts/linq/linq-to-objects).
//
// See also:
//
// https://en.wikipedia.org/wiki/Language_Integrated_Query
//
// https://docs.microsoft.com/dotnet/csharp/programming-guide/concepts/linq/
//
// https://docs.microsoft.com/dotnet/api/system.linq.enumerable
//
// Inspired by: https://codeblog.jonskeet.uk/category/edulinq/
//
// Methods involving two Enumerator parameters
// (Concat…, Except…, GroupJoin…, Intersect…, Join…, SequenceEqualMust…, Union…, ZipErr…)
// are not safe to use the arguments based on the same Enumerator instance
// (see Test_ZipSelf for such examples).
// The problem arises from the fact that calling MoveNext on one Enumerator will affect the other too.
// So if you need to use Enumerators based on the same instance
// (such as performing operations on adjacent elements (see Test_ZipSelf/AdjacentElements)),
// use corresponding …Self… counterpart methods instead.
package go2linq
