package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type GameController struct {
	Sprites       []rl.Texture2D
	player        Player
	rooms         []Room
	currentRoom   int
	enemies       []Enemy
	attackManager AttackManager
	menu          []Menu
	gamestate     GameState
	camera        rl.Camera2D
	saveslot      []SaveSlot
	currentSave   *SaveSlot
}
type GameState int

var shot bool = false

const (
	menu GameState = iota
	selectsave
	playing
	home
	gameover
	pause
	upgrade
	buildmode
	quit
)

func NewGameController() GameController {
	rooms := make([]Room, 0)
	enemy := make([]Enemy, 0)
	AllMenu := make([]Menu, 0)
	camera := rl.NewCamera2D(
		rl.Vector2{X: 0, Y: 0}, // offset
		rl.Vector2{X: 0, Y: 0}, // target
		0,                      // rotation
		1,                      // zoom
	)
	GameController := GameController{
		rooms:     rooms,
		enemies:   enemy,
		menu:      AllMenu,
		gamestate: menu,
		camera:    camera,
	}
	return GameController
}
func (gc *GameController) Initialize() {
	gc.setstate(menu)

	gc.mainmenu()                                        // index 0
	gc.savemenu(gc.LoadSaves())                          // index 1
	gc.PlayingUI()                                       // index 2
	gc.menu = append(gc.menu, newMenu(rl.Vector2Zero())) // home, index 3
	gc.gameovermenu()                                    // gameover, index 4
	gc.pausemenu()                                       // pause, index 5
	gc.upgrademenu()                                     // index 6
	gc.menu = append(gc.menu, newMenu(rl.Vector2Zero())) // buildmode, index 7
	gc.menu = append(gc.menu, newMenu(rl.Vector2Zero())) // quit, index 8

	gc.attackManager = AttackManager{}
}
func (gc *GameController) GameState() {

}
func (gc *GameController) Update() {
	if gc.gamestate == playing && rl.IsKeyPressed(rl.KeyEscape) {
		gc.setstate(pause)
		return
	}

	if gc.gamestate == pause && rl.IsKeyPressed(rl.KeyEscape) {
		gc.setstate(playing)
		return
	}

	if int(gc.gamestate) < len(gc.menu) {
		gc.menu[gc.gamestate].Update(gc)
	}

	if gc.gamestate == playing {
		gc.camera.Target = gc.player.pos
		gc.camera.Offset = rl.Vector2{X: 500, Y: 500}

		gc.player.velocity = rl.Vector2{X: 0, Y: 0}

		if rl.IsKeyDown(rl.KeyD) {
			gc.player.velocity.X += gc.player.Speed
		}
		if rl.IsKeyDown(rl.KeyA) {
			gc.player.velocity.X -= gc.player.Speed
		}
		if rl.IsKeyDown(rl.KeyW) {
			gc.player.velocity.Y -= gc.player.Speed
		}
		if rl.IsKeyDown(rl.KeyS) {
			gc.player.velocity.Y += gc.player.Speed
		}

		mouse := rl.GetMousePosition()
		mouseWorld := rl.GetScreenToWorld2D(mouse, gc.camera)
		dir := rl.Vector2Subtract(mouseWorld, gc.player.pos)

		if rl.IsKeyPressed(rl.KeySpace) {
			gc.attackManager.SpawnProjectile(gc.player.pos, dir, &gc.player, gc.player.GunDamage, textures[bullet], 0)

		}

		if rl.IsKeyPressed(rl.KeyQ) {
			gc.attackManager.SpawnMelee(gc.player.pos, dir, &gc.player, gc.player.MeleeDamage, 120, textures[melee], 1)
		}

		gc.player.Update(&gc.rooms[gc.currentRoom])
		gc.CheckEncounterTriggers()
		gc.UpdateEncounters()
			if gc.CurrentRoomCompleted() {
				gc.GoToNextRoom()
				return
			}
		gc.attackManager.Update(rl.GetFrameTime())
		gc.CheckEncounterHits()
		if gc.player.CurrentHealth <= 0 || !gc.player.Alive {
			gc.TriggerGameOver()
		}
	}
}
func (gc *GameController) Draw() {
	drawWorld := gc.gamestate == playing || gc.gamestate == pause || gc.gamestate == upgrade

	if drawWorld {
		rl.BeginMode2D(gc.camera)
		gc.rooms[gc.currentRoom].Draw()
		gc.player.Draw()
		gc.player.DrawHealthBar()
		gc.attackManager.Draw()
		gc.DrawEncounters()
		rl.EndMode2D()
	}

	if gc.gamestate == pause {
		rl.DrawRectangle(0, 0, 1000, 1000, rl.Color{R: 0, G: 0, B: 0, A: 120})
	}

	if gc.gamestate == upgrade {
		rl.DrawRectangle(0, 0, 1000, 1000, rl.Color{R: 0, G: 0, B: 0, A: 160})
		gc.DrawUpgradeStats()
	}

	if int(gc.gamestate) < len(gc.menu) {
		gc.menu[gc.gamestate].Draw()
	}
}
func (gc *GameController) PlayingUI() {
	playingMenu := newMenu(rl.Vector2Zero())
	playingMenu.CustomDraw = func() {
		gc.DrawPlayerGUI()
	}
	upgradeButton := playingMenu.newButton(func() { gc.setstate(upgrade) }, rl.Vector2{X: 720, Y: 20}, 250, 50, "UPGRADE", int(smallbutton), false)
	pauseButton := playingMenu.newButton(func() { gc.setstate(pause) }, rl.Vector2{X: 720, Y: 80}, 250, 50, "PAUSE", int(smallbutton), false)
	playingMenu.buttons = append(playingMenu.buttons, upgradeButton, pauseButton)
	gc.menu = append(gc.menu, playingMenu)
}
func (gc *GameController) gameovermenu() {
	gameOverMenu := newMenu(rl.Vector2Zero())

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

	gc.menu = append(gc.menu, gameOverMenu)
}
func (gc *GameController) TriggerGameOver() {
	if gc.gamestate == gameover {
		return
	}

	gc.DeleteCurrentSave()
	gc.setstate(gameover)
}

