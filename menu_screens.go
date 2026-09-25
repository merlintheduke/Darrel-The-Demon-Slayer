package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func (gc *GameController) PlayingUI() {
	playingMenu := gc.newMenu(rl.Vector2Zero())
	playingMenu.CustomDraw = func() {
		gc.DrawPlayerGUI()
	}
	upgradeButton := playingMenu.newButton(func() { gc.setstate(upgrade) }, rl.Vector2{X: 720, Y: 20}, 250, 50, "UPGRADE", int(smallbutton), false)
	pauseButton := playingMenu.newButton(func() { gc.setstate(pause) }, rl.Vector2{X: 720, Y: 80}, 250, 50, "PAUSE", int(smallbutton), false)
	playingMenu.buttons = append(playingMenu.buttons, upgradeButton, pauseButton)
	gc.menus[playing] = playingMenu
}

func (gc *GameController) gameovermenu() {
	gameOverMenu := gc.newMenu(rl.Vector2Zero())

	title := gameOverMenu.newTextBox(rl.Vector2{X: 330, Y: 220}, "GAME OVER", 70)

	message := gameOverMenu.newTextBox(rl.Vector2{X: 260, Y: 310}, "Your save has been deleted.", 30)

	mainMenu := gameOverMenu.newButton(
		func() {
			gc.setstate(menu)
			gc.ResetAfterGameOver()
		}, rl.Vector2{X: 350, Y: 430}, 300, 80, "MAIN MENU", int(button), false,
	)

	gameOverMenu.textbox = append(gameOverMenu.textbox, title, message)
	gameOverMenu.buttons = append(gameOverMenu.buttons, mainMenu)

	gc.menus[gameover] = gameOverMenu
}

func (gc *GameController) pausemenu() {
	pauseMenu := gc.newMenu(rl.Vector2Zero())

	title := pauseMenu.newTextBox(rl.Vector2{X: 410, Y: 180}, "PAUSED", 60)
	resume := pauseMenu.newButton(func() { gc.setstate(playing) }, rl.Vector2{X: 380, Y: 300}, 260, 70, "RESUME", int(button), false)
	upgrades := pauseMenu.newButton(func() { gc.setstate(upgrade) }, rl.Vector2{X: 380, Y: 390}, 260, 70, "UPGRADES", int(button), false)
	mainMenu := pauseMenu.newButton(func() { gc.setstate(menu) }, rl.Vector2{X: 380, Y: 480}, 260, 70, "MENU", int(button), false)
	saveButton := pauseMenu.newButton(func() { gc.SaveGame(gc.currentSave.PlayerName) }, rl.Vector2{X: 380, Y: 570}, 260, 70, "SAVE", int(button), false)
	confirmationText := pauseMenu.newTextBox(rl.Vector2{X: 350, Y: 560}, "", 40)
	var yesButton *Button
	var noButton *Button
	var closeButton *Button
	closeButton = pauseMenu.newButton(func() {
		resume.Showbutton = false
		upgrades.Showbutton = false
		mainMenu.Showbutton = false
		saveButton.Showbutton = false
		closeButton.Showbutton = false
		confirmationText.contents = "ARE YOU SURE?"
		yesButton.Showbutton = true
		noButton.Showbutton = true
	}, rl.Vector2{X: 380, Y: 660}, 260, 70, "CLOSE", int(button), false)
	yesButton = pauseMenu.newButton(func() { gc.setstate(quit) }, rl.Vector2{X: 300, Y: 660}, 120, 70, "YES", int(button), false)
	noButton = pauseMenu.newButton(func() {
		resume.Showbutton = true
		upgrades.Showbutton = true
		mainMenu.Showbutton = true
		saveButton.Showbutton = true
		closeButton.Showbutton = true
		confirmationText.contents = ""
		yesButton.Showbutton = false
		noButton.Showbutton = false
	}, rl.Vector2{X: 580, Y: 660}, 120, 70, "NO", int(button), false)
	yesButton.Showbutton = false
	noButton.Showbutton = false
	pauseMenu.textbox = append(pauseMenu.textbox, title)
	pauseMenu.textbox = append(pauseMenu.textbox, confirmationText)
	pauseMenu.buttons = append(pauseMenu.buttons, resume, upgrades, mainMenu, saveButton, closeButton, yesButton, noButton)

	gc.menus[pause] = pauseMenu
}

func (gc *GameController) upgrademenu() {
	upgradeMenu := gc.newMenu(rl.Vector2Zero())

	title := upgradeMenu.newTextBox(rl.Vector2{X: 330, Y: 120}, "UPGRADES", 60)

	health := upgradeMenu.newButton(func() { gc.player.UpgradeHealth() }, rl.Vector2{X: 300, Y: 260}, 400, 60, "+ HEALTH", int(button), false)

	speed := upgradeMenu.newButton(func() { gc.player.UpgradeSpeed() }, rl.Vector2{X: 300, Y: 340}, 400, 60, "+ SPEED", int(button), false)

	gun := upgradeMenu.newButton(func() { gc.player.UpgradeGunDamage() }, rl.Vector2{X: 300, Y: 420}, 400, 60, "+ GUN DAMAGE", int(button), false)

	melee := upgradeMenu.newButton(func() { gc.player.UpgradeMeleeDamage() }, rl.Vector2{X: 300, Y: 500}, 400, 60, "+ MELEE DAMAGE", int(button), false)

	back := upgradeMenu.newButton(func() { gc.setstate(playing) }, rl.Vector2{X: 300, Y: 620}, 400, 60, "BACK", int(button), false)

	upgradeMenu.textbox = append(upgradeMenu.textbox, title)
	upgradeMenu.buttons = append(upgradeMenu.buttons, health, speed, gun, melee, back)

	gc.menus[upgrade] = upgradeMenu
}

func (gc *GameController) mainmenu() {
	homeMenu := gc.newMenu(rl.Vector2{X: 0, Y: 0})
	homeMenu.background = NewSpriteRenderer(gc.assets.Textures[menubackground], rl.White, rl.Vector2{X: 500, Y: 450}, .6, 0)
	start := homeMenu.newButton(func() { gc.gamestate = selectsave }, rl.Vector2{X: 400, Y: 300}, 200, 80 /*no sprite yet*/, "START", 0, false)
	quit := homeMenu.newButton(func() { gc.setstate(quit) }, rl.Vector2{X: 400, Y: 420}, 200, 80, "QUIT", 0, false)
	homeMenu.buttons = append(homeMenu.buttons, start, quit)
	gc.menus[menu] = homeMenu
}
