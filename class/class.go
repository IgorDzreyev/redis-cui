package class

import (
	"log"
	"strconv"
	"strings"

	"github.com/kyai/gocui"
	"github.com/kyai/redis-cui/common"
)

type Class interface {
	AddRow(s ...string)
	Render(v *gocui.View)
	Column() []string
	Widths() []int
}

func NewClass(t string, g *gocui.Gui) Class {
	base := Base{
		g:    g,
		Type: t,
	}
	switch t {
	case "STRING":
		return &String{base}
	case "LIST", "SET":
		return &List{base}
	case "HASH":
		return &Hash{base}
	case "ZSET":
		return &Zset{base}
	case "STREAM":
		return &Stream{base}
	default:
		log.Panic("Unknown type: ", t)
		return nil
	}
}

type Base struct {
	g    *gocui.Gui
	Type string
	Rows [][]string
}

func (b *Base) AddRow(s ...string) {
	for k, v := range s {
		s[k] = strings.ReplaceAll(v, "\n", "")
	}
	b.Rows = append(b.Rows, s)
}

func (b *Base) Widths(c []string) []int {
	w := make([]int, len(c))
	for _, row := range b.Rows {
		for k, v := range row {
			if l := len(v); l > w[k] {
				w[k] = l
			}
		}
	}
	for k, v := range c {
		if l := len(v); l > w[k] {
			w[k] = l
		}
	}
	return w
}

func (b *Base) rownum() {
	for i, row := range b.Rows {
		b.Rows[i] = append([]string{strconv.Itoa(i + 1)}, row...)
	}
}

func (b *Base) render(view *gocui.View, rows [][]string, c []string, w []int) {
	var title string
	for k, v := range c {
		if k < len(c)-1 {
			title += common.FillRight(v, '─', w[k])
		} else {
			title += v
		}
	}
	view.Title = title

	view.Clear()
	for i, row := range rows {
		line := " "
		for k, v := range row {
			line += common.FillRight(v, ' ', w[k])
		}
		view.Write([]byte(line))
		if i < len(rows)-1 {
			view.Write([]byte("\n"))
		}
	}
}
