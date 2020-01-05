// SPDX-License-Identifier: Apache-2.0

package game

import (
	"fmt"
	"image/color"
	"strings"

	"github.com/hajimehoshi/bitmapfont"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/hajimehoshi/ebiten/text"
	"golang.org/x/image/font"
)

type Grid byte

const (
	GridFloor Grid = iota
	GridGoal
	GridWall
)

const GridSize = 16

var (
	tmpFloorImage, _  = ebiten.NewImage(GridSize, GridSize, ebiten.FilterDefault)
	tmpWallImage, _   = ebiten.NewImage(GridSize, GridSize, ebiten.FilterDefault)
	tmpGoalImage, _   = ebiten.NewImage(GridSize, GridSize, ebiten.FilterDefault)
	tmpBoxImage, _    = ebiten.NewImage(GridSize, GridSize, ebiten.FilterDefault)
	tmpPlayerImage, _ = ebiten.NewImage(GridSize, GridSize, ebiten.FilterDefault)
)

func init() {
	tmpFloorImage.Fill(color.RGBA{0xff, 0xcc, 0x99, 0xff})
	tmpGoalImage.Fill(color.RGBA{0xff, 0xcc, 0x99, 0xff})
	// TODO: Use vector.
	ebitenutil.DrawRect(tmpGoalImage, 1, 1, 14, 14, color.RGBA{0x66, 0x00, 0x00, 0x66})
	tmpWallImage.Fill(color.RGBA{0x66, 0x33, 0x00, 0xff})

	ebitenutil.DrawRect(tmpBoxImage, 1, 1, 14, 14, color.RGBA{0x99, 0x99, 0x33, 0xff})
	ebitenutil.DrawRect(tmpPlayerImage, 1, 1, 14, 14, color.RGBA{0x33, 0x33, 0xcc, 0xff})
}

type Player struct {
	x int
	y int
}

type Box struct {
	x int
	y int
}

type Field struct {
	initGrids  [][]Grid
	initPlayer *Player
	initBoxes  []*Box

	grids  [][]Grid
	player *Player
	boxes  []*Box
	steps  int
}

func ParseField(str string) *Field {
	str = strings.TrimSpace(str)
	f := &Field{}
	for j, line := range strings.Split(str, "\n") {
		var gs []Grid
		for i, ch := range line {
			switch ch {
			case ' ':
				gs = append(gs, GridFloor)
			case '#':
				gs = append(gs, GridWall)
			case '.':
				gs = append(gs, GridGoal)
			case '@':
				gs = append(gs, GridFloor)
				f.initPlayer = &Player{
					x: i,
					y: j,
				}
			case '$':
				gs = append(gs, GridFloor)
				f.initBoxes = append(f.initBoxes, &Box{
					x: i,
					y: j,
				})
			}

		}
		f.initGrids = append(f.initGrids, gs)
	}
	f.reset()
	return f
}

func (f *Field) reset() {
	f.steps = 0
	f.grids = nil
	f.player = nil
	f.boxes = nil

	for _, gs := range f.initGrids {
		newGs := make([]Grid, len(gs))
		copy(newGs, gs)
		f.grids = append(f.grids, newGs)
	}
	f.player = &Player{
		x: f.initPlayer.x,
		y: f.initPlayer.y,
	}
	for _, b := range f.initBoxes {
		f.boxes = append(f.boxes, &Box{
			x: b.x,
			y: b.y,
		})
	}
}

func (f *Field) gridAt(x, y int) Grid {
	if x < 0 {
		return GridFloor
	}
	if y < 0 {
		return GridFloor
	}
	if len(f.grids) <= y {
		return GridFloor
	}
	if len(f.grids[y]) <= x {
		return GridFloor
	}
	return f.grids[y][x]
}

func (f *Field) boxAt(x, y int) *Box {
	for _, c := range f.boxes {
		if c.x == x && c.y == y {
			return c
		}
	}
	return nil
}

func (f *Field) canMove(dx, dy int) (ok bool, box *Box) {
	x, y := f.player.x+dx, f.player.y+dy
	if f.gridAt(x, y) == GridWall {
		return false, nil
	}
	c := f.boxAt(x, y)
	if c == nil {
		return true, nil
	}
	if f.gridAt(x+dx, y+dy) == GridWall {
		return false, nil
	}
	if f.boxAt(x+dx, y+dy) != nil {
		return false, nil
	}
	return true, c
}

func (f *Field) tryMove(dx, dy int) {
	ok, box := f.canMove(dx, dy)
	if !ok {
		return
	}
	f.player.x += dx
	f.player.y += dy
	if box != nil {
		box.x += dx
		box.y += dy
	}
	f.steps++
}

func (f *Field) Update() {
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		f.tryMove(0, -1)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		f.tryMove(0, 1)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		f.tryMove(-1, 0)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		f.tryMove(1, 0)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		f.reset()
	}
}

func (f *Field) sizeInGrids() (int, int) {
	h := len(f.grids)
	w := 0
	for _, gs := range f.grids {
		if len(gs) > w {
			w = len(gs)
		}
	}
	return w, h
}

func (f *Field) Draw(screen *ebiten.Image) {
	w, h := f.sizeInGrids()
	offsetX, offsetY := (GameWidth-w*GridSize)/2, (GameHeight-h*GridSize)/2

	for j, gs := range f.grids {
		for i, g := range gs {
			var img *ebiten.Image
			switch g {
			case GridWall:
				img = tmpWallImage
			case GridGoal:
				img = tmpGoalImage
			default:
				img = tmpFloorImage
			}
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(i)*GridSize, float64(j)*GridSize)
			op.GeoM.Translate(float64(offsetX), float64(offsetY))
			screen.DrawImage(img, op)
		}
	}
	for _, c := range f.boxes {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(c.x)*GridSize, float64(c.y)*GridSize)
		op.GeoM.Translate(float64(offsetX), float64(offsetY))
		screen.DrawImage(tmpBoxImage, op)
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(f.player.x)*GridSize, float64(f.player.y)*GridSize)
	op.GeoM.Translate(float64(offsetX), float64(offsetY))
	screen.DrawImage(tmpPlayerImage, op)

	// Render steps
	bf := bitmapfont.Gothic12r

	stepsMsg := fmt.Sprintf("Steps: %d", f.steps)
	b, _ := font.BoundString(bf, stepsMsg)
	x := -b.Min.X.Round() + 8
	y := -b.Min.Y.Round() + 8
	text.Draw(screen, stepsMsg, bf, x, y, color.White)

	// Render instructions
	instMsg := "Press R to reset"
	b, _ = font.BoundString(bf, instMsg)
	x = -b.Min.X.Round() + 8
	y = -b.Min.Y.Round() + 24
	text.Draw(screen, instMsg, bf, x, y, color.White)
}
