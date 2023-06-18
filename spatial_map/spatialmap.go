package spatialmap

import "math"

// For a possibly for scalable solution, reference:
// https://carmencincotti.com/2022-10-31/spatial-hash-maps-part-one/#what-is-a-spatial-hash-table

type SpatialMap[T any] struct {
	cells map[coordinate][]*T

	spacing int // constant
}
type coordinate struct {
	x, y int
}

func NewSpatialMap[T any](spacing int) SpatialMap[T] {
	return SpatialMap[T]{
		cells: make(map[coordinate][]*T),

		spacing: spacing,
	}
}

func (s *SpatialMap[T]) Add(x, y int, data *T) {
	coord := coordinate{
		int(math.Floor(float64(x) / float64(s.spacing))),
		int(math.Floor(float64(y) / float64(s.spacing))),
	}
	s.cells[coord] = append(s.cells[coord], data)
}

func (s *SpatialMap[T]) Nearby(x, y, radius int) []*T {
	nearby := []*T{}
	for i := x - radius; i <= x+radius; i++ {
		for j := x - radius; j <= x+radius; j++ {
			nearby = append(nearby, s.cells[coordinate{i, j}]...)
		}
	}
	return nearby
}

func (s *SpatialMap[T]) Get(x, y int) []*T {
	return s.cells[coordinate{x, y}]
}
