package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func (gc *GameController) savemenu(savesmeta []SaveMeta) {
	gc.menu = append(gc.menu, gc.buildSaveMenu(savesmeta))
}

func (gc *GameController) buildSaveMenu(savesmeta []SaveMeta) Menu {
	savesMenu := newMenu(rl.Vector2{X: 0, Y: 0})

	choosesave := savesMenu.newTextBox(rl.Vector2{X: 300, Y: 150}, "Choose Save", 50)
	savesMenu.textbox = append(savesMenu.textbox, choosesave)

	for i := range savesmeta {
		name := savesmeta[i].PlayerName

		if name == "" {
			continue
		}
		x := float32(100 + 400*(i%2))
		y := float32(350 + 220*(i/2))
		loadButton := savesMenu.newButton(func() {
			err := gc.LoadGame(name)
			if err == nil {
				gc.StartGame()
			}
		}, rl.Vector2{X: x, Y: y}, 300, 150, name, int(smallbutton), false)

		deleteButton := savesMenu.newButton(func() { gc.DeleteSave(name) }, rl.Vector2{X: x + 310, Y: y}, 40, 40, "", int(smallbutton), true)

		savesMenu.buttons = append(savesMenu.buttons, loadButton, deleteButton)
	}
	newsave := savesMenu.newButton(func() { gc.newsavebutton() }, rl.Vector2{X: 50, Y: 800}, 100, 100, "NEW", int(smallbutton), false)
	backbutton := savesMenu.newButton(func() { gc.setstate(menu) }, rl.Vector2{X: 50, Y: 20}, 100, 100, "Back", int(smallbutton), false)
	savesMenu.buttons = append(savesMenu.buttons, backbutton, newsave)

	return savesMenu
}

func (gc *GameController) RefreshSaveMenu() {
	if int(selectsave) < len(gc.menu) {
		gc.menu[selectsave] = gc.buildSaveMenu(gc.LoadSaves())
	}
}

func (gc *GameController) newsavebutton() {
	b1 := gc.menu[selectsave].newButton(func() {}, rl.Vector2{X: 300, Y: 275}, 300, 50, "", int(button), false)
	b1.OnClick = func() { b1.allowUserInput = true }
	gc.menu[selectsave].buttons = append(gc.menu[selectsave].buttons, b1)
}
