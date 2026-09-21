package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func TestSaveGameDoesNotPersistRuntimeTextureData(t *testing.T) {
	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	workingDir := t.TempDir()
	if err := os.Chdir(workingDir); err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = os.Chdir(oldWD)
	}()

	gc := NewGameController()
	gc.player = newplayer(rl.Texture2D{ID: 42}, "Alice")
	gc.player.Level = 3
	gc.player.XpToNext = 200
	gc.player.GunDamage = 15
	gc.player.MeleeDamage = 35
	gc.player.Position = rl.Vector2{X: 123, Y: 456}
	gc.player.pos = gc.player.Position
	gc.currentSave = &SaveSlot{PlayerName: "Alice"}

	if err := gc.SaveGame("Alice"); err != nil {
		t.Fatalf("SaveGame returned error: %v", err)
	}

	data, err := os.ReadFile(filepath.Join("saves", "Alice.json"))
	if err != nil {
		t.Fatalf("Read save file: %v", err)
	}

	jsonText := string(data)
	if strings.Contains(jsonText, "\"Sprite\"") {
		t.Fatal("save file contains runtime Sprite data")
	}
	if strings.Contains(jsonText, "\"collisionBox\"") {
		t.Fatal("save file contains runtime collision data")
	}
	if !strings.Contains(jsonText, "\"PlayerName\": \"Alice\"") {
		t.Fatal("save file did not preserve player metadata")
	}
}
