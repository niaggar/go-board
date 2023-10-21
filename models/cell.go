package models

type Cell struct {
	objects []*int
}

func NewCell() *Cell {
	return &Cell{}
}

func (c *Cell) GetObjects() []*int {
	return c.objects
}

func (c *Cell) AddObject(id int) {
	c.objects = append(c.objects, &id)
}

func (c *Cell) Clear() {
	c.objects = c.objects[:0]
}
