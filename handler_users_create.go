package main

import (
	"encoding/json"
	"net/http"
	"time"
	"github.com/austinwilson1296/Chirpy/internal/auth"
	"github.com/austinwilson1296/Chirpy/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	IsChirpyRed bool    `json:"is_chirpy_red"`
}

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email string `json:"email"`
	}
	type response struct {
		User
	}


	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	if params.Password == ""{
		respondWithError(w, http.StatusBadRequest, "password required", nil)
		return
	}

	hashedPassword, err:= auth.HashPassword(params.Password)
	if err != nil{
		respondWithError(w, http.StatusBadRequest, "could not set password",err)
		return
	}

	

	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		HashedPassword: hashedPassword,
		Email: params.Email,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
	})
}