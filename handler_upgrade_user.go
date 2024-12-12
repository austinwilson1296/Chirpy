package main

import(
	"net/http"
	"encoding/json"
	"database/sql"
	"github.com/google/uuid"
	"github.com/austinwilson1296/Chirpy/internal/database"
	"github.com/austinwilson1296/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerUpgradeUser(w http.ResponseWriter, r *http.Request) {
	key,err := auth.GetAPIKey(r.Header)
		if err != nil{
			respondWithError(w,http.StatusUnauthorized, "no key found in header",err)
			return
		}
	if key != cfg.polkaKey{
		respondWithError(w,http.StatusUnauthorized, "invalid api key",err)
		return
	}

	type parameters struct {
		Event string `json:"event"`
		Data struct {
			UserID uuid.UUID `json:"user_id"`
		}`json:"data"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}
	if params.Event != "user.upgraded"{
		w.WriteHeader(http.StatusNoContent)
		return
	}
	_,err = cfg.db.UpgradeUser(r.Context(),database.UpgradeUserParams{
		IsChirpyRed: sql.NullBool{
			Bool: true,
			Valid: true,
		},
		ID: params.Data.UserID,
	})
	if err != nil{
		respondWithError(w, http.StatusNotFound,"user not found", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
	return
}
