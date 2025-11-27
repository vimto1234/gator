package main

import (
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
		//fmt.Printf("error :%v\n", err)
		//os.Exit(1)
		log.Fatal(err)
	}
}
