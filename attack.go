package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type AttackType int

const (
	Projectile AttackType = iota
	Melee
)

type Attack struct {
	PhysicsBody
	SpriteRenderer

	Type AttackType

	Damage int
	Level  int

	Lifetime  float32
	TimeAlive float32
	Active    bool

	Radius float32

	Owner interface{}
	soundeffect func()
	OnHit func(target interface{})
}

type AttackManager struct {
	Attacks []Attack
}

func (am *AttackManager) SpawnProjectile(pos, dir rl.Vector2, owner interface{}, dmg int, sprite rl.Texture2D, soundeffect int) {
	PlaySoundEffect(soundeffect)
	if dir.X == 0 && dir.Y == 0 {
		return
	}

	dir = rl.Vector2Normalize(dir)

	angle := float32(math.Atan2(float64(dir.Y), float64(dir.X)) * 180 / math.Pi)

	a := Attack{
		Type: Projectile,

		PhysicsBody: PhysicsBody{
			pos:      pos,
			velocity: rl.Vector2Scale(dir, 600),
		},

		SpriteRenderer: SpriteRenderer{
			Sprite:       sprite,
			Color:        rl.White,
			Position:     pos,
			Angle:        angle,
			Scale:        .25,
			totalFrames:  1,
			currentFrame: 1,
		},

		Damage:   dmg,
		Lifetime: 3,
		Radius:   8,
		Active:   true,
		Owner:    owner,
	}

	am.Attacks = append(am.Attacks, a)
}

func (am *AttackManager) SpawnMelee(pos, dir rl.Vector2, owner interface{}, dmg int, radius float32,sprite rl.Texture2D,soundeffect int) {
	PlaySoundEffect(soundeffect)
	if dir.X == 0 && dir.Y == 0 {
		dir = rl.Vector2{X: 1, Y: 0}
	}
	angle := float32(math.Atan2(float64(dir.Y), float64(dir.X)) * 180 / math.Pi)
	dir = rl.Vector2Normalize(dir)
	
	hitPos := rl.Vector2Add(pos, rl.Vector2Scale(dir, radius))

	a := Attack{
		Type: Melee,

		PhysicsBody: PhysicsBody{
			pos: hitPos ,
		},
		SpriteRenderer: SpriteRenderer{
			Sprite:       sprite,
			Color:        rl.White,
			Position:     hitPos,
			Angle:        angle,
			Scale:        1,
			totalFrames:  1,
			currentFrame: 1,
		},

		Damage:   dmg,
		Lifetime: 0.12,
		Radius:   radius,
		Active:   true,
		Owner:    owner,
	}

	am.Attacks = append(am.Attacks, a)
}

func (am *AttackManager) Update(dt float32) {
	for i := range am.Attacks {
		a := &am.Attacks[i]

		if !a.Active {
			continue
		}

		if a.Type == Projectile {
			a.pos = rl.Vector2Add(a.pos, rl.Vector2Scale(a.velocity, dt))
			a.SpriteRenderer.Position = a.pos
		}

		a.TimeAlive += dt

		if a.TimeAlive >= a.Lifetime {
			a.Active = false
		}
	}

	am.RemoveDeadAttacks()
}

func (am *AttackManager) Draw() {
	for i := range am.Attacks {
		a := &am.Attacks[i]

		if !a.Active {
			continue
		}

		if a.Type == Projectile {
			a.SpriteRenderer.Draw()
		} else if a.Type == Melee {
			a.SpriteRenderer.Draw()
		}
	}
}

func (am *AttackManager) CheckHits(enemies []Enemy, player *Player, room *Room) {
	for i := range am.Attacks {
		a := &am.Attacks[i]

		if !a.Active {
			continue
		}

		// Projectile wall collision.
		if room != nil && a.Type == Projectile {
			tile := room.coordinatesToTile(a.pos)
			x := int(tile.X)
			y := int(tile.Y)

			if y < 0 || y >= len(room.Tiles) || x < 0 || x >= len(room.Tiles[0]) {
				a.Active = false
				continue
			}

			if room.Tiles[y][x] == wall || room.Tiles[y][x] == corner {
				a.Active = false
				continue
			}
		}

		// Enemy-owned attacks hit the player.
		if a.Owner != player {
			if player != nil && player.Alive {
				if rl.CheckCollisionCircles(a.pos, a.Radius, player.pos, player.radius) {
					player.TakeDamage(a.Damage)
					a.Active = false
				}
			}

			continue
		}

		// Player-owned attacks hit enemies.
		for j := range enemies {
			e := &enemies[j]

			if !e.Alive {
				continue
			}

			hit := rl.CheckCollisionCircles(a.pos, a.Radius, e.pos, e.radius)

			if hit {
				killed := e.TakeDamage(a.Damage)

				if killed {
					player.GainXP(e.XPReward)
				}

				a.Active = false
				break
			}
		}
	}

	am.RemoveDeadAttacks()
}

func (am *AttackManager) RemoveDeadAttacks() {
	alive := make([]Attack, 0)

	for _, a := range am.Attacks {
		if a.Active {
			alive = append(alive, a)
		}
	}

	am.Attacks = alive
}