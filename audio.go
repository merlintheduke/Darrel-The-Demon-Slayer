package main

import rl "github.com/gen2brain/raylib-go/raylib"

type SoundEffect int

const (
	shootSound SoundEffect = iota
	meleeSound
	enemyHitSound
	playerHitSound
	roomClearSound
	gameOverSound
	buttonClickSound
)

func LoadSounds() []rl.Sound {
	sound0 := rl.LoadSound("Assets/sounds/shot.wav")
	sound1 := rl.LoadSound("Assets/sounds/whoosh.wav")

	// SET VOLUME HERE
	rl.SetSoundVolume(sound0, 2) // quieter gun
	rl.SetSoundVolume(sound1, 2)
	return []rl.Sound{sound0, sound1}
}
func PlaySoundEffect(sounds []rl.Sound, sound int) {
	if sound < 0 || sound >= len(sounds) {
		return
	}

	rl.PlaySound(sounds[sound])
}

func UnloadSounds(sounds []rl.Sound) {
	for i := range sounds {
		rl.UnloadSound(sounds[i])
	}
}
