package models

type Cell struct {
	objects   []*int
	obstacles []*int
}

func NewCell() *Cell {
	return &Cell{}
}

func (c *Cell) GetObjects() ([]*int, []*int) {
	return c.objects, c.obstacles
}

func (c *Cell) AddObject(id int) {
	c.objects = append(c.objects, &id)
}

func (c *Cell) AddObstacle(id int) {
	c.obstacles = append(c.obstacles, &id)
}

func (c *Cell) Clear() {
	c.objects = c.objects[:0]
	c.obstacles = c.obstacles[:0]
}
