package main

import rl "github.com/gen2brain/raylib-go/raylib"

type SpriteRenderer struct {
	Sprite         rl.Texture2D
	Color          rl.Color
	Position       rl.Vector2
	Angle          float32
	Scale          float32
	CurrentFrame   int
	FacingRow      int
	FrameTimer     float32
	Animator       Animator
	SpriteAnimator SpriteAnimator
}

func NewSpriteRenderer(sprite rl.Texture2D, newColor rl.Color, newPosition rl.Vector2, scale float32, angle float32) SpriteRenderer {
	sr := SpriteRenderer{
		Sprite:   sprite,
		Color:    newColor,
		Position: newPosition,
		Angle:    angle,
		Scale:    scale,
		Animator: NewAnimator(map[string]AnimationClip{
			"default": {Texture: sprite, FrameCount: 1, FramesPerSecond: 1, Loop: true},
		}),
	}
	sr.Animator.Play("default")
	return sr
}

func NewAnimatedSpriteRenderer(sprite rl.Texture2D, newColor rl.Color, newPosition rl.Vector2, scale float32, angle float32, frameCount int, framesPerSecond float32) SpriteRenderer {
	sr := NewSpriteRenderer(sprite, newColor, newPosition, scale, angle)
	sr.Animator.Clips["default"] = AnimationClip{
		Texture:         sprite,
		FrameCount:      frameCount,
		FramesPerSecond: framesPerSecond,
		Loop:            true,
	}
	sr.Animator.Play("default")
	return sr
}

func (sr SpriteRenderer) Draw() {
	if sr.Sprite.ID == 0 {
		return
	}
	if sr.SpriteAnimator.Ready() {
		sr.SpriteAnimator.Draw(sr.Position, sr.Scale, sr.Color)
		return
	}

	if sr.Sprite.Width == 480 && sr.Sprite.Height == 960 {
		source := rl.Rectangle{
			X:      float32(sr.CurrentFrame * 120),
			Y:      float32(sr.FacingRow * 120),
			Width:  120,
			Height: 120,
		}
		destination := rl.Vector2{X: sr.Position.X - 60, Y: sr.Position.Y - 60}
		rl.DrawTextureRec(sr.Sprite, source, destination, sr.Color)
		return
	}

	clip, ok := sr.Animator.CurrentClip()
	if !ok {
		clip = AnimationClip{Texture: sr.Sprite, FrameCount: 1}
	}

	frameCount := clip.FrameCount
	if frameCount <= 0 {
		frameCount = 1
	}
	frame := sr.Animator.Frame
	if frame < 0 || frame >= frameCount {
		frame = 0
	}

	frameHeight := float32(clip.Texture.Height) / float32(frameCount)
	sourceRect := rl.NewRectangle(0, float32(frame)*frameHeight, float32(clip.Texture.Width), frameHeight)
	destRect := rl.NewRectangle(sr.Position.X, sr.Position.Y, float32(clip.Texture.Width)*sr.Scale, frameHeight*sr.Scale)
	origin := rl.Vector2Scale(
		rl.NewVector2(
			float32(clip.Texture.Width)/2,
			frameHeight/2,
		),
		sr.Scale,
	)
	rl.DrawTexturePro(clip.Texture, sourceRect,
		destRect,
		origin, sr.Angle, sr.Color)
}
