package api

import (
	"fmt"
	"net/http"

	"github.com/SoWieMarkus/wheelgpt/internal/auth"
	"github.com/SoWieMarkus/wheelgpt/internal/db"
	"github.com/google/uuid"
)

func (s *Server) handleTwitchRedirect(w http.ResponseWriter, r *http.Request) {
	params := fmt.Sprintf(
		"client_id=%s&redirect_uri=%s&response_type=code&scope=user:read:email",
		s.cfg.TwitchClientID, s.cfg.TwitchRedirectURL,
	)
	http.Redirect(w, r, "https://id.twitch.tv/oauth2/authorize?"+params, http.StatusFound)
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Code string `json:"code"`
	}
	if err := decodeJSON(r, &body); err != nil || body.Code == "" {
		jsonError(w, http.StatusBadRequest, "missing code")
		return
	}

	ctx := r.Context()
	userToken, err := s.tc.ExchangeCode(ctx, body.Code)
	if err != nil {
		jsonError(w, http.StatusBadGateway, "token exchange failed")
		return
	}

	twitchUser, err := s.tc.GetUser(ctx, userToken)
	if err != nil {
		jsonError(w, http.StatusBadGateway, "user lookup failed")
		return
	}

	isNew := false
	_, err = s.store.GetChannel(ctx, twitchUser.ID)
	if err != nil {
		isNew = true
	}

	if err := s.store.UpsertChannel(ctx, db.UpsertChannelParams{
		ID:           twitchUser.ID,
		Login:        twitchUser.Login,
		Token:        uuid.NewString(),
		DisplayName:  twitchUser.DisplayName,
		ProfileImage: twitchUser.ProfileImageURL,
	}); err != nil {
		jsonError(w, http.StatusInternalServerError, "db error")
		return
	}

	if isNew {
		ch, _ := s.store.GetChannel(ctx, twitchUser.ID)
		s.bot.Join(botCfg(ch))
		if stream, err := s.tc.GetStream(ctx, twitchUser.ID); err == nil {
			live := 0
			if stream.Type == "live" {
				live = 1
			}
			_ = s.store.SetChannelLive(ctx, db.SetChannelLiveParams{
				IsLive: int64(live),
				ID:     twitchUser.ID,
			})
		}
		if s.cfg.UpdateWebHooks {
			_ = s.tc.AddWebhooks(ctx, twitchUser.ID, s.cfg.EventSubCallbackURL, s.cfg.EventSubSecret)
		}
	}

	webToken, err := auth.IssueWebToken(twitchUser.ID, s.cfg.JWTSecretWeb)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "token issue failed")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"webToken": webToken})
}

func (s *Server) handleRenewToken(w http.ResponseWriter, r *http.Request) {
	channelID := channelIDFrom(r.Context())
	newToken := uuid.NewString()
	if err := s.store.RotateChannelToken(r.Context(), db.RotateChannelTokenParams{
		Token: newToken,
		ID:    channelID,
	}); err != nil {
		jsonError(w, http.StatusInternalServerError, "db error")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"token": newToken})
}
