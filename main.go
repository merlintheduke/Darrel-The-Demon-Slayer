package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Texture int

const (
	button Texture = iota
	smallbutton
	cowboy
	tile
	littledemon
	bullet
	melee
	menubackground
)

func main() {
	//Screensize and window setup
	rl.SetConfigFlags(rl.FlagWindowUndecorated)

	rl.InitWindow(gameWindowWidth, gameWindowHeight, "GAME")

	defer rl.CloseWindow()
	rl.SetTargetFPS(60)
	//screenwidth := rl.GetScreenWidth()
	//screenheight := rl.GetScreenHeight()

	rl.InitAudioDevice()
	defer rl.CloseAudioDevice()

	music := rl.LoadMusicStream("Assets/sounds/backgroundsong.wav")
	rl.SetMasterVolume(0.5)
	music.Looping = true
	rl.PlayMusicStream(music)
	defer rl.UnloadMusicStream(music)

	BuildMode := false
	assets := LoadAssets()
	defer assets.Unload()
	LoadPlayerAnimatorTemplate()

	gc := NewGameController()
	gc.assets = assets
	gc.player = newplayer(assets.Textures[cowboy], "")
	room := newRoom(defaultTileSize, assets.Textures[tile])
	room.BuildRoom()
	gc.rooms = append(gc.rooms, room)

	gc.Initialize()

	//gc.rooms[0].GenerateDungeon(20)

	for !rl.WindowShouldClose() && gc.gamestate != quit {
		rl.BeginDrawing()
		rl.ClearBackground(rl.Orange)
		gc.Update()
		gc.Draw()
		rl.UpdateMusicStream(music)
		camera := gc.camera

		if rl.IsKeyPressed(rl.KeyP) && BuildMode == false {
			BuildMode = true
		} else if rl.IsKeyPressed(rl.KeyP) && BuildMode == true {
			BuildMode = false
		}
		if !BuildMode {
			camera.Zoom = 1
		} else {
			room.RoomBuilderView()
			camera.Target = rl.Vector2{X: 0, Y: 0}
			camera.Zoom = .1
		}
		gc.camera = camera
		/*rl.BeginMode2D(camera)
		room.Draw()
		//enemy.Draw()

		gc.player.Draw()

		rl.EndMode2D()
		rl.DrawText("Hello, World!", 10, 10, 30, rl.RayWhite)
		*/
		//player.collisionBox = rl.NewRectangle(player.Position.X-20*player.Scale, player.Position.Y-42*player.Scale, 40*player.Scale, 90*player.Scale)
		rl.EndDrawing()
	}
}
