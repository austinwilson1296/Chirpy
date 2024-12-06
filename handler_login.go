package main

import (
	"encoding/json"
	"net/http"
	"github.com/austinwilson1296/Chirpy/internal/auth"
)


func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
    type parameters struct {
        Password string `json:"password"`
        Email    string `json:"email"`
    }
    type response struct {
        User
    }

    // First decode the parameters
    decoder := json.NewDecoder(r.Body)
    params := parameters{}
    err := decoder.Decode(&params)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
        return
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
	
	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
	})
}