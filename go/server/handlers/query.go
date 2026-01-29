package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/tehnerd/vatran/go/server/models"
)

func decodeVIPRequest(r *http.Request, req *models.VIPRequest) error {
	query := r.URL.Query()
	_, hasAddress := query["address"]
	_, hasPort := query["port"]
	_, hasProto := query["proto"]
	if hasAddress || hasPort || hasProto {
		address := query.Get("address")
		portStr := query.Get("port")
		protoStr := query.Get("proto")
		if address == "" || portStr == "" || protoStr == "" {
			return fmt.Errorf("missing required query parameters")
		}
		port, err := strconv.ParseUint(portStr, 10, 16)
		if err != nil {
			return fmt.Errorf("invalid port: %w", err)
		}
		proto, err := strconv.ParseUint(protoStr, 10, 8)
		if err != nil {
			return fmt.Errorf("invalid proto: %w", err)
		}
		req.Address = address
		req.Port = uint16(port)
		req.Proto = uint8(proto)
		return nil
	}
	return json.NewDecoder(r.Body).Decode(req)
}

func decodeRealIndexRequest(r *http.Request, req *models.GetRealIndexRequest) error {
	query := r.URL.Query()
	if values, ok := query["address"]; ok {
		if len(values) == 0 || values[0] == "" {
			return fmt.Errorf("missing required query parameters")
		}
		req.Address = values[0]
		return nil
	}
	return json.NewDecoder(r.Body).Decode(req)
}

func decodeRealStatsRequest(r *http.Request, req *models.GetRealStatsRequest) error {
	query := r.URL.Query()
	if values, ok := query["index"]; ok {
		if len(values) == 0 || values[0] == "" {
			return fmt.Errorf("missing required query parameters")
		}
		index, err := strconv.ParseUint(values[0], 10, 32)
		if err != nil {
			return fmt.Errorf("invalid index: %w", err)
		}
		req.Index = uint32(index)
		return nil
	}
	return json.NewDecoder(r.Body).Decode(req)
}

func decodeBPFMapStatsRequest(r *http.Request, req *models.GetBPFMapStatsRequest) error {
	query := r.URL.Query()
	if values, ok := query["map_name"]; ok {
		if len(values) == 0 || values[0] == "" {
			return fmt.Errorf("missing required query parameters")
		}
		req.MapName = values[0]
		return nil
	}
	return json.NewDecoder(r.Body).Decode(req)
}
