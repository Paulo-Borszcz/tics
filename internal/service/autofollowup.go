package service

import (
	"log"
	"strings"

	"github.com/pauloborszcz/tics/internal/config"
	"github.com/pauloborszcz/tics/internal/glpi"
)

// AutoFollowup checks tickets in "processing (assigned)" status and sends
// the auto-followup HTML if it hasn't been sent yet.
func AutoFollowup(client *glpi.Client, tickets []glpi.Ticket) {
	for _, t := range tickets {
		// Only process tickets with status "Processing (assigned)" = 2
		if t.StatusCode != glpi.StatusProcessing {
			continue
		}

		// Check if the auto-followup was already sent
		followups, err := client.GetFollowups(t.ID)
		if err != nil {
			log.Printf("autofollowup: error getting followups for ticket #%d: %v", t.ID, err)
			continue
		}

		alreadySent := false
		for _, f := range followups {
			if strings.Contains(f.Content, config.AutoFollowupMarker) {
				alreadySent = true
				break
			}
		}

		if alreadySent {
			continue
		}

		// Send the auto-followup
		log.Printf("autofollowup: sending to ticket #%d - %s", t.ID, t.Name)
		if err := client.AddFollowup(t.ID, config.AutoFollowupHTML); err != nil {
			log.Printf("autofollowup: error sending to ticket #%d: %v", t.ID, err)
		} else {
			log.Printf("autofollowup: sent successfully to ticket #%d", t.ID)
		}
	}
}
