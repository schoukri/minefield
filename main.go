package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {

	filePath := flag.String("file", "input.txt", "file containing the input data")
	flag.Parse()

	file, err := os.Open(*filePath)
	if err != nil {
		log.Fatalf("error: cannot open file %s: %v", *filePath, err)
	}

	//nolint: errcheck
	defer file.Close()

	// parse the input data and load all the mines into the minefield
	field, err := ParseData(file)
	if err != nil {
		log.Fatalf("error: cannot parse input data into a minefield: %v", err)
	}

	// for each mine, start a chain reacion and record the peak time interval with the most explosions
	for _, mine := range field.Mines {
		peak := field.StartExplosions(mine)
		mine.Peak = peak
	}

	// sort the mines by its peak time interval
	sort.Sort(byPeakExplosions(field.Mines))

	// print the winning mine with the highest number of explosions in its peak interval
	// (if there are multiple mines with the same highest explosion count, print them all)
	maxExplosions := field.Mines[0].Peak.Explosions
	for i, m := range field.Mines {
		if m.Peak.Explosions < maxExplosions {
			break
		}
		fmt.Printf("Winner (%d): Mine ID=%d, X=%f, Y=%f, Peak Time=%v, Peak Explosions=%v\n", i, m.ID, m.X, m.Y, m.Peak.Time, m.Peak.Explosions)
	}

}

// Mine describes the position and explosive power of a mine (and other information).
type Mine struct {
	ID        int
	X, Y      float32
	Power     float32
	Exploded  bool
	Neighbors []*Mine
	Peak      *Interval
}

// Field contains the list of mines on the minefield.
type Field struct {
	Mines []*Mine
}

// Distance returns the straight line distance between two mines.
func (m *Mine) Distance(m2 *Mine) float32 {
	xDiffSquared := math.Pow(float64(m2.X-m.X), 2)
	yDiffSquared := math.Pow(float64(m2.Y-m.Y), 2)
	return float32(math.Sqrt(xDiffSquared + yDiffSquared))
}

// Explode explodes the given mine and returns true on success.
// It returns false if the mine was already exploded.
func (m *Mine) Explode() bool {

	// cannot explode a mine that's already been exploded (return false)
	if m.Exploded {
		return false
	}

	// otherwise, explode the mine and return true for success
	m.Exploded = true

	return true
}

// AddMine adds the given mine to the mine field (returns an error on failure)
func (f *Field) AddMine(mine *Mine) error {

	if mine.Power < 0 {
		return fmt.Errorf("error: cannot add mine with negative explosive Power=%f, ID=%d", mine.Power, mine.ID)
	}

	// pre-calculate the neighbors for each mine as it is added
	mine.Neighbors = make([]*Mine, 0)
	for _, neighbor := range f.Mines {

		if mine.ID == neighbor.ID {
			return fmt.Errorf("error: cannot add two mines with same ID = %d", mine.ID)
		}

		if mine.X == neighbor.X && mine.Y == neighbor.Y {
			return fmt.Errorf("error: cannot add two mines with same coordinates X=%f, Y=%f", mine.X, mine.Y)
		}

		// distance between mine and neighbor
		distance := mine.Distance(neighbor)

		// is the distance within the mine's explosive power?
		// (if so, add the neighbor to the mines's list of neighbors)
		if distance <= mine.Power {
			mine.Neighbors = append(mine.Neighbors, neighbor)
		}

		// is the distance within the neighbors's explosive power?
		// (if so, add the mine to the neighbor's list of neighbors)
		if distance <= neighbor.Power {
			neighbor.Neighbors = append(neighbor.Neighbors, mine)
		}

	}

	f.Mines = append(f.Mines, mine)

	return nil
}

// Interval contains a time value and the number of explosions during that time.
type Interval struct {
	Time       int
	Explosions int
}

// byPeakExplosions implements sort.Interface for a slice of Mines
type byPeakExplosions []*Mine

// Len returns the number of mines in the slice.
func (m byPeakExplosions) Len() int {
	return len(m)
}

// Swap swaps the Mines in the slice specified by the indices i and j.
func (m byPeakExplosions) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

// Less returns true if the Mine at index i should be ordered before the Mine at index j.
func (m byPeakExplosions) Less(i, j int) bool {

	// primary sort: higher peak explosion count
	if m[i].Peak.Explosions != m[j].Peak.Explosions {
		return m[i].Peak.Explosions > m[j].Peak.Explosions
	}

	// secondary sort: lower X coordinate
	if m[i].X != m[j].X {
		return m[i].X < m[j].X
	}

	// tertiary sort: lower Y coordinate
	return m[i].Y < m[j].Y

}

// StartExplosions explodes the input mine and starts a chain reaction of exploding nearby mines.
// It returns the time interval with the peak number of explosions.
func (f *Field) StartExplosions(m *Mine) *Interval {

	// make sure all the mines are unexploded at the start of a new chain reaction
	for _, m := range f.Mines {
		m.Exploded = false
	}

	// keep track of the number of explosions per time interval
	explosions := make(map[int]int)

	// start the chain reaction of explosions with the input mine
	mines := []*Mine{m}
	for time := 0; len(mines) > 0; time++ {
		neighbors := make([]*Mine, 0)
		for _, m := range mines {

			// skip any mines that have already been exploded.
			if ok := m.Explode(); !ok {
				continue
			}

			explosions[time]++

			// build up a list of unexploded neighboring mines to explode in the next time interval
			for _, neighbor := range m.Neighbors {
				if !neighbor.Exploded {
					neighbors = append(neighbors, neighbor)
				}
			}

		}
		mines = neighbors
	}

	// find the peak interval
	intervals := make([]*Interval, 0)
	for time, count := range explosions {
		intervals = append(intervals, &Interval{time, count})
	}

	// sort the intervals by explosion count in desc order
	sort.Slice(intervals, func(i, j int) bool {
		// primary sort: higher explosion count
		if intervals[i].Explosions != intervals[j].Explosions {
			return intervals[i].Explosions > intervals[j].Explosions
		}

		// secondary sort: earlier time interval
		return intervals[i].Time < intervals[j].Time

	})

	// return the first (peak) interval
	return intervals[0]
}

// ParseData builds a mine field from the input lines of data.
func ParseData(data io.Reader) (*Field, error) {

	id := 0
	field := &Field{}

	scanner := bufio.NewScanner(data)
	for scanner.Scan() {

		id++
		line := scanner.Text()

		fields := strings.Fields(line)
		if len(fields) != 3 {
			return nil, fmt.Errorf("error: line line '%s' does not have 3 fields: %d", line, len(fields))
		}

		x, err := strconv.ParseFloat(fields[0], 32)
		if err != nil {
			return nil, fmt.Errorf("error: line line '%s' invalid float for X coord: %v", line, fields[0])
		}

		y, err := strconv.ParseFloat(fields[1], 32)
		if err != nil {
			return nil, fmt.Errorf("error: line line '%s' invalid float for Y coord: %v", line, fields[1])
		}

		power, err := strconv.ParseFloat(fields[2], 32)
		if err != nil {
			return nil, fmt.Errorf("error: line line '%s' invalid float for power: %v", line, fields[2])
		}

		mine := &Mine{
			ID:    id,
			X:     float32(x),
			Y:     float32(y),
			Power: float32(power),
		}

		if err := field.AddMine(mine); err != nil {
			return nil, err
		}

	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return field, nil
}
