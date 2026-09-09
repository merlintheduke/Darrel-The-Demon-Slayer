package main
import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type roomdescriptor struct {
	firsttile  rl.Vector2
	lasttile rl.Vector2
	doortile rl.Vector2
	lock      bool
	Button
	*Encounter
}
func newroomdescriptor(first, last, door rl.Vector2) roomdescriptor {
	return roomdescriptor{
		firsttile: rl.Vector2{
			X: float32(int(first.X)), 
			Y: float32(int(first.Y)),
		},
		lasttile: rl.Vector2{
			X: float32(int(last.X)),  
			Y: float32(int(last.Y)),
		},
		doortile: rl.Vector2{
			X: float32(int(door.X)), 
			Y: float32(int(door.Y)),
		},
	}
}
