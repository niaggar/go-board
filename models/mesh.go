package models

import "go-board/utils/gmath"

type Mesh struct {
	Columns int
	Rows    int
	MaxSize float32
	Cells   []*Cell
}

func NewMesh(bounds Bounds, maxSize float32) *Mesh {
	columns := int(bounds.Width/maxSize) + 1
	rows := int(bounds.Height/maxSize) + 1

	cells := make([]*Cell, columns*rows)

	for i := 0; i < len(cells); i++ {
		cells[i] = NewCell()
	}

	return &Mesh{
		MaxSize: maxSize,
		Columns: columns,
		Rows:    rows,
		Cells:   cells,
	}
}

func GetCell(m *Mesh, position gmath.Vector) *Cell {
	x := int(position.X / m.MaxSize)
	y := int(position.Y / m.MaxSize)

	return m.Cells[y*m.Columns+x]
}

func (m *Mesh) AddBall(position gmath.Vector, id int) {
	cell := GetCell(m, position)
	cell.AddObject(id)
}

func (m *Mesh) AddObstacle(position gmath.Vector, id int) {
	cell := GetCell(m, position)
	cell.AddObstacle(id)
}

func (m *Mesh) GetElementsAround(col, row int) ([]*int, []*int) {
	var balls []*int
	var obstacles []*int

	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			x := col + i
			y := row + j

			if x >= 0 && x < m.Columns && y >= 0 && y < m.Rows {
				cell := m.Cells[y*m.Columns+x]
				ballCell, obsCell := cell.GetObjects()

				balls = append(balls, ballCell...)
				obstacles = append(obstacles, obsCell...)
			}
		}
	}

	return balls, obstacles
}

func (m *Mesh) Clear() {
	for i := 0; i < len(m.Cells); i++ {
		m.Cells[i].Clear()
	}
}
