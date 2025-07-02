package model

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/TicketsBot-cloud/gdl/objects/interaction/component"
	"github.com/TicketsBot-cloud/status-updates/internal/db"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Incident represents the structure of the given JSON object

// Incident represents the structure of the provided JSON
type Incident struct {
	ID                                        string           `json:"id"`
	Components                                []Component      `json:"components"`
	CreatedAt                                 time.Time        `json:"created_at"`
	Impact                                    string           `json:"impact"`
	ImpactOverride                            string           `json:"impact_override"`
	IncidentUpdates                           []IncidentUpdate `json:"incident_updates"`
	Metadata                                  Metadata         `json:"metadata"`
	MonitoringAt                              time.Time        `json:"monitoring_at"`
	Name                                      string           `json:"name"`
	PageID                                    string           `json:"page_id"`
	PostmortemBody                            string           `json:"postmortem_body"`
	PostmortemBodyLastUpdatedAt               time.Time        `json:"postmortem_body_last_updated_at"`
	PostmortemIgnored                         bool             `json:"postmortem_ignored"`
	PostmortemNotifiedSubscribers             bool             `json:"postmortem_notified_subscribers"`
	PostmortemNotifiedTwitter                 bool             `json:"postmortem_notified_twitter"`
	PostmortemPublishedAt                     time.Time        `json:"postmortem_published_at"`
	ResolvedAt                                time.Time        `json:"resolved_at"`
	ScheduledAutoCompleted                    bool             `json:"scheduled_auto_completed"`
	ScheduledAutoInProgress                   bool             `json:"scheduled_auto_in_progress"`
	ScheduledFor                              time.Time        `json:"scheduled_for"`
	AutoTransitionDeliverNotificationsAtEnd   bool             `json:"auto_transition_deliver_notifications_at_end"`
	AutoTransitionDeliverNotificationsAtStart bool             `json:"auto_transition_deliver_notifications_at_start"`
	AutoTransitionToMaintenanceState          bool             `json:"auto_transition_to_maintenance_state"`
	AutoTransitionToOperationalState          bool             `json:"auto_transition_to_operational_state"`
	ScheduledRemindPrior                      bool             `json:"scheduled_remind_prior"`
	ScheduledRemindedAt                       time.Time        `json:"scheduled_reminded_at"`
	ScheduledUntil                            time.Time        `json:"scheduled_until"`
	Shortlink                                 string           `json:"shortlink"`
	Status                                    string           `json:"status"`
	UpdatedAt                                 time.Time        `json:"updated_at"`
	ReminderIntervals                         string           `json:"reminder_intervals"` // JSON string representation of intervals
}

// IncidentUpdate represents the structure of an incident update
type IncidentUpdate struct {
	ID                   string              `json:"id"`
	IncidentID           string              `json:"incident_id"`
	AffectedComponents   []AffectedComponent `json:"affected_components"`
	Body                 string              `json:"body"`
	CreatedAt            time.Time           `json:"created_at"`
	CustomTweet          string              `json:"custom_tweet"`
	DeliverNotifications bool                `json:"deliver_notifications"`
	DisplayAt            time.Time           `json:"display_at"`
	Status               string              `json:"status"`
	TweetID              string              `json:"tweet_id"`
	TwitterUpdatedAt     time.Time           `json:"twitter_updated_at"`
	UpdatedAt            time.Time           `json:"updated_at"`
	WantsTwitterUpdate   bool                `json:"wants_twitter_update"`
}

func (iu IncidentUpdate) AsString() string {
	return fmt.Sprintf("[<t:%d:t>] %s", iu.DisplayAt.Unix(), iu.Body)
}

