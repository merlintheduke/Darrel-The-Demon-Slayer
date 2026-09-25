package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var textBuffer = ""

type SaveGame struct {
	Meta        SaveMeta
	Player      SavePlayer
	CurrentRoom int
	GameState   GameState
}

type SaveMeta struct {
	PlayerName string
	Lvl        int
	CreatedAt  string
	LastPlayed int64
	Index      int
}

type SavePlayer struct {
	Name          string
	Position      rl.Vector2
	Angle         float32
	Scale         float32
	Speed         float32
	MaxHealth     float64
	CurrentHealth float64
	Alive         bool
	Xp            int
	Level         int
	XpToNext      int
	UpgradePoints int
	GunDamage     int
	MeleeDamage   int
	Roomscleared  int
	Difficulty    int
}

type SaveSlot struct {
	Index      int
	PlayerName string
	Exists     bool
}

func newSavePlayer(p Player) SavePlayer {
	return SavePlayer{
		Name:          p.Name,
		Position:      p.Position,
		Angle:         p.Angle,
		Scale:         p.Scale,
		Speed:         p.Speed,
		MaxHealth:     p.MaxHealth,
		CurrentHealth: p.CurrentHealth,
		Alive:         p.Alive,
		Xp:            p.Xp,
		Level:         p.Level,
		XpToNext:      p.XpToNext,
		UpgradePoints: p.UpgradePoints,
		GunDamage:     p.GunDamage,
		MeleeDamage:   p.MeleeDamage,
		Roomscleared:  p.Roomscleared,
		Difficulty:    p.Difficulty,
	}
}

func (sp SavePlayer) toPlayer() Player {
	p := Player{
		Entity: Entity{
			Name: sp.Name,
			PhysicsBody: PhysicsBody{
				pos:      sp.Position,
				velocity: rl.Vector2Zero(),
			},
			SpriteRenderer: SpriteRenderer{
				Color:        rl.White,
				Position:     sp.Position,
				Angle:        sp.Angle,
				Scale:        sp.Scale,
				frameSize:    120,
				totalFrames:  7,
				currentFrame: 1,
			},
			HealthBar:     HealthBar{Width: 80, Height: 10},
			Speed:         sp.Speed,
			MaxHealth:     sp.MaxHealth,
			CurrentHealth: sp.CurrentHealth,
			Alive:         sp.Alive,
		},
		Xp:            sp.Xp,
		Level:         sp.Level,
		XpToNext:      sp.XpToNext,
		UpgradePoints: sp.UpgradePoints,
		GunDamage:     sp.GunDamage,
		MeleeDamage:   sp.MeleeDamage,
		Roomscleared:  sp.Roomscleared,
		Difficulty:    sp.Difficulty,
	}
	if p.Position == (rl.Vector2{}) {
		p.Position = rl.Vector2{X: 250, Y: 250}
	}
	p.pos = p.Position
	p.collisionBox = rl.NewRectangle(p.pos.X-25, p.pos.Y-50, 50, 100)
	p.Position = p.pos
	return p
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
	playerSnapshot := newSavePlayer(gc.player)
	if playerSnapshot.Name == "" {
		playerSnapshot.Name = filename
	}

	savegame := SaveGame{
		Meta: SaveMeta{
			PlayerName: filename,
			Lvl:        gc.player.Level,
			CreatedAt:  time.Now().Format("2006-01-02"),
			LastPlayed: time.Now().Unix(),
			Index:      0,
		},
		Player:      playerSnapshot,
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
	files, err := os.ReadDir("saves")
	if err != nil {
		return nil
	}

	var savemeta []SaveMeta
	for i, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".json") {
			continue
		}

		data, err := os.ReadFile(filepath.Join("saves", file.Name()))
		if err != nil {
			continue
		}

		var save SaveGame
		if err := json.Unmarshal(data, &save); err != nil {
			continue
		}
		if save.Meta.PlayerName == "" {
			save.Meta.PlayerName = save.Player.Name
		}
		if save.Meta.PlayerName == "" {
			continue
		}

		savemeta = append(savemeta, SaveMeta{
			PlayerName: save.Meta.PlayerName,
			Lvl:        save.Player.Level,
			CreatedAt:  save.Meta.CreatedAt,
			LastPlayed: save.Meta.LastPlayed,
			Index:      i,
		})
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
		Player:      newSavePlayer(gc.player),
		CurrentRoom: 0,
		GameState:   playing,
	}
	if savegame.Player.Name == "" {
		savegame.Player.Name = filename
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

	gc.player = save.Player.toPlayer()
	if gc.player.Name == "" {
		gc.player.Name = playerName
	}
	gc.FixPlayerAfterLoad()
	gc.currentRoom = save.CurrentRoom
	gc.currentSave = &SaveSlot{
		PlayerName: playerName,
		Index:      0,
		Exists:     true,
	}
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
	gc.player.Sprite = textures[cowboy]
	gc.player.Color = rl.White
	gc.player.Scale = 1
	gc.player.frameSize = 120
	gc.player.totalFrames = 7
	gc.player.currentFrame = 1

	if gc.player.Position.X != 0 || gc.player.Position.Y != 0 {
		gc.player.pos = gc.player.Position
	} else {
		gc.player.pos = rl.Vector2{X: 250, Y: 250}
		gc.player.Position = gc.player.pos
	}

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
