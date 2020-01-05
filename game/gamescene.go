// SPDX-License-Identifier: Apache-2.0

package game

import (
	"github.com/hajimehoshi/ebiten"
)

type GameScene struct{
	field *Field
}

// See http://www.sokobano.de/wiki/index.php?title=Level_format
// TODO: Add other types

var testField = `
#######
#. . .#
# $$$ #
#.$@$.#
# $$$ #
#. . .#
#######
`

func (g *GameScene) Update(sceneManager SceneManager) error {
	if g.field == nil {
		g.field = ParseField(testField)
	}
	g.field.Update()
	return nil
}

func (g *GameScene) Draw(screen *ebiten.Image) {
	g.field.Draw(screen)
}
