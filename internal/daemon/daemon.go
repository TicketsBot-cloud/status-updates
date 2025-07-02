package daemon

import (
	"context"
	"fmt"
	"time"

	"github.com/TicketsBot-cloud/gdl/objects/channel"
	"github.com/TicketsBot-cloud/gdl/objects/channel/message"
	"github.com/TicketsBot-cloud/gdl/objects/interaction/component"
	"github.com/TicketsBot-cloud/gdl/rest"
	"github.com/TicketsBot-cloud/status-updates/internal/config"
	"github.com/TicketsBot-cloud/status-updates/internal/model"
	"github.com/TicketsBot-cloud/status-updates/internal/statuspage"
	"go.uber.org/zap"
)

type Daemon struct {
	logger           *zap.Logger
	config           config.Config
	statusPageClient statuspage.StatusPageClient
}

func NewDaemon(logger *zap.Logger, spc statuspage.StatusPageClient) *Daemon {
	return &Daemon{
		logger:           logger,
		config:           config.Conf,
		statusPageClient: spc,
	}
}

func (d *Daemon) Start(ctx context.Context) error {
	d.logger.Info("Starting daemon")

	// Run once immediately to avoid waiting for the first timer tick
	if err := d.runOnce(ctx); err != nil {
		d.logger.Error("Failed to run initial check", zap.Error(err))
	}

	timer := time.NewTimer(d.config.Daemon.Frequency)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			start := time.Now()
			d.logger.Info("Run started", zap.Time("start_time", start))
			if err := d.runOnce(ctx); err != nil {
				d.logger.Error("Failed to run", zap.Error(err))
			}
			d.logger.Info("Run completed", zap.Time("end_time", time.Now()), zap.Duration("duration", time.Since(start)))
			timer.Reset(d.config.Daemon.Frequency)
		case <-ctx.Done():
			d.logger.Info("Shutting down daemon")
			return nil

		}
	}
}

