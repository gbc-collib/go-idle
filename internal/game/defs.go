package game

type BuildingDef struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Produces    string  `json:"produces"` // e.g. "cpu"
	BaseProd    float64 `json:"base_prod"`
	ProdScaling float64 `json:"prod_scaling"` // per level
	BaseCost    float64 `json:"base_cost"`
	CostScaling float64 `json:"cost_scaling"`
}

var BuildingsDefs map[string]BuildingDef
