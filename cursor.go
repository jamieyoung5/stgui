package stgui

type Cursor struct {
	x, y, gridX, gridY int
	controls           map[string]string
	selectedColour     string
}

func NewCursor(x, y, gridX, gridY int, controls map[string]string) *Cursor {
	return &Cursor{
		x:        x,
		y:        y,
		gridX:    gridX,
		gridY:    gridY,
		controls: controls,
	}
}

func (c *Cursor) Coords() (x int, y int) {
	return c.x, c.y
}

func (c *Cursor) GridCoords() (x int, y int) {
	return c.gridX, c.gridY
}

func (c *Cursor) Up() {
	c.y--
}

func (c *Cursor) Down() {
	c.y++
}

func (c *Cursor) Right() {
	c.x++
}

func (c *Cursor) Left() {
	c.x--
}
