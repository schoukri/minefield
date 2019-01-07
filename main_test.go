package main

import (
	"io"
	"math"
	"reflect"
	"sort"
	"strings"
	"testing"
)

func TestMine_Distance(t *testing.T) {
	tests := []struct {
		name string
		m1   *Mine
		m2   *Mine
		want float32
	}{
		{
			name: "sqrt 2",
			m1:   &Mine{X: 1, Y: 1},
			m2:   &Mine{X: 2, Y: 2},
			want: float32(math.Sqrt(2)), // 1.4142135,
		},
		{
			name: "sqrt 25 (3-4-5 triangle)",
			m1:   &Mine{X: 0, Y: 0},
			m2:   &Mine{X: 3, Y: 4},
			want: float32(math.Sqrt(25)), // 5
		},
		{
			name: "sqrt 32 (zero negative coordinates)",
			m1:   &Mine{X: 2, Y: 2},
			m2:   &Mine{X: 6, Y: 6},
			want: float32(math.Sqrt(32)), // 5.656854,
		},
		{
			name: "sqrt 32 (one negative coordinate)",
			m1:   &Mine{X: 2, Y: 2},
			m2:   &Mine{X: -2, Y: -2},
			want: float32(math.Sqrt(32)), // 5.656854,
		},
		{
			name: "sqrt 32 (two negative coordinates)",
			m1:   &Mine{X: -6, Y: -6},
			m2:   &Mine{X: -2, Y: -2},
			want: float32(math.Sqrt(32)), // 5.656854,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m1.Distance(tt.m2); !almostEqual(got, tt.want) {
				t.Errorf("Mine.Distance() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestField_AddMine(t *testing.T) {

	field := &Field{}

	tests := []struct {
		name    string
		mine    *Mine
		wantErr bool
	}{
		{
			name:    "one",
			mine:    &Mine{ID: 1, X: 1, Y: 1, Power: 1},
			wantErr: false,
		},
		{
			name:    "two",
			mine:    &Mine{ID: 2, X: 2, Y: 2, Power: 2},
			wantErr: false,
		},
		{
			name:    "two again (same ID)",
			mine:    &Mine{ID: 2, X: 3, Y: 3, Power: 3},
			wantErr: true,
		},
		{
			name:    "two again (same X,Y)",
			mine:    &Mine{ID: 3, X: 2, Y: 2, Power: 2},
			wantErr: true,
		},
		{
			name:    "negative power",
			mine:    &Mine{ID: 3, X: 3, Y: 3, Power: -3},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := field.AddMine(tt.mine); (err != nil) != tt.wantErr {
				t.Errorf("Field.AddMine() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMine_Neighbors(t *testing.T) {

	m1 := &Mine{ID: 1, X: 1, Y: 1, Power: 2}
	m2 := &Mine{ID: 2, X: 2, Y: 2, Power: 2}
	m3 := &Mine{ID: 3, X: 3, Y: 3, Power: 4}
	m4 := &Mine{ID: 4, X: 4, Y: 4, Power: 5}

	f := &Field{}
	for _, m := range []*Mine{m1, m2, m3, m4} {
		f.AddMine(m)
	}

	tests := []struct {
		name string
		mine *Mine
		want []*Mine
	}{
		{
			name: "one",
			mine: m1,
			want: []*Mine{m2},
		},
		{
			name: "two",
			mine: m2,
			want: []*Mine{m1, m3},
		},
		{
			name: "three",
			mine: m3,
			want: []*Mine{m1, m2, m4},
		},
		{
			name: "four",
			mine: m4,
			want: []*Mine{m1, m2, m3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.mine.Neighbors; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Mine.Neighbors = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestField_StartExplosions(t *testing.T) {

	m1 := &Mine{ID: 1, X: 1, Y: 1, Power: 1.5} // Moderate
	m2 := &Mine{ID: 2, X: 1, Y: 2, Power: 1.5} // Moderate
	m3 := &Mine{ID: 3, X: 1, Y: 3, Power: 1.5} // Moderate
	m4 := &Mine{ID: 4, X: 2, Y: 1, Power: 1.5} // Moderate
	m5 := &Mine{ID: 5, X: 2, Y: 2, Power: 1.5} // Moderate
	m6 := &Mine{ID: 6, X: 2, Y: 3, Power: 1.5} // Moderate
	m7 := &Mine{ID: 7, X: 3, Y: 1, Power: 5.0} // Powerful
	m8 := &Mine{ID: 8, X: 3, Y: 2, Power: 0.9} // Dud
	m9 := &Mine{ID: 9, X: 3, Y: 3, Power: 1.1} // Weak

	f := &Field{}
	for _, m := range []*Mine{m1, m2, m3, m4, m5, m6, m7, m8, m9} {
		f.AddMine(m)
	}

	tests := []struct {
		name  string
		field *Field
		mine  *Mine
		want  *Interval
	}{
		{
			name:  "one",
			field: f,
			mine:  m1,
			want:  &Interval{Time: 2, Explosions: 5},
		},
		{
			name:  "two",
			field: f,
			mine:  m2,
			want:  &Interval{Time: 1, Explosions: 5},
		},
		{
			name:  "three",
			field: f,
			mine:  m3,
			want:  &Interval{Time: 2, Explosions: 5},
		},
		{
			name:  "four",
			field: f,
			mine:  m4,
			want:  &Interval{Time: 1, Explosions: 5},
		},
		{
			name:  "five",
			field: f,
			mine:  m5,
			want:  &Interval{Time: 1, Explosions: 8},
		},
		{
			name:  "six",
			field: f,
			mine:  m6,
			want:  &Interval{Time: 1, Explosions: 5},
		},
		{
			name:  "seven",
			field: f,
			mine:  m7,
			want:  &Interval{Time: 1, Explosions: 8},
		},
		{
			name:  "eight",
			field: f,
			mine:  m8,
			want:  &Interval{Time: 0, Explosions: 1},
		},
		{
			name:  "nine",
			field: f,
			mine:  m9,
			want:  &Interval{Time: 2, Explosions: 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.field.StartExplosions(tt.mine)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Field.StartExplosions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMines_SortByExplosions(t *testing.T) {

	m1 := &Mine{ID: 1, X: 1, Y: 1, Peak: &Interval{Explosions: 6}}
	m2 := &Mine{ID: 2, X: 2, Y: 2, Peak: &Interval{Explosions: 5}}
	m3 := &Mine{ID: 3, X: 3, Y: 3, Peak: &Interval{Explosions: 5}}
	m4 := &Mine{ID: 4, X: 4, Y: 4, Peak: &Interval{Explosions: 4}}
	m5 := &Mine{ID: 5, X: 4, Y: 5, Peak: &Interval{Explosions: 4}}
	m6 := &Mine{ID: 6, X: 5, Y: 2, Peak: &Interval{Explosions: 4}}
	m7 := &Mine{ID: 7, X: 6, Y: 6, Peak: &Interval{Explosions: 3}}
	m8 := &Mine{ID: 8, X: 7, Y: 7, Peak: &Interval{Explosions: 2}}
	m9 := &Mine{ID: 9, X: 8, Y: 8, Peak: &Interval{Explosions: 1}}

	tests := []struct {
		name  string
		mines []*Mine
		want  []*Mine
	}{
		{
			name:  "all different explosions",
			mines: []*Mine{m9, m4, m8, m7, m2, m1},
			want:  []*Mine{m1, m2, m4, m7, m8, m9},
		},
		{
			name:  "same explosions, tiebreaker with X coord",
			mines: []*Mine{m3, m2},
			want:  []*Mine{m2, m3},
		},
		{
			name:  "same explosions, some with same X coord, tiebreaker with Y coord",
			mines: []*Mine{m5, m6, m4},
			want:  []*Mine{m4, m5, m6},
		},
		{
			name:  "overall",
			mines: []*Mine{m3, m9, m5, m6, m4, m7, m1, m8, m2},
			want:  []*Mine{m1, m2, m3, m4, m5, m6, m7, m8, m9},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sort.Sort(byPeakExplosions(tt.mines))
			if !reflect.DeepEqual(tt.mines, tt.want) {
				t.Errorf("sort.Sort(byPeakExplosions(mines))\n got = %+v\nwant = %+v", tt.mines, tt.want)
			}
		})
	}
}
func TestParseData(t *testing.T) {

	m1 := &Mine{ID: 1, X: 1, Y: 1, Power: 1}
	m2 := &Mine{ID: 2, X: 2, Y: 2, Power: 2}

	f := &Field{}
	for _, m := range []*Mine{m1, m2} {
		f.AddMine(m)
	}

	tests := []struct {
		name    string
		data    io.Reader
		want    *Field
		wantErr bool
	}{
		{
			name: "good data",
			data: strings.NewReader(
				`1 1 1
				 2 2 2`,
			),
			want:    f,
			wantErr: false,
		},
		{
			name: "incorrect number of fields",
			data: strings.NewReader(
				`1 1 1
				 2 2`,
			),
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid float in field 1",
			data: strings.NewReader(
				`- 1 1
				 2 2 2`,
			),
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid float in field 2",
			data: strings.NewReader(
				`1 - 1
				 2 2 2`,
			),
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid float in field 3",
			data: strings.NewReader(
				`1 1 -
				 2 2 2`,
			),
			want:    nil,
			wantErr: true,
		},
		{
			name: "duplicate coordinates",
			data: strings.NewReader(
				`1 1 1
				 1 1 2`,
			),
			want:    nil,
			wantErr: true,
		},
		{
			name: "negative power",
			data: strings.NewReader(
				`1 1 1
				 2 2 -2`,
			),
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseData(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.want == nil && got == nil {
				return
			}

			for i, want := range tt.want.Mines {
				// `Neighbors` should be the same logical mines, but DeepEqual() returns false because they are different pointers.
				// Set them to nil since it isn't the point of this test to validate `Neighbors`.
				got.Mines[i].Neighbors = nil
				want.Neighbors = nil
				if !reflect.DeepEqual(got.Mines[i], want) {
					t.Errorf("ParseData() = %+v, want %+v", got.Mines[i], want)
				}
			}
		})
	}
}

// almostEqual returns true if the input floats are within a small epsilon value of each other
func almostEqual(a, b float32) bool {
	epsilon := math.Nextafter(1, 2) - 1
	return math.Abs(float64(a-b)) <= epsilon
}
