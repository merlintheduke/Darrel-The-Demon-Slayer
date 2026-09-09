package main

import (
	"fmt"
	"math"
	"math/rand/v2"

	rl "github.com/gen2brain/raylib-go/raylib"
)
type TileType int
const (
	nothing TileType = iota
	floor
	wall
	corner
	door
	completeddoor
)
type Room struct {
	Tiles [][]TileType
	Roomdescriptor []roomdescriptor
	walls []rl.Vector2
	//rooms
	roomspawn rl.Vector2
	pos      rl.Vector2
	tileSize float32
	sprite rl.Texture2D 
}

func newRoom(tileSize float32, sprite rl.Texture2D) Room {
	tile := make([][]TileType, 300)
	for i := range tile {
		tile[i] = make([]TileType, 300)
	}
	roomdesc := make([]roomdescriptor,0)
	walls := make([]rl.Vector2,0)
	room := Room{
		Tiles:    tile,
		Roomdescriptor: roomdesc,
		walls: walls,
		pos:      rl.Vector2Zero(),
		tileSize: tileSize,
		sprite:   sprite,
	}
	return room
}

func (r Room) Draw() {
	
	for i := 0; i < len(r.Tiles); i++ {
		for j := 0; j < len(r.Tiles[i]); j++ {
				x := r.pos.X + float32(j)*r.tileSize
				y := r.pos.Y + float32(i)*r.tileSize
			switch TileType(r.Tiles[i][j]) {
			case nothing: 
				continue
			case wall:
				rl.DrawRectangle(int32(x), int32(y), int32(r.tileSize), int32(r.tileSize), rl.DarkGreen)
				rl.DrawTextureV(r.sprite, rl.Vector2{X: x, Y: y}, rl.White)
			case floor:
				rl.DrawRectangle(int32(x), int32(y), int32(r.tileSize), int32(r.tileSize), rl.DarkBrown)
				rl.DrawTextureV(r.sprite, rl.Vector2{X: x, Y: y}, rl.DarkBrown)
			case corner:
				rl.DrawTextureV(r.sprite, rl.Vector2{X: x, Y: y}, rl.White)
			case door:
				rl.DrawRectangle(int32(x), int32(y), int32(r.tileSize), int32(r.tileSize), rl.Black)
				case completeddoor:
				rl.DrawTextureV(r.sprite, rl.Vector2{X: x, Y: y}, rl.Green)
			}
			
		}
	}
}
func (r *Room) SetTile(x int, y int, value TileType) {
	if r.Tiles[y][x] == nothing {
		r.Tiles[y][x] = value
	} else {
		return
	}
}
func (r *Room) SetRoom(pos, doorway rl.Vector2, width int, height int) {
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			if i == 0 || i == height-1 || j == 0 || j == width-1 {
				if i == 0 && j == 0 {
					r.Tiles[i+int(pos.Y)][j+int(pos.X)] = corner
				} else if i == 0 && j == width-1 {
					r.Tiles[i+int(pos.Y)][j+int(pos.X)] = corner
				} else if i == height-1 && j == 0 {
					r.Tiles[i+int(pos.Y)][j+int(pos.X)] = corner
				} else if i == height-1 && j == width-1 {
					r.Tiles[i+int(pos.Y)][j+int(pos.X)] = corner
				} else {
					r.Tiles[i+int(pos.Y)][j+int(pos.X)] = wall
					r.walls = append(r.walls, rl.Vector2{X: float32(j+int(pos.X)), Y: float32(i + int(pos.Y))})
				}
			} else {
				r.Tiles[i+int(pos.Y)][j+int(pos.X)] = floor
			}
		}
	}
}
func (r *Room) SetRoomV(pos, pos2, doorway rl.Vector2) {

	width := int(pos2.X - pos.X)
	height := int(pos2.Y - pos.Y)

	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			if i == 0 || i == height-1 || j == 0 || j == width-1 {
				if i == 0 && j == 0 {
					r.Tiles[i+int(pos.Y)][j+int(pos.X)] = corner
				} else if i == 0 && j == width-1 {
					r.Tiles[i+int(pos.Y)][j+int(pos.X)] = corner
				} else if i == height-1 && j == 0 {
					r.Tiles[i+int(pos.Y)][j+int(pos.X)] = corner
				} else if i == height-1 && j == width-1 {
					r.Tiles[i+int(pos.Y)][j+int(pos.X)] = corner
				} else {
					r.Tiles[i+int(pos.Y)][j+int(pos.X)] = wall
					r.walls = append(r.walls, rl.Vector2{X: float32(j+int(pos.X)), Y: float32(i + int(pos.Y))})
				}
			} else {
				r.Tiles[i+int(pos.Y)][j+int(pos.X)] = floor
			}
		}
	}
}

