package main

import rl "github.com/gen2brain/raylib-go/raylib"

var sounds []rl.Sound

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

func LoadSounds() {
	sound0 := rl.LoadSound("Assets/sounds/shot.wav")
	sound1 := rl.LoadSound("Assets/sounds/whoosh.wav")
	
	// SET VOLUME HERE
	rl.SetSoundVolume(sound0, 2) // quieter gun
	rl.SetSoundVolume(sound1, 2)
	sounds = append(sounds, sound0, sound1)
}
func PlaySoundEffect(sound int) {
	if sound < 0 || sound >= len(sounds) {
		return
	}

	rl.PlaySound(sounds[sound])
}

func UnloadSounds() {
	for i := range sounds {
		rl.UnloadSound(sounds[i])
	}
}
