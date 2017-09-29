package main

import (
	"flag"
	"fmt"

	"github.com/takashabe/go-isucon-exercise/portal/models"
)

// Register Team
func main() {
	var (
		teamID   int
		name     string
		email    string
		password string
		instance string
	)
	flag.IntVar(&teamID, "team_id", 1, "team_id")
	flag.StringVar(&name, "name", "gopher", "team name")
	flag.StringVar(&email, "email", "gopher@example.com", "email use by login")
	flag.StringVar(&password, "password", "gopher", "password use by login")
	flag.StringVar(&instance, "instance", "http://localhost:8000", "webapp instance url")
	flag.Parse()

	err := models.NewTeam().Register(teamID, name, email, password, instance)
	if err != nil {
		panic(err)
	}
	fmt.Printf("succeed register team: %d", teamID)
}
