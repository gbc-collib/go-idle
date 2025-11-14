package game;
import (
	"os"
	"encoding/json"
)

func InitGameData() error {
	bytes, err := os.ReadFile("data/buildings.json")
	if err != nil{
		return err
	}

	var defs []BuildingDef
    if err := json.Unmarshal(bytes, &defs); err != nil {
        return err
    }

    BuildingsDefs = make(map[string]BuildingDef)
    for _, d := range defs {
        BuildingsDefs[d.ID] = d
    }
	return nil
}
