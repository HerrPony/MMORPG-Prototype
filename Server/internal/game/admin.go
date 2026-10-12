package game

import (
	"encoding/json"
	"log"
	"path/filepath"
	"sort"
	"strings"

	"mmorpg/internal/config"
	"mmorpg/internal/gamedata"
	"mmorpg/internal/protocol"
)

// Map editor back end. Only the gateway calls these, and only for accounts on the admin list.

func editsPath(cfg config.Config) string { return filepath.Join(cfg.DataDir, "world_edits.json") }

// AdminGet sends everything the editor needs: the effective world plus the edits that are applied.
func (w *World) AdminGet(out Sender) {
	w.Do(func() {
		items := map[string]protocol.M{}
		for id, it := range gamedata.Items {
			items[id] = protocol.M{"name": it.Name, "icon": it.Icon}
		}
		out.Send(protocol.M{
			"t": "adminData", "scenery": w.sceneryJSON, "edits": gamedata.CurrentEdits(),
			"spawns": gamedata.AllSpawns(), "objects": gamedata.ObjectSpawns,
			"npcs": gamedata.NPCs, "mobs": gamedata.Mobs, "quests": gamedata.Quests,
			"items": items, "objDefs": gamedata.Objects, "kinds": gamedata.PropKinds, "bkinds": gamedata.BuildingKinds,
			"zones": gamedata.Zones, "online": len(w.players),
		})
	})
}

type adminSave struct {
	Spawns    *[]gamedata.Spawn          `json:"spawns"`
	Objects   *[]gamedata.ObjectSpawn    `json:"objects"`
	Buildings *[]gamedata.Building       `json:"buildings"`
	PropsAdd  *[]gamedata.PropEdit       `json:"propsAdd"`
	PropsDel  *[][2]float64              `json:"propsDel"`
	Heights   [][3]int                   `json:"heights"` // changed cells: [i, j, absolute height*20]
	Biomes    [][3]int                   `json:"biomes"`  // changed cells: [i, j, biome]
	NPCs      map[string]json.RawMessage `json:"npcs"`
	Mobs      map[string]json.RawMessage `json:"mobs"`
	Quests    map[string]json.RawMessage `json:"quests"`
	Reset     []string                   `json:"reset"` // "spawns" "objects" "buildings" "props" "terrain" "npcs:id" "mobs:id" "quests:id"
}

// AdminSave merges the editor's changes into the stored edits, applies them to the running world and tells
// all players to reload.
func (w *World) AdminSave(out Sender, data []byte) {
	var in adminSave
	if err := json.Unmarshal(data, &in); err != nil {
		out.Send(protocol.M{"t": "adminError", "error": "Ungültige Daten"})
		return
	}
	w.Do(func() {
		old := gamedata.CurrentEdits()
		ed := mergeEdits(old, &in)
		if err := gamedata.ApplyEdits(ed); err != nil {
			_ = gamedata.ApplyEdits(old)
			out.Send(protocol.M{"t": "adminError", "error": err.Error()})
			return
		}
		if err := gamedata.SaveEdits(editsPath(w.cfg), ed); err != nil {
			log.Printf("world_edits.json: %v", err)
			out.Send(protocol.M{"t": "adminError", "error": "Speichern fehlgeschlagen: " + err.Error()})
			return
		}
		w.reloadContent()
		log.Printf("Map-Editor: Änderungen gespeichert (%d Spieler online)", len(w.players))
		out.Send(protocol.M{"t": "adminSaved"})
	})
}

func mergeEdits(old *gamedata.Edits, in *adminSave) *gamedata.Edits {
	ed := *old
	cp := func(m map[string]json.RawMessage) map[string]json.RawMessage {
		n := map[string]json.RawMessage{}
		for k, v := range m {
			n[k] = v
		}
		return n
	}
	ed.NPCs, ed.Mobs, ed.Quests = cp(old.NPCs), cp(old.Mobs), cp(old.Quests)
	for _, r := range in.Reset {
		kind, id, _ := strings.Cut(r, ":")
		switch kind {
		case "spawns":
			ed.Spawns = nil
		case "objects":
			ed.Objects = nil
		case "buildings":
			ed.Buildings = nil
		case "props":
			ed.PropsAdd, ed.PropsDel = nil, nil
		case "terrain":
			ed.Heights, ed.Biomes = nil, nil
		case "npcs":
			delete(ed.NPCs, id)
		case "mobs":
			delete(ed.Mobs, id)
		case "quests":
			delete(ed.Quests, id)
		}
	}
	if in.Spawns != nil {
		ed.Spawns = in.Spawns
	}
	if in.Objects != nil {
		ed.Objects = in.Objects
	}
	if in.Buildings != nil {
		ed.Buildings = in.Buildings
	}
	if in.PropsAdd != nil {
		ed.PropsAdd = *in.PropsAdd
	}
	if in.PropsDel != nil {
		ed.PropsDel = *in.PropsDel
	}
	for k, v := range in.NPCs {
		ed.NPCs[k] = v
	}
	for k, v := range in.Mobs {
		ed.Mobs[k] = v
	}
	for k, v := range in.Quests {
		ed.Quests[k] = v
	}
	ed.Heights = mergeCells(ed.Heights, in.Heights, true)
	ed.Biomes = mergeCells(ed.Biomes, in.Biomes, false)
	return &ed
}

// mergeCells folds changed grid cells into the stored sparse edit list. For heights the editor sends absolute
// values (stored as the difference to the generated terrain), for biomes the new biome.
func mergeCells(old, changed [][3]int, height bool) [][3]int {
	type cell struct{ i, j int }
	m := map[cell]int{}
	for _, c := range old {
		m[cell{c[0], c[1]}] = c[2]
	}
	for _, c := range changed {
		if c[0] < 0 || c[0] >= gamedata.GridNX || c[1] < 0 || c[1] >= gamedata.GridNZ {
			continue
		}
		bh, bb := gamedata.BaseGrid(c[0], c[1])
		if height {
			if d := c[2] - bh; d != 0 {
				m[cell{c[0], c[1]}] = d
			} else {
				delete(m, cell{c[0], c[1]})
			}
		} else if c[2] != bb {
			m[cell{c[0], c[1]}] = c[2]
		} else {
			delete(m, cell{c[0], c[1]})
		}
	}
	out := make([][3]int, 0, len(m))
	for k, v := range m {
		out = append(out, [3]int{k.i, k.j, v})
	}
	sort.Slice(out, func(a, b int) bool {
		if out[a][1] != out[b][1] {
			return out[a][1] < out[b][1]
		}
		return out[a][0] < out[b][0]
	})
	return out
}

// reloadContent replaces all static world content (mobs, NPCs, objects, ferry) with the freshly edited data
// and tells every player to reload the client, which fetches the new scenery on login.
func (w *World) reloadContent() {
	for _, p := range w.players {
		w.savePlayer(p)
	}
	for id, e := range w.ents {
		if e.Inst != 0 {
			continue
		}
		switch e.Kind {
		case KindMob, KindNPC, KindObj, KindShip:
			delete(w.ents, id)
		}
	}
	for _, p := range w.players {
		p.P.Target, p.P.Auto, p.P.Cast, p.P.Aboard = 0, false, nil, 0
	}
	w.refreshStatic()
	w.spawnStatic()
	w.broadcast(protocol.M{"t": "worldChanged"})
}
