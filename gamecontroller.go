package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type GameController struct {
	player         Player
	rooms          []Room
	currentRoom    int
	enemies        []Enemy
	attackManager  AttackManager
	assets         Assets
	menus          map[GameState]Menu
	saveRepository SaveRepository
	gamestate      GameState
	camera         rl.Camera2D
	saveslot       []SaveSlot
	currentSave    *SaveSlot
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
	camera := rl.NewCamera2D(
		rl.Vector2{X: 0, Y: 0}, // offset
		rl.Vector2{X: 0, Y: 0}, // target
		0,                      // rotation
		1,                      // zoom
	)
	GameController := GameController{
		rooms:          rooms,
		enemies:        enemy,
		menus:          make(map[GameState]Menu),
		saveRepository: SaveRepository{Directory: "saves"},
		gamestate:      menu,
		camera:         camera,
	}
	return GameController
}
func (gc *GameController) Initialize() {
	gc.setstate(menu)

	gc.mainmenu()
	gc.savemenu(gc.LoadSaves())
	gc.PlayingUI()
	gc.menus[home] = gc.newMenu(rl.Vector2Zero())
	gc.gameovermenu()
	gc.pausemenu()
	gc.upgrademenu()
	gc.menus[buildmode] = gc.newMenu(rl.Vector2Zero())
	gc.menus[quit] = gc.newMenu(rl.Vector2Zero())

	gc.attackManager = AttackManager{Sounds: gc.assets.Sounds}
}
func (gc *GameController) GameState() {

}
func (gc *GameController) Update() {
	if gc.updateStateTransitions() {
		return
	}

	if currentMenu, ok := gc.menus[gc.gamestate]; ok {
		currentMenu.Update(gc)
	}

	if gc.gamestate == playing {
		gc.updatePlaying()
	}
}

func (gc *GameController) updateStateTransitions() bool {
	if gc.gamestate == playing && rl.IsKeyPressed(rl.KeyEscape) {
		gc.setstate(pause)
		return true
	}

	if gc.gamestate == pause && rl.IsKeyPressed(rl.KeyEscape) {
		gc.setstate(playing)
		return true
	}

	return false
}

func (gc *GameController) updatePlaying() {
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

	mouses := rl.GetMousePosition()
	mouseWorld := rl.GetScreenToWorld2D(mouses, gc.camera)
	dir := rl.Vector2Subtract(mouseWorld, gc.player.pos)

	if rl.IsKeyPressed(rl.KeySpace) {
		gc.attackManager.SpawnProjectile(gc.player.pos, dir, &gc.player, gc.player.GunDamage, gc.assets.Textures[bullet], 0)
	}

	if rl.IsKeyPressed(rl.KeyQ) {
		gc.attackManager.SpawnMelee(gc.player.pos, dir, &gc.player, gc.player.MeleeDamage, 120, gc.assets.Textures[melee], 1)
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

	if currentMenu, ok := gc.menus[gc.gamestate]; ok {
		currentMenu.Draw()
	}
}
func (gc *GameController) TriggerGameOver() {
	if gc.gamestate == gameover {
		return
	}

	gc.DeleteCurrentSave()
	gc.setstate(gameover)
}

func (gc *GameController) ResetAfterGameOver() {
	gc.player = newplayer(gc.assets.Textures[cowboy], "")

	gc.rooms = make([]Room, 0)
	room := newRoom(100, gc.assets.Textures[tile])
	room.BuildRoom()
	gc.rooms = append(gc.rooms, room)

	gc.currentRoom = 0
	gc.attackManager = AttackManager{Sounds: gc.assets.Sounds}
	gc.currentSave = nil

	gc.RefreshSaveMenu()
}
func (gc *GameController) StartGame() {
	gc.currentRoom = 0

	if len(gc.rooms) == 0 {
		room := newRoom(100, gc.assets.Textures[tile])
		room.BuildRoom()
		gc.rooms = append(gc.rooms, room)
	}

	gc.rooms[gc.currentRoom].SetupEncounters(&gc.player, gc.player.Difficulty, gc.assets.Textures)

	gc.MovePlayerToRoomSpawn()

	gc.camera.Target = gc.player.pos
	gc.camera.Offset = rl.Vector2{X: 500, Y: 500}
	gc.camera.Zoom = 1

	gc.setstate(playing)
}

func (gc *GameController) setstate(state GameState) {
	gc.gamestate = state
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
	gc.player.Renderer.Position = spawn
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
	gc.attackManager = AttackManager{Sounds: gc.assets.Sounds}

	// make a brand new room
	newRoom := newRoom(100, gc.assets.Textures[tile])
	newRoom.BuildRoom()
	newRoom.SetupEncounters(&gc.player, gc.player.Difficulty, gc.assets.Textures)

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
