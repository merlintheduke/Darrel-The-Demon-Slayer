package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Player struct {
	Entity
	Xp            int
	Level         int
	XpToNext      int
	UpgradePoints int

	GunDamage   int
	MeleeDamage int

	Roomscleared  int
	Difficulty    int
	Sprinting     bool
	shootCooldown float32
}

func newPlayerAnimator(sprite rl.Texture2D) Animator {
	if playerAnimatorTemplate.Ready() {
		clips := map[string]AnimationClip{}
		if clip, ok := playerAnimatorTemplate.Clip(AnimationIdle); ok {
			clips[playerIdleAnimation] = AnimationClip{
				Texture:         clip.Texture,
				FrameCount:      int(clip.Frames),
				FramesPerSecond: clip.FPS,
				Loop:            clip.Loop,
			}
		}
		if clip, ok := playerAnimatorTemplate.Clip(AnimationWalk); ok {
			clips[playerWalkAnimation] = AnimationClip{
				Texture:         clip.Texture,
				FrameCount:      int(clip.Frames),
				FramesPerSecond: clip.FPS,
				Loop:            clip.Loop,
			}
		}
		if clip, ok := playerAnimatorTemplate.Clip(AnimationGun); ok {
			clips[playerAttackAnimation] = AnimationClip{
				Texture:         clip.Texture,
				FrameCount:      int(clip.Frames),
				FramesPerSecond: clip.FPS,
				Loop:            clip.Loop,
			}
		}
		if clip, ok := playerAnimatorTemplate.Clip(AnimationHurt); ok {
			clips[playerHurtAnimation] = AnimationClip{
				Texture:         clip.Texture,
				FrameCount:      int(clip.Frames),
				FramesPerSecond: clip.FPS,
				Loop:            clip.Loop,
			}
		}
		if clip, ok := playerAnimatorTemplate.Clip(AnimationDeath); ok {
			clips[playerDeathAnimation] = AnimationClip{
				Texture:         clip.Texture,
				FrameCount:      int(clip.Frames),
				FramesPerSecond: clip.FPS,
				Loop:            clip.Loop,
			}
		}
		if len(clips) >= 5 {
			animator := NewAnimator(clips)
			animator.Play(playerIdleAnimation)
			return animator
		}
	}

	clips := map[string]AnimationClip{}
	for _, name := range []string{playerIdleAnimation, playerWalkAnimation} {
		clips[name] = playerAnimationClip(sprite, true)
	}
	for _, name := range []string{playerAttackAnimation, playerHurtAnimation, playerDeathAnimation} {
		clips[name] = playerAnimationClip(sprite, false)
	}

	animator := NewAnimator(clips)
	animator.Play(playerIdleAnimation)
	return animator
}

func playerAnimationClip(sprite rl.Texture2D, loop bool) AnimationClip {
	return AnimationClip{
		Texture:         sprite,
		FrameCount:      7,
		FramesPerSecond: 1,
		Loop:            loop,
	}
}

func newPlayerSpriteAnimator() SpriteAnimator {
	if playerAnimatorTemplate.Ready() {
		return NewPlayerAnimator()
	}
	return SpriteAnimator{}
}

func (p *Player) UpdateAnimation(attacking bool) {
	animator := &p.Renderer.Animator
	spriteAnimator := &p.Renderer.SpriteAnimator
	if !p.Alive {
		animator.Play(playerDeathAnimation)
		spriteAnimator.Play(AnimationDeath)
		return
	}

	if attacking {
		animator.Play(playerAttackAnimation)
		spriteAnimator.Play(AnimationGun)
		return
	}

	if (animator.Current == playerAttackAnimation || animator.Current == playerHurtAnimation) && !animator.Finished {
		return
	}
	if spriteAnimator.IsLocked() {
		return
	}

	if p.velocity.X != 0 || p.velocity.Y != 0 {
		if p.Sprinting {
			animator.Play(playerWalkAnimation)
			spriteAnimator.Play(AnimationSprint)
		} else {
			animator.Play(playerWalkAnimation)
			spriteAnimator.Play(AnimationWalk)
		}
		spriteAnimator.SetFacing(p.velocity)
		return
	}

	animator.Play(playerIdleAnimation)
	spriteAnimator.Play(AnimationIdle)
}

func (p *Player) SetFacing(direction rl.Vector2) {
	p.Renderer.SpriteAnimator.SetFacing(direction)
}

func (p *Player) UpdateShootCooldown(deltaTime float32) {
	p.shootCooldown -= deltaTime
	if p.shootCooldown < 0 {
		p.shootCooldown = 0
	}
}

func (p *Player) CanShoot() bool {
	return p.shootCooldown <= 0
}

func (p *Player) StartGunAttack() {
	p.shootCooldown = playerGunCooldown
	p.Renderer.Animator.Replay(playerAttackAnimation)
	p.Renderer.SpriteAnimator.Replay(AnimationGun)
}

