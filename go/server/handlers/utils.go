package handlers

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/tehnerd/vatran/go/katran"
	"github.com/tehnerd/vatran/go/server/lb"
	"github.com/tehnerd/vatran/go/server/models"
)

// UtilsHandler handles utility operations.
type UtilsHandler struct {
	manager *lb.Manager
}

// NewUtilsHandler creates a new UtilsHandler.
//
// Returns a new UtilsHandler instance.
func NewUtilsHandler() *UtilsHandler {
	return &UtilsHandler{
		manager: lb.GetManager(),
	}
}

// HandleMAC handles MAC address operations based on HTTP method.
//
// Parameters:
//   - w: The http.ResponseWriter to write to.
//   - r: The HTTP request.
func (h *UtilsHandler) HandleMAC(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.handleGetMAC(w, r)
	case http.MethodPut:
		h.handleChangeMAC(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleGetMAC handles GET /utils/mac - gets the current MAC address.
func (h *UtilsHandler) handleGetMAC(w http.ResponseWriter, r *http.Request) {
	lbInstance, ok := h.manager.Get()
	if !ok {
		models.WriteError(w, http.StatusServiceUnavailable, models.NewLBNotInitializedError())
		return
	}

	mac, err := lbInstance.GetMAC()
	if err != nil {
		models.WriteKatranError(w, err)
		return
	}

	// Format MAC as hex string with colons
	macStr := fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x",
		mac[0], mac[1], mac[2], mac[3], mac[4], mac[5])

	models.WriteSuccess(w, models.MACResponse{MAC: macStr})
}

// handleChangeMAC handles PUT /utils/mac - changes the MAC address.
func (h *UtilsHandler) handleChangeMAC(w http.ResponseWriter, r *http.Request) {
	lbInstance, ok := h.manager.Get()
	if !ok {
		models.WriteError(w, http.StatusServiceUnavailable, models.NewLBNotInitializedError())
		return
	}

	var req models.ChangeMACRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		models.WriteError(w, http.StatusBadRequest,
			models.NewInvalidRequestError("invalid request body: "+err.Error()))
		return
	}

	mac, err := parseMAC(req.MAC)
	if err != nil {
		models.WriteError(w, http.StatusBadRequest,
			models.NewInvalidRequestError("invalid MAC address: "+err.Error()))
		return
	}

	if err := lbInstance.ChangeMAC(mac); err != nil {
		models.WriteKatranError(w, err)
		return
	}

	models.WriteSuccess(w, nil)
}

// HandleRealForFlow handles GET /utils/real-for-flow - gets the real for a flow.
//
// Parameters:
//   - w: The http.ResponseWriter to write to.
//   - r: The HTTP request.
func (h *UtilsHandler) HandleRealForFlow(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	lbInstance, ok := h.manager.Get()
	if !ok {
		models.WriteError(w, http.StatusServiceUnavailable, models.NewLBNotInitializedError())
		return
	}

	var req models.FlowRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		models.WriteError(w, http.StatusBadRequest,
			models.NewInvalidRequestError("invalid request body: "+err.Error()))
		return
	}

	flow := katran.Flow{
		Src:     req.Src,
		Dst:     req.Dst,
		SrcPort: req.SrcPort,
		DstPort: req.DstPort,
		Proto:   req.Proto,
	}

	realAddr, err := lbInstance.GetRealForFlow(flow)
	if err != nil {
		models.WriteKatranError(w, err)
		return
	}

	models.WriteSuccess(w, models.RealForFlowResponse{Address: realAddr})
}

// HandleSimulatePacket handles POST /utils/simulate-packet - simulates a packet.
//
// Parameters:
//   - w: The http.ResponseWriter to write to.
//   - r: The HTTP request.
func (h *UtilsHandler) HandleSimulatePacket(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	lbInstance, ok := h.manager.Get()
	if !ok {
		models.WriteError(w, http.StatusServiceUnavailable, models.NewLBNotInitializedError())
		return
	}

	var req models.SimulatePacketRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		models.WriteError(w, http.StatusBadRequest,
			models.NewInvalidRequestError("invalid request body: "+err.Error()))
		return
	}

	// Decode base64 packet
	inPacket, err := base64.StdEncoding.DecodeString(req.Packet)
	if err != nil {
		models.WriteError(w, http.StatusBadRequest,
			models.NewInvalidRequestError("invalid base64 packet: "+err.Error()))
		return
	}

	outPacket, err := lbInstance.SimulatePacket(inPacket)
	if err != nil {
		models.WriteKatranError(w, err)
		return
	}

	// Encode output as base64
	outPacketB64 := base64.StdEncoding.EncodeToString(outPacket)

	models.WriteSuccess(w, models.SimulatePacketResponse{Packet: outPacketB64})
}

