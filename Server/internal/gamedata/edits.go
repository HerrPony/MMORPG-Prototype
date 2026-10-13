package gamedata

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"sync"
)

// Map editor edits: a JSON file that is layered on top of the built-in world data. Everything in it is
// optional; nil/empty means "use the defaults from the code". ApplyEdits can be called at any time (from the
// world goroutine) and always starts again from the defaults, so removing an edit restores the original.
type Edits struct {
	Scale     float64        `json:"scale,omitempty"`     // WorldScale the edits were made for (the grid and all positions depend on it)
	Spawns    *[]Spawn       `json:"spawns,omitempty"`    // complete open-world mob spawn list
	Objects   *[]ObjectSpawn `json:"objects,omitempty"`   // complete list of gather nodes / portals
	Buildings *[]Building    `json:"buildings,omitempty"` // complete building list
	PropsAdd  []PropEdit     `json:"propsAdd,omitempty"`
	PropsDel  [][2]float64   `json:"propsDel,omitempty"`
	Heights   [][3]int       `json:"heights,omitempty"` // [grid i, grid j, delta in 1/20 m] against the generated terrain
	Biomes    [][3]int       `json:"biomes,omitempty"`  // [grid i, grid j, biome]

	NPCs   map[string]json.RawMessage `json:"npcs,omitempty"` // full definitions (replace or add)
	Mobs   map[string]json.RawMessage `json:"mobs,omitempty"`
	Quests map[string]json.RawMessage `json:"quests,omitempty"`
}

var BuildingKinds = []string{"house", "hut", "tower", "wall", "keep", "fountain", "pad", "tent", "cave", "pagoda", "gate", "archgate", "pond", "lantern", "dock"}

var (
	defOnce   sync.Once
	defSpawns []Spawn
	defObjs   []ObjectSpawn
	defBlds   []Building
	defHm     [GridNX * GridNZ]float32
	defBm     [GridNX * GridNZ]uint8
	defNPCs   []json.RawMessage
	defMobs   map[string]json.RawMessage
	defQuests map[string]json.RawMessage
	curEdits  = &Edits{}
	editMu    sync.Mutex
)

func snapshotDefaults() {
	defOnce.Do(func() {
		Init()
		ensureWorld()
		defSpawns = append([]Spawn(nil), spawnsAll...)
		defObjs = append([]ObjectSpawn(nil), ObjectSpawns...)
		defBlds = append([]Building(nil), Buildings()...)
		defHm, defBm = hm, bm
		for _, n := range NPCs {
			b, _ := json.Marshal(n)
			defNPCs = append(defNPCs, b)
		}
		defMobs, defQuests = map[string]json.RawMessage{}, map[string]json.RawMessage{}
		for id, m := range Mobs {
			defMobs[id], _ = json.Marshal(m)
		}
		for id, q := range Quests {
			defQuests[id], _ = json.Marshal(q)
		}
	})
}

// CurrentEdits returns the edits that are currently applied.
func CurrentEdits() *Edits { editMu.Lock(); defer editMu.Unlock(); return curEdits }

// BaseGrid returns the generated terrain without edits (int16 height*20 and biome) of one grid cell.
func BaseGrid(i, j int) (int, int) {
	snapshotDefaults()
	return int(math.Round(float64(defHm[j*GridNX+i]) * 20)), int(defBm[j*GridNX+i])
}

func LoadEdits(path string) (*Edits, error) {
	b, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return &Edits{}, nil
	}
	if err != nil {
		return nil, err
	}
	var e Edits
	if err := json.Unmarshal(b, &e); err != nil {
		return nil, err
	}
	if e.Scale != WorldScale && string(b) != "{}" {
		// made for a world of another size: the positions would be wrong, keep the file aside
		_ = os.Rename(path, path+".unscaled")
		return &Edits{}, nil
	}
	return &e, nil
}