func (d *Daemon) runOnce(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, d.config.Daemon.ExecutionTimeout)
	defer cancel()

	incidents, err := d.statusPageClient.GetIncidents()
	if err != nil {
		d.logger.Error("Failed to fetch incidents", zap.Error(err))
		return err
	}

	for _, incident := range incidents {
		exists, err := incident.Exists()
		if err != nil {
			d.logger.Error("Failed to check if incident exists", zap.Error(err))
			continue
		}

		// Order updates in reverse order
		incident.OrderUpdates()
		container := incident.GenerateContainer()
		msgComponents := []component.Component{
			component.BuildTextDisplay(component.TextDisplay{
				Content: fmt.Sprintf("-# A new incident has been reported <@&%d>", d.config.Discord.UpdateRoleId),
			}),
			container,
		}

		if !exists {
			d.logger.Info("New incident detected. Sending Discord message...", zap.String("incident_id", incident.ID), zap.String("status", incident.Status))
			msg, err := rest.CreateMessage(ctx, d.config.Discord.Token, nil, d.config.Discord.ChannelId, rest.CreateMessageData{
				Components: msgComponents,
				Flags:      message.SumFlags(message.FlagComponentsV2),
				AllowedMentions: message.AllowedMention{
					Roles: []uint64{d.config.Discord.UpdateRoleId},
				},
			})
			if err != nil {
				d.logger.Error("Error sending message", zap.Error(err))
				continue
			}

			channelInfo, err := rest.GetChannel(ctx, d.config.Discord.Token, nil, d.config.Discord.ChannelId)
			if err != nil {
				d.logger.Error("Error retrieving channel info", zap.Error(err))
				continue
			}

			if channelInfo.Type == channel.ChannelTypeGuildNews && config.Conf.Discord.ShouldCrosspost {
				if err := rest.CrosspostMessage(ctx, d.config.Discord.Token, nil, d.config.Discord.ChannelId, msg.Id); err != nil {
					d.logger.Error("Error crossposting message", zap.Error(err))
				}
			}

			d.logger.Info("Discord message sent for incident", zap.String("incident_id", incident.ID), zap.Uint64("message_id", msg.Id))

			// Create role & thread
			role, err := rest.CreateGuildRole(ctx, d.config.Discord.Token, nil, d.config.Discord.GuildId, rest.GuildRoleData{
				Name: fmt.Sprintf("Incident Updates: %s", incident.ID),
			})
			if err != nil {
				d.logger.Error("Error creating role", zap.Error(err))
				continue
			}

			thread, err := rest.StartThreadWithMessage(ctx, d.config.Discord.Token, nil, d.config.Discord.ChannelId, msg.Id, rest.StartThreadWithMessageData{
				Name:                fmt.Sprintf("Incident Updates: %s", incident.ID),
				AutoArchiveDuration: 1440, // 24 hours
			})
			if err != nil {
				d.logger.Error("Error starting thread", zap.Error(err))
				continue
			}

			incidentInfo := model.IncidentInfo{
				Id:            incident.ID,
				RoleId:        role.Id,
				MessageId:     msg.Id,
				ThreadId:      thread.Id,
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
				CurrentStatus: incident.Status,
			}

			if err := incidentInfo.Save(); err != nil {
				d.logger.Error("Error saving incident info", zap.Error(err))
				continue
			}
			d.logger.Info("Incident info saved", zap.String("incident_id", incident.ID), zap.Uint64("role_id", role.Id), zap.Uint64("message_id", msg.Id), zap.Uint64("thread_id", thread.Id))

		} else {
			incidentInfo, err := incident.Get()
			if err != nil {
				d.logger.Error("Error retrieving incident info", zap.Error(err))
				continue
			}
			if incident.IncidentUpdates[len(incident.IncidentUpdates)-1].DisplayAt.After(incidentInfo.UpdatedAt) {
				d.logger.Info("Update detected for incident. Editing Discord message...",
					zap.String("incident_id", incident.ID),
					zap.Uint64("message_id", incidentInfo.MessageId),
				)
				// Update the message if the last update is newer
				_, err := rest.EditMessage(ctx, d.config.Discord.Token, nil, d.config.Discord.ChannelId, incidentInfo.MessageId, rest.EditMessageData{
					Components: msgComponents,
					Flags:      message.SumFlags(message.FlagComponentsV2),
				})
				if err != nil {
					d.logger.Error("Error editing message", zap.Error(err))
					continue
				}

				d.logger.Info("Discord message updated for incident", zap.String("incident_id", incident.ID))

				// Send update to thread
				mostRecentUpdate := incident.IncidentUpdates[len(incident.IncidentUpdates)-1]
				updateContainer := incident.GenerateUpdateContainer(mostRecentUpdate)
				_, err = rest.CreateMessage(ctx, d.config.Discord.Token, nil, incidentInfo.ThreadId, rest.CreateMessageData{
					Components: []component.Component{
						component.BuildTextDisplay(component.TextDisplay{
							Content: fmt.Sprintf("-# A new update has been posted <@&%d>", incidentInfo.RoleId),
						}),
						updateContainer,
					},
					Flags: message.SumFlags(message.FlagComponentsV2),
					AllowedMentions: message.AllowedMention{
						Roles: []uint64{incidentInfo.RoleId},
					},
				})
				if err != nil {
					d.logger.Error("Error creating message in thread", zap.Error(err))
					continue
				}

				d.logger.Info("Update message sent in thread", zap.String("incident_id", incident.ID), zap.Uint64("thread_id", incidentInfo.ThreadId))

				// Check if its resolved, if it is, close everything down
				if incident.Status == "resolved" || incident.Status == "completed" {
					d.logger.Info("Incident resolved, closing thread and removing role", zap.String("incident_id", incident.ID))
					archive := true

					// Close the thread
					if _, err := rest.ModifyChannel(ctx, d.config.Discord.Token, nil, incidentInfo.ThreadId, rest.ModifyChannelData{
						ThreadMetadataModifyData: &rest.ThreadMetadataModifyData{
							Archived: &archive,
							Locked:   &archive,
						},
					}); err != nil {
						d.logger.Error("Error closing thread", zap.Error(err))
					}

					// Delete the role
					if err := rest.DeleteGuildRole(ctx, d.config.Discord.Token, nil, d.config.Discord.GuildId, incidentInfo.RoleId); err != nil {
						d.logger.Error("Error deleting role", zap.Error(err))
					}

					d.logger.Info("Thread closed and role deleted for incident", zap.String("incident_id", incident.ID))
				}

				incidentInfo.CurrentStatus = incident.Status
				incidentInfo.UpdatedAt = time.Now()

				if err := incidentInfo.Save(); err != nil {
					d.logger.Error("Error saving incident info", zap.Error(err))
					continue
				}
			}
		}

		continue
	}

	return nil
}