func (r *Room) RoomBuilderView() {
	room := r.Roomdescriptor[0]
	mousepos := rl.GetMousePosition()
	mouse := rl.Vector2Scale(mousepos, float32(r.tileSize/1000))
	mouse.X = float32(math.Ceil(float64(mouse.X)))
	mouse.Y = float32(math.Ceil(float64(mouse.Y)))
	rl.DrawText(fmt.Sprintf("%.0f, %.0f", mouse.X, mouse.Y), int32(mousepos.X), int32(mousepos.Y), 20, rl.RayWhite)
	if rl.CheckCollisionPointRec(mousepos, rl.Rectangle{X: mouse.X*r.tileSize, Y: mouse.Y*r.tileSize, Width: 12, Height: 12}) && room.lock == false {
		room.lock = true
		room.firsttile = mouse
		rl.DrawRectangle(int32(room.firsttile.X*r.tileSize), int32(room.firsttile.Y*r.tileSize), int32(r.tileSize/10), int32(r.tileSize/10), rl.Green)
	}
	if rl.IsMouseButtonUp(rl.MouseLeftButton) && room.lock == true {
		//room.lasttile = mouse
	}

}
func (r *Room) BuildRoom() {
	r.GenerateDungeon(20)
}


func (r *Room) RandomRoom() bool {
	width := rand.IntN(8) + 6
	height := rand.IntN(8) + 6

	if len(r.walls) == 0 {
		return false
	}

	for tries := 0; tries < len(r.walls); tries++ {
		door1 := r.walls[rand.IntN(len(r.walls))]

		if r.Tiles[int(door1.Y)][int(door1.X)] != wall {
			continue
		}

		doorX := int(door1.X)
		doorY := int(door1.Y)

		halfW := width / 2
		halfH := height / 2

		leftroom := newroomdescriptor(
			rl.Vector2{X: float32(doorX - width + 1), Y: float32(doorY - halfH)},
			rl.Vector2{X: float32(doorX + 1), Y: float32(doorY - halfH + height)},
			door1,
		)

		rightroom := newroomdescriptor(
			rl.Vector2{X: float32(doorX), Y: float32(doorY - halfH)},
			rl.Vector2{X: float32(doorX + width), Y: float32(doorY - halfH + height)},
			door1,
		)

		toproom := newroomdescriptor(
			rl.Vector2{X: float32(doorX - halfW), Y: float32(doorY - height + 1)},
			rl.Vector2{X: float32(doorX - halfW + width), Y: float32(doorY + 1)},
			door1,
		)

		bottomroom := newroomdescriptor(
			rl.Vector2{X: float32(doorX - halfW), Y: float32(doorY)},
			rl.Vector2{X: float32(doorX - halfW + width), Y: float32(doorY + height)},
			door1,
		)

		choices := []roomdescriptor{rightroom, bottomroom, leftroom, toproom}

		for _, room := range choices {
			if r.CheckRoom(room) {
				r.SetRoomV(room.firsttile, room.lasttile, door1)
				r.Tiles[doorY][doorX] = door
				r.Roomdescriptor = append(r.Roomdescriptor, room)
				return true
			}
		}
	}

	return false
}

func (r *Room) CheckRoom(room roomdescriptor) bool {
	firstX := int(room.firsttile.X)
	firstY := int(room.firsttile.Y)
	lastX := int(room.lasttile.X)
	lastY := int(room.lasttile.Y)

	doorX := int(room.doortile.X)
	doorY := int(room.doortile.Y)

	if firstX < 1 || firstY < 1 || lastX >= len(r.Tiles[0]) || lastY >= len(r.Tiles) {
		return false
	}

	checkStartX := firstX - 1
	checkEndX := lastX + 1
	checkStartY := firstY - 1
	checkEndY := lastY + 1

	// Room extends right, door is on left edge.
	if doorX == firstX {
		checkStartX = firstX + 1
	}

	// Room extends left, door is on right edge.
	if doorX == lastX-1 {
		checkEndX = lastX - 1
	}

	// Room extends down, door is on top edge.
	if doorY == firstY {
		checkStartY = firstY + 1
	}

	// Room extends up, door is on bottom edge.
	if doorY == lastY-1 {
		checkEndY = lastY - 1
	}

	for y := checkStartY; y < checkEndY; y++ {
		for x := checkStartX; x < checkEndX; x++ {
			if r.Tiles[y][x] != nothing {
				return false
			}
		}
	}

	return true
}

