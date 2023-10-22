package models

import "go-board/gmath"

type Mesh struct {
	Columns int
	Rows    int
	MaxSize float32
	Cells   []*Cell
}

func NewMesh(bounds gmath.Vector, maxSize float32) *Mesh {
	columns := int(bounds.X/maxSize) + 1
	rows := int(bounds.Y/maxSize) + 1

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

func (m *Mesh) AddObject(position gmath.Vector, id int) {
	cell := m.getCell(position)
	cell.AddObject(id)
}

func (m *Mesh) AddObstacle(position gmath.Vector, id int) {
	cell := m.getCell(position)
	cell.AddObstacle(id)
}

func (m *Mesh) getCell(position gmath.Vector) *Cell {
	x := int(position.X / m.MaxSize)
	y := int(position.Y / m.MaxSize)

	return m.Cells[y*m.Columns+x]
}

func (m *Mesh) GetElements(col, row int) ([]*int, []*int) {
	cell := m.Cells[row*m.Columns+col]
	return cell.GetObjects()
}

func (m *Mesh) GetElementsAround(col, row int) ([]*int, []*int) {
	var objects []*int
	var obstacles []*int

	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			x := col + i
			y := row + j

			if x >= 0 && x < m.Columns && y >= 0 && y < m.Rows {
				cell := m.Cells[y*m.Columns+x]
				objCell, obsCell := cell.GetObjects()

				objects = append(objects, objCell...)
				obstacles = append(obstacles, obsCell...)
			}
		}
	}

	return objects, obstacles
}

func (m *Mesh) Clear() {
	for i := 0; i < len(m.Cells); i++ {
		m.Cells[i].Clear()
	}
}