func SaveEdits(path string, e *Edits) error {
	e.Scale = WorldScale
	b, err := json.Marshal(e)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, b, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

func inWorld(x, z float64) bool {
	return !math.IsNaN(x) && !math.IsNaN(z) && x > WorldMinX+2 && x < WorldMaxX-2 && math.Abs(z) < WorldHalf-2
}

// ValidateEdits rejects edits that would break the world.
func ValidateEdits(e *Edits) error {
	snapshotDefaults()
	mobIDs := map[string]bool{}
	for id := range defMobs {
		mobIDs[id] = true
	}
	for id, raw := range e.Mobs {
		var m Mob
		if err := json.Unmarshal(raw, &m); err != nil {
			return fmt.Errorf("Mob %s: %v", id, err)
		}
		if m.ID != id || m.Name == "" || m.Model == "" || m.LevelMin < 1 || m.LevelMax < m.LevelMin || m.HPBase < 0 || m.DmgMax < m.DmgMin {
			return fmt.Errorf("Mob %s: ungültige Werte (ID, Name, Modell, Level, Schaden)", id)
		}
		mobIDs[id] = true
	}
	if e.Spawns != nil {
		if len(*e.Spawns) > 4000 {
			return errors.New("Zu viele Spawns")
		}
		for _, s := range *e.Spawns {
			if !mobIDs[s.Mob] || !inWorld(s.X, s.Z) {
				return fmt.Errorf("Spawn %q bei %.0f,%.0f ist ungültig", s.Mob, s.X, s.Z)
			}
		}
	}
	if e.Objects != nil {
		if len(*e.Objects) > 2000 {
			return errors.New("Zu viele Objekte")
		}
		for _, o := range *e.Objects {
			if Objects[o.Obj] == nil || !inWorld(o.X, o.Z) {
				return fmt.Errorf("Objekt %q bei %.0f,%.0f ist ungültig", o.Obj, o.X, o.Z)
			}
		}
	}
	if e.Buildings != nil {
		if len(*e.Buildings) > 4000 {
			return errors.New("Zu viele Gebäude")
		}
		for _, b := range *e.Buildings {
			ok := false
			for _, k := range BuildingKinds {
				ok = ok || k == b.Kind
			}
			if !ok || !inWorld(b.X, b.Z) || b.W < 0.5 || b.W > 120 || b.D < 0.5 || b.D > 120 || b.H < 0.1 || b.H > 60 {
				return fmt.Errorf("Gebäude %q bei %.0f,%.0f ist ungültig", b.Kind, b.X, b.Z)
			}
		}
	}
	if len(e.PropsAdd) > 60000 || len(e.PropsDel) > 60000 || len(e.Heights) > 400000 || len(e.Biomes) > 400000 {
		return errors.New("Zu viele Änderungen")
	}
	for _, p := range e.PropsAdd {
		if kindIndex2(p.Kind) < 0 || !inWorld(p.X, p.Z) {
			return fmt.Errorf("Prop %q ist ungültig", p.Kind)
		}
	}
	for _, c := range e.Heights {
		if c[0] < 0 || c[0] >= GridNX || c[1] < 0 || c[1] >= GridNZ || c[2] < -2000 || c[2] > 2000 {
			return errors.New("Höhenänderung außerhalb der Karte")
		}
	}
	for _, c := range e.Biomes {
		if c[0] < 0 || c[0] >= GridNX || c[1] < 0 || c[1] >= GridNZ || c[2] < 0 || c[2] >= int(numBiomes) {
			return errors.New("Biom-Änderung ungültig")
		}
	}
	for id, raw := range e.NPCs {
		var n NPC
		if err := json.Unmarshal(raw, &n); err != nil || n.ID != id || n.Name == "" || !inWorld(n.X, n.Z) {
			return fmt.Errorf("NPC %s ist ungültig", id)
		}
	}
	for id, raw := range e.Quests {
		var q Quest
		if err := json.Unmarshal(raw, &q); err != nil {
			return fmt.Errorf("Quest %s: %v", id, err)
		}
		if q.ID != id || q.Title == "" || len(q.Objectives) == 0 {
			return fmt.Errorf("Quest %s: ID, Titel und mindestens ein Ziel nötig", id)
		}
		for _, o := range q.Objectives {
			if o.Count < 1 {
				return fmt.Errorf("Quest %s: Zielanzahl muss >= 1 sein", id)
			}
		}
	}
	return nil
}

// ApplyEdits rebuilds the world data from the defaults plus e. Call it on the world goroutine.
func ApplyEdits(e *Edits) error {
	snapshotDefaults()
	if err := ValidateEdits(e); err != nil {
		return err
	}
	editMu.Lock()
	defer editMu.Unlock()
	curEdits = e

	// terrain
	hm, bm = defHm, defBm
	for _, c := range e.Heights {
		idx := c[1]*GridNX + c[0]
		hm[idx] = float32(float64(int16(math.Round(float64(defHm[idx])*20))+int16(c[2])) / 20)
	}
	for _, c := range e.Biomes {
		bm[c[1]*GridNX+c[0]] = uint8(c[2])
	}

	// definitions: mobs and quests are replaced in place (pointers stay valid), NPCs are rebuilt
	restore := func(id string, raw json.RawMessage, into any) { _ = json.Unmarshal(raw, into) }
	for id, raw := range defMobs {
		m := Mobs[id]
		*m = Mob{}
		restore(id, raw, m)
	}
	for id, raw := range defQuests {
		q := Quests[id]
		*q = Quest{}
		restore(id, raw, q)
	}
	for id := range Mobs {
		if _, isDef := defMobs[id]; !isDef {
			delete(Mobs, id)
		}
	}
	for id := range Quests {
		if _, isDef := defQuests[id]; !isDef {
			delete(Quests, id)
		}
	}
	for id, raw := range e.Mobs {
		m := Mobs[id]
		if m == nil {
			m = &Mob{}
			Mobs[id] = m
		}
		*m = Mob{}
		_ = json.Unmarshal(raw, m)
	}
	for id, raw := range e.Quests {
		q := Quests[id]
		if q == nil {
			q = &Quest{}
			Quests[id] = q
		}
		*q = Quest{}
		_ = json.Unmarshal(raw, q)
	}
	NPCs = NPCs[:0:0]
	byID := map[string]int{}
	for _, raw := range defNPCs {
		n := &NPC{}
		_ = json.Unmarshal(raw, n)
		byID[n.ID] = len(NPCs)
		NPCs = append(NPCs, n)
	}
	ids := make([]string, 0, len(e.NPCs))
	for id := range e.NPCs {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	for _, id := range ids {
		n := &NPC{}
		_ = json.Unmarshal(e.NPCs[id], n)
		if i, ok := byID[id]; ok {
			NPCs[i] = n
		} else {
			NPCs = append(NPCs, n)
		}
	}

	// placement lists
	spawnsAll = append([]Spawn(nil), defSpawns...)
	if e.Spawns != nil {
		spawnsAll = append([]Spawn(nil), *e.Spawns...)
	}
	ObjectSpawns = append([]ObjectSpawn(nil), defObjs...)
	if e.Objects != nil {
		ObjectSpawns = append([]ObjectSpawn(nil), *e.Objects...)
	}
	buildingList = append([]Building(nil), defBlds...)
	if e.Buildings != nil {
		buildingList = append([]Building(nil), *e.Buildings...)
	}
	propAdd, propDel = e.PropsAdd, e.PropsDel
	resetProps()
	computeQuestAreas()
	return nil
}
