//go:build go1.18

package go2linq

import (
	"testing"
)

// https://github.com/jskeet/edulinq/blob/master/src/Edulinq.Tests/CountTest.cs

func Test_Count_int(t *testing.T) {
	type args struct {
		source Enumerator[int]
	}
	tests := []struct {
		name        string
		args        args
		want        int
		wantErr     bool
		expectedErr error
	}{
		{name: "NonCollectionCount",
			args: args{
				source: RangeMust(2, 5),
			},
			want: 5,
		},
		{name: "0",
			args: args{
				source: Empty[int](),
			},
			want: 0,
		},
		{name: "NullSourceThrowsArgumentNullException",
			wantErr:     true,
			expectedErr: ErrNilSource,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Count(tt.args.source)
			if (err != nil) != tt.wantErr {
				t.Errorf("Count() error = '%v', wantErr '%v'", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if err != tt.expectedErr {
					t.Errorf("Count() error = '%v', expectedErr '%v'", err, tt.expectedErr)
				}
				return
			}
			if got != tt.want {
				t.Errorf("Count() = '%v', want '%v'", got, tt.want)
			}
		})
	}
}

func Test_Count_string(t *testing.T) {
	type args struct {
		source Enumerator[string]
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{name: "1",
			args: args{
				source: NewOnSlice("zero", "one", "two", "three", "four", "five"),
			},
			want: 6,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := Count(tt.args.source); got != tt.want {
				t.Errorf("Count() = '%v', want '%v'", got, tt.want)
			}
		})
	}
}

func Test_CountPred_int(t *testing.T) {
	type args struct {
		source    Enumerator[int]
		predicate func(int) bool
	}
	tests := []struct {
		name        string
		args        args
		want        int
		wantErr     bool
		expectedErr error
	}{
		{name: "PredicatedNullSourceThrowsArgumentNullException",
			args: args{
				predicate: func(x int) bool { return x == 1 },
			},
			wantErr:     true,
			expectedErr: ErrNilSource,
		},
		{name: "PredicatedNullPredicateThrowsArgumentNullException",
			args: args{
				source: NewOnSlice(3, 5, 20, 15),
			},
			wantErr:     true,
			expectedErr: ErrNilPredicate,
		},
		{name: "PredicatedCount",
			args: args{
				source:    RangeMust(2, 5),
				predicate: func(x int) bool { return x%2 == 0 },
			},
			want: 3,
		},
		{name: "11",
			args: args{
				source:    NewOnSlice(1, 2, 3, 4),
				predicate: func(int) bool { return false },
			},
			want: 0,
		},
		{name: "12",
			args: args{
				source:    NewOnSlice(1, 2, 3, 4),
				predicate: func(int) bool { return true },
			},
			want: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CountPred(tt.args.source, tt.args.predicate)
			if (err != nil) != tt.wantErr {
				t.Errorf("CountPred() error = '%v', wantErr '%v'", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if err != tt.expectedErr {
					t.Errorf("CountPred() error = '%v', expectedErr '%v'", err, tt.expectedErr)
				}
				return
			}
			if got != tt.want {
				t.Errorf("CountPred() = '%v', want '%v'", got, tt.want)
			}
		})
	}
}

func Test_CountPred_string(t *testing.T) {
	type args struct {
		source    Enumerator[string]
		predicate func(string) bool
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{name: "21",
			args: args{
				source:    NewOnSlice("one", "two", "three", "four"),
				predicate: func(s string) bool { return len(s) == 3 },
			},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := CountPred(tt.args.source, tt.args.predicate); got != tt.want {
				t.Errorf("CountPred() = '%v', want '%v'", got, tt.want)
			}
		})
	}
}
