package cursor

import (
    "math"
)

func clamp(n, lower, upper int) int {
    if n < lower {
        return lower
    }
    if n > upper {
        return upper
    }
    return n
}

//struct for handling cursor movement & text scrolling
type Cursor struct {
    x int
    targetX int //the target x of the cursor (for moving up/down through text)
    y int
    termY int //the vertical offset that the text will be drawn from (view scrolling)
}

func (c *Cursor) X() int {
    return c.x
}

func (c *Cursor) Y() int {
    return c.y
}

func (c *Cursor) TermY() int {
    return c.termY
}

func (c *Cursor) SetDefaultValues() {
    c.x = -1
    c.targetX = c.x
}

func (c *Cursor) scrollVertical(text []string, termHeight, n int) {
    c.y = clamp(c.y+n, 0, len(text)-1)
    c.x = clamp(c.targetX, -1, len(text[c.y])-1)

    //scrolling view
    if (n > 0 && c.y >= c.termY+termHeight) || (n < 0 && c.y < c.termY) {
        c.termY = clamp(c.termY+n, 0, len(text)-1)
    }
}

func (c *Cursor) ScrollVertical(text []string, termHeight, n int) {
    if n != 0 {
        c.scrollVertical(text, termHeight, n)
    }
}

//horizontal scrolling requires iteration: large numbers could require moving over multiple lines
func (c *Cursor) scrollHorizontal(text []string, termHeight, n int) {
    sign := n / int(math.Abs(float64(n)))
    for i := 0; i < int(math.Abs(float64(n))); i++ {
        c.x += sign
        //moving onto another line
        if c.x < -1 || (c.x > len(text[c.y])-1) {
            c.ScrollVertical(text, termHeight, sign)
            c.x = -1
            if sign < 0 {
                c.x = len(text[c.y]) - 1
            }
        }
        c.targetX = c.x
    }
}

func (c *Cursor) ScrollHorizontal(text []string, termheight, n int) {
    if n != 0 {
        c.scrollHorizontal(text, termheight, n)
    }
}

//very basic method - should only be used externally, not for internal methods where clamping is preferred
func (c *Cursor) ScrollToLineStart() {
    c.x = -1
    c.targetX = c.x
}

//very basic method - should only be used externally, not for internal methods where clamping is preferred
func (c *Cursor) ScrollToLineEnd(text []string) {
    c.x = len(text[c.y]) - 1
    c.targetX = c.x
}
