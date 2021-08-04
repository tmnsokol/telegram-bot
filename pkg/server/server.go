package server

import (
	"log"
	"net/http"
	"strconv"

	"github.com/youtube-telegram-bot/youtube-telegram-bot/pkg/telegram/repository"
	"github.com/zhashkevych/go-pocket-sdk"
)

type AuthorizationServer struct {
	server          *http.Server
	pocketClient    *pocket.Client
	tokenRepository repository.TokenRepository
	redirectURL     string
}

func NewAuthorizationServer(pocketClient *pocket.Client,
	tokenRepository repository.TokenRepository,
	redirectURL string) *AuthorizationServer {
	return &AuthorizationServer{pocketClient: pocketClient, tokenRepository: tokenRepository, redirectURL: redirectURL}
}

func (s *AuthorizationServer) Start() error {
	s.server = &http.Server{
		Addr:    ":80",
		Handler: s,
	}

	return s.server.ListenAndServe()
}

func (s *AuthorizationServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	chatIDParam := r.URL.Query().Get("chat_id")
	if chatIDParam == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	chatID, err := strconv.ParseInt(chatIDParam, 10, 64)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	requestToken, err := s.tokenRepository.Get(chatID, repository.RequestTokens)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
	}

	authResponse, err := s.pocketClient.Authorize(r.Context(), requestToken)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	err = s.tokenRepository.Save(chatID, authResponse.AccessToken, repository.AccessTokens)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	log.Printf("chat_id: %d\nrequest_token: %s\naccess_token: %s", chatID, requestToken, authResponse.AccessToken)

	w.Header().Add("Location", s.redirectURL)
	w.WriteHeader(http.StatusMovedPermanently)
}
