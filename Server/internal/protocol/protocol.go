// Package protocol defines the messages exchanged over the WebSocket.
// Every message is a JSON object with a type field "t".
package protocol

import "mmorpg/internal/storage"

// M is a generic outgoing message.
type M map[string]any

// Envelope is used to read the type of an incoming message.
type Envelope struct {
	T string `json:"t"`
}

// ---- Client -> Server ----

const (
	CRegister = "register"
	CLogin    = "login"
	CCreate   = "createChar"
	CDelete   = "deleteChar"
	CEnter    = "enterWorld"
	CMove     = "move"
	CTarget   = "target"
	CCast     = "cast"
	CAuto     = "auto"
	CInteract = "interact"
	CAccept   = "acceptQuest"
	CTurnIn   = "turnInQuest"
	CAbandon  = "abandonQuest"
	CLoot     = "loot"
	CChat     = "chat"
	CRevive   = "revive"
	CLogout   = "logout"
	CTrain    = "train"
	CSell     = "sell"
	CBuy      = "buy"
	CSellJunk = "sellJunk"
	CEquip    = "equip"
	CUnequip  = "unequip"
	CUse      = "use"
	CSaveSet  = "saveSet"
	CEquipSet = "equipSet"
	CDelSet   = "deleteSet"
	CInvite   = "invite"
	CInviteRe = "inviteReply"
	CLeaveP   = "leaveParty"
	CKick     = "kick"
	CPromote  = "promote"
	CFlight   = "flight"
	CPetCmd   = "petCmd"
)

type InviteMsg struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type InviteReplyMsg struct {
	Accept bool `json:"accept"`
}

type IDMsg struct {
	ID int `json:"id"`
}

type FlightMsg struct {
	NPC  int    `json:"npc"`
	Dest string `json:"dest"`
}

type EquipMsg struct {
	Idx  int    `json:"idx"`
	Slot string `json:"slot"` // optional target slot
}

type UnequipMsg struct {
	Slot string `json:"slot"`
}

type UseMsg struct {
	Idx int `json:"idx"`
}

type BuyMsg struct {
	NPC   int    `json:"npc"`
	Item  string `json:"item"`
	Count int    `json:"count"`
}

type SellJunkMsg struct {
	NPC int `json:"npc"`
}

type SetMsg struct {
	Name string `json:"name"`
}

type TrainMsg struct {
	NPC   int    `json:"npc"`
	Skill string `json:"skill"`
}

// SellMsg sells the inventory stack at Idx to a vendor.
type SellMsg struct {
	NPC int `json:"npc"`
	Idx int `json:"idx"`
}

type Credentials struct {
	User string `json:"user"`
	Pass string `json:"pass"`
}

type CreateChar struct {
	Name  string       `json:"name"`
	Race  string       `json:"race"`
	Class string       `json:"class"`
	Look  storage.Look `json:"look"`
}

type CharID struct {
	ID string `json:"id"`
}

type Move struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
	R float64 `json:"r"`
}

type TargetMsg struct {
	ID int `json:"id"`
}

type CastMsg struct {
	Skill string `json:"skill"`
}

type AutoMsg struct {
	On bool `json:"on"`
}

type InteractMsg struct {
	ID int `json:"id"`
}

type QuestMsg struct {
	Quest string `json:"quest"`
}

// LootMsg takes one item (Idx >= 0), the copper (Idx == -2) or everything (Idx == -1).
type LootMsg struct {
	ID  int `json:"id"`
	Idx int `json:"idx"`
}

type ChatMsg struct {
	Text string `json:"text"`
}
