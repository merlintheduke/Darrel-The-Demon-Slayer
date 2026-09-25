package main

import rl "github.com/gen2brain/raylib-go/raylib"

type Assets struct {
	Textures []rl.Texture2D
	Sounds   []rl.Sound
}

func LoadAssets() Assets {
	textures := []rl.Texture2D{
		rl.LoadTexture("Assets/imgs/buttonbox.png"),
		rl.LoadTexture("Assets/imgs/smallXbutton.png"),
		rl.LoadTexture("Assets/imgs/CowboySpriteSheet.png"),
		rl.LoadTexture("Assets/imgs/tilesart.jpg"),
		rl.LoadTexture("Assets/imgs/demonV1F1.png"),
		rl.LoadTexture("Assets/imgs/bullet.png"),
		rl.LoadTexture("Assets/imgs/melee.png"),
		rl.LoadTexture("Assets/imgs/Darrell_the_Demon_Slayer.png"),
	}

	return Assets{Textures: textures, Sounds: LoadSounds()}
}

func (a Assets) Unload() {
	for i := range a.Textures {
		rl.UnloadTexture(a.Textures[i])
	}
	UnloadSounds(a.Sounds)
}
