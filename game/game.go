// SPDX-License-Identifier: Apache-2.0

package game

import (
	"github.com/hajimehoshi/ebiten"
)

const (
	GameTitle  = "Boxes"
	GameWidth  = 320
	GameHeight = 240
)

type Game struct {
	scene     Scene
	nextScene Scene
}

func (*Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return GameWidth, GameHeight
}

func (g *Game) Update(screen *ebiten.Image) error {
	if g.nextScene != nil {
		g.scene = g.nextScene
		g.nextScene = nil
	}
	if g.scene == nil {
		g.scene = &TitleScene{}
	}
	if err := g.scene.Update(g); err != nil {
		return err
	}
	if ebiten.IsDrawingSkipped() {
		return nil
	}
	g.scene.Draw(screen)
	return nil
}

func (g *Game) GoToGameScene() {
	g.nextScene = &GameScene{}
}

type Scene interface {
	Update(sceneManager SceneManager) error
	Draw(screen *ebiten.Image)
}

type SceneManager interface {
	GoToGameScene()
}
