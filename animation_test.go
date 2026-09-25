package main

import (
	"testing"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func TestPlayerAnimatorDefinesNamedClips(t *testing.T) {
	animator := newPlayerAnimator(rl.Texture2D{ID: 1})

	for _, name := range []string{
		playerIdleAnimation,
		playerWalkAnimation,
		playerAttackAnimation,
		playerHurtAnimation,
		playerDeathAnimation,
	} {
		clip, ok := animator.Clips[name]
		if !ok {
			t.Fatalf("missing player animation clip %q", name)
		}
		if clip.FrameCount != 7 || clip.FramesPerSecond != 1 {
			t.Fatalf("unexpected %q clip: %+v", name, clip)
		}
	}

	if animator.Current != playerIdleAnimation {
		t.Fatalf("current animation = %q, want %q", animator.Current, playerIdleAnimation)
	}
}

func TestPlayerUpdateAnimationChoosesMovementStates(t *testing.T) {
	player := newplayer(rl.Texture2D{ID: 1}, "Test")

	player.UpdateAnimation(false)
	if player.Renderer.Animator.Current != playerIdleAnimation {
		t.Fatalf("idle animation = %q, want %q", player.Renderer.Animator.Current, playerIdleAnimation)
	}

	player.velocity.X = 1
	player.UpdateAnimation(false)
	if player.Renderer.Animator.Current != playerWalkAnimation {
		t.Fatalf("moving animation = %q, want %q", player.Renderer.Animator.Current, playerWalkAnimation)
	}

	player.UpdateAnimation(true)
	if player.Renderer.Animator.Current != playerAttackAnimation {
		t.Fatalf("attacking animation = %q, want %q", player.Renderer.Animator.Current, playerAttackAnimation)
	}
}

func TestPlayerWalkAnimationAdvancesFrames(t *testing.T) {
	player := newplayer(rl.Texture2D{ID: 1}, "Test")
	player.Renderer.SpriteAnimator = NewSpriteAnimator(map[AnimationID]AnimationClip{
		AnimationIdle: {Texture: rl.Texture2D{ID: 1}, Columns: 4, Rows: 8, Frames: 4, FPS: 5, Loop: true},
		AnimationWalk: {Texture: rl.Texture2D{ID: 1}, Columns: 8, Rows: 8, Frames: 8, FPS: 6, Loop: true, Directional: true},
		AnimationGun:  {Texture: rl.Texture2D{ID: 1}, Columns: 4, Rows: 8, Frames: 4, FPS: 16, Loop: false, Directional: true},
	}, AnimationIdle)

	player.velocity.X = 1
	player.UpdateAnimation(false)
	if player.Renderer.SpriteAnimator.Current != AnimationWalk {
		t.Fatalf("atlas animation = %q, want %q", player.Renderer.SpriteAnimator.Current, AnimationWalk)
	}

	player.Renderer.Update(1.0 / 6.0)
	if player.Renderer.SpriteAnimator.Frame != 1 {
		t.Fatalf("walk frame = %d, want 1", player.Renderer.SpriteAnimator.Frame)
	}

	player.UpdateAnimation(false)
	if player.Renderer.SpriteAnimator.Frame != 1 {
		t.Fatalf("walk frame reset during state refresh: frame=%d", player.Renderer.SpriteAnimator.Frame)
	}
}

func TestPlayerSprintAnimationUsesSprintClip(t *testing.T) {
	player := newplayer(rl.Texture2D{ID: 1}, "Test")
	player.Renderer.SpriteAnimator = NewSpriteAnimator(map[AnimationID]AnimationClip{
		AnimationIdle:   {Texture: rl.Texture2D{ID: 1}, Columns: 4, Rows: 8, Frames: 4, FPS: 5, Loop: true},
		AnimationSprint: {Texture: rl.Texture2D{ID: 1}, Columns: 8, Rows: 8, Frames: 8, FPS: 9, Loop: true, Directional: true},
	}, AnimationIdle)

	player.velocity.X = 1
	player.Sprinting = true
	player.UpdateAnimation(false)
	if player.Renderer.SpriteAnimator.Current != AnimationSprint {
		t.Fatalf("sprint animation = %q, want %q", player.Renderer.SpriteAnimator.Current, AnimationSprint)
	}
}

func TestPlayerAttackUsesGunAnimation(t *testing.T) {
	player := newplayer(rl.Texture2D{ID: 1}, "Test")
	player.Renderer.SpriteAnimator = NewSpriteAnimator(map[AnimationID]AnimationClip{
		AnimationIdle: {Texture: rl.Texture2D{ID: 1}, Columns: 4, Rows: 8, Frames: 4, FPS: 5, Loop: true},
		AnimationGun:  {Texture: rl.Texture2D{ID: 1}, Columns: 4, Rows: 8, Frames: 4, FPS: 16, Loop: false, Directional: true},
	}, AnimationIdle)

	player.UpdateAnimation(true)
	if player.Renderer.SpriteAnimator.Current != AnimationGun {
		t.Fatalf("attack animation = %q, want %q", player.Renderer.SpriteAnimator.Current, AnimationGun)
	}
}

func TestPlayerGunCooldownMatchesAnimation(t *testing.T) {
	player := newplayer(rl.Texture2D{ID: 1}, "Test")
	if !player.CanShoot() {
		t.Fatal("new player should be able to shoot")
	}

	player.StartGunAttack()
	if player.CanShoot() {
		t.Fatal("player shot again before the gun animation duration")
	}

	player.UpdateShootCooldown(playerGunCooldown - 0.001)
	if player.CanShoot() {
		t.Fatal("player cooldown ended before the gun animation duration")
	}

	player.UpdateShootCooldown(0.002)
	if !player.CanShoot() {
		t.Fatal("player cooldown did not end with the gun animation duration")
	}
}

func TestDeathAnimationMustFinishBeforeGameOver(t *testing.T) {
	player := newplayer(rl.Texture2D{ID: 1}, "Test")
	player.Renderer.SpriteAnimator = NewSpriteAnimator(map[AnimationID]AnimationClip{
		AnimationDeath: {Texture: rl.Texture2D{ID: 1}, Columns: 4, Rows: 8, Frames: 4, FPS: 8, Loop: false, Directional: true},
	}, AnimationDeath)

	player.Alive = false
	player.Renderer.SpriteAnimator.Play(AnimationDeath)
	if player.DeathAnimationFinished() {
		t.Fatal("death animation finished before any frames elapsed")
	}

	player.Renderer.SpriteAnimator.Update(0.5)
	if !player.DeathAnimationFinished() {
		t.Fatal("death animation did not finish after its full duration")
	}
}

func TestAnimatorPlayResetsAnimation(t *testing.T) {
	animator := NewAnimator(map[string]AnimationClip{
		"walk": {FrameCount: 4, FramesPerSecond: 8, Loop: true},
	})
	animator.Frame = 2
	animator.Elapsed = 0.1
	animator.Finished = true

	animator.Play("walk")

	if animator.Current != "walk" || animator.Frame != 0 || animator.Elapsed != 0 || animator.Finished {
		t.Fatalf("Play did not reset animation: %+v", animator)
	}
}

func TestAnimatorLoops(t *testing.T) {
	animator := NewAnimator(map[string]AnimationClip{
		"walk": {FrameCount: 3, FramesPerSecond: 10, Loop: true},
	})
	animator.Play("walk")

	animator.Update(0.1)
	animator.Update(0.1)
	animator.Update(0.1)

	if animator.Frame != 0 || animator.Finished {
		t.Fatalf("looping animation did not wrap: frame=%d finished=%v", animator.Frame, animator.Finished)
	}
}

func TestAnimatorStopsOnLastFrame(t *testing.T) {
	animator := NewAnimator(map[string]AnimationClip{
		"death": {FrameCount: 3, FramesPerSecond: 10, Loop: false},
	})
	animator.Play("death")

	animator.Update(0.3)

	if animator.Frame != 2 || !animator.Finished {
		t.Fatalf("non-looping animation did not finish on last frame: frame=%d finished=%v", animator.Frame, animator.Finished)
	}

	animator.Update(1)
	if animator.Frame != 2 {
		t.Fatalf("finished animation advanced: frame=%d", animator.Frame)
	}
}

func TestAnimatorIgnoresUnknownClip(t *testing.T) {
	animator := NewAnimator(map[string]AnimationClip{})
	animator.Play("missing")
	animator.Update(1)

	if animator.Current != "" || animator.Frame != 0 {
		t.Fatalf("unknown clip changed animator: %+v", animator)
	}
}
