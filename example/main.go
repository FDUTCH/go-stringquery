package main

import (
	"fmt"
	"go-stringquery/query"
)

func main() {
	formatter := query.NewFormatter[*Player]().
		WithOption("health", playerHealth).
		WithOption("name", playerName).
		WithOption("food", playerFood).
		WithOption("strength", playerStrength).
		WithOption("speed", playerSpeed).
		WithOption("regeneration", playerRegeneration)

	p := &Player{
		health:       100,
		name:         "Steve",
		food:         10,
		strength:     2,
		speed:        1.75,
		regeneration: 7.5,
	}

	textQuery := "your health: ${health}"
	statsQuery := "Strength: ${strength}, Speed: ${speed}, Regeneration speed: ${regeneration}"
	jsonQuery := `{"name":"${name}","health":${health},"food":${food}}`

	fmt.Println(formatter.Query(textQuery, p))
	fmt.Println(formatter.Query(statsQuery, p))
	fmt.Println(formatter.Query(jsonQuery, p))
}

func playerHealth(p *Player) any {
	return p.Health()
}

func playerName(p *Player) any {
	return p.Name()
}

func playerFood(p *Player) any {
	return p.Food()
}

func playerStrength(p *Player) any {
	return p.Strength()
}
func playerSpeed(p *Player) any {
	return p.Speed()
}

func playerRegeneration(p *Player) any {
	return p.Regeneration()
}

type Player struct {
	health       float64
	name         string
	food         int
	strength     int
	speed        float64
	regeneration float64
}

func (p *Player) Regeneration() float64 {
	return p.regeneration
}

func (p *Player) Strength() int {
	return p.strength
}

func (p *Player) Speed() float64 {
	return p.speed
}

func (p *Player) Health() float64 {
	return p.health
}

func (p *Player) Name() string {
	return p.name
}

func (p *Player) Food() int {
	return p.food
}
