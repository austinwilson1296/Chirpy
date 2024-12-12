package main

import (
	"net/http"
	"github.com/google/uuid"
	"github.com/austinwilson1296/Chirpy/internal/auth"
	"github.com/austinwilson1296/Chirpy/internal/database"
	"database/sql"
)

func (cfg *apiConfig) handlerChirpDelete(w http.ResponseWriter, r *http.Request) {
	chirpIDString := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
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
	chirp, err := cfg.db.GetChirp(r.Context(), chirpID)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, "chirp not found", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "error fetching chirp", err)
		return
	}

	// Check ownership
	if chirp.UserID != id {
		respondWithError(w, http.StatusForbidden, "you don't own this chirp", nil)
		return
	}
	err = cfg.db.DeleteChirp(r.Context(),database.DeleteChirpParams{
		ID: chirpID,
		UserID: id,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			// What status code should go here?
			respondWithError(w, http.StatusNotFound, "chirp not found", err)
			return
		}
		// Handle other types of errors
		respondWithError(w, http.StatusForbidden, "chirp not associated with user", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
	return
}