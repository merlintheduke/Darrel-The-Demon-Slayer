package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	
)
type Entity struct {
	Name string
	PhysicsBody
	SpriteRenderer
	HealthBar

	radius        float32
	Atkrange     float32
	Speed        float32
	MaxHealth    float64
	CurrentHealth float64
	Alive        bool
}
type HealthBar struct {
	Width  float32
	Height float32
}
type PhysicsBody struct {
	collisionBox rl.Rectangle
	directionfaced rl.Vector2
	velocity     rl.Vector2
	pos          rl.Vector2
}
func (p *PhysicsBody) UpdateCollision(room *Room) {

	for j := range room.Tiles{
		for i := range room.Tiles {

			// changed: correct world-space tile position
			tileRect := rl.NewRectangle(
				room.pos.X + float32(i)*room.tileSize,
				room.pos.Y + float32(j)*room.tileSize,
				room.tileSize,
				room.tileSize,
			)
			if room.Tiles[j][i] != wall {
				continue
			}
			if rl.CheckCollisionRecs(p.collisionBox, tileRect) {

				dxLeft := (p.collisionBox.X + p.collisionBox.Width) - tileRect.X
				dxRight := (tileRect.X + tileRect.Width) - p.collisionBox.X
				dyTop := (p.collisionBox.Y + p.collisionBox.Height) - tileRect.Y
				dyBottom := (tileRect.Y + tileRect.Height) - p.collisionBox.Y

				if dxLeft < dxRight && dxLeft < dyTop && dxLeft < dyBottom {
					p.pos.X -= dxLeft/2
					p.velocity.X = 0

				} else if dxRight < dyTop && dxRight < dyBottom {
					p.pos.X += dxRight/2
					p.velocity.X = 0

				} else if dyTop < dyBottom {
					p.pos.Y -= dyTop/2
					p.velocity.Y = 0

				} else {
					p.pos.Y += dyBottom/2
					p.velocity.Y = 0
				}
			}
		}
	}
}
func (p *PhysicsBody) drawCollisionBox() {
	rl.DrawRectangleRec(p.collisionBox,rl.Color{R: 255,G: 0,B: 0,A: 100})
}
func (e *Entity) DrawHealthBar() {
	if e.MaxHealth <= 0 {
		return
	}

	barWidth := e.HealthBar.Width
	barHeight := e.HealthBar.Height

	// fallback defaults
	if barWidth <= 0 {
		barWidth = 70
	}
	if barHeight <= 0 {
		barHeight = 8
	}

	// health percent clamp
	percent := float32(e.CurrentHealth / e.MaxHealth)
	if percent < 0 {
		percent = 0
	}
	if percent > 1 {
		percent = 1
	}

	// position above entity
	x := e.pos.X - barWidth/2
	y := e.pos.Y - e.collisionBox.Height/2 - 15

	// background
	rl.DrawRectangle(
		int32(x),
		int32(y),
		int32(barWidth),
		int32(barHeight),
		rl.DarkGray,
	)

	// health fill
	rl.DrawRectangle(
		int32(x),
		int32(y),
		int32(barWidth*percent),
		int32(barHeight),
		rl.Red,
	)

	// outline
	rl.DrawRectangleLines(
		int32(x),
		int32(y),
		int32(barWidth),
		int32(barHeight),
		rl.Black,
	)
}