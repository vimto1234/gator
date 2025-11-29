package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/vimto1234/gator/internal/config"
	"github.com/vimto1234/gator/internal/database"
)

type state struct {
	db      *database.Queries
	configP *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	commandsMap map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	command, ok := c.commandsMap[cmd.name]
	if !ok {
		return fmt.Errorf("no such command exists")
	}

	return command(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.commandsMap[name] = f
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("login command expects 1 arg")
	}

	loggedInName := cmd.args[0]

	_, err := s.db.GetUser(context.Background(), loggedInName)

	if err != nil {
		return fmt.Errorf("user has not been registered")
	}

	s.configP.SetUser(loggedInName)

	fmt.Println("Username set")

	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("register command expects 1 arg")
	}

	newUserName := cmd.args[0]

	_, err := s.db.GetUser(context.Background(), newUserName)

	if err == nil {
		return fmt.Errorf("user already registered")
	}

	newUser := database.CreateUserParams{
		ID:        uuid.New(),
		Name:      newUserName,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	createdUser, err := s.db.CreateUser(context.Background(), newUser)

	if err != nil {
		return err
	}

	s.configP.SetUser(createdUser.Name)

	fmt.Println("Username set")

	return nil
}

func handlerClear(s *state, cmd command) error {

	err := s.db.ClearUsers(context.Background())

	if err != nil {
		return err
	}

	return nil
}

func handlerGetAllUsers(s *state, cmd command) error {

	allUsers, err := s.db.FetchUsers(context.Background())

	if err != nil {
		return err
	}

	for _, user := range allUsers {
		if user.Name == s.configP.CurrentUserName {
			fmt.Printf("* %v (current)\n", user.Name)
		} else {
			fmt.Printf("* %v\n", user.Name)
		}
	}

	return nil
}

func handlerAgg(s *state, cmd command) error {

	address := "https://www.wagslane.dev/index.xml"

	rssFeed, err := fetchFeed(context.Background(), address)
	if err != nil {
		return err
	}

	fmt.Print(rssFeed)
	return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("add feed command expects 2 args")
	}

	newFeedName := cmd.args[0]
	newFeedURL := cmd.args[1]

	newFeed := database.CreateFeedParams{
		ID:        uuid.New(),
		Name:      newFeedName,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Url:       newFeedURL,
		UserID:    user.ID,
	}

	createdFeed, err := s.db.CreateFeed(context.Background(), newFeed)

	if err != nil {
		return err
	}

	followFeedsParams := database.FollowFeedsParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    createdFeed.ID,
	}

	_, err = s.db.FollowFeeds(context.Background(), followFeedsParams)

	if err != nil {
		return err
	}

	fmt.Printf("Added feed :%v\n", createdFeed.Name)

	return nil
}

func handlerFeeds(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("feeds command expects no args")
	}

	allFeeds, err := s.db.GetFeeds(context.Background())

	if err != nil {
		return err
	}

	for _, user := range allFeeds {
		fmt.Printf("%v\n", user.Name)
		fmt.Printf(" *%v\n", user.Url)
		fmt.Printf(" *%v\n", user.Username)
	}

	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("feeds command expects one arg")
	}

	feedURL := cmd.args[0]

	feed, err := s.db.GetFeedbyURL(context.Background(), feedURL)

	if err != nil {
		return err
	}

	followFeedsParams := database.FollowFeedsParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	addedFollowFeed, err := s.db.FollowFeeds(context.Background(), followFeedsParams)

	if err != nil {
		return err
	}

	fmt.Printf("Added feed to user, name :%v, username:%v\n", addedFollowFeed.FeedName, addedFollowFeed.UserName)

	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("feeds command expects no arg")
	}

	feeds, err := s.db.GetAllFeedsUserFollows(context.Background(), user.ID)

	if err != nil {
		return err
	}

	fmt.Printf("%v is following the below feeds \n", user.Name)

	for _, feed := range feeds {
		println(feed.FeedsName)
	}

	return nil
}
