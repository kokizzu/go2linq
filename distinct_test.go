//go:build go1.18

package go2linq

import (
	"testing"
)

// https://github.com/jskeet/edulinq/blob/master/src/Edulinq.Tests/DistinctTest.cs

var (
	testString1 = "test"
	testString2 = "test"
)

func Test_Distinct_int(t *testing.T) {
	type args struct {
		source Enumerator[int]
	}
	tests := []struct {
		name        string
		args        args
		want        Enumerator[int]
		wantErr     bool
		expectedErr error
	}{
		{name: "NullSourceNoComparer",
			wantErr:     true,
			expectedErr: ErrNilSource,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Distinct(tt.args.source)
			if (err != nil) != tt.wantErr {
				t.Errorf("Distinct() error = '%v', wantErr '%v'", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if err != tt.expectedErr {
					t.Errorf("Distinct() error = '%v', expectedErr '%v'", err, tt.expectedErr)
				}
				return
			}
			if !SequenceEqualMust(got, tt.want) {
				got.Reset()
				tt.want.Reset()
				t.Errorf("Distinct() = '%v', want '%v'", String(got), String(tt.want))
			}
		})
	}
}

func Test_Distinct_string(t *testing.T) {
	type args struct {
		source Enumerator[string]
	}
	tests := []struct {
		name string
		args args
		want Enumerator[string]
	}{
		{name: "1",
			args: args{
				source: NewOnSlice("A", "a", "b", "c", "b"),
			},
			want: NewOnSlice("A", "a", "b", "c"),
		},
		{name: "2",
			args: args{
				source: NewOnSlice("b", "a", "d", "a"),
			},
			want: NewOnSlice("b", "a", "d"),
		},
		// https://docs.microsoft.com/dotnet/csharp/programming-guide/concepts/linq/set-operations#distinct-and-distinctby
		{name: "Distinct",
			args: args{
				source: NewOnSlice("Mercury", "Venus", "Venus", "Earth", "Mars", "Earth"),
			},
			want: NewOnSlice("Mercury", "Venus", "Earth", "Mars"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := Distinct(tt.args.source)
			if !SequenceEqualMust(got, tt.want) {
				got.Reset()
				tt.want.Reset()
				t.Errorf("Distinct() = %v, want %v", String(got), String(tt.want))
			}
		})
	}
}

func Test_DistinctEq_string(t *testing.T) {
	type args struct {
		source  Enumerator[string]
		equaler Equaler[string]
	}
	tests := []struct {
		name        string
		args        args
		want        Enumerator[string]
		wantErr     bool
		expectedErr error
	}{
		{name: "NullSourceWithComparer",
			args: args{
				equaler: CaseInsensitiveEqualer,
			},
			wantErr:     true,
			expectedErr: ErrNilSource,
		},
		{name: "NullComparerUsesDefault",
			args: args{
				source: NewOnSlice("xyz", testString1, "XYZ", testString2, "def"),
			},
			want: NewOnSlice("xyz", testString1, "XYZ", "def"),
		},
		{name: "NonNullEqualer",
			args: args{
				source:  NewOnSlice("xyz", testString1, "XYZ", testString2, "def"),
				equaler: CaseInsensitiveEqualer,
			},
			want: NewOnSlice("xyz", testString1, "def"),
		},
		{name: "DistinctStringsWithCaseInsensitiveComparer",
			args: args{
				source:  NewOnSlice("xyz", testString1, "XYZ", testString2, "def"),
				equaler: CaseInsensitiveEqualer,
			},
			want: NewOnSlice("xyz", testString1, "def"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DistinctEq(tt.args.source, tt.args.equaler)
			if (err != nil) != tt.wantErr {
				t.Errorf("DistinctEq() error = '%v', wantErr '%v'", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if err != tt.expectedErr {
					t.Errorf("DistinctEq() error = '%v', expectedErr '%v'", err, tt.expectedErr)
				}
				return
			}
			if !SequenceEqualMust(got, tt.want) {
				got.Reset()
				tt.want.Reset()
				t.Errorf("DistinctEq() = '%v', want '%v'", String(got), String(tt.want))
			}
		})
	}
}

func Test_DistinctCmp_string(t *testing.T) {
	type args struct {
		source Enumerator[string]
		cmp    Comparer[string]
	}
	tests := []struct {
		name string
		args args
		want Enumerator[string]
	}{
		{name: "DistinctStringsWithCaseInsensitiveComparer",
			args: args{
				source: NewOnSlice("xyz", testString1, "XYZ", testString2, "def"),
				cmp:    CaseInsensitiveComparer,
			},
			want: NewOnSlice("xyz", testString1, "def"),
		},
		{name: "3",
			args: args{
				source: NewOnSlice("A", "a", "b", "c", "b"),
				cmp:    CaseInsensitiveComparer,
			},
			want: NewOnSlice("A", "b", "c"),
		},
		{name: "4",
			args: args{
				source: NewOnSlice("b", "a", "d", "a"),
				cmp:    CaseInsensitiveComparer,
			},
			want: NewOnSlice("b", "a", "d"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := DistinctCmp(tt.args.source, tt.args.cmp); !SequenceEqualMust(got, tt.want) {
				got.Reset()
				tt.want.Reset()
				t.Errorf("DistinctCmp() = '%v', want '%v'", String(got), String(tt.want))
			}
		})
	}
}

func Test_DistinctCmp_int(t *testing.T) {
	type args struct {
		source Enumerator[int]
		cmp    Comparer[int]
	}
	tests := []struct {
		name string
		args args
		want Enumerator[int]
	}{
		{name: "EmptyEnumerable",
			args: args{
				source: Empty[int](),
				cmp:    Order[int]{},
			},
			want: Empty[int](),
		},
		{name: "1",
			args: args{
				source: NewOnSlice(1, 2, 3, 4),
				cmp:    Order[int]{},
			},
			want: NewOnSlice(1, 2, 3, 4),
		},
		{name: "2",
			args: args{
				source: ConcatMust(NewOnSliceEn(1, 2, 3, 4), NewOnSliceEn(1, 2, 3, 4)),
				cmp:    Order[int]{},
			},
			want: NewOnSlice(1, 2, 3, 4),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := DistinctCmp(tt.args.source, tt.args.cmp); !SequenceEqualMust(got, tt.want) {
				got.Reset()
				tt.want.Reset()
				t.Errorf("DistinctCmp() = '%v', want '%v'", String(got), String(tt.want))
			}
		})
	}
}

func Test_DistinctEq_Reset(t *testing.T) {
	source, _ := DistinctEq(
		NewOnSliceEn("xyz", testString1, "XYZ", testString2, "def"),
		CaseInsensitiveEqualer)
	got1 := NewOnSliceEn(Slice(source)...)
	source.Reset()
	got2 := NewOnSliceEn(Slice(source)...)
	if !SequenceEqualMust(got1, got2) {
		got1.Reset()
		got2.Reset()
		t.Errorf("Reset error: '%v' != '%v'", String(got1), String(got2))
	}
}
