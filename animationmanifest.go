package main

import (
	"encoding/json"
	"fmt"
	"os"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	playerAnimationManifestPath = "Assets/imgs/player/animations.json"
	playerIdleSpritePath        = "Assets/imgs/player/idle.png"
)

type AnimationManifest struct {
	FormatVersion  int                            `json:"format_version"`
	FrameWidth     int32                          `json:"frame_width"`
	FrameHeight    int32                          `json:"frame_height"`
	DirectionOrder []string                       `json:"direction_order"`
	Animations     map[AnimationID]AnimationEntry `json:"animations"`
}

type AnimationEntry struct {
	Sheet       string  `json:"sheet"`
	Columns     int32   `json:"columns"`
	Rows        int32   `json:"rows"`
	Frames      int32   `json:"frames"`
	FPS         float32 `json:"fps"`
	Loop        bool    `json:"loop"`
	Directional bool    `json:"directional"`
}

// The template owns the texture handles. Individual player/enemy animators
// copy its clip map but maintain their own frame, timer, and facing state.
var playerAnimatorTemplate SpriteAnimator

func LoadSpriteAnimatorManifest(path string, initial AnimationID) (SpriteAnimator, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return SpriteAnimator{}, fmt.Errorf("read animation manifest %q: %w", path, err)
	}

	var manifest AnimationManifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return SpriteAnimator{}, fmt.Errorf("parse animation manifest %q: %w", path, err)
	}
	if manifest.FormatVersion != 1 {
		return SpriteAnimator{}, fmt.Errorf("unsupported animation manifest version %d", manifest.FormatVersion)
	}
	if manifest.FrameWidth <= 0 || manifest.FrameHeight <= 0 {
		return SpriteAnimator{}, fmt.Errorf("animation manifest %q has invalid frame dimensions", path)
	}

	clips := make(map[AnimationID]AnimationClip, len(manifest.Animations))
	loaded := make([]rl.Texture2D, 0, len(manifest.Animations))
	for id, entry := range manifest.Animations {
		if entry.Sheet == "" || entry.Columns <= 0 || entry.Rows <= 0 || entry.Frames <= 0 || entry.FPS <= 0 {
			unloadAnimationTextures(loaded)
			return SpriteAnimator{}, fmt.Errorf("animation %q in %q has invalid metadata", id, path)
		}

		texture := rl.LoadTexture(entry.Sheet)
		if texture.ID == 0 {
			unloadAnimationTextures(loaded)
			return SpriteAnimator{}, fmt.Errorf("load sprite sheet %q for animation %q", entry.Sheet, id)
		}
		if texture.Width != manifest.FrameWidth*entry.Columns || texture.Height != manifest.FrameHeight*entry.Rows {
			rl.UnloadTexture(texture)
			unloadAnimationTextures(loaded)
			return SpriteAnimator{}, fmt.Errorf(
				"sheet %q is %dx%d; manifest expects %dx%d",
				entry.Sheet,
				texture.Width,
				texture.Height,
				manifest.FrameWidth*entry.Columns,
				manifest.FrameHeight*entry.Rows,
			)
		}

		loaded = append(loaded, texture)
		clips[id] = AnimationClip{
			Texture:     texture,
			Columns:     entry.Columns,
			Rows:        entry.Rows,
			Frames:      entry.Frames,
			FPS:         entry.FPS,
			Loop:        entry.Loop,
			Directional: entry.Directional,
		}
	}

	animator := NewSpriteAnimator(clips, initial)
	if !animator.Ready() {
		unloadAnimationTextures(loaded)
		return SpriteAnimator{}, fmt.Errorf("animation manifest %q has no usable %q clip", path, initial)
	}
	return animator, nil
}

func (anim *SpriteAnimator) Unload() {
	textures := make([]rl.Texture2D, 0, len(anim.Clips))
	for _, clip := range anim.Clips {
		textures = append(textures, clip.Texture)
	}
	unloadAnimationTextures(textures)
	anim.Clips = nil
}

func unloadAnimationTextures(textures []rl.Texture2D) {
	seen := make(map[uint32]bool)
	for _, texture := range textures {
		if texture.ID == 0 || seen[texture.ID] {
			continue
		}
		seen[texture.ID] = true
		rl.UnloadTexture(texture)
	}
}

func LoadPlayerAnimatorTemplate() {
	if !rl.IsWindowReady() {
		return
	}

	template, err := LoadSpriteAnimatorManifest(playerAnimationManifestPath, AnimationIdle)
	if err != nil {
		fmt.Println("player animation setup:", err)
		return
	}
	playerAnimatorTemplate = template
}

func LoadPlayerIdleTexture() rl.Texture2D {
	if !rl.IsWindowReady() {
		return rl.Texture2D{}
	}

	texture := rl.LoadTexture(playerIdleSpritePath)
	if texture.ID == 0 {
		fmt.Println("failed to load player idle sheet:", playerIdleSpritePath)
	}
	return texture
}

func NewPlayerAnimator() SpriteAnimator {
	return NewSpriteAnimator(playerAnimatorTemplate.Clips, AnimationIdle)
}
