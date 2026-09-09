package main

import rl "github.com/gen2brain/raylib-go/raylib"

type EnemyState int

const (
	idle EnemyState = iota
	chasing
	attacking
	fleeing
	dead
	stunned
)

type EnemyType int

const (
	MeleeDemon EnemyType = iota
	RangedDemon
	TankDemon
	FastDemon
)

type EnemyConfig struct {
	EnemyType EnemyType

	Sprite rl.Texture2D

	MaxHealth float64
	Speed     float32
	Radius    float32

	SenseRange      float32
	PreferredRange  float32
	TooCloseRange   float32
	AttackRange     float32
	AttackCooldown  float32
	AttackDamage    int
	XPReward        int

	UsesProjectile  bool
	ProjectileSprite rl.Texture2D
	MeleeRadius     float32
}

type Enemy struct {
	Entity

	Config EnemyConfig

	target *Player
	state  EnemyState

	attackcooldown float32

	AttackDamage int
	XPReward     int
}

func GetEnemyConfig(enemyType EnemyType) EnemyConfig {
	switch enemyType {
	case RangedDemon:
		return EnemyConfig{
			EnemyType: enemyType,
			Sprite:    textures[littledemon], // later change to textures[rangeddemon]

			MaxHealth: 70,
			Speed:     75,
			Radius:    35,

			SenseRange:     1000,
			PreferredRange: 350,
			TooCloseRange:  220,

			AttackRange:    500,
			AttackCooldown: 1.4,
			AttackDamage:   8,
			XPReward:       35,

			UsesProjectile:   true,
			ProjectileSprite: textures[bullet],
			MeleeRadius:      0,
		}

	case TankDemon:
		return EnemyConfig{
			EnemyType: enemyType,
			Sprite:    textures[littledemon],

			MaxHealth: 180,
			Speed:     55,
			Radius:    45,

			SenseRange:     900,
			PreferredRange: 65,
			TooCloseRange:  0,

			AttackRange:    90,
			AttackCooldown: 1.6,
			AttackDamage:   18,
			XPReward:       50,

			UsesProjectile: false,
			MeleeRadius:    70,
		}

	case FastDemon:
		return EnemyConfig{
			EnemyType: enemyType,
			Sprite:    textures[littledemon],

			MaxHealth: 60,
			Speed:     150,
			Radius:    30,

			SenseRange:     1000,
			PreferredRange: 55,
			TooCloseRange:  0,

			AttackRange:    75,
			AttackCooldown: 0.8,
			AttackDamage:   7,
			XPReward:       30,

			UsesProjectile: false,
			MeleeRadius:    55,
		}

	default:
		return EnemyConfig{
			EnemyType: enemyType,
			Sprite:    textures[littledemon],

			MaxHealth: 100,
			Speed:     90,
			Radius:    35,

			SenseRange:     1000,
			PreferredRange: 60,
			TooCloseRange:  0,

			AttackRange:    80,
			AttackCooldown: 1.2,
			AttackDamage:   10,
			XPReward:       25,

			UsesProjectile: false,
			MeleeRadius:    60,
		}
	}
}

func newEnemy(enemyType EnemyType, position rl.Vector2, target *Player) Enemy {
	config := GetEnemyConfig(enemyType)

	enemy := Enemy{
		Entity: newEntity(config.Sprite),
		Config: config,
		target: target,
		state:  idle,

		AttackDamage: config.AttackDamage,
		XPReward:     config.XPReward,
	}

	enemy.pos = position
	enemy.Position = position

	enemy.radius = config.Radius
	enemy.Atkrange = config.AttackRange
	enemy.Speed = config.Speed
	enemy.MaxHealth = config.MaxHealth
	enemy.CurrentHealth = config.MaxHealth
	enemy.Alive = true

	enemy.HealthBar = HealthBar{
		Width:  70,
		Height: 8,
	}

	enemy.totalFrames = 1
	enemy.currentFrame = 1
	enemy.Scale = 1
	enemy.Color = rl.White

	enemy.SyncCollisionBox()

	return enemy
}

func (e *Enemy) update(room *Room, attackManager *AttackManager) {
	if !e.Alive {
		return
	}

	dt := rl.GetFrameTime()

	if e.attackcooldown > 0 {
		e.attackcooldown -= dt
	}

	if e.CurrentHealth <= 0 {
		e.Die()
		return
	}

	if e.target == nil || !e.target.Alive {
		e.velocity = rl.Vector2Zero()
		e.state = idle
		e.Entity.Update(room)
		return
	}

	e.UpdateAI(attackManager)

	if room != nil {
		e.Entity.Update(room)
	} else {
		e.move()
	}
}

func (e *Enemy) UpdateAI(attackManager *AttackManager) {
	distance := rl.Vector2Distance(e.pos, e.target.pos)

	dirToPlayer := rl.Vector2Subtract(e.target.pos, e.pos)

	if dirToPlayer.X != 0 || dirToPlayer.Y != 0 {
		dirToPlayer = rl.Vector2Normalize(dirToPlayer)
		e.directionfaced = dirToPlayer
	}

	if distance > e.Config.SenseRange {
		e.state = idle
		e.velocity = rl.Vector2Zero()
		return
	}

	// Ranged enemies back away if player gets too close.
	if e.Config.TooCloseRange > 0 && distance < e.Config.TooCloseRange {
		e.state = fleeing
		e.velocity = rl.Vector2Scale(dirToPlayer, -e.Speed)
		return
	}

	// Move toward player until enemy reaches preferred distance.
	if distance > e.Config.PreferredRange {
		e.state = chasing
		e.velocity = rl.Vector2Scale(dirToPlayer, e.Speed)
		return
	}

	// Stop and attack.
	e.state = attacking
	e.velocity = rl.Vector2Zero()

	if distance <= e.Config.AttackRange {
		e.TryAttack(attackManager)
	}
}

func (e *Enemy) TryAttack(attackManager *AttackManager) {
	if attackManager == nil {
		return
	}

	if e.attackcooldown > 0 {
		return
	}

	if e.target == nil {
		return
	}

	e.attackcooldown = e.Config.AttackCooldown

	dir := rl.Vector2Subtract(e.target.pos, e.pos)

	if dir.X == 0 && dir.Y == 0 {
		dir = rl.Vector2{X: 1, Y: 0}
	}

	if e.Config.UsesProjectile {
		attackManager.SpawnProjectile(
			e.pos,
			dir,
			e,
			e.Config.AttackDamage,
			e.Config.ProjectileSprite,
			1,
		)
		return
	}

	attackManager.SpawnMelee(
		e.pos,
		dir,
		e,
		e.Config.AttackDamage,
		e.Config.MeleeRadius,
		textures[Melee],
		1,
	)
}

func (e *Enemy) Die() {
	e.CurrentHealth = 0
	e.Alive = false
	e.state = dead
	e.velocity = rl.Vector2Zero()
}

func (e *Enemy) setpositiontile(x, y int) {
	e.pos = rl.Vector2{X: float32(x-1)*120 + 60, Y: float32(y-1)*120 + 60}
	e.Position = e.pos
	e.SyncCollisionBox()
}

func (e *Enemy) settexture(t rl.Texture2D) {
	e.Sprite = t
}

func (e *Enemy) TakeDamage(dmg int) bool {
	if !e.Alive {
		return false
	}

	e.CurrentHealth -= float64(dmg)

	if e.CurrentHealth <= 0 {
		e.Die()
		return true
	}

	return false
}