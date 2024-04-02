package routes

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/AdityaP1502/Instant-Messanging/api/database"
	httpx "github.com/AdityaP1502/Instant-Messanging/api/http"
	"github.com/AdityaP1502/Instant-Messanging/api/http/middleware"
	"github.com/AdityaP1502/Instant-Messanging/api/http/responseerror"
	"github.com/AdityaP1502/Instant-Messanging/api/jsonutil"
	"github.com/AdityaP1502/Instant-Messanging/api/service/games/config"
	"github.com/AdityaP1502/Instant-Messanging/api/service/games/payload"
	"github.com/gorilla/mux"
)

var querynator = database.Querynator{
	DriverName: "postgres",
}

type GenericResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func syncUserGamesHandler(db *sql.DB, conf interface{}, w http.ResponseWriter, r *http.Request) responseerror.HTTPCustomError {
	body := r.Context().Value(middleware.PayloadKey).(*payload.Collections)

	// read username from path parameter
	vars := mux.Vars(r)
	username := vars["username"]

	filterExecutor := querynator.PrepareFilterOperation()
	err := filterExecutor.UseTransaction(db)

	if err != nil {
		return responseerror.CreateInternalServiceError(err)
	}

	userGames := []*payload.UserGames{}

	for _, game := range body.Games {
		userGames = append(userGames, &payload.UserGames{
			Username: username,
			GameID:   game.GameID,
		})
	}

	filterExecutor.AddTableSource("games", "game_id", "game_id")
	filterExecutor.UseExplicitCast()
	err = filterExecutor.BatchInsert(userGames, db, "user_games")

	if err != nil {
		if filterExecutor.Tx != nil {
			filterExecutor.Rollback()
		}

		return responseerror.CreateInternalServiceError(err)
	}

	json, err := jsonutil.EncodeToJson(&GenericResponse{
		Status:  "success",
		Message: "user game collections has been successfully synced",
	})

	if err != nil {
		return responseerror.CreateInternalServiceError(err)
	}

	w.Write(json)
	w.WriteHeader(200)

	filterExecutor.Commit()

	return nil
}

func SetGamesRoute(r *mux.Router, db *sql.DB, conf *config.Config) {
	certMiddleware := middleware.CertMiddleware(conf.RootCAs)

	syncGamesPayloadMiddleware, err := middleware.PayloadCheckMiddleware(
		&payload.Collections{},
		"Games",
	)

	if err != nil {
		log.Fatal(err)
	}

	subrouter := r.PathPrefix("/games").Subrouter()

	syncGames := httpx.CreateHTTPHandler(db, conf, syncUserGamesHandler)
	subrouter.Handle("/{username}/sync", middleware.UseMiddleware(db, conf, syncGames,
		certMiddleware,
		syncGamesPayloadMiddleware,
	))
	// subrouter.Handle("/{username}/sync", middleware.UseMiddleware(db, conf, syncGames, syncGamesPayloadMiddleware)).Methods("POST")

}
