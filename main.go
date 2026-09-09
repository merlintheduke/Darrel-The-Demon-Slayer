package main

import (
	

	rl "github.com/gen2brain/raylib-go/raylib"
)

var textures []rl.Texture2D

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
	
	rl.InitWindow(1000, 1000, "GAME")

	defer rl.CloseWindow()
	rl.SetTargetFPS(60)
	//screenwidth := rl.GetScreenWidth()
	//screenheight := rl.GetScreenHeight()

	rl.InitAudioDevice()
	defer rl.CloseAudioDevice()
	LoadSounds()
	defer UnloadSounds()

	music := rl.LoadMusicStream("Assets/sounds/backgroundsong.wav")
	rl.SetMasterVolume(0.5) 
	music.Looping = true
	rl.PlayMusicStream(music)
	defer rl.UnloadMusicStream(music)

	BuildMode := false
	texture0 := rl.LoadTexture("Assets/imgs/buttonbox.png")
	texture1 := rl.LoadTexture("Assets/imgs/smallXbutton.png")
	texture2 := rl.LoadTexture("Assets/imgs/CowboySpriteSheet.png")
	texture3 := rl.LoadTexture("Assets/imgs/tilesart.jpg")
	texture4 := rl.LoadTexture("Assets/imgs/demonV1F1.png")
	texture5 := rl.LoadTexture("Assets/imgs/bullet.png")
	texture6 := rl.LoadTexture("Assets/imgs/melee.png")
	texture7 := rl.LoadTexture("Assets/imgs/Darrell_the_Demon_Slayer.png")
	textures = append(textures, texture0, texture1, texture2, texture3, texture4, texture5, texture6, texture7)
	defer UnloadTextures()

	gc := NewGameController()
	gc.player = newplayer(texture2, "")
	room := newRoom(100, textures[tile])
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
		if rl.IsKeyPressed(rl.KeySpace) {
			player.nextFrame()
		}
		*/
		//player.collisionBox = rl.NewRectangle(player.Position.X-20*player.Scale, player.Position.Y-42*player.Scale, 40*player.Scale, 90*player.Scale)
		rl.EndDrawing()
	}
}



func UnloadTextures() {
	for i := range textures {
		rl.UnloadTexture(textures[i])
	}
}
