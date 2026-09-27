package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/kartFr/Asset-Reuploader/internal/app/assets"
	"github.com/kartFr/Asset-Reuploader/internal/app/config"
	"github.com/kartFr/Asset-Reuploader/internal/app/request"
	"github.com/kartFr/Asset-Reuploader/internal/app/response"
	"github.com/kartFr/Asset-Reuploader/internal/color"
	"github.com/kartFr/Asset-Reuploader/internal/files"
	"github.com/kartFr/Asset-Reuploader/internal/roblox"
)

var CompatiblePluginVersion = ""

const (
	routeRoot     = "GET /"
	routeReupload = "POST /reupload"
)

func getOutputFileName(reuploadType string) string {
	t := time.Now()
	return fmt.Sprintf("Output_%s_%s.json", reuploadType, t.Format("2006-01-02_15-04-05"))
}

func serve(c *roblox.Client) error {
	var mu sync.Mutex
	var exportedJSONName string
	var exportJSON bool
	var busy bool
	finished := true

	respHistory := make([]response.ResponseItem, 0)
	isBusy := func() (bool, bool) {
		mu.Lock()
		defer mu.Unlock()
		return busy, finished
	}
	setBusy := func(b, f bool) {
		mu.Lock()
		defer mu.Unlock()
		busy, finished = b, f
	}
	setExport := func(e bool, name string) {
		mu.Lock()
		defer mu.Unlock()
		exportJSON, exportedJSONName = e, name
	}
	appendHistory := func(item response.ResponseItem) {
		mu.Lock()
		if !exportJSON {
			mu.Unlock()
			return
		}
		respHistory = append(respHistory, item)
		name := exportedJSONName
		history := append([]response.ResponseItem(nil), respHistory...)
		j, err := json.Marshal(history)
		if err != nil {
			mu.Unlock()
			color.Error.Println("encode history:", err)
			return
		}
		// Hold lock through file write to keep history file ordered
		// (writes are per-upload, small JSON, local disk — negligible hold).
		if err := files.Write(name, string(j)); err != nil {
			mu.Unlock()
			color.Error.Println("write history:", err)
			return
		}
		mu.Unlock()
	}
	resetHistory := func() {
		mu.Lock()
		defer mu.Unlock()
		respHistory = make([]response.ResponseItem, 0)
	}
	resp := response.New(func(i response.ResponseItem) {
		appendHistory(i)
	})

	http.HandleFunc(routeRoot, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		busyNow, finishedNow := isBusy()
		if resp.Len() == 0 && !busyNow {
			if !finishedNow {
				setBusy(false, true)
				setExport(false, "")
				resp.Clear()
				resetHistory()

				fmt.Fprint(w, "done")
				fmt.Println("Finished reuploading. (you can rerun without restarting)")
			}

			return
		}

		if err := resp.EncodeJSON(json.NewEncoder(w)); err != nil {
			color.Error.Println("encode response:", err)
		} else {
			resp.Clear()
		}
	})

	http.HandleFunc(routeReupload, func(w http.ResponseWriter, r *http.Request) {
		if b, f := isBusy(); b || !f {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		var req request.RawRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			color.Error.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if CompatiblePluginVersion != "" && req.PluginVersion != CompatiblePluginVersion {
			w.WriteHeader(http.StatusConflict)
			return
		}

		if exists := assets.DoesModuleExist(req.AssetType); !exists {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		startReupload, err := assets.NewReuploadHandlerWithType(req.AssetType, c, &req, resp)
		if err != nil {
			color.Error.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if req.ExportJSON {
			setExport(true, getOutputFileName(req.AssetType))
		} else {
			setExport(false, "")
		}

		setBusy(true, false)

		go func() {
			start := time.Now()
			err := startReupload()
			if err != nil {
				setBusy(false, true)
				color.Error.Println("Failed to start reuploading: ", err)
				return
			}

			duration := time.Since(start)
			fmt.Printf("Reuploading took %d hours, %d minutes, and %d seconds\n", int(duration.Hours()), int(duration.Minutes())%60, int(duration.Seconds())%60)
			fmt.Println("Waiting for client to finish changing ids...")
		}()

		w.WriteHeader(http.StatusOK)
	})

	return http.ListenAndServe(":"+config.Get("port"), nil)
}
