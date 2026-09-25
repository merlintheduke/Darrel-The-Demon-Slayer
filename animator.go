package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// AnimationID is intentionally shared by the player and enemies. An actor
// only needs to provide clips for the states it actually supports.
type AnimationID string

const (
	AnimationIdle   AnimationID = "idle"
	AnimationWalk   AnimationID = "walk"
	AnimationSprint AnimationID = "sprint"
	AnimationGun    AnimationID = "gun"
	AnimationSword  AnimationID = "sword"
	AnimationFist   AnimationID = "fist"
	AnimationHurt   AnimationID = "hurt"
	AnimationHeal   AnimationID = "heal"
	AnimationDeath  AnimationID = "death"
)

// FacingDirection doubles as the row index for directional sprite sheets.
// Its ordering is reusable by every 8-way sheet in the project.
type FacingDirection int32

const (
	FacingSouth FacingDirection = iota
	FacingSouthEast
	FacingEast
	FacingNorthEast
	FacingNorth
	FacingNorthWest
	FacingWest
	FacingSouthWest
)

type AnimationClip struct {
	Texture         rl.Texture2D
	FrameCount      int
	FramesPerSecond float32
	Loop            bool
	Columns         int32
	Rows            int32
	Frames          int32
	FPS             float32
	Directional     bool
}

// SpriteAnimator is the reusable atlas animator for every game actor.
// A directional clip uses Frame as the X column and Facing as the Y row.
// A non-directional clip consumes Frames left-to-right, then top-to-bottom.
type SpriteAnimator struct {
	Clips    map[AnimationID]AnimationClip
	Current  AnimationID
	Facing   FacingDirection
	Frame    int32
	elapsed  float32
	complete bool
}

func NewSpriteAnimator(clips map[AnimationID]AnimationClip, initial AnimationID) SpriteAnimator {
	return SpriteAnimator{
		Clips:   clips,
		Current: initial,
		Facing:  FacingSouth,
	}
}

func NewStaticAnimator(texture rl.Texture2D) SpriteAnimator {
	return NewSpriteAnimator(map[AnimationID]AnimationClip{
		AnimationIdle: {
			Texture: texture,
			Columns: 1,
			Rows:    1,
			Frames:  1,
			FPS:     1,
			Loop:    true,
		},
	}, AnimationIdle)
}

func (anim SpriteAnimator) Clip(animation AnimationID) (AnimationClip, bool) {
	clip, ok := anim.Clips[animation]
	if !ok || clip.Texture.ID == 0 || clip.Columns <= 0 || clip.Rows <= 0 || clip.Frames <= 0 {
		return AnimationClip{}, false
	}
	return clip, true
}

func (anim SpriteAnimator) Ready() bool {
	_, ok := anim.Clip(anim.Current)
	return ok
}

func (anim SpriteAnimator) IsComplete() bool {
	return anim.complete
}

func (anim SpriteAnimator) IsLocked() bool {
	clip, ok := anim.Clip(anim.Current)
	return ok && !clip.Loop && !anim.complete
}

func (anim *SpriteAnimator) Play(animation AnimationID) {
	if _, ok := anim.Clip(animation); !ok {
		return
	}
	if anim.Current == animation {
		return
	}

	anim.Current = animation
	anim.Frame = 0
	anim.elapsed = 0
	anim.complete = false
}

func (anim *SpriteAnimator) Replay(animation AnimationID) {
	if _, ok := anim.Clip(animation); !ok {
		return
	}

	anim.Current = animation
	anim.Frame = 0
	anim.elapsed = 0
	anim.complete = false
}

func (anim *SpriteAnimator) SetFacing(direction rl.Vector2) {
	if direction.X == 0 && direction.Y == 0 {
		return
	}

	angle := math.Atan2(float64(direction.Y), float64(direction.X))
	sector := int(math.Floor((angle+math.Pi/8)/(math.Pi/4))) % 8
	if sector < 0 {
		sector += 8
	}

	// Screen-space Y grows down, so this maps the mathematical angle to the
	// fixed sheet row ordering: S, SE, E, NE, N, NW, W, SW.
	switch sector {
	case 0:
		anim.Facing = FacingEast
	case 1:
		anim.Facing = FacingSouthEast
	case 2:
		anim.Facing = FacingSouth
	case 3:
		anim.Facing = FacingSouthWest
	case 4:
		anim.Facing = FacingWest
	case 5:
		anim.Facing = FacingNorthWest
	case 6:
		anim.Facing = FacingNorth
	default:
		anim.Facing = FacingNorthEast
	}
}

func (anim *SpriteAnimator) Update(dt float32) {
	clip, ok := anim.Clip(anim.Current)
	if !ok || clip.FPS <= 0 {
		return
	}

	anim.elapsed += dt
	frameDuration := 1 / clip.FPS
	for anim.elapsed >= frameDuration {
		anim.elapsed -= frameDuration
		anim.Frame++

		if anim.Frame < clip.Frames {
			continue
		}

		if clip.Loop {
			anim.Frame = 0
			continue
		}

		anim.Frame = clip.Frames - 1
		anim.complete = true
		return
	}
}

func (anim SpriteAnimator) Draw(position rl.Vector2, scale float32, tint rl.Color) {
	clip, ok := anim.Clip(anim.Current)
	if !ok {
		return
	}

	frameWidth := clip.Texture.Width / clip.Columns
	frameHeight := clip.Texture.Height / clip.Rows
	column := anim.Frame % clip.Columns
	row := anim.Frame / clip.Columns
	if clip.Directional {
		row = int32(anim.Facing)
	}

	if row >= clip.Rows {
		return
	}

	source := rl.NewRectangle(
		float32(column*frameWidth),
		float32(row*frameHeight),
		float32(frameWidth),
		float32(frameHeight),
	)
	destination := rl.NewRectangle(
		position.X,
		position.Y,
		float32(frameWidth)*scale,
		float32(frameHeight)*scale,
	)
	origin := rl.NewVector2(
		float32(frameWidth)*scale/2,
		float32(frameHeight)*scale/2,
	)
	rl.DrawTexturePro(clip.Texture, source, destination, origin, 0, tint)
}
