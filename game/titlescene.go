// SPDX-License-Identifier: Apache-2.0

package game

import (
	"image/color"

	"github.com/hajimehoshi/bitmapfont"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/hajimehoshi/ebiten/text"
	"golang.org/x/image/font"
)

type TitleScene struct{}

func (*TitleScene) Update(sceneManager SceneManager) error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		sceneManager.GoToGameScene()
	}
	return nil
}

func (*TitleScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x33, 0x66, 0x99, 0xff})

	f := bitmapfont.Gothic12r
	b, _ := font.BoundString(f, GameTitle)
	x := -b.Min.X.Round() + (GameWidth-(b.Max.X-b.Min.X).Round())/2
	y := -b.Min.Y.Round() + (GameHeight-(b.Max.Y-b.Min.Y).Round())/2 - 32
	text.Draw(screen, GameTitle, f, x, y, color.White)

	msg := "Press Enter"
	b, _ = font.BoundString(f, msg)
	x = -b.Min.X.Round() + (GameWidth-(b.Max.X-b.Min.X).Round())/2
	y = -b.Min.Y.Round() + (GameHeight-(b.Max.Y-b.Min.Y).Round())/2 + 32
	text.Draw(screen, msg, f, x, y, color.White)
}