func (i Incident) GenerateContainer() component.Component {
	statusCaser := cases.Title(language.English)
	color := i.GetColor()
	var msgFormat string
	for _, update := range i.IncidentUpdates {
		msgFormat += fmt.Sprintf("%s\n\n", update.AsString())
	}
	buttons := []component.Component{component.BuildButton(component.Button{
		Label: "Status Page",
		Style: component.ButtonStyleLink,
		Url:   &i.Shortlink,
	})}
	if i.Status != "resolved" && i.Status != "completed" {
		buttons = append(buttons, component.BuildButton(component.Button{
			Label:    "Receive Updates",
			Style:    component.ButtonStyleSecondary,
			CustomId: fmt.Sprintf("incident-role-%s", i.ID),
		}))
	}
	return component.BuildContainer(component.Container{
		Components: []component.Component{
			component.BuildTextDisplay(component.TextDisplay{
				Content: fmt.Sprintf("## %s - %s", i.GetSeverity(), i.Name),
			}),
			component.BuildSeparator(component.Separator{}),
			component.BuildTextDisplay(component.TextDisplay{
				Content: msgFormat,
			}),
			component.BuildTextDisplay(component.TextDisplay{
				Content: "Status: " + statusCaser.String(i.Status),
			}),
			component.BuildActionRow(buttons...),
		},
		AccentColor: &color,
	})
}

func (i Incident) GetSeverity() string {
	severity := ""
	if len(i.Components) == 0 {
		return "Unknown"
	}
	switch i.Components[0].Status {
	case "partial_outage":
		severity = "Partial Outage"
	case "major_outage":
		severity = "Major Outage"
	case "degraded_performance":
		severity = "Degraded Performance"
	case "operational":
		severity = "Operational"
	default:
		severity = "Unknown"
	}

	return severity
}

func (i Incident) GetColor() int {
	color := 0x00CD00
	if len(i.Components) == 0 {
		return color
	}
	switch i.Components[0].Status {
	case "partial_outage":
		color = 0xFFA500
	case "major_outage":
		color = 0xFF0000
	case "degraded_performance":
		color = 0xFFA500
	case "operational":
		color = 0x00CD00
	default:
		color = 0x00CD00
	}

	return color
}

func (i Incident) GenerateUpdateContainer(u IncidentUpdate) component.Component {
	color := i.GetColor()
	return component.BuildContainer(component.Container{
		Components: []component.Component{
			component.BuildTextDisplay(component.TextDisplay{
				Content: fmt.Sprintf("## %s - %s", i.GetSeverity(), i.Name),
			}),
			component.BuildSeparator(component.Separator{}),
			component.BuildTextDisplay(component.TextDisplay{
				Content: fmt.Sprintf("[<t:%d:t>] %s", u.DisplayAt.Unix(), u.Body),
			}),
			component.BuildTextDisplay(component.TextDisplay{
				Content: "Status: " + cases.Title(language.English).String(u.Status),
			}),
		},
		AccentColor: &color,
	})
}

func (i *Incident) OrderUpdates() {
	newIncidentUpdates := make([]IncidentUpdate, len(i.IncidentUpdates))
	for j, k := len(i.IncidentUpdates)-1, 0; j >= 0; j, k = j-1, k+1 {
		newIncidentUpdates[k] = i.IncidentUpdates[j]
	}
	i.IncidentUpdates = newIncidentUpdates
}

func (i Incident) Exists() (bool, error) {
	var exists bool
	if err := db.Client.Get(&exists, "SELECT EXISTS(SELECT 1 FROM incidents WHERE id = $1)", i.ID); err != nil {
		fmt.Printf("Error checking if incident exists: %v\n", err)
		return false, err
	}

	return exists, nil
}

func (i Incident) Get() (IncidentInfo, error) {
	var info IncidentInfo
	err := db.Client.Get(&info, "SELECT * FROM incidents WHERE id = $1", i.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return IncidentInfo{}, fmt.Errorf("incident with ID %s not found", i.ID)
		}
		fmt.Printf("Error retrieving incident info: %v\n", err)
		return IncidentInfo{}, err
	}

	return info, nil
}