func (gc *GameController) ResetAfterGameOver() {
	gc.player = newplayer(textures[cowboy], "")

	gc.rooms = make([]Room, 0)
	room := newRoom(100, textures[tile])
	room.BuildRoom()
	gc.rooms = append(gc.rooms, room)

	gc.currentRoom = 0
	gc.attackManager = AttackManager{}
	gc.currentSave = nil

	gc.RefreshSaveMenu()
}
func (gc *GameController) pausemenu() {
	pauseMenu := newMenu(rl.Vector2Zero())

	title := pauseMenu.newTextBox(rl.Vector2{X: 410, Y: 180}, "PAUSED", 60)
	resume := pauseMenu.newButton(func() { gc.setstate(playing) }, rl.Vector2{X: 380, Y: 300}, 260, 70, "RESUME", int(button), false)
	upgrades := pauseMenu.newButton(func() { gc.setstate(upgrade) }, rl.Vector2{X: 380, Y: 390}, 260, 70, "UPGRADES", int(button), false)
	mainMenu := pauseMenu.newButton(func() { gc.setstate(menu) }, rl.Vector2{X: 380, Y: 480}, 260, 70, "MENU", int(button), false)
	quitButton := pauseMenu.newButton(func() { gc.setstate(quit) }, rl.Vector2{X: 380, Y: 570}, 260, 70, "QUIT", int(button), false)
	saveButton := pauseMenu.newButton(func() { gc.SaveGame(gc.currentSave.PlayerName) }, rl.Vector2{X: 380, Y: 570}, 260, 70, "SAVE", int(button), false)
	pauseMenu.textbox = append(pauseMenu.textbox, title)
	pauseMenu.buttons = append(pauseMenu.buttons, resume, upgrades, mainMenu, quitButton, saveButton)

	gc.menu = append(gc.menu, pauseMenu)
}

