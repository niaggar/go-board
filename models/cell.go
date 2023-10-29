package models

type Cell struct {
	balls     []*int
	obstacles []*int
}

func NewCell() *Cell {
	return &Cell{}
}

func (c *Cell) GetObjects() ([]*int, []*int) {
	return c.balls, c.obstacles
}

func (c *Cell) AddObject(id int) {
	c.balls = append(c.balls, &id)
}

func (c *Cell) AddObstacle(id int) {
	c.obstacles = append(c.obstacles, &id)
}

func (c *Cell) Clear() {
	c.balls = c.balls[:0]
	c.obstacles = c.obstacles[:0]
}
