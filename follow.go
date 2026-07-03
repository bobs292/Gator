package main

import (
	"fmt"
	"time"
	"github.com/google/uuid"
	"gator.com/m/v2/internal/database"
	"context"
)

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.arg) != 1{
		return fmt.Errorf("Requires a url argument")
	}
	feed, err := s.db.GetFeedsByURL(context.Background(),cmd.arg[0])
	if err != nil {
		return err
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

func handlerFollowing(s *state, cmd command, user database.User) error {
	feedFollows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}
	for _, feedFollow := range feedFollows{
		fmt.Println(feedFollow.FeedName)
	}
	return nil
}

func handlerUnFollow(s *state, cmd command, user database.User) error {
	feed, err := s.db.GetFeedsByURL(context.Background(),cmd.arg[0])
	if err != nil {
		return err
	}

	err = s.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		UserID:		user.ID,
		FeedID:		feed.ID,
	})
	if err != nil {
		return err
	}
	return nil
}