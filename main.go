// SPDX-License-Identifier: Apache-2.0

package main

import (
	"github.com/hajimehoshi/ebiten"

	"github.com/hajimehoshi/boxes/game"
)

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle(game.GameTitle)
	if err := ebiten.RunGame(&game.Game{}); err != nil {
		panic(err)
	}
}