func (p *Player) DeathAnimationFinished() bool {
	if p.Renderer.SpriteAnimator.Ready() {
		return p.Renderer.SpriteAnimator.Current == AnimationDeath && p.Renderer.SpriteAnimator.IsComplete()
	}
	return p.Renderer.Animator.Current == playerDeathAnimation && p.Renderer.Animator.Finished
}

func newplayer(cowboySprite rl.Texture2D, Name string) Player {
	if cowboySprite.ID == 0 && rl.IsWindowReady() {
		cowboySprite = LoadPlayerIdleTexture()
	}

	player := Player{
		Entity: Entity{
			Name: Name,
			PhysicsBody: PhysicsBody{
				pos:      rl.Vector2{X: playerSpawnX, Y: playerSpawnY},
				velocity: rl.Vector2{X: 0, Y: 0},
			},
			Renderer: EntityRenderer{SpriteRenderer: SpriteRenderer{
				Sprite:         cowboySprite,
				Color:          rl.White,
				Position:       rl.Vector2{X: playerSpawnX, Y: playerSpawnY},
				Scale:          1,
				CurrentFrame:   0,
				FacingRow:      0,
				Animator:       newPlayerAnimator(cowboySprite),
				SpriteAnimator: newPlayerSpriteAnimator(),
			}, HealthBar: HealthBar{
				Width:  80,
				Height: 10,
			}},
			Speed:         300,
			MaxHealth:     100,
			CurrentHealth: 100,
			Alive:         true,
		},

		Xp:            0,
		Level:         1,
		XpToNext:      100,
		UpgradePoints: 0,
		GunDamage:     10,
		MeleeDamage:   10,
	}
	player.Position = player.pos
	player.collisionBox = rl.NewRectangle(player.Position.X, player.Position.Y, 120*player.Renderer.Scale, 120*player.Renderer.Scale)
	return player
}

func (e *Entity) move() {
	e.pos = rl.Vector2Add(e.pos, rl.Vector2Scale(e.velocity, rl.GetFrameTime()))
	e.SyncCollisionBox()

}

func newEntity(sprite rl.Texture2D) Entity {
	entity := Entity{
		PhysicsBody: PhysicsBody{
			pos:      rl.Vector2{X: playerSpawnX, Y: playerSpawnY},
			velocity: rl.Vector2{X: 0, Y: 0},
		},
		Renderer:      EntityRenderer{SpriteRenderer: NewSpriteRenderer(sprite, rl.White, rl.Vector2{X: playerSpawnX, Y: playerSpawnY}, 1, 0)},
		Atkrange:      100,
		Speed:         50,
		MaxHealth:     100,
		CurrentHealth: 100,
		Alive:         true,
	}
	return entity
}

func (e *Entity) Update(room *Room) {
	e.Renderer.Update(rl.GetFrameTime())
	e.move()
	e.UpdateCollision(room)
	e.SyncCollisionBox() // keep it tight after correction
}
func (e *Entity) SyncCollisionBox() {
	e.collisionBox.X = e.pos.X - e.collisionBox.Width/2
	e.collisionBox.Y = e.pos.Y - e.collisionBox.Height/2
	e.collisionBox.Height = entityCollisionHeight
	e.collisionBox.Width = entityCollisionWidth
	e.Position = e.pos
	e.Renderer.Position = e.pos
}
func (p *Player) TakeDamage(dmg int) {
	if !p.Alive {
		return
	}

	p.CurrentHealth -= float64(dmg)

	if p.CurrentHealth <= 0 {
		p.CurrentHealth = 0
		p.Alive = false
		p.Renderer.Animator.Play(playerDeathAnimation)
		p.Renderer.SpriteAnimator.Play(AnimationDeath)
		return
	}

	p.Renderer.Animator.Play(playerHurtAnimation)
	p.Renderer.SpriteAnimator.Play(AnimationHurt)
}

func (p *Player) GainXP(amount int) {
	if amount <= 0 {
		return
	}

	if p.Level <= 0 {
		p.Level = 1
	}

	if p.XpToNext <= 0 {
		p.XpToNext = 100
	}

	p.Xp += amount

	for p.Xp >= p.XpToNext {
		p.Xp -= p.XpToNext
		p.Level++
		p.UpgradePoints++
		p.XpToNext = int(float32(p.XpToNext) * 1.25)

		p.MaxHealth += 10
		p.CurrentHealth = p.MaxHealth

		p.Speed += 5
	}
}
func (p *Player) SpendUpgradePoint() bool {
	if p.UpgradePoints <= 0 {
		return false
	}

	p.UpgradePoints--
	return true
}

func (p *Player) UpgradeHealth() {
	if !p.SpendUpgradePoint() {
		return
	}

	p.MaxHealth += 25
	p.CurrentHealth += 25
}

func (p *Player) UpgradeSpeed() {
	if !p.SpendUpgradePoint() {
		return
	}

	p.Speed += 25
}

func (p *Player) UpgradeGunDamage() {
	if !p.SpendUpgradePoint() {
		return
	}

	p.GunDamage += 5
}

func (p *Player) UpgradeMeleeDamage() {
	if !p.SpendUpgradePoint() {
		return
	}

	p.MeleeDamage += 5
}
