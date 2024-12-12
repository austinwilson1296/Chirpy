package main

import(
	"encoding/json"
	"net/http"
	"github.com/austinwilson1296/Chirpy/internal/auth"
	"github.com/austinwilson1296/Chirpy/internal/database"
)

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	type response struct {
		Email string `json:"email"`
		
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find token", err)
		return
	}

	id,err := auth.ValidateJWT(token,cfg.jwtSecret)
	if err != nil{
		respondWithError(w, http.StatusUnauthorized, "couldn't validate user",err)
		return
	}

	hashedPass,err:= auth.HashPassword(params.Password)
	if err != nil{
		respondWithError(w, http.StatusBadRequest, "unable to hash password", err)
		return
	}
	
	_,err = cfg.db.UpdateUser(r.Context(),database.UpdateUserParams{
		Email: params.Email,
		HashedPassword: hashedPass,
		ID: id,
	})
	if err != nil{
		respondWithError(w, http.StatusBadRequest, "error updating specified user", err)
		return
	}

	respondWithJSON(w, http.StatusOK,response{
		Email: params.Email,
	})
}