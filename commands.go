package main

import(
	"context"
	"fmt"
	"errors"
	"time"
	"github.com/google/uuid"
	"gator.com/m/v2/internal/database"
)

type command struct{
	name string
	arg []string
}

type commands struct {
	registeredCommands map[string]func(*state, command) error
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arg) != 1 {
		return errors.New("Add username")
	}
	user, err := s.db.GetUser(context.Background(),cmd.arg[0])
	if err != nil{
		return err
	}
	ok := s.cfg.SetUser(user.Name)
	if ok != nil{
		return ok
	}
	fmt.Println("Login success for", cmd.arg[0])
	return nil
}

func (c *commands) run(s *state, cmd command) error {
	handler, ok := c.registeredCommands[cmd.name]
	if !ok {
		return fmt.Errorf("unknown command: %s", cmd.name)
	}

	return handler(s, cmd)
}
func (c *commands) register(name string, f func(*state, command) error) {
	c.registeredCommands[name] = f
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.arg) != 1 {
		return fmt.Errorf("usage: %v <name>", cmd.name)
	}

	name := cmd.arg[0]

	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      name,
	})

	if err != nil {
		return fmt.Errorf("couldn't create user: %w", err)
	}

	fmt.Println(user.ID, user.Name)
	s.cfg.SetUser(user.Name)
	fmt.Println("The users has been created and set")
	return nil
}

func handlerReset(s *state, cmd command) error{
	err := s.db.DeleteAllUsers(context.Background())
	if err != nil {
		return err
	}
	fmt.Println("All users deleted")
	return nil
}
func handlerGetUsers(s *state, cmd command) error{
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return nil
	}
	for _, user := range users{
		fmt.Print(user.Name)
		if user.Name == s.cfg.CurrentUserName{
			fmt.Print(" (current)")
		}
		fmt.Println("")
	}
	return nil
}

func handlerAddfeed(s *state, cmd command, user database.User) error {
	if len(cmd.arg) != 2 {
		return fmt.Errorf("Requires name of the feed and URL")
	}
	feed, err := s.db.Createfeed(context.Background(), database.CreatefeedParams{
		ID: 		uuid.New(),
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
		Name: 		cmd.arg[0], 
		Url: 		cmd.arg[1], 
		UserID: 	user.ID,
		})
	if err != nil {
		return fmt.Errorf("Couldn't create feed: %w", err)
	}
	
	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        	uuid.New(),
		CreatedAt: 	time.Now().UTC(),
		UpdatedAt: 	time.Now().UTC(),
		UserID:    	user.ID,
		FeedID:		feed.ID,
	})
	if err != nil {
		return err
	}
	fmt.Println("%v is now following %v", user.Name, feed.Name)
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	if len(cmd.arg) != 0 {
		return fmt.Errorf("Requires no arg values")
	}
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil{
		return fmt.Errorf("Couldn't fetch feeds list: %w", err)
	}
	for _, feed := range feeds {
		user, err := s.db.GetUserbyId(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("Couldn't get user name: %w", err)
		}
		fmt.Println(feed.Name, feed.Url, user.Name)
	}
	return nil
}