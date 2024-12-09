package main

import (
	"encoding/json"
	"net/http"
	"github.com/austinwilson1296/Chirpy/internal/auth"
	"time"
)


func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
    type parameters struct {
        Password string `json:"password"`
        Email    string `json:"email"`
		Expires time.Duration `json:"expires_in_seconds"`
    }
    type response struct {
        User
		Token string `json:"token"`
    }

    // First decode the parameters
    decoder := json.NewDecoder(r.Body)
    params := parameters{}
    err := decoder.Decode(&params)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
        return
    }

	params.Expires = time.Duration(params.Expires) * time.Second

	if params.Expires == 0{
		params.Expires = 1 * time.Hour
	}else if params.Expires > time.Hour * 1 {
		params.Expires = time.Hour * 1
	}
	
    // Then get the user
    user, err := cfg.db.GetUser(r.Context(), params.Email)
    if err != nil {
        respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", nil)
        return
    }

    // Now what would you add here to:
    // 1. Check the password hash
	err = auth.CheckPasswordHash(user.HashedPassword, params.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", nil)
		return
	}

	key,err := auth.MakeJWT(user.ID,cfg.tokenSecret,params.Expires)
	if err != nil{
		respondWithError(w, http.StatusUnauthorized, "error generating token", nil)
		return 
	}
	
	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
			
		},
		Token: key,
		
	})
}