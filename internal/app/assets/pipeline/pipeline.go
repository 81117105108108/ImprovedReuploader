package pipeline

import (
	"time"

	"github.com/kartFr/Asset-Reuploader/internal/app/config"
	"github.com/kartFr/Asset-Reuploader/internal/app/request"
	"github.com/kartFr/Asset-Reuploader/internal/roblox/develop"
)

// RatePolicy centralizes per-asset-type tuning (was scattered consts).
type RatePolicy struct {
	StartsPerMinute   int
	MaxConcurrent     int
	MaxParallelChunks int
	UploadTries       int
	UploadDelay       time.Duration
}

// Default policies (config-overridable, fallback to previous hardcoded behavior).
// Global keys: starts_per_minute/max_concurrent/max_parallel_chunks/upload_tries.
// Per-type keys (optional): animation_starts_per_minute, sound_starts_per_minute, mesh_starts_per_minute, etc.
func AnimationPolicy() RatePolicy {
	return RatePolicy{
		StartsPerMinute:   config.GetInt("animation_starts_per_minute", config.StartsPerMinute()),
		MaxConcurrent:     config.GetInt("animation_max_concurrent", config.MaxConcurrent()),
		MaxParallelChunks: config.MaxParallelChunks(),
		UploadTries:       config.GetInt("animation_upload_tries", 5),
		UploadDelay:       900 * time.Millisecond,
	}
}

func SoundPolicy() RatePolicy {
	return RatePolicy{
		StartsPerMinute:   config.GetInt("sound_starts_per_minute", 120),
		MaxConcurrent:     config.GetInt("sound_max_concurrent", 8),
		MaxParallelChunks: config.MaxParallelChunks(),
		UploadTries:       config.UploadTries(),
		UploadDelay:       time.Second,
	}
}

func MeshPolicy() RatePolicy {
	return RatePolicy{
		StartsPerMinute:   config.GetInt("mesh_starts_per_minute", 3000),
		MaxConcurrent:     config.GetInt("mesh_max_concurrent", 16),
		MaxParallelChunks: config.MaxParallelChunks(),
		UploadTries:       config.UploadTries(),
		UploadDelay:       time.Second,
	}
}

// GetDefaultPlaceIDs copies r.DefaultPlaceIDs and merges r.PlaceID (copy semantics, no aliasing).
func GetDefaultPlaceIDs(r *request.Request) ([]int64, map[int64]struct{}) {
	ids := append([]int64(nil), r.DefaultPlaceIDs...)
	m := make(map[int64]struct{}, len(ids)+1)
	for _, id := range ids {
		m[id] = struct{}{}
	}
	if r.PlaceID > 0 {
		if _, ok := m[r.PlaceID]; !ok {
			ids = append(ids, r.PlaceID)
			m[r.PlaceID] = struct{}{}
		}
	}
	return ids, m
}

// GroupByCreator groups filtered assets by Creator.Type -> Creator.TargetID.
func GroupByCreator(infos []*develop.AssetInfo) map[string]map[int64][]*develop.AssetInfo {
	out := make(map[string]map[int64][]*develop.AssetInfo)
	for _, info := range infos {
		byType, ok := out[info.Creator.Type]
		if !ok {
			byType = make(map[int64][]*develop.AssetInfo)
			out[info.Creator.Type] = byType
		}
		byType[info.Creator.TargetID] = append(byType[info.Creator.TargetID], info)
	}
	return out
}

// CalcBatchSize returns 50 except last chunk (idsToUpload%50, 50 if divisible).
func CalcBatchSize(idsToUpload, taskIndex, taskCount int) int {
	if taskIndex != taskCount-1 {
		return 50
	}
	if n := idsToUpload % 50; n != 0 {
		return n
	}
	return 50
}
