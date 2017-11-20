package hand

func (h *Hand) HasGame(id string) bool {
	_, ok := h.Games[id]
	return ok
}

func (h *Hand) GetGame(id string) *Game {
	g, _ := h.Games[id]
	return g
}

func (h *Hand) CreateGame(id string) (g *Game) {
	g = &Game{}
	h.Games[id] = g
	return
}

func (h *Hand) RemoveGame(id string) {
	if h.HasGame(id) {
		delete(h.Games, id)
	}
}
