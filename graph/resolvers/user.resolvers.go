package resolvers

import (
	"context"
	"errors"

	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"

	"time"

	"miikka.xyz/rs/graph/model"
)

func (r *queryResolver) GetUsersWithPrefix(ctx context.Context, prefix string) ([]*model.User, error) {
	// TODO: Add these as params to function
	const pageNumber int = 0
	const pageSize int = 10

	started := time.Now()

	const maxeywordLen = 128
	if len(prefix) > maxeywordLen {
		return nil, errors.New("keyword was too long")
	}

	users, err := r.db.GetUsersWithPrefix(prefix, pageNumber, pageSize)
	if err != nil {
		log.Error("GetUsersWithPrefix failed: ", err)
		return nil, errors.New("fetching users failed")
	}

	ended := time.Now()
	log.Printf("GetUsersWithPrefix [%s] took [%v] to run", prefix, ended.Sub(started))

	return users, nil
}

func (r *mutationResolver) CreateRFIDForUser(ctx context.Context, token string, tag string) (*string, error) {
	username, err := r.db.GetUsernameFromToken(token)
	if err != nil {
		log.Error("validating token failed ", err)
		return nil, errors.New("authentication failed")
	}

	rfid, err := r.db.AddRFID(*username, tag)
	if err != nil {
		log.Error("error while adding rfid ", err)
		return nil, errors.New("rfid failed")
	}
	return rfid, nil
}

// Get tags that are readed by device and free to take
func (r *queryResolver) GetLatestAvailableTags(ctx context.Context, token string, count int) ([]*model.TagAvailable, error) {
	log.Info("latest tags with ", token)
	if _, err := r.db.GetUsernameFromToken(token); err != nil {
		log.Error("validating token failed ", err)
		return nil, errors.New("authentication failed")
	}

	tags, err := r.db.GetLatestTags(3)
	if err != nil {
		log.Error("Error while fetching available tags ", err)
		return nil, errors.New("tags failed")
	}
	return tags, nil
}

// Deprecated, using Github-authentication only. Left as example
// Login compares password with password in database and if succesful write token to
// session table
func (r *mutationResolver) Login(ctx context.Context, username string, password string) (*string, error) {
	const loginErr = "login failed"

	pswd, err := r.db.GetUserPassword(username)
	if err != nil {
		log.Error("error while fetching password ", err)
		return nil, errors.New(loginErr)
	}

	// Comparing the password with the hash
	if err := bcrypt.CompareHashAndPassword([]byte(*pswd), []byte(password)); err != nil {
		log.Error("comparing passwords failed ", err)
		return nil, errors.New(loginErr)
	}

	// Create token
	token, err := bcrypt.GenerateFromPassword([]byte(*pswd), bcrypt.DefaultCost)
	if err != nil {
		log.Error("generating token failed ", err)
		return nil, errors.New(loginErr)
	}

	// Write token to database
	err = r.db.WriteSession(username, string(token))
	if err != nil {
		log.Error("writing token failed ", err)
		return nil, errors.New(loginErr)
	}

	ret := string(token)

	return &ret, nil
}

// Deprecated, using Github-authentication only. Left as example
// Logout removes token from session table
func (r *mutationResolver) Logout(ctx context.Context, token string) (*string, error) {
	err := r.db.DeleteSession(token)
	if err != nil {
		log.Error("logout failed ", err)
		return nil, errors.New("logout failed")
	}

	ret := "Logged out"
	return &ret, nil
}
