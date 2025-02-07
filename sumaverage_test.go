//go:build go1.18

package go2linq

import (
	"math"
	"testing"
)

// https://github.com/jskeet/edulinq/blob/master/src/Edulinq.Tests/SumTest.cs
// https://github.com/jskeet/edulinq/blob/master/src/Edulinq.Tests/AverageTest.cs

func Test_Sum_string_int(t *testing.T) {
	type args struct {
		source   Enumerator[string]
		selector func(string) int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{name: "EmptySequenceIntWithSelector",
			args: args{
				source:   Empty[string](),
				selector: func(s string) int { return len(s) },
			},
			want: 0,
		},
		{name: "SimpleSumIntWithSelector",
			args: args{
				source:   NewOnSlice("x", "abc", "de"),
				selector: func(s string) int { return len(s) },
			},
			want: 6,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SumMust(tt.args.source, tt.args.selector); got != tt.want {
				t.Errorf("Sum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_Sum_string_float64(t *testing.T) {
	type args struct {
		source   Enumerator[string]
		selector func(string) float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{name: "EmptySequenceFloat64WithSelector",
			args: args{
				source:   Empty[string](),
				selector: func(s string) float64 { return float64(len(s)) },
			},
			want: 0,
		},
		{name: "SimpleSumFloat64WithSelector",
			args: args{
				source:   NewOnSlice("x", "abc", "de"),
				selector: func(s string) float64 { return float64(len(s)) },
			},
			want: 6,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SumMust(tt.args.source, tt.args.selector); got != tt.want {
				t.Errorf("Sum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_Sum_string_float64IsNaN(t *testing.T) {
	type args struct {
		source   Enumerator[string]
		selector func(string) float64
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "SimpleSumFloat64WithSelectorWithNan",
			args: args{
				source: NewOnSlice("x", "abc", "de"),
				selector: func(s string) float64 {
					l := len(s)
					if l == 3 {
						return math.NaN()
					}
					return float64(l)
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if want := math.IsNaN(SumMust(tt.args.source, tt.args.selector)); want != tt.want {
				t.Errorf("IsNaN(Sum()) = %v, want %v", want, tt.want)
			}
		})
	}
}

func Test_Sum_float64_float64IsInf(t *testing.T) {
	type args struct {
		source   Enumerator[float64]
		selector func(float64) float64
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "OverflowToNegInfinityFloat64",
			args: args{
				source:   NewOnSlice(-math.MaxFloat64, -math.MaxFloat64),
				selector: Identity[float64],
			},
			want: true,
		},
		{name: "OverflowToInfinityFloat64",
			args: args{
				source:   NewOnSlice(math.MaxFloat64, math.MaxFloat64),
				selector: Identity[float64],
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SumMust(tt.args.source, tt.args.selector)
			var want bool
			switch tt.name {
			case "OverflowToNegInfinityFloat64":
				want = math.IsInf(got, -1)
			case "OverflowToInfinityFloat64":
				want = math.IsInf(got, +1)
			}
			if want != tt.want {
				t.Errorf("IsInf(Sum()) = %v, want %v", want, tt.want)
			}
		})
	}
}

func Test_Sum_string_float64IsInf(t *testing.T) {
	type args struct {
		source   Enumerator[string]
		selector func(string) float64
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "OverflowToInfinityFloat64WithSelector",
			args: args{
				source:   NewOnSlice("x", "y"),
				selector: func(string) float64 { return math.MaxFloat64 },
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SumMust(tt.args.source, tt.args.selector)
			want := math.IsInf(got, +1)
			if want != tt.want {
				t.Errorf("IsInf(Sum()) = %v, want %v", want, tt.want)
			}
		})
	}
}

func Test_Average_int_int(t *testing.T) {
	type args struct {
		source   Enumerator[int]
		selector func(int) int
	}
	tests := []struct {
		name        string
		args        args
		want        float64
		wantErr     bool
		expectedErr error
	}{
		{name: "EmptySequenceIntNoSelector",
			args: args{
				source: Empty[int](),
			},
			wantErr:     true,
			expectedErr: ErrNilSelector,
		},
		{name: "SimpleAverageInt",
			args: args{
				source:   NewOnSlice(5, 10, 0, 15),
				selector: Identity[int],
			},
			want: 7.5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Average(tt.args.source, tt.args.selector)
			if (err != nil) != tt.wantErr {
				t.Errorf("Average() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if err != tt.expectedErr {
					t.Errorf("Average() error = %v, expectedErr %v", err, tt.expectedErr)
				}
				return
			}
			if got != tt.want {
				t.Errorf("Average() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_Average_string_int(t *testing.T) {
	type args struct {
		source   Enumerator[string]
		selector func(string) int
	}
	tests := []struct {
		name        string
		args        args
		want        float64
		wantErr     bool
		expectedErr error
	}{
		{name: "SourceStrNilSelector",
			args: args{
				source: NewOnSlice("one", "two", "three", "four"),
			},
			wantErr:     true,
			expectedErr: ErrNilSelector,
		},
		{name: "EmptySequenceIntWithSelector",
			args: args{
				source:   Empty[string](),
				selector: func(s string) int { return len(s) },
			},
			wantErr:     true,
			expectedErr: ErrEmptySource,
		},
		{name: "SimpleAverageIntWithSelector",
			args: args{
				source:   NewOnSlice("", "abcd", "a", "b"),
				selector: func(s string) int { return len(s) },
			},
			want: 1.5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Average(tt.args.source, tt.args.selector)
			if (err != nil) != tt.wantErr {
				t.Errorf("Average() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if err != tt.expectedErr {
					t.Errorf("Average() error = %v, expectedErr %v", err, tt.expectedErr)
				}
				return
			}
			if got != tt.want {
				t.Errorf("Average() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnum_Average_string_float64IsNaN(t *testing.T) {
	type args struct {
		source   Enumerator[string]
		selector func(string) float64
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "SequenceContainingNan",
			args: args{
				source: NewOnSlice("x", "abc", "de"),
				selector: func(s string) float64 {
					l := len(s)
					if l == 3 {
						return math.NaN()
					}
					return float64(l)
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AverageMust(tt.args.source, tt.args.selector)
			want := math.IsNaN(got)
			if want != tt.want {
				t.Errorf("IsNaN(Average()) = %v, want %v", want, tt.want)
			}
		})
	}
}

func TestEnum_Average_float64_float64IsInf(t *testing.T) {
	type args struct {
		source   Enumerator[float64]
		selector func(float64) float64
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "Float64OverflowsToInfinity",
			args: args{
				source:   NewOnSlice(math.MaxFloat64, math.MaxFloat64, -math.MaxFloat64, -math.MaxFloat64),
				selector: Identity[float64],
			},
			want: true,
		},
		{name: "Float64OverflowsToNegInfinity",
			args: args{
				source:   NewOnSlice(-math.MaxFloat64, -math.MaxFloat64, math.MaxFloat64, math.MaxFloat64),
				selector: Identity[float64],
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AverageMust(tt.args.source, tt.args.selector)
			var want bool
			switch tt.name {
			case "Float64OverflowsToInfinity":
				want = math.IsInf(got, +1)
			case "Float64OverflowsToNegInfinity":
				want = math.IsInf(got, -1)
			}
			if want != tt.want {
				t.Errorf("IsInf(Average()) = %v, want %v", want, tt.want)
			}
		})
	}
}
