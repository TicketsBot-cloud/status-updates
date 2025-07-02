package httpserver

import (
	"database/sql"
	"strings"

	"github.com/TicketsBot-cloud/gdl/objects/channel/message"
	"github.com/TicketsBot-cloud/gdl/objects/interaction"
	"github.com/TicketsBot-cloud/gdl/objects/interaction/component"
	"github.com/TicketsBot-cloud/gdl/rest"
	"github.com/TicketsBot-cloud/status-updates/internal/db"
	"github.com/TicketsBot-cloud/status-updates/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func (s *Server) HandleInteraction(ctx *gin.Context) {
	var body interaction.Interaction
	if err := ctx.ShouldBindBodyWith(&body, binding.JSON); err != nil {
		ctx.JSON(400, errorJson("Failed to parse body"))
		return
	}

	switch body.Type {
	case interaction.InteractionTypePing:
		ctx.JSON(200, interaction.NewResponsePong())
	case interaction.InteractionTypeMessageComponent:
		var commandData interaction.MessageComponentInteraction
		if err := ctx.ShouldBindBodyWith(&commandData, binding.JSON); err != nil {
			_ = ctx.Error(errors.Wrap(err, "failed to parse application command payload"))
			return
		}

		if strings.HasPrefix(commandData.Data.AsButton().CustomId, "incident-role-") {
			incidentId := strings.TrimPrefix(commandData.Data.AsButton().CustomId, "incident-role-")
			var incident model.IncidentInfo

			if err := db.Client.Get(&incident, "SELECT * FROM incidents WHERE id = $1", incidentId); err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					ctx.JSON(404, errorJson("Incident not found"))
					return
				}
				_ = ctx.Error(errors.Wrap(err, "failed to fetch incident from database"))
				return
			}

			// Add incident updates role
			if err := rest.AddGuildMemberRole(ctx, s.config.Discord.Token, nil, s.config.Discord.GuildId, commandData.Member.User.Id, incident.RoleId); err != nil {
				ctx.JSON(500, errorJson("Failed to add role"))
				s.logger.Error("Failed to add role", zap.Error(err))
				return
			}

			// Add to thread
			if err := rest.AddThreadMember(ctx, s.config.Discord.Token, nil, incident.ThreadId, commandData.Member.User.Id); err != nil {
				ctx.JSON(500, errorJson("Failed to add to thread"))
				s.logger.Error("Failed to add to thread", zap.Error(err))
				return
			}

			ctx.JSON(200, interaction.NewResponseChannelMessage(interaction.ApplicationCommandCallbackData{
				Flags: message.SumFlags(message.FlagEphemeral, message.FlagComponentsV2),
				Components: []component.Component{
					component.BuildContainer(component.Container{
						Components: []component.Component{
							component.BuildTextDisplay(component.TextDisplay{
								Content: "You have been added to the incident updates role and thread.",
							}),
						},
					}),
				},
			}))

			return
		}

		ctx.JSON(400, gin.H{"error": "not found"})
	}

}
