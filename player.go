package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	
)

type Player struct {
	Entity
	Xp           int
	Level        int
	XpToNext     int
	UpgradePoints int

	GunDamage    int
	MeleeDamage  int

	Roomscleared int
	Difficulty   int
}

func newplayer(cowboySprite rl.Texture2D, Name string) Player {
	
	player := Player{
		Entity: Entity{
			Name: Name,
			PhysicsBody: PhysicsBody{
				pos:      rl.Vector2{X: 250, Y: 250},
				velocity: rl.Vector2{X: 0, Y: 0},
			},
			SpriteRenderer: SpriteRenderer{
			Sprite:       cowboySprite,
			Color:        rl.White,
			Position:     rl.Vector2{X: 250, Y: 250},
			Angle:        0,
			Scale:        1,
			frameSize:    120,
			totalFrames:  7,
			currentFrame: 1,
		},
		HealthBar: HealthBar{
					Width:  80,
					Height: 10,
				},
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
	player.collisionBox = rl.NewRectangle(player.Position.X, player.Position.Y, 120*player.Scale, 120*player.Scale)
	return player
}

func (e *Entity) move() {
	e.pos = rl.Vector2Add(e.pos, rl.Vector2Scale(e.velocity, rl.GetFrameTime()))
	e.SyncCollisionBox()
	
	
}

func newEntity(sprite rl.Texture2D) Entity {
	Entity := Entity{
			PhysicsBody: PhysicsBody{
				pos:      rl.Vector2{X: 250, Y: 250},
				velocity: rl.Vector2{X: 0, Y: 0},
			},
			SpriteRenderer: SpriteRenderer{
			Sprite:       sprite,
			Color:        rl.White,
			Position:     rl.Vector2{X: 250, Y: 250},
			Angle:        0,
			Scale:        1,
			frameSize:    120,
			totalFrames:  7,
			currentFrame: 1,
		},
		Atkrange:	  100,
		Speed:         50,
		MaxHealth:     100,
		CurrentHealth: 100,
		Alive:         true,
		}
		return Entity
	}

func (e *Entity) Update(room *Room) {
	e.move()
	e.UpdateCollision(room)
	e.SyncCollisionBox() // keep it tight after correction
}
func (e *Entity) SyncCollisionBox() {
	e.collisionBox.X = e.pos.X - e.collisionBox.Width / 2
	e.collisionBox.Y = e.pos.Y - e.collisionBox.Height / 2
	e.collisionBox.Height = 100
	e.collisionBox.Width = 50
	e.Position = e.pos
}
func (p *Player) TakeDamage(dmg int) {
	if !p.Alive {
		return
	}

	p.CurrentHealth -= float64(dmg)

	if p.CurrentHealth <= 0 {
		p.CurrentHealth = 0
		p.Alive = false
	}
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