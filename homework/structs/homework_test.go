package main

import (
	"fmt"
	"math"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

type Option func(*GamePerson)

func WithName(name string) func(*GamePerson) {
	return func(person *GamePerson) {
		copy(person.name[:], name)
	}
}

func WithCoordinates(x, y, z int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.x = int32(x)
		person.y = int32(y)
		person.z = int32(z)
	}
}

func WithGold(gold int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.gold = uint32(gold)
	}
}

func WithMana(mana int) func(*GamePerson) {
	return func(person *GamePerson) {
		setField32(&person.other, uint32(mana), 0, 10)
	}
}

func WithHealth(health int) func(*GamePerson) {
	return func(person *GamePerson) {
		setField32(&person.other, uint32(health), 10, 10)
	}
}

func WithRespect(respect int) func(*GamePerson) {
	return func(person *GamePerson) {
		setField16(&person.stats, uint16(respect), 0, 4)
	}
}

func WithStrength(strength int) func(*GamePerson) {
	return func(person *GamePerson) {
		setField16(&person.stats, uint16(strength), 4, 4)
	}
}

func WithExperience(experience int) func(*GamePerson) {
	return func(person *GamePerson) {
		setField16(&person.stats, uint16(experience), 8, 4)
	}
}

func WithLevel(level int) func(*GamePerson) {
	return func(person *GamePerson) {
		setField16(&person.stats, uint16(level), 12, 4)
	}
}

func WithHouse() func(*GamePerson) {
	return func(person *GamePerson) {
		setField32(&person.other, 1, 20, 1)
	}
}

func WithGun() func(*GamePerson) {
	return func(person *GamePerson) {
		setField32(&person.other, 1, 21, 1)
	}
}

func WithFamily() func(*GamePerson) {
	return func(person *GamePerson) {
		setField32(&person.other, 1, 22, 1)
	}
}

func WithType(personType int) func(*GamePerson) {
	return func(person *GamePerson) {
		setField32(&person.other, uint32(personType), 23, 2)
	}
}

func setField32(dst *uint32, value uint32, offset, bits uint) {
	mask := uint32((1<<bits - 1) << offset)
	*dst = (*dst & ^mask) | ((value << offset) & mask)
}

func setField16(dst *uint16, value uint16, offset, bits uint) {
	mask := uint16((1<<bits - 1) << offset)
	*dst = (*dst & ^mask) | ((value << offset) & mask)
}

func getField32(src uint32, offset, bits uint) uint32 {
	mask := uint32((1<<bits - 1) << offset)
	return (src & mask) >> offset
}

func getField16(src uint16, offset, bits uint) uint16 {
	mask := uint16((1<<bits - 1) << offset)
	return (src & mask) >> offset
}

const (
	BuilderGamePersonType = iota
	BlacksmithGamePersonType
	WarriorGamePersonType
)

type GamePerson struct {
	x     int32
	y     int32
	z     int32
	gold  uint32
	other uint32 // тут mana health house weapon family playerTypes
	stats uint16 // 0000| 0000 | 0000| 0000 тут respect strength experience level
	name  [42]byte
}

func NewGamePerson(options ...Option) GamePerson {
	var gamePerson GamePerson
	for _, option := range options {
		option(&gamePerson)
	}

	return gamePerson
}

func (p *GamePerson) Name() string {
	name := p.name[:]
	for i, b := range name {
		if b == 0 {
			return string(name[:i])
		}
	}
	return string(name)
}

func (p *GamePerson) X() int {
	return int(p.x)
}

func (p *GamePerson) Y() int {
	return int(p.y)
}

func (p *GamePerson) Z() int {
	return int(p.z)
}

func (p *GamePerson) Gold() int {
	return int(p.gold)
}

func (p *GamePerson) Mana() int {
	return int(getField32(p.other, 0, 10))
}

func (p *GamePerson) Health() int {
	return int(getField32(p.other, 10, 10))
}

func (p *GamePerson) Respect() int {
	return int(getField16(p.stats, 0, 4))
}

func (p *GamePerson) Strength() int {
	return int(getField16(p.stats, 4, 4))
}

func (p *GamePerson) Experience() int {
	return int(getField16(p.stats, 8, 4))
}

func (p *GamePerson) Level() int {
	return int(getField16(p.stats, 12, 4))
}

func (p *GamePerson) HasHouse() bool {
	return getField32(p.other, 20, 1) != 0
}

func (p *GamePerson) HasGun() bool {
	return getField32(p.other, 21, 1) != 0
}

func (p *GamePerson) HasFamilty() bool {
	return getField32(p.other, 22, 1) != 0
}

func (p *GamePerson) Type() int {
	return int(getField32(p.other, 23, 2))
}

func TestGamePerson(t *testing.T) {
	fmt.Println(unsafe.Sizeof(GamePerson{}))
	assert.LessOrEqual(t, unsafe.Sizeof(GamePerson{}), uintptr(64))

	const x, y, z = math.MinInt32, math.MaxInt32, 0
	const name = "aaaaaaaaaaaaa_bbbbbbbbbbbbb_cccccccccccccc"
	const personType = BuilderGamePersonType
	const gold = math.MaxInt32
	const mana = 1000
	const health = 1000
	const respect = 10
	const strength = 10
	const experience = 10
	const level = 10

	options := []Option{
		WithName(name),
		WithCoordinates(x, y, z),
		WithGold(gold),
		WithMana(mana),
		WithHealth(health),
		WithRespect(respect),
		WithStrength(strength),
		WithExperience(experience),
		WithLevel(level),
		WithHouse(),
		WithFamily(),
		WithType(personType),
	}

	person := NewGamePerson(options...)
	assert.Equal(t, name, person.Name())
	assert.Equal(t, x, person.X())
	assert.Equal(t, y, person.Y())
	assert.Equal(t, z, person.Z())
	assert.Equal(t, gold, person.Gold())
	assert.Equal(t, mana, person.Mana())
	assert.Equal(t, health, person.Health())
	assert.Equal(t, respect, person.Respect())
	assert.Equal(t, strength, person.Strength())
	assert.Equal(t, experience, person.Experience())
	assert.Equal(t, level, person.Level())
	assert.True(t, person.HasHouse())
	assert.True(t, person.HasFamilty())
	assert.False(t, person.HasGun())
	assert.Equal(t, personType, person.Type())
}