func (gc *GameController) upgrademenu() {
	upgradeMenu := newMenu(rl.Vector2Zero())

	title := upgradeMenu.newTextBox(rl.Vector2{X: 330, Y: 120}, "UPGRADES", 60)

	health := upgradeMenu.newButton(func() { gc.player.UpgradeHealth() }, rl.Vector2{X: 300, Y: 260}, 400, 60, "+ HEALTH", int(button), false)

	speed := upgradeMenu.newButton(func() { gc.player.UpgradeSpeed() }, rl.Vector2{X: 300, Y: 340}, 400, 60, "+ SPEED", int(button), false)

	gun := upgradeMenu.newButton(func() { gc.player.UpgradeGunDamage() }, rl.Vector2{X: 300, Y: 420}, 400, 60, "+ GUN DAMAGE", int(button), false)

	melee := upgradeMenu.newButton(func() { gc.player.UpgradeMeleeDamage() }, rl.Vector2{X: 300, Y: 500}, 400, 60, "+ MELEE DAMAGE", int(button), false)

	back := upgradeMenu.newButton(func() { gc.setstate(playing) }, rl.Vector2{X: 300, Y: 620}, 400, 60, "BACK", int(button), false)

	upgradeMenu.textbox = append(upgradeMenu.textbox, title)
	upgradeMenu.buttons = append(upgradeMenu.buttons, health, speed, gun, melee, back)

	gc.menu = append(gc.menu, upgradeMenu)
}

func (gc *GameController) StartGame() {
	gc.currentRoom = 0

	if len(gc.rooms) == 0 {
		room := newRoom(100, textures[tile])
		room.BuildRoom()
		gc.rooms = append(gc.rooms, room)
	}

	gc.rooms[gc.currentRoom].SetupEncounters(&gc.player, gc.player.Difficulty)

	gc.MovePlayerToRoomSpawn()

	gc.camera.Target = gc.player.pos
	gc.camera.Offset = rl.Vector2{X: 500, Y: 500}
	gc.camera.Zoom = 1

	gc.setstate(playing)
}

func (gc *GameController) mainmenu() {
	homeMenu := newMenu(rl.Vector2{X: 0, Y: 0})
	homeMenu.background = NewSpriteRenderer(7, rl.White, rl.Vector2{X: 500, Y: 450}, .6, 0)
	start := homeMenu.newButton(func() { gc.gamestate = selectsave }, rl.Vector2{X: 400, Y: 300}, 200, 80 /*no sprite yet*/, "START", 0, false)
	quit := homeMenu.newButton(func() { gc.setstate(quit) }, rl.Vector2{X: 400, Y: 420}, 200, 80, "QUIT", 0, false)
	homeMenu.buttons = append(homeMenu.buttons, start, quit)
	gc.menu = append(gc.menu, homeMenu)
}
func (gc *GameController) setstate(state GameState) {
	gc.gamestate = state
}
func (gc *GameController) UnloadTextures() {
	for i := range gc.Sprites {
		rl.UnloadTexture(gc.Sprites[i])
	}
}
func (gc *GameController) MovePlayerToRoomSpawn() {
	if len(gc.rooms) == 0 {
		return
	}

	spawn := gc.rooms[gc.currentRoom].roomspawn

	if spawn.X == 0 && spawn.Y == 0 {
		return
	}

	gc.player.pos = spawn
	gc.player.Position = spawn
	gc.player.SyncCollisionBox()

	gc.camera.Target = gc.player.pos
}
func (gc *GameController) CheckEncounterTriggers() {
	if len(gc.rooms) == 0 {
		return
	}

	room := &gc.rooms[gc.currentRoom]

	for i := range room.Roomdescriptor {
		// Skip start room.
		if i == 0 {
			continue
		}

		roomdesc := &room.Roomdescriptor[i]

		if roomdesc.Encounter == nil {
			continue
		}

		if roomdesc.Encounter.started || roomdesc.Encounter.ended {
			continue
		}

		// This triggers after the player has passed through the door
		// and entered the actual room area.
		if room.playerPastDoor(gc.player.pos, roomdesc) {
			room.blockDoor(roomdesc)
			roomdesc.Encounter.startEncounter(room, roomdesc)
		}
	}
}

func (gc *GameController) UpdateEncounters() {
	if len(gc.rooms) == 0 {
		return
	}

	room := &gc.rooms[gc.currentRoom]

	for i := range room.Roomdescriptor {
		roomdesc := &room.Roomdescriptor[i]

		if roomdesc.Encounter == nil {
			continue
		}

		wasEnded := roomdesc.Encounter.ended

		roomdesc.Encounter.updateEncounter(room, &gc.attackManager)

		if !wasEnded && roomdesc.Encounter.ended {
			room.completeDoor(roomdesc)
		}
	}
}

func (gc *GameController) DrawEncounters() {
	if len(gc.rooms) == 0 {
		return
	}

	room := &gc.rooms[gc.currentRoom]

	for i := range room.Roomdescriptor {
		roomdesc := &room.Roomdescriptor[i]

		if roomdesc.Encounter != nil {
			roomdesc.Encounter.drawEncounter()
		}
	}
}

