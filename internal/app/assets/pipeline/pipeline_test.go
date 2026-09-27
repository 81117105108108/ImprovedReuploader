package pipeline

import (
	"testing"

	"github.com/kartFr/Asset-Reuploader/internal/app/request"
	"github.com/kartFr/Asset-Reuploader/internal/roblox/develop"
)

func mkInfo(typ string, id int64) *develop.AssetInfo {
	a := &develop.AssetInfo{}
	a.Creator.Type = typ
	a.Creator.TargetID = id
	return a
}

func TestCalcBatchSize(t *testing.T) {
	if got := CalcBatchSize(120, 0, 3); got != 50 {
		t.Fatalf("first %d", got)
	}
	if got := CalcBatchSize(120, 2, 3); got != 20 {
		t.Fatalf("last 120->20, got %d", got)
	}
	if got := CalcBatchSize(100, 1, 2); got != 50 {
		t.Fatalf("divisible %d", got)
	}
}

func TestGroupByCreator(t *testing.T) {
	infos := []*develop.AssetInfo{mkInfo("User", 1), mkInfo("User", 1), mkInfo("Group", 2)}
	g := GroupByCreator(infos)
	if len(g["User"][1]) != 2 || len(g["Group"][2]) != 1 {
		t.Fatalf("grouping %+v", g)
	}
}

func TestGetDefaultPlaceIDs(t *testing.T) {
	r := &request.Request{DefaultPlaceIDs: []int64{1, 2}, PlaceID: 2}
	ids, m := GetDefaultPlaceIDs(r)
	if len(ids) != 2 || len(m) != 2 {
		t.Fatalf("dedupe %+v %+v", ids, m)
	}
	r2 := &request.Request{DefaultPlaceIDs: []int64{1}, PlaceID: 3}
	ids2, _ := GetDefaultPlaceIDs(r2)
	if len(ids2) != 2 || ids2[1] != 3 {
		t.Fatalf("merge %+v", ids2)
	}
	if len(r2.DefaultPlaceIDs) != 1 {
		t.Fatalf("aliased caller slice")
	}
}