// HandleKatranProgFD handles GET /utils/prog-fd - gets the katran program FD.
//
// Parameters:
//   - w: The http.ResponseWriter to write to.
//   - r: The HTTP request.
func (h *UtilsHandler) HandleKatranProgFD(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	lbInstance, ok := h.manager.Get()
	if !ok {
		models.WriteError(w, http.StatusServiceUnavailable, models.NewLBNotInitializedError())
		return
	}

	fd, err := lbInstance.GetKatranProgFD()
	if err != nil {
		models.WriteKatranError(w, err)
		return
	}

	models.WriteSuccess(w, models.ProgFDResponse{FD: fd})
}

// HandleHealthcheckerProgFD handles GET /utils/hc-prog-fd - gets the healthchecker program FD.
//
// Parameters:
//   - w: The http.ResponseWriter to write to.
//   - r: The HTTP request.
func (h *UtilsHandler) HandleHealthcheckerProgFD(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	lbInstance, ok := h.manager.Get()
	if !ok {
		models.WriteError(w, http.StatusServiceUnavailable, models.NewLBNotInitializedError())
		return
	}

	fd, err := lbInstance.GetHealthcheckerProgFD()
	if err != nil {
		models.WriteKatranError(w, err)
		return
	}

	models.WriteSuccess(w, models.ProgFDResponse{FD: fd})
}

// HandleBPFMapFD handles GET /utils/map-fd - gets a BPF map FD by name.
//
// Parameters:
//   - w: The http.ResponseWriter to write to.
//   - r: The HTTP request.
func (h *UtilsHandler) HandleBPFMapFD(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	lbInstance, ok := h.manager.Get()
	if !ok {
		models.WriteError(w, http.StatusServiceUnavailable, models.NewLBNotInitializedError())
		return
	}

	var req models.GetBPFMapStatsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		models.WriteError(w, http.StatusBadRequest,
			models.NewInvalidRequestError("invalid request body: "+err.Error()))
		return
	}

	fd, err := lbInstance.GetBPFMapFDByName(req.MapName)
	if err != nil {
		models.WriteKatranError(w, err)
		return
	}

	models.WriteSuccess(w, models.ProgFDResponse{FD: fd})
}

// HandleGlobalLRUMapsFDs handles GET /utils/global-lru-map-fds - gets global LRU map FDs.
//
// Parameters:
//   - w: The http.ResponseWriter to write to.
//   - r: The HTTP request.
func (h *UtilsHandler) HandleGlobalLRUMapsFDs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	lbInstance, ok := h.manager.Get()
	if !ok {
		models.WriteError(w, http.StatusServiceUnavailable, models.NewLBNotInitializedError())
		return
	}

	fds, err := lbInstance.GetGlobalLRUMapsFDs()
	if err != nil {
		models.WriteKatranError(w, err)
		return
	}

	models.WriteSuccess(w, models.MapFDsResponse{FDs: fds})
}

// HandleAddSrcIPForPcktEncap handles POST /utils/src-ip-encap - adds source IP for packet encapsulation.
//
// Parameters:
//   - w: The http.ResponseWriter to write to.
//   - r: The HTTP request.
func (h *UtilsHandler) HandleAddSrcIPForPcktEncap(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	lbInstance, ok := h.manager.Get()
	if !ok {
		models.WriteError(w, http.StatusServiceUnavailable, models.NewLBNotInitializedError())
		return
	}

	var req models.AddSrcIPForPcktEncapRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		models.WriteError(w, http.StatusBadRequest,
			models.NewInvalidRequestError("invalid request body: "+err.Error()))
		return
	}

	if err := lbInstance.AddSrcIPForPcktEncap(req.Src); err != nil {
		models.WriteKatranError(w, err)
		return
	}

	models.WriteSuccess(w, nil)
}

// parseMAC parses a MAC address string (e.g., "aa:bb:cc:dd:ee:ff") to bytes.
func parseMAC(mac string) ([]byte, error) {
	// Remove common separators
	mac = strings.ReplaceAll(mac, ":", "")
	mac = strings.ReplaceAll(mac, "-", "")
	return hex.DecodeString(mac)
}