func (gc *GameController) CheckEncounterHits() {
	if len(gc.rooms) == 0 {
		return
	}

	room := &gc.rooms[gc.currentRoom]

	for i := range room.Roomdescriptor {
		roomdesc := &room.Roomdescriptor[i]

		if roomdesc.Encounter == nil {
			continue
		}
		if !roomdesc.Encounter.started || roomdesc.Encounter.ended {
			continue
		}
		gc.attackManager.CheckHits(roomdesc.Encounter.enemies, &gc.player, room)
	}
}
func (gc *GameController) DrawPlayerGUI() {
	healthPercent := float32(0)

	if gc.player.MaxHealth > 0 {
		healthPercent = float32(gc.player.CurrentHealth / gc.player.MaxHealth)
	}

	if healthPercent < 0 {
		healthPercent = 0
	}
	if healthPercent > 1 {
		healthPercent = 1
	}

	rl.DrawText("HP", 25, 20, 24, rl.White)
	rl.DrawRectangle(70, 22, 250, 22, rl.DarkGray)
	rl.DrawRectangle(70, 22, int32(250*healthPercent), 22, rl.Red)
	rl.DrawRectangleLines(70, 22, 250, 22, rl.Black)

	xpToNext := gc.player.XpToNext
	if xpToNext <= 0 {
		xpToNext = 100
	}

	xpPercent := float32(gc.player.Xp) / float32(xpToNext)

	if xpPercent < 0 {
		xpPercent = 0
	}
	if xpPercent > 1 {
		xpPercent = 1
	}

	rl.DrawText("XP", 25, 55, 24, rl.White)
	rl.DrawRectangle(70, 57, 250, 18, rl.DarkGray)
	rl.DrawRectangle(70, 57, int32(250*xpPercent), 18, rl.Blue)
	rl.DrawRectangleLines(70, 57, 250, 18, rl.Black)
	rl.DrawText(fmt.Sprintf("Level: %d", gc.player.Level), 25, 90, 24, rl.White)
	rl.DrawText(fmt.Sprintf("Upgrade Points: %d", gc.player.UpgradePoints), 25, 120, 22, rl.White)
}

func (gc *GameController) DrawUpgradeStats() {
	rl.DrawText(fmt.Sprintf("Upgrade Points: %d", gc.player.UpgradePoints), 300, 200, 30, rl.White)
	rl.DrawText(fmt.Sprintf("Health: %.0f / %.0f", gc.player.CurrentHealth, gc.player.MaxHealth), 730, 275, 22, rl.White)
	rl.DrawText(fmt.Sprintf("Speed: %.0f", gc.player.Speed), 730, 355, 22, rl.White)
	rl.DrawText(fmt.Sprintf("Gun Damage: %d", gc.player.GunDamage), 730, 435, 22, rl.White)
	rl.DrawText(fmt.Sprintf("Melee Damage: %d", gc.player.MeleeDamage), 730, 515, 22, rl.White)
}
func (gc *GameController) CurrentRoomCompleted() bool {
	if len(gc.rooms) == 0 {
		return false
	}

	room := &gc.rooms[gc.currentRoom]

	// If there are no real encounter rooms, don't advance.
	hasEncounter := false

	for i := range room.Roomdescriptor {
		// skip start room
		if i == 0 {
			continue
		}

		enc := room.Roomdescriptor[i].Encounter
		if enc == nil {
			continue
		}

		hasEncounter = true

		if !enc.ended {
			return false
		}
	}

	return hasEncounter
}

func (gc *GameController) GoToNextRoom() {
	// clear bullets/melee attacks from old room
	gc.attackManager = AttackManager{}

	// make a brand new room
	newRoom := newRoom(100, textures[tile])
	newRoom.BuildRoom()
	newRoom.SetupEncounters(&gc.player, gc.player.Difficulty)

	gc.rooms = append(gc.rooms, newRoom)

	// move to the new room
	gc.currentRoom++

	// optional difficulty increase
	gc.player.Difficulty++

	gc.MovePlayerToRoomSpawn()

	gc.camera.Target = gc.player.pos
	gc.camera.Offset = rl.Vector2{X: 500, Y: 500}
	gc.camera.Zoom = 1
}