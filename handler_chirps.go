package main

import (
	"encoding/json"
	"net/http"
	"github.com/google/uuid"
	"time"
	"github.com/austinwilson1296/Chirpy/internal/database"
	
)
type Chirp struct{
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body     string    `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerChirpCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}
	type returnVals struct {
		Chirp
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}
	chirp,err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body: params.Body,
		UserID: params.UserID,
	})
	if err != nil{
		respondWithError(w, http.StatusBadRequest, "Error creating chirp",err)
		return
	}


	respondWithJSON(w, http.StatusCreated, returnVals{
		Chirp: Chirp{
			ID: chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body: chirp.Body,
			UserID: params.UserID,
		},
	})
}
func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
    dbChirps, err := cfg.db.SelectAllChirps(r.Context())
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Error retrieving chirps", err)
        return
    }

    chirps := []Chirp{}
    for _, dbChirp := range dbChirps {
        chirps = append(chirps, Chirp{
            ID:        dbChirp.ID,
            CreatedAt: dbChirp.CreatedAt,
            UpdatedAt: dbChirp.UpdatedAt,
            Body:      dbChirp.Body,
            UserID:    dbChirp.UserID,
        })
    }

    respondWithJSON(w, http.StatusOK, chirps)
}


func (cfg *apiConfig) handlerGetSingleChirp(w http.ResponseWriter, r *http.Request){
	paramID := r.PathValue("chirpID")
	
	parseUUID,err := uuid.Parse(paramID)
	if err != nil {
		respondWithError(w,http.StatusNotFound,"invalid path", err)
	}
	type returnVals struct {
		Chirp
	}
	
	chirp,err := cfg.db.SelectSingleChirp(r.Context(),parseUUID)
	if err != nil{
		respondWithError(w, http.StatusNotFound, "Error retrieving chirp", err)
        return
	}
	respondWithJSON(w, http.StatusOK, returnVals{
		Chirp: Chirp{
			ID: chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body: chirp.Body,
			UserID: chirp.UserID,
		},
	})

}