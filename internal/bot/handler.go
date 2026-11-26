package bot

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/SoWieMarkus/gokit/httpx"
)

type Handler struct {
	bot *Bot
}

func NewHandler(bot *Bot) *Handler {
	return &Handler{bot: bot}
}

func (h *Handler) Mount(r chi.Router) {
	r.Get("/commands", h.list)
}

type CommandResponse struct {
	Name    string   `json:"name"    example:"guess"`
	Aliases []string `json:"aliases" example:"g"`
	// everyone, subscriber, vip, moderator or broadcaster
	Access          string `json:"access"          example:"everyone" enums:"everyone,subscriber,vip,moderator,broadcaster"`
	CooldownSeconds int64  `json:"cooldownSeconds" example:"0"`
	Usage           string `json:"usage"           example:"!guess 1:23.456"`
	Description     string `json:"description"     example:"Guess the streamer's next personal best on the current map."`
}

// list returns the chat commands the bot understands.
//
//	@Summary	Chat commands
//	@Tags		bot
//	@Produce	json
//	@Success	200	{object}	httpx.ListDoc{items=[]CommandResponse}
//	@Router		/commands [get]
func (h *Handler) list(w http.ResponseWriter, _ *http.Request) {
	cmds := h.bot.Commands()
	items := make([]CommandResponse, len(cmds))
	for i, c := range cmds {
		aliases := c.Aliases
		if aliases == nil {
			aliases = []string{}
		}
		items[i] = CommandResponse{
			Name: c.Name, Aliases: aliases, Access: c.Access.String(),
			CooldownSeconds: int64(c.Cooldown.Seconds()), Usage: c.Usage, Description: c.Description,
		}
	}
	httpx.JSON(w, http.StatusOK, httpx.List[CommandResponse]{Items: items})
}
