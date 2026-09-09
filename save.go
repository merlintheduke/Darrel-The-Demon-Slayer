package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)
var textBuffer = ""
type SaveGame struct {
	Meta        SaveMeta
	Player      Player
	CurrentRoom int
	GameState   GameState
}
type SaveMeta struct {
	PlayerName  string
	Lvl int
	CreatedAt string
	LastPlayed int64
	Index int
}
type SaveSlot struct {
	Index      int
	PlayerName string
	Exists     bool
}
func (gc *GameController) SaveGame(filename string) error {
	if filename == "" {
		if gc.currentSave == nil || gc.currentSave.PlayerName == "" {
			return fmt.Errorf("no current save selected")
		}

		filename = gc.currentSave.PlayerName
	}

	err := os.MkdirAll("saves", 0755)
	if err != nil {
		return err
	}

	file := filepath.Join("saves", filename+".json")

	savegame := SaveGame{
		Meta: SaveMeta{
			PlayerName:  filename,
			Lvl:         gc.player.Level,
			CreatedAt:   time.Now().Format("2006-01-02"),
			LastPlayed:  time.Now().Unix(),
			Index:       0,
		},
		Player:      gc.player,
		CurrentRoom: gc.currentRoom,
		GameState:   gc.gamestate,
	}

	data, err := json.MarshalIndent(savegame, "", "  ")
	if err != nil {
		fmt.Println(err)
		return err
	}

	return os.WriteFile(file, data, 0644)
}
func (gc *GameController) SaveCurrentGame() {
	if gc.currentSave == nil || gc.currentSave.PlayerName == "" {
		fmt.Println("No current save selected")
		return
	}

	err := gc.SaveGame(gc.currentSave.PlayerName)
	if err != nil {
		fmt.Println("Save error:", err)
	}
}
func (gc *GameController) DeleteSave(playerName string) {
	if playerName == "" {
		return
	}

	filename := filepath.Join("saves", playerName+".json")

	err := os.Remove(filename)
	if err != nil {
		fmt.Println("Delete save error:", err)
	}

	if gc.currentSave != nil && gc.currentSave.PlayerName == playerName {
		gc.currentSave = nil
	}

	gc.RefreshSaveMenu()
}

func (gc *GameController) DeleteCurrentSave() {
	if gc.currentSave == nil || gc.currentSave.PlayerName == "" {
		return
	}

	gc.DeleteSave(gc.currentSave.PlayerName)
}
func (gc *GameController) LoadSaves() []SaveMeta {
    files, _ := os.ReadDir("saves")
	var savemeta []SaveMeta
	for i, file := range files {
		data, _ := os.ReadFile("saves/" + file.Name())
		var save SaveGame
		json.Unmarshal(data, &save)
		savemeta = append(savemeta, SaveMeta{
			PlayerName: save.Meta.PlayerName,
			Lvl: save.Meta.Lvl,
			CreatedAt:  save.Meta.CreatedAt,
			LastPlayed: save.Meta.LastPlayed,
			Index:      i,})
		}
		return savemeta
}

func (gc *GameController) NewSaveGame(filename string) error {
	filesname := filename + ".json"
	folder := "saves"
	filepath := filepath.Join(folder, filesname)

	savemeta := SaveMeta{
		PlayerName: filename,
		Lvl:        0,
		CreatedAt:  time.Now().Format("2006-01-02"),
	}

	savegame := SaveGame{
		Meta:        savemeta,
		Player:      gc.player,
		CurrentRoom: 0,
		GameState:   playing,
	}

	data, err := json.MarshalIndent(savegame, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath, data, 0644)
	if err != nil {
		return err
	}

	gc.currentSave = &SaveSlot{
		PlayerName: filename,
		Index:      0,
		Exists:     true,
	}
	err = gc.LoadGame(filename)
	if err != nil {
	fmt.Println("New save load error:", err)
	return err
	}

	gc.StartGame()

	return nil
}
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

		deleteButton := savesMenu.newButton(func() {gc.DeleteSave(name)}, rl.Vector2{X: x + 310, Y: y}, 40, 40, "", int(smallbutton), true)

		savesMenu.buttons = append(savesMenu.buttons, loadButton, deleteButton)
	}
	newsave := savesMenu.newButton(func() { gc.newsavebutton() },rl.Vector2{X: 50, Y: 800},100,100,"NEW",int(smallbutton),false,)
	backbutton := savesMenu.newButton(func() { gc.setstate(menu) },rl.Vector2{X: 50, Y: 20},100,100,"Back",int(smallbutton),false,)
	savesMenu.buttons = append(savesMenu.buttons, backbutton, newsave)

	return savesMenu
}

func (gc *GameController) RefreshSaveMenu() {
	if int(selectsave) < len(gc.menu) {
		gc.menu[selectsave] = gc.buildSaveMenu(gc.LoadSaves())
	}
}
func (gc *GameController) newsavebutton() {
        //currentsave := 1
        b1 := gc.menu[selectsave].newButton(func() {}, rl.Vector2{X: 300, Y: 275 }, 300, 50, "", int(button),false)
		b1.OnClick = func() {b1.allowUserInput = true}
		gc.menu[selectsave].buttons = append(gc.menu[selectsave].buttons, b1)
		
}
func (gc *GameController) LoadGame(playerName string) error {
	filename := filepath.Join("saves", playerName+".json")

	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("LoadGame read error:", err)
		return err
	}

	var save SaveGame
	err = json.Unmarshal(data, &save)
	if err != nil {
		fmt.Println("LoadGame json error:", err)
		return err
	}

	gc.player = save.Player
	gc.FixPlayerAfterLoad()
	gc.currentRoom = save.CurrentRoom
	gc.currentSave = &SaveSlot{
		PlayerName: playerName,
		Index:      0,
		Exists:     true,
	}
	// IMPORTANT:
	// Do not trust saved Room yet because Room fields are lowercase
	// and textures cannot be restored from JSON.
	if len(gc.rooms) == 0 {
		room := newRoom(100, textures[tile])
		room.BuildRoom()
		gc.rooms = append(gc.rooms, room)
	}

	gc.camera.Target = gc.player.pos
	gc.camera.Offset = rl.Vector2{X: 500, Y: 500}
	gc.camera.Zoom = 1

	gc.setstate(playing)

	return nil
}
func (gc *GameController) FixPlayerAfterLoad() {
	// restore texture/runtime sprite data
	gc.player.Sprite = textures[cowboy]
	gc.player.Color = rl.White
	gc.player.Scale = 1
	gc.player.frameSize = 120
	gc.player.totalFrames = 7
	gc.player.currentFrame = 1

	// JSON saved Position, but not pos.
	// So copy visible sprite position back into physics position.
	if gc.player.Position.X != 0 || gc.player.Position.Y != 0 {
		gc.player.pos = gc.player.Position
	} else {
		gc.player.pos = rl.Vector2{X: 250, Y: 250}
		gc.player.Position = gc.player.pos
	}

	// restore collision/runtime values
	gc.player.velocity = rl.Vector2Zero()
	gc.player.collisionBox = rl.NewRectangle(
		gc.player.pos.X-25,
		gc.player.pos.Y-50,
		50,
		100,
	)

	if gc.player.GunDamage <= 0 {
	gc.player.GunDamage = 10
	}
	if gc.player.MeleeDamage <= 0 {
		gc.player.MeleeDamage = 25
	}
	if gc.player.XpToNext <= 0 {
		gc.player.XpToNext = 100
	}
	if gc.player.Level <= 0 {
		gc.player.Level = 1
	}
	gc.player.Alive = true

	gc.player.SyncCollisionBox()
}