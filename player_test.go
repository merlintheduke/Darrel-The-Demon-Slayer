package main

import (
	"testing"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func testPlayer() Player {
	return newplayer(rl.Texture2D{}, "Test")
}

func TestPlayerTakeDamage(t *testing.T) {
	player := testPlayer()

	player.TakeDamage(30)
	if player.CurrentHealth != 70 || !player.Alive {
		t.Fatalf("partial damage changed player incorrectly: health=%v alive=%v", player.CurrentHealth, player.Alive)
	}

	player.TakeDamage(100)
	if player.CurrentHealth != 0 || player.Alive {
		t.Fatalf("lethal damage changed player incorrectly: health=%v alive=%v", player.CurrentHealth, player.Alive)
	}

	player.TakeDamage(10)
	if player.CurrentHealth != 0 || player.Alive {
		t.Fatalf("dead player took additional damage: health=%v alive=%v", player.CurrentHealth, player.Alive)
	}
}

func TestPlayerGainXPLevelsUp(t *testing.T) {
	player := testPlayer()

	player.GainXP(100)

	if player.Level != 2 {
		t.Fatalf("level = %d, want 2", player.Level)
	}
	if player.Xp != 0 {
		t.Fatalf("xp = %d, want 0", player.Xp)
	}
	if player.XpToNext != 125 {
		t.Fatalf("xp to next = %d, want 125", player.XpToNext)
	}
	if player.UpgradePoints != 1 {
		t.Fatalf("upgrade points = %d, want 1", player.UpgradePoints)
	}
	if player.MaxHealth != 110 || player.CurrentHealth != 110 {
		t.Fatalf("health after level up = %v/%v, want 110/110", player.CurrentHealth, player.MaxHealth)
	}
	if player.Speed != 305 {
		t.Fatalf("speed = %v, want 305", player.Speed)
	}
}

func TestPlayerUpgradesSpendPoints(t *testing.T) {
	player := testPlayer()
	player.UpgradePoints = 1

	player.UpgradeHealth()
	if player.UpgradePoints != 0 || player.MaxHealth != 125 || player.CurrentHealth != 125 {
		t.Fatalf("health upgrade failed: points=%d health=%v/%v", player.UpgradePoints, player.CurrentHealth, player.MaxHealth)
	}

	player.UpgradeSpeed()
	if player.Speed != 300 {
		t.Fatalf("speed changed without an upgrade point: %v", player.Speed)
	}

	player.UpgradePoints = 1
	player.UpgradeSpeed()
	if player.UpgradePoints != 0 || player.Speed != 325 {
		t.Fatalf("speed upgrade failed: points=%d speed=%v", player.UpgradePoints, player.Speed)
	}
}

func TestPlayerDamageUpgrades(t *testing.T) {
	player := testPlayer()
	player.UpgradePoints = 2

	player.UpgradeGunDamage()
	player.UpgradeMeleeDamage()

	if player.UpgradePoints != 0 {
		t.Fatalf("upgrade points = %d, want 0", player.UpgradePoints)
	}
	if player.GunDamage != 15 {
		t.Fatalf("gun damage = %d, want 15", player.GunDamage)
	}
	if player.MeleeDamage != 15 {
		t.Fatalf("melee damage = %d, want 15", player.MeleeDamage)
	}
}
