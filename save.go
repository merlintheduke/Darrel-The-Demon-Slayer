package main

import (
	"fmt"
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
		Angle:         p.Renderer.Angle,
		Scale:         p.Renderer.Scale,
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
			Renderer:      EntityRenderer{SpriteRenderer: SpriteRenderer{Sprite: rl.Texture2D{}, Color: rl.White, Position: sp.Position, Angle: sp.Angle, Scale: sp.Scale}, HealthBar: HealthBar{Width: 80, Height: 10}},
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
		p.Position = rl.Vector2{X: playerSpawnX, Y: playerSpawnY}
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

	return gc.saveRepository.Save(savegame)
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

	err := gc.saveRepository.Delete(playerName)
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
	saveMetadata, err := gc.saveRepository.List()
	if err != nil {
		return nil
	}
	return saveMetadata
}

func (gc *GameController) NewSaveGame(filename string) error {
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

	if err := gc.saveRepository.Save(savegame); err != nil {
		return err
	}

	gc.currentSave = &SaveSlot{
		PlayerName: filename,
		Index:      0,
		Exists:     true,
	}
	err := gc.LoadGame(filename)
	if err != nil {
		fmt.Println("New save load error:", err)
		return err
	}

	gc.StartGame()

	return nil
}

func (gc *GameController) LoadGame(playerName string) error {
	save, err := gc.saveRepository.Load(playerName)
	if err != nil {
		fmt.Println("LoadGame error:", err)
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
		room := newRoom(defaultTileSize, gc.assets.Textures[tile])
		room.BuildRoom()
		gc.rooms = append(gc.rooms, room)
	}

	gc.camera.Target = gc.player.pos
	gc.camera.Offset = rl.Vector2{X: cameraCenterX, Y: cameraCenterY}
	gc.camera.Zoom = 1

	gc.setstate(playing)

	return nil
}

func (gc *GameController) FixPlayerAfterLoad() {
	gc.player.Renderer.Sprite = gc.assets.Textures[cowboy]
	gc.player.Renderer.Color = rl.White
	gc.player.Renderer.Scale = 1
	gc.player.Renderer.Animator = newPlayerAnimator(gc.assets.Textures[cowboy])
	gc.player.Renderer.SpriteAnimator = newPlayerSpriteAnimator()

	if gc.player.Position.X != 0 || gc.player.Position.Y != 0 {
		gc.player.pos = gc.player.Position
	} else {
		gc.player.pos = rl.Vector2{X: playerSpawnX, Y: playerSpawnY}
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
