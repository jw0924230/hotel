package handlers

import "testing"

func TestInferTownship(t *testing.T) {
	candidates := []TownshipCandidate{
		{ID: 1, Name: "信義區", CityName: "台北"},
		{ID: 2, Name: "信義區", CityName: "基隆"},
		{ID: 3, Name: "大安區", CityName: "台北"},
		{ID: 4, Name: "竹北市", CityName: "新竹"},
		{ID: 5, Name: "東區", CityName: "嘉義"},
		{ID: 6, Name: "南竿鄉", CityName: "馬祖"},
	}

	tests := []struct {
		name, area, address string
		wantID              int
		wantReason          string
	}{
		{"normalizes tai character", "台北", "106 臺北市 大安區信義路三段", 3, ""},
		{"same township name uses city", "基隆", "基隆市信義區東明路", 2, ""},
		{"merged hsinchu parent", "新竹", "新竹縣竹北市光明六路", 4, ""},
		{"merged chiayi parent", "嘉義", "嘉義市東區中山路", 5, ""},
		{"matsu alias", "", "連江縣南竿鄉介壽村", 6, ""},
		{"conflicting city", "台北", "基隆市信義區東明路", 0, "city_conflict"},
		{"missing address", "台北", "", 0, "missing_address"},
		{"unknown township", "台北", "台北市不存在路一號", 0, "township_not_found"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			match := InferTownship(test.area, test.address, candidates)
			if test.wantID > 0 {
				if match.Candidate == nil || match.Candidate.ID != test.wantID {
					t.Fatalf("got %+v, want township ID %d", match, test.wantID)
				}
				return
			}
			if match.Candidate != nil || match.Reason != test.wantReason {
				t.Fatalf("got %+v, want reason %s", match, test.wantReason)
			}
		})
	}
}