func (r *Room) coordinatesToTile(pos rl.Vector2) rl.Vector2 {
	tileX := float32(math.Floor(float64(pos.X / r.tileSize)))
	tileY := float32(math.Floor(float64(pos.Y / r.tileSize)))
	return rl.Vector2{X: tileX, Y: tileY}
}
func (r *Room) GenerateDungeon(goal int) {
	startX := len(r.Tiles[0]) / 2
	startY := len(r.Tiles) / 2

	startW := 8
	startH := 8

	r.SetRoom(
		rl.Vector2{X: float32(startX), Y: float32(startY)},
		rl.Vector2{},
		startW,
		startH,
	)

	r.roomspawn = rl.Vector2{
		X: float32(startX+startW/2)*r.tileSize + r.tileSize/2,
		Y: float32(startY+startH/2)*r.tileSize + r.tileSize/2,
	}

	startDoor := rl.Vector2{
	X: float32(startX + startW/2),
	Y: float32(startY),
	}

	r.Roomdescriptor = append(r.Roomdescriptor,
		newroomdescriptor(
			rl.Vector2{X: float32(startX), Y: float32(startY)},
			rl.Vector2{X: float32(startX + startW), Y: float32(startY + startH)},
			startDoor,
		),
	)

	roomsMade := 1
	failCount := 0

	for roomsMade < goal {
		if r.RandomRoom() {
			roomsMade++
			failCount = 0
		} else {
			failCount++
		}

		if failCount > 50 {
			break
		}
	}
}
func (r *Room) SetupEncounters(p *Player, difficulty int) {
	for i := range r.Roomdescriptor {
		// Start room should not have an encounter.
		if i == 0 {
			r.Roomdescriptor[i].Encounter = nil
			continue
		}

		if r.Roomdescriptor[i].Encounter == nil {
			encounter := newencounter(difficulty, p)
			r.Roomdescriptor[i].Encounter = &encounter
		} else {
			r.Roomdescriptor[i].Encounter.p = p
		}
	}
}

func (r *Room) randomFloorPosInRoom(roomdesc *roomdescriptor) rl.Vector2 {
	firstX := int(roomdesc.firsttile.X) + 1
	firstY := int(roomdesc.firsttile.Y) + 1
	lastX := int(roomdesc.lasttile.X) - 1
	lastY := int(roomdesc.lasttile.Y) - 1

	for tries := 0; tries < 100; tries++ {
		x := rand.IntN(lastX-firstX) + firstX
		y := rand.IntN(lastY-firstY) + firstY

		if r.Tiles[y][x] == floor {
			return rl.Vector2{
				X: float32(x)*r.tileSize + r.tileSize/2,
				Y: float32(y)*r.tileSize + r.tileSize/2,
			}
		}
	}

	// fallback: center of the room
	centerX := (firstX + lastX) / 2
	centerY := (firstY + lastY) / 2

	return rl.Vector2{
		X: float32(centerX)*r.tileSize + r.tileSize/2,
		Y: float32(centerY)*r.tileSize + r.tileSize/2,
	}
}

func (r *Room) blockDoor(roomdesc *roomdescriptor) {
	x := int(roomdesc.doortile.X)
	y := int(roomdesc.doortile.Y)

	if y >= 0 && y < len(r.Tiles) && x >= 0 && x < len(r.Tiles[0]) {
		r.Tiles[y][x] = wall
	}
}

func (r *Room) openDoor(roomdesc *roomdescriptor) {
	x := int(roomdesc.doortile.X)
	y := int(roomdesc.doortile.Y)

	if y >= 0 && y < len(r.Tiles) && x >= 0 && x < len(r.Tiles[0]) {
		r.Tiles[y][x] = door
	}
}
func (r *Room) completeDoor(roomdesc *roomdescriptor) {
	x := int(roomdesc.doortile.X)
	y := int(roomdesc.doortile.Y)

	if y >= 0 && y < len(r.Tiles) && x >= 0 && x < len(r.Tiles[0]) {
		r.Tiles[y][x] = completeddoor
	}
}

func (r *Room) playerInsideRoomdesc(playerPos rl.Vector2, roomdesc *roomdescriptor) bool {
	tile := r.coordinatesToTile(playerPos)

	return tile.X >= roomdesc.firsttile.X &&
		tile.X < roomdesc.lasttile.X &&
		tile.Y >= roomdesc.firsttile.Y &&
		tile.Y < roomdesc.lasttile.Y
}
func (r *Room) playerPastDoor(playerPos rl.Vector2, roomdesc *roomdescriptor) bool {
	if !r.playerInsideRoomdesc(playerPos, roomdesc) {
		return false
	}

	doorX := int(roomdesc.doortile.X)
	doorY := int(roomdesc.doortile.Y)

	firstX := int(roomdesc.firsttile.X)
	firstY := int(roomdesc.firsttile.Y)
	lastX := int(roomdesc.lasttile.X)
	lastY := int(roomdesc.lasttile.Y)

	doorCenterX := float32(doorX)*r.tileSize + r.tileSize/2
	doorCenterY := float32(doorY)*r.tileSize + r.tileSize/2

	// How far past the door the player must be before the door closes.
	margin := r.tileSize * 0.75

	// Room is to the right of the door.
	if doorX == firstX {
		return playerPos.X > doorCenterX+margin
	}

	// Room is to the left of the door.
	if doorX == lastX-1 {
		return playerPos.X < doorCenterX-margin
	}

	// Room is below the door.
	if doorY == firstY {
		return playerPos.Y > doorCenterY+margin
	}

	// Room is above the door.
	if doorY == lastY-1 {
		return playerPos.Y < doorCenterY-margin
	}

	return false
}