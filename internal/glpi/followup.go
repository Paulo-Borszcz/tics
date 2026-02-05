package glpi

import (
	"encoding/json"
	"fmt"
)

type Followup struct {
	ID           int    `json:"id"`
	Content      string `json:"content"`
	DateCreation string `json:"date_creation"`
	UsersID      int    `json:"users_id"`
	IsPrivate    int    `json:"is_private"`
}

// GetFollowups returns all followups for a given ticket.
func (c *Client) GetFollowups(ticketID int) ([]Followup, error) {
	data, err := c.get(fmt.Sprintf("/Ticket/%d/ITILFollowup", ticketID))
	if err != nil {
		return nil, fmt.Errorf("getting followups for ticket %d: %w", ticketID, err)
	}

	var followups []Followup
	if err := json.Unmarshal(data, &followups); err != nil {
		return nil, fmt.Errorf("decoding followups: %w", err)
	}

	return followups, nil
}

// AddFollowup adds a new followup to a ticket.
func (c *Client) AddFollowup(ticketID int, content string) error {
	payload := fmt.Sprintf(`{"input":{"items_id":%d,"itemtype":"Ticket","content":%q}}`, ticketID, content)
	_, err := c.post("/ITILFollowup/", payload)
	if err != nil {
		return fmt.Errorf("adding followup to ticket %d: %w", ticketID, err)
	}
	return nil
}
