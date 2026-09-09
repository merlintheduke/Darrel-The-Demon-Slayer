package main

import (
	"math"
	"math/rand/v2"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Encounter struct {
	enemies    []Enemy
	started    bool
	generated  bool
	ended      bool
	quantity   int
	difficulty int
	timer      float32
	p          *Player
}

func newencounter(difficulty int, p *Player) Encounter {
	return Encounter{
		difficulty: difficulty,
		p:          p,
	}
}

func randomEnemyType(difficulty int, roomsCleared int) EnemyType {
	roll := rand.IntN(100)

	// First couple rooms are easier.
	if roomsCleared < 2 {
		return MeleeDemon
	}

	if roll < 55 {
		return MeleeDemon
	}

	if roll < 80 {
		return RangedDemon
	}

	if roll < 92 {
		return FastDemon
	}

	return TankDemon
}

func (e *Encounter) startEncounter(r *Room, roomdesc *roomdescriptor) {
	if e.started || e.ended {
		return
	}

	e.started = true
	e.generated = true

	e.quantity = e.p.quantityBasedOnDifficulty(e.difficulty)

	if e.quantity < 1 {
		e.quantity = 1
	}

	for i := 0; i < e.quantity; i++ {
		spawnPos := r.randomFloorPosInRoom(roomdesc)

		enemyType := randomEnemyType(e.difficulty, e.p.Roomscleared)
		enemy := newEnemy(enemyType, spawnPos, e.p)

		e.enemies = append(e.enemies, enemy)
	}
}

func (e *Encounter) updateEncounter(room *Room, attackManager *AttackManager) {
	if !e.started || e.ended {
		return
	}

	e.timer += rl.GetFrameTime()

	aliveCount := 0

	for i := range e.enemies {
		if e.enemies[i].Alive {
			e.enemies[i].update(room, attackManager)
			aliveCount++
		}
	}

	if e.generated && aliveCount == 0 {
		e.endEncounter()
	}
}

func (e *Encounter) drawEncounter() {
	if !e.started || e.ended {
		return
	}

	for i := range e.enemies {
		if e.enemies[i].Alive {
			e.enemies[i].Draw()
			e.enemies[i].Entity.DrawHealthBar()
		}
	}
}

func (e *Encounter) endEncounter() {
	if e.ended {
		return
	}

	e.ended = true

	if e.p != nil {
		e.p.Roomscleared++
	}
}

func (p *Player) quantityBasedOnDifficulty(difficulty int) int {
	total := difficulty + p.Roomscleared/10

	enemies := rand.IntN(3+total) + total + int(math.Floor(float64(p.Roomscleared/50.0)))
	enemies += int(math.Floor(float64(p.Roomscleared/100.0)))

	if enemies < 1 {
		enemies = 1
	}

	return enemies
}