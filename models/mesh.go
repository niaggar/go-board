package models

import "go-board/gmath"

type Mesh struct {
	columns int
	rows    int
	cells   []Cell
}

func NewMesh(columns, rows int) *Mesh {
	cells := make([]Cell, columns*rows)

	for i := range cells {
		cells[i] = NewCell()
	}

	return &Mesh{
		columns: columns,
		rows:    rows,
		cells:   cells,
	}
}

func (m *Mesh) AddObject(position gmath.Vector, id int) {
	cell := m.getCell(position)
	cell.AddObject(id)
}

func (m *Mesh) getCell(position gmath.Vector) *Cell {
	x := int(position.X)
	y := int(position.Y)

	return &m.cells[y*m.columns+x]
}

func (m *Mesh) GetObjects(position gmath.Vector) []*int {
	return m.getCell(position).GetObjects()
}

func (m *Mesh) GetCellsAround(col, row int) []*int {
	cells := make([]*int, 0, 9)

	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			x := col + i
			y := row + j

			if x >= 0 && x < m.columns && y >= 0 && y < m.rows {
				cell := m.cells[y*m.columns+x]
				cells = append(cells, cell.GetObjects()...)
			}
		}
	}

	return cells
}

func (m *Mesh) Clear() {
	for _, cell := range m.cells {
		cell.Clear()
	}
}
