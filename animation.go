package main

import rl "github.com/gen2brain/raylib-go/raylib"

const (
	playerIdleAnimation   = "idle"
	playerWalkAnimation   = "walk"
	playerAttackAnimation = "attack"
	playerHurtAnimation   = "hurt"
	playerDeathAnimation  = "death"
)

type AnimationClip struct {
	Texture         rl.Texture2D
	FrameCount      int
	FramesPerSecond float32
	Loop            bool
}

type Animator struct {
	Clips    map[string]AnimationClip
	Current  string
	Frame    int
	Elapsed  float32
	Finished bool
}

func NewAnimator(clips map[string]AnimationClip) Animator {
	return Animator{Clips: clips}
}

func (a *Animator) Play(name string) {
	if _, ok := a.Clips[name]; !ok {
		return
	}
	if a.Current == name {
		return
	}

	a.Current = name
	a.Frame = 0
	a.Elapsed = 0
	a.Finished = false
}

func (a *Animator) Update(deltaTime float32) {
	clip, ok := a.Clips[a.Current]
	if !ok || clip.FrameCount <= 0 || clip.FramesPerSecond <= 0 || a.Finished {
		return
	}

	frameDuration := 1 / clip.FramesPerSecond
	a.Elapsed += deltaTime

	for a.Elapsed >= frameDuration {
		a.Elapsed -= frameDuration
		a.Frame++

		if a.Frame < clip.FrameCount {
			continue
		}

		if clip.Loop {
			a.Frame = 0
		} else {
			a.Frame = clip.FrameCount - 1
			a.Finished = true
			a.Elapsed = 0
			return
		}
	}
}

func (a Animator) CurrentClip() (AnimationClip, bool) {
	clip, ok := a.Clips[a.Current]
	return clip, ok
}
