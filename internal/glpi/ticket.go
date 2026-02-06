package glpi

import (
	"encoding/json"
	"fmt"
	"html"
	"net/url"
)

type Ticket struct {
	ID           int             `json:"id"`
	Name         string          `json:"name"`
	Content      string          `json:"content"`
	RawStatus    json.RawMessage `json:"status"`
	Status       string          `json:"-"`
	StatusCode   int             `json:"-"`
	Priority     int             `json:"priority"`
	DateCreation string          `json:"date_creation"`
	DateMod      string          `json:"date_mod"`
	Entity       string          `json:"entities_id"`
	Category     string          `json:"itilcategories_id"`
	Requester    string          `json:"users_id_recipient"`
}

func (t *Ticket) parseStatus() {
	var code int
	if err := json.Unmarshal(t.RawStatus, &code); err == nil {
		t.StatusCode = code
		t.Status = statusName(code)
		return
	}
	var s string
	if err := json.Unmarshal(t.RawStatus, &s); err == nil {
		t.Status = s
	}
}

// Status codes
const (
	StatusNew              = 1
	StatusProcessing       = 2 // Assigned
	StatusProcessingPlaned = 3
	StatusPending          = 4
	StatusSolved           = 5
	StatusClosed           = 6
)

// SearchMyTickets returns tickets: new tickets OR (assigned to current user AND not solved/closed).
func (c *Client) SearchMyTickets() ([]Ticket, error) {
	params := url.Values{}
	// criteria[0]: status = New (1)
	params.Set("criteria[0][link]", "AND")
	params.Set("criteria[0][field]", "12")
	params.Set("criteria[0][searchtype]", "equals")
	params.Set("criteria[0][value]", "1")
	// criteria[1]: OR technician = current logged-in user
	params.Set("criteria[1][link]", "OR")
	params.Set("criteria[1][field]", "5")
	params.Set("criteria[1][searchtype]", "equals")
	params.Set("criteria[1][value]", "myself")
	// criteria[2]: AND status = notold (not solved, not closed)
	params.Set("criteria[2][link]", "AND")
	params.Set("criteria[2][field]", "12")
	params.Set("criteria[2][searchtype]", "equals")
	params.Set("criteria[2][value]", "notold")
	params.Set("forcedisplay[0]", "1")  // id
	params.Set("forcedisplay[1]", "2")  // name
	params.Set("forcedisplay[2]", "12") // status
	params.Set("forcedisplay[3]", "3")  // priority
	params.Set("forcedisplay[4]", "15") // date_creation
	params.Set("forcedisplay[5]", "19") // date_mod
	params.Set("range", "0-50")
	params.Set("sort[0]", "19")
	params.Set("order[0]", "DESC")

	data, err := c.get("/search/Ticket?" + params.Encode())
	if err != nil {
		return nil, fmt.Errorf("searching tickets: %w", err)
	}

	var result struct {
		Data []map[string]interface{} `json:"data"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("decoding search results: %w", err)
	}

	var tickets []Ticket
	for _, item := range result.Data {
		statusCode := intFromSearchField(item, "12")
		t := Ticket{
			ID:         intFromSearchField(item, "2"),
			Name:       unescapeAll(stringFromSearchField(item, "1")),
			Status:     statusName(statusCode),
			StatusCode: statusCode,
			Priority:   intFromSearchField(item, "3"),
		}
		if v, ok := item["15"].(string); ok {
			t.DateCreation = v
		}
		if v, ok := item["19"].(string); ok {
			t.DateMod = v
		}
		tickets = append(tickets, t)
	}

	return tickets, nil
}

// GetTicket returns the full details of a single ticket.
func (c *Client) GetTicket(id int) (*Ticket, error) {
	data, err := c.get(fmt.Sprintf("/Ticket/%d?expand_dropdowns=true", id))
	if err != nil {
		return nil, fmt.Errorf("getting ticket %d: %w", id, err)
	}

	var t Ticket
	if err := json.Unmarshal(data, &t); err != nil {
		return nil, fmt.Errorf("decoding ticket: %w", err)
	}
	t.parseStatus()
	t.Name = unescapeAll(t.Name)
	return &t, nil
}

func statusName(code int) string {
	switch code {
	case 1:
		return "Novo"
	case 2:
		return "Processando (atribuído)"
	case 3:
		return "Processando (planejado)"
	case 4:
		return "Pendente"
	case 5:
		return "Solucionado"
	case 6:
		return "Fechado"
	default:
		return fmt.Sprintf("Status %d", code)
	}
}

// PriorityName returns a human-readable priority label.
func PriorityName(p int) string {
	switch p {
	case 1:
		return "Muito baixa"
	case 2:
		return "Baixa"
	case 3:
		return "Média"
	case 4:
		return "Alta"
	case 5:
		return "Muito alta"
	case 6:
		return "Crítica"
	default:
		return fmt.Sprintf("Prioridade %d", p)
	}
}

func intFromSearchField(item map[string]interface{}, key string) int {
	v, ok := item[key]
	if !ok {
		return 0
	}
	switch val := v.(type) {
	case float64:
		return int(val)
	case string:
		return 0
	default:
		return 0
	}
}

func unescapeAll(s string) string {
	for {
		decoded := html.UnescapeString(s)
		if decoded == s {
			return s
		}
		s = decoded
	}
}

func stringFromSearchField(item map[string]interface{}, key string) string {
	v, ok := item[key]
	if !ok {
		return ""
	}
	switch val := v.(type) {
	case string:
		return val
	case float64:
		return fmt.Sprintf("%.0f", val)
	default:
		return fmt.Sprintf("%v", val)
	}
}
