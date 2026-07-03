package main
import _ "github.com/lib/pq"
import (
	"context"
	"log"
	"os"
	"gator.com/m/v2/internal/config"
	"gator.com/m/v2/internal/database"
	"database/sql"
)

type state struct{
	db *database.Queries
	cfg *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}
	db, err := sql.Open("postgres", cfg.DBURL)
	dbQueries := database.New(db)
	if err != nil {
		log.Fatal(err)
	}
	s := &state{
		db: dbQueries,
		cfg: &cfg,
	}

	cmds := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerGetUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddfeed))
	cmds.register("feeds", handlerFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnFollow))
	cmds.register("browse", middlewareLoggedIn(handlerBrowse))

	if len(os.Args) < 2 {
		log.Fatal("not enough arguments")
	}

	cmd := command{
		name: os.Args[1],
		arg: os.Args[2:],
	}

	if err := cmds.run(s, cmd); err != nil {
		log.Fatal(err)
	}
}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return err
		}

		return handler(s, cmd, user)
	}
}
