package handlers

import "net/http"

type VIPHandler struct{}

func NewVIPHandler() *VIPHandler { return &VIPHandler{} }

// VIP exists, but feature is not implemented
func (h *VIPHandler) VIPStatus(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"status":  "in processing",
		"message": "VIP features are planned but not implemented in Milestone 2",
	})
}
