package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"snowflake-id/internal/snowflake"
)

func main() {
	machineID, err := readMachineID()
	exitOnUsageError(err)

	port, err := readPort()
	exitOnUsageError(err)

	g, err := snowflake.New(machineID)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = w.Write([]byte("ok\n"))
	})
	mux.HandleFunc("/id", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		id, err := g.NextID()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(idResponse{
			ID:        id,
			MachineID: machineID,
		})
	})

	addr := ":" + strconv.Itoa(port)
	srv := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("machine_id=%d listening_on=%s\n", machineID, addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

type idResponse struct {
	ID        uint64 `json:"id"`
	MachineID uint16 `json:"machine_id"`
}

func exitOnUsageError(err error) {
	if err == nil {
		return
	}
	fmt.Fprintln(os.Stderr, err)
	os.Exit(2)
}

func readMachineID() (uint16, error) {
	machineIDStr, ok := os.LookupEnv("MACHINE_ID")
	if !ok || machineIDStr == "" {
		return 0, fmt.Errorf("MACHINE_ID env var is required (0-1023)")
	}

	machineID64, err := strconv.ParseUint(machineIDStr, 10, 16)
	if err != nil {
		return 0, fmt.Errorf("invalid MACHINE_ID %q: %v", machineIDStr, err)
	}

	return uint16(machineID64), nil
}

func readPort() (int, error) {
	portStr, ok := os.LookupEnv("PORT")
	if !ok || strings.TrimSpace(portStr) == "" {
		return 8080, nil
	}

	port64, err := strconv.ParseInt(portStr, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid PORT %q: %v", portStr, err)
	}
	if port64 <= 0 || port64 > 65535 {
		return 0, fmt.Errorf("invalid PORT %q: must be in 1..65535", portStr)
	}

	return int(port64), nil
}
