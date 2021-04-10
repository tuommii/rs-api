package github

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	log "github.com/sirupsen/logrus"

	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"miikka.xyz/rs/config"
	"miikka.xyz/rs/store"
)

const (
	cookieUsername = "user_id"
	cookieSession  = "session"
)

// Auth injects database
type Auth struct {
	DB *store.DB
}

// GetStartHandler returns function that handles starting the authentication
func (a *Auth) GetStartHandler(config *config.Config) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// If user has valid existing cookie don't send request to GitHub
		if cookie, _ := r.Cookie(cookieSession); cookie != nil {
			username, err := a.DB.GetUsernameFromToken(cookie.Value)
			if err != nil {
				w.Write([]byte("invalid token in cookie"))
				return
			}

			log.Info("user ", *username, " is already authenticated (not secure yet)")
			http.Redirect(w, r, config.OauthRedirectURL, http.StatusTemporaryRedirect)
			return
		}

		// Otherwise go to GitHub
		log.Info("starting authentication with GitHub")
		query := "?client_id=" + config.GithubID + "&state=" + config.GithubState + "&redirect_uri=" + config.GithubCallbackURL
		http.Redirect(w, r, "https://github.com/login/oauth/authorize"+query, http.StatusTemporaryRedirect)
	})
}

// GetCallbackHandler returns function that handles getting info about user
func (a *Auth) GetCallbackHandler(config *config.Config) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Info("URL params: ", r.URL.Query())

		if err := checkState(r, config.GithubState); err != nil {
			w.Write([]byte("error"))
			return
		}

		var token string
		if err := getAccessToken(r, &token, config.GithubID, config.GithubSecret); err != nil {
			w.Write([]byte("error"))
			return
		}

		user := User{}
		if err := getUserInfo(r, token, &user); err != nil {
			w.Write([]byte("error"))
			return
		}

		// ===================================================================
		// Now we are done with Github-part, rest is for internal use
		// ===================================================================

		log.Printf("user logged in, %+v\n", user)

		isUsernameAvailable := a.DB.IsUsernameAvailable(user.Login)
		log.Info("username: ", user.Login, " available:", isUsernameAvailable)
		// TODO: Handle if Github name has not " "
		names := strings.SplitN(user.Name, " ", 2)
		if isUsernameAvailable {
			id, err := a.DB.CreateUser(user.Login, names[0], names[1], user.Email)
			if err != nil {
				log.Error("error while creating new user ", err)
			} else {
				log.Info("User ", user.Login, " created with id: ", id)
			}
		}

		var hash string
		if err := saveToken(r, token, &user, a.DB, &hash); err != nil {
			log.Error(err)
			return
		}

		expires := time.Now().Add(time.Minute * 30)
		sessionCookie := &http.Cookie{
			Name:    cookieSession,
			Value:   hash,
			Expires: expires,
		}
		usernameCookie := &http.Cookie{
			Name:    cookieUsername,
			Value:   user.Login,
			Expires: expires,
		}
		http.SetCookie(w, sessionCookie)
		http.SetCookie(w, usernameCookie)
		http.Redirect(w, r, config.OauthRedirectURL, http.StatusTemporaryRedirect)
	})
}

// ===========================================================================
// HELPER FUNCTIONS
// ===========================================================================

func checkState(r *http.Request, wanted string) error {
	got := r.URL.Query().Get("state")
	if got != wanted {
		log.Error("state mismatch")
		return errors.New("state mismatch")
	}
	return nil
}

func getAccessToken(r *http.Request, token *string, clientID string, secret string) error {
	code := r.URL.Query().Get("code")
	tokenQuery := "?client_id=" + clientID + "&client_secret=" + secret + "&code=" + code
	resp, err := http.Post("https://github.com/login/oauth/access_token"+tokenQuery, "application/json", nil)
	if err != nil {
		log.Error("post failed ", err)
		return errors.New("post failed")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("reading body failed ", err)
		return errors.New("reading body failed")
	}

	values, err := url.ParseQuery(string(body))
	if err != nil {
		log.Error("parsing access token failed ", err)
		return errors.New("parsing token failed")
	}

	*token = values.Get("access_token")
	return nil
}

func getUserInfo(r *http.Request, token string, user *User) error {
	userInfoReq, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		log.Error("creating user info request failed ", err)
		return errors.New("creating user info request failed")
	}

	userInfoReq.Header.Set("Authorization", "token "+token)
	client := http.Client{}
	resp, err := client.Do(userInfoReq)
	if err != nil {
		log.Error("reading user response failed ", err)
		return errors.New("reading user response failed")
	}
	defer resp.Body.Close()

	// Read user info response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("reading user body failed ", err)
		return errors.New("reading user response body failed")
	}

	err = json.Unmarshal(body, user)
	if err != nil {
		log.Error("parsing user json failed ", err)
		return errors.New("parsing user json failed")
	}
	return nil
}

func saveToken(r *http.Request, token string, user *User, db *store.DB, ret *string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Login+user.Email), bcrypt.DefaultCost)
	if err != nil {
		log.Error("generating token failed ", err)
		return errors.New("generating token failed")
	}
	*ret = string(hash)
	err = db.WriteSession(user.Login, *ret)
	if err != nil {
		log.Error("saving session failed ", err)
		return errors.New("saving token failed")
	}
	return nil
}
