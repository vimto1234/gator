package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/vimto1234/gator/internal/config"
	"github.com/vimto1234/gator/internal/database"
)

func main() {
	gatorConfig, err := config.Read()
	if err != nil {
		fmt.Printf("Error reading config: %v", err)
	}

	db, err := sql.Open("postgres", gatorConfig.DBURL)
	if err != nil {
		fmt.Printf("Error connecting to DB: %v", err)
	}

	dbQueries := database.New(db)

	currentState := state{
		db:      dbQueries,
		configP: &gatorConfig,
	}

	commands := commands{
		commandsMap: make(map[string]func(*state, command) error),
	}

	commands.register("login", handlerLogin)
	commands.register("register", handlerRegister)
	commands.register("reset", handlerClear)
	commands.register("users", handlerGetAllUsers)
	commands.register("agg", handlerAgg)
	commands.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	commands.register("feeds", handlerFeeds)
	commands.register("follow", middlewareLoggedIn(handlerFollow))
	commands.register("following", middlewareLoggedIn(handlerFollowing))

	args := os.Args

	if len(args) < 2 {
		fmt.Print("not enough args provided")
		os.Exit(1)
	}

	commandName := args[1]

	commandArgs := []string{}

	if len(args) >= 3 {
		commandArgs = args[2:]
	}

	commandsToRun := command{
		name: commandName,
		args: commandArgs,
	}

	err = commands.run(&currentState, commandsToRun)
	if err != nil {
		log.Fatal(err)
	}

}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.configP.CurrentUserName)
		if err != nil {
			return err
		}
		return handler(s, cmd, user)
	}
}
