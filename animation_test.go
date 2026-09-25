package main

import "testing"

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
