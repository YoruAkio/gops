package functions

import (
	"encoding/binary"
	"fmt"
	"strconv"

	"github.com/codecat/go-libs/log"
	enet "github.com/eikarna/gotops"
	"github.com/eikarna/gotps/functions/variants"
	variant "github.com/eikarna/gotps/functions/variants"
	pkt "github.com/eikarna/gotps/packet"
	tankpacket "github.com/eikarna/gotps/packet/TankPacket"
	player "github.com/eikarna/gotps/players"
	"github.com/eikarna/gotps/utils"
	"github.com/eikarna/gotps/worlds"
)

var ListActiveWorld = make(map[string]int)

func OnRemove(peer enet.Peer, netid uint32) {
	log.Warn("OnRemove netID: %s", fmt.Sprint(netid))
	variant := variant.NewVariant(0, -1)
	variant.InsertString("OnRemove")
	variant.InsertString("netID|" + fmt.Sprint(netid))
	variant.Send(peer)
}

func OnDialogRequest(peer enet.Peer, dialog string, delay int) {
	variant := variant.NewVariant(delay, -1)
	variant.InsertString("OnDialogRequest")
	variant.InsertString(dialog)
	variant.Send(peer)
}

func ModifyInventory(peer enet.Peer, itemId int, count int, pl *player.Player) {
	/*if player.NotSafePlayer(peer) {
		return
	}*/
	inv := pl.Inventory
	for i := range inv {
		if inv[i].ID == itemId {
			if inv[i].Qty > 0 {
				inv[i].Qty += int16(count)
				log.Info("%d", inv[i].Qty)
			} else {
				inv = append(inv[:i-1], inv[i+1:]...)
				log.Info("Del Inv: %#v", inv)
			}
		}
	}
	pl.Inventory = inv
	ReducePack := &tankpacket.TankPacket{
		PacketType: 13,
		Value:      uint32(itemId),
	}
	ReducedPacket := ReducePack.Serialize(56, true)
	Packet, err := enet.NewPacket(ReducedPacket, enet.PacketFlagReliable)
	if err != nil {
		log.Error("Packet type 13:", err.Error())
	}
	peer.SendPacket(Packet, 0)
	pl, ReducePack, ReducedPacket, Packet = nil, nil, nil, nil
}

func OnWrench(peer enet.Peer, Tank *tankpacket.TankPacket, name string) {

}

func AddTile(peer enet.Peer, Tank *tankpacket.TankPacket) {
	Coords := Tank.PunchX + (Tank.PunchY * uint32(worlds.Worlds[player.PInfo(peer).CurrentWorld].SizeX))
	PlacePack := &tankpacket.TankPacket{
		PacketType:     3,
		NetID:          player.PInfo(peer).NetID,
		CharacterState: Tank.CharacterState,
		Value:          Tank.Value,
		X:              Tank.X,
		Y:              Tank.Y,
		XSpeed:         Tank.XSpeed,
		YSpeed:         Tank.YSpeed,
		PunchX:         Tank.PunchX,
		PunchY:         Tank.PunchY,
	}
	PlacePacket := PlacePack.Serialize(56, true)
	worlds.Worlds[player.PInfo(peer).CurrentWorld].Tiles[Coords].Fg = int16(Tank.Value)
	worlds.Worlds[player.PInfo(peer).CurrentWorld].Tiles[Coords].Bg = int16(Tank.Value)
	TalkBubble(peer, player.PInfo(peer).NetID, 100, false, "ID: %d, Qty: %d", Tank.Value, player.GetCountItemFromInventory(peer, int(Tank.Value)))
	Packet, err := enet.NewPacket(PlacePacket, enet.PacketFlagReliable)
	if err != nil {
		log.Error("Error Packet 3:", err)
	}
	for _, currentPeer := range player.GetPeers(player.PlayerMap) {
		if player.NotSafePlayer(currentPeer) {
			continue
		}
		if player.PInfo(peer).CurrentWorld == player.PInfo(currentPeer).CurrentWorld {
			currentPeer.SendPacket(Packet, 0)
		}
	}
	Tank, PlacePacket, Packet = nil, nil, nil
}

/*
func OnPunch(peer enet.Peer, Tank *tankpacket.TankPacket, world *worlds.World) {
	//player.PlayerMap[peer].PunchX = Tank.PunchX
	//player.PlayerMap[peer].PunchY = Tank.PunchY
	/* test, err := worlds.GetWorld(player.PlayerMap[peer].CurrentWorld)
	if err != nil {
		return
	}*/
//test, ok := worlds.Worlds[player.PlayerMap[peer].CurrentWorld]
/*world, err := worlds.LoadWorld(db, name)
	if err != nil {
		peer.DisconnectLater(0)
		log.Error("Worlds with name: %s is not found in our database!", name)
	}
	Coords := Tank.PunchX + (Tank.PunchY * uint32(world.SizeX))
	ConsoleMsg(peer, 0, "PunchX: %d, PunchY: %d, TotalXY: %d", Tank.PunchX, Tank.PunchY, Coords)
	switch world.Tiles[Coords].Fg {

	case 6, 8:
		{
			TalkBubble(peer, player.PlayerMap[peer].NetID, 0, false, "It's too strong to break.")
			return
			break
		}
	default:
		{
			break
		}
	}
	switch world.Tiles[Coords].Bg {
	default:
		{
			break
		}
	}
	// TalkBubble(peer, player.PInfo(peer).NetID, 100, false, "Error! Tile Data is not valid at X:%d, Y:%d, Tiles:%d", Tank.PunchX, Tank.PunchY, Coords)
	testt := &tankpacket.TankPacket{
		PacketType:     3,
		NetID:          player.PInfo(peer).NetID,
		CharacterState: Tank.CharacterState,
		Value:          Tank.Value,
		X:              Tank.X,
		Y:              Tank.Y,
		XSpeed:         Tank.XSpeed,
		YSpeed:         Tank.YSpeed,
		PunchX:         Tank.PunchX,
		PunchY:         Tank.PunchY,
	}
	bbb := testt.Serialize(56, true)
	aaa, err := enet.NewPacket(bbb, enet.PacketFlagReliable)
	if err != nil {
		log.Error("Error Packet 3:", err)
	}
	for _, currentPeer := range player.GetPeers(player.PlayerMap) {
		if player.NotSafePlayer(currentPeer) {
			continue
		}
		if player.PlayerMap[peer].CurrentWorld == player.PlayerMap[currentPeer].CurrentWorld {
			currentPeer.SendPacket(aaa, 0)
		}
	}
	worlds.Worlds[world.Name] = world
	//UpdateWorld(peer, player.PInfo(peer).CurrentWorld)
	LogMsg(peer, "[Punch/Place] X:%d, Y:%d, Value:%d, NetID:%d", Tank.PunchX, Tank.PunchY, Tank.Value, Tank.NetID)
	// worlds.Worlds[name].Tiles[Coords] = world.Tiles[Coords]
	return
}*/

/*func OnPunch(peer enet.Peer, Tank *tankpacket.TankPacket, world *worlds.World) {
	Coords := Tank.PunchX + (Tank.PunchY * uint32(world.SizeX))
	ConsoleMsg(peer, 0, "PunchX: %d, PunchY: %d, TotalXY: %d", Tank.PunchX, Tank.PunchY, Coords)

	// Check if the tile being punched is not blank
	if world.Tiles[Coords].Fg != 0 {
		// Check if the block is unbreakable
		if world.Tiles[Coords].Fg == 6 || world.Tiles[Coords].Fg == 8 {
			TalkBubble(peer, player.PlayerMap[peer].NetID, 0, false, "It's too strong to break.")
			return
		}

		// Update the tile to be blank
		world.Tiles[Coords].Fg = 0
		// Optionally, you may want to update other properties of the tile like Bg, Flags, Label, etc. if needed.

		// Inform players about the change
		testt := &tankpacket.TankPacket{
			PacketType:     3,
			NetID:          player.PInfo(peer).NetID,
			CharacterState: Tank.CharacterState,
			Value:          Tank.Value,
			X:              Tank.X,
			Y:              Tank.Y,
			XSpeed:         Tank.XSpeed,
			YSpeed:         Tank.YSpeed,
			PunchX:         Tank.PunchX,
			PunchY:         Tank.PunchY,
		}
		bbb := testt.Serialize(56, true)
		aaa, err := enet.NewPacket(bbb, enet.PacketFlagReliable)
		if err != nil {
			log.Error("Error Packet 3:", err)
		}
		for _, currentPeer := range player.GetPeers(player.PlayerMap) {
			if player.NotSafePlayer(currentPeer) {
				continue
			}
			if player.PlayerMap[peer].CurrentWorld == player.PlayerMap[currentPeer].CurrentWorld {
				currentPeer.SendPacket(aaa, 0)
			}
		}

		// Update the world
		worlds.Worlds[world.Name] = world

		// Log the action
		LogMsg(peer, "[Punch/Place] X:%d, Y:%d, Value:%d, NetID:%d", Tank.PunchX, Tank.PunchY, Tank.Value, Tank.NetID)
	}
}*/

func PunchLoop(peer enet.Peer, Tank *tankpacket.TankPacket, world *worlds.World) {
	//coordsListX := make([]uint32, far)
	//coordsListY := make([]uint32, far)
	//coordsList := make([]uint32, far)
	//listPacket := make([]Tank, far)
	/*for i := range far {
	farX := uint32((Tank.X - float32(world.SizeX)) / 32)
	if farX < Tank.PunchX {
		ConsoleMsg(peer, 0, "PunchLoop at X: %d, PunchX: %d", farX, Tank.PunchX)
		coordsListX = append(coordsListX, Tank.PunchX-uint32(i))
	} else if farX > Tank.PunchX {
		ConsoleMsg(peer, 0, "PunchLoop at X: %d, PunchX: %d", farX, Tank.PunchX)
		coordsListX = append(coordsListX, Tank.PunchX+uint32(i))
	} /*else {
		coordsListX = append(coordsListX, Tank.PunchX*uint32(i))
	}
	/*if uint32(Tank.Y) < Tank.PunchY {
		coordsListY = append(coordsListY, Tank.PunchY+uint32(i))
	} else if uint32(Tank.Y) > Tank.PunchY {
		coordsListY = append(coordsListY, Tank.PunchY-uint32(i))
	} else {
		coordsListY = append(coordsListY, Tank.PunchY*uint32(i))
	}*/
	// }
	// for i := range coordsListX {
	// PunchY: coordsListY[coordsListX[i]],
	// }
	/*for i := range coordsListY {
		world.Tiles[i].Fg = 0
		world.Tiles[i].Bg = 0
		testt := &tankpacket.TankPacket{
			PacketType:     3,
			NetID:          player.PInfo(peer).NetID,
			CharacterState: Tank.CharacterState,
			Value:          Tank.Value,
			X:              Tank.X,
			Y:              Tank.Y,
			XSpeed:         Tank.XSpeed,
			YSpeed:         Tank.YSpeed,
			// PunchX:         Tank.PunchX,
			PunchX: Tank.PunchX - uint32(i),
			PunchY: coordsListY[i],
		}
		bbb := testt.Serialize(56, true)
		aaa, err := enet.NewPacket(bbb, enet.PacketFlagReliable)
		if err != nil {
			log.Error("Error Packet 3:", err)
		}
		for _, currentPeer := range player.GetPeers(player.PlayerMap) {
			if player.NotSafePlayer(currentPeer) {
				continue
			}
			if player.PInfo(peer).CurrentWorld == player.PInfo(currentPeer).CurrentWorld {
				currentPeer.SendPacket(aaa, 0)
				break
			}
		}
	}*/
}

func OnPunch(peer enet.Peer, Tank *tankpacket.TankPacket, world *worlds.World) {
	Coords := Tank.PunchX + (Tank.PunchY * uint32(world.SizeX))
	//ConsoleMsg(peer, 0, "PunchX: %d, PunchY: %d, TotalXY: %d", Tank.PunchX, Tank.PunchY, Coords)
	if world.Tiles[Coords].Fg == 6 || world.Tiles[Coords].Fg == 8 {
		TalkBubble(peer, player.PInfo(peer).NetID, 0, false, "It's too strong to break.")
		return
	}
	worlds.Worlds[world.Name].Tiles[Coords].Fg = 0
	worlds.Worlds[world.Name].Tiles[Coords].Bg = 0
	Tank.NetID = player.PInfo(peer).NetID
	Tank.X = player.PInfo(peer).PosX
	Tank.Y = player.PInfo(peer).PosY
	bbb := Tank.Serialize(56, true)
	aaa, err := enet.NewPacket(bbb, enet.PacketFlagReliable)
	if err != nil {
		log.Error("Error Packet 3:", err)
	}
	for _, currentPeer := range player.GetPeers(player.PlayerMap) {
		if player.NotSafePlayer(currentPeer) {
			continue
		}
		if player.PInfo(peer).CurrentWorld == player.PInfo(currentPeer).CurrentWorld {
			currentPeer.SendPacket(aaa, 0)
		}
	}
	Tank, bbb, aaa, world = nil, nil, nil, nil
}

func SendWorldMenu(peer enet.Peer) {
	var world_packet string
	// Add World Start as default
	ListActiveWorld["START"] = 65
	world_packet += "add_filter|\n"
	world_packet += "add_heading|Top Worlds<ROW2>|\n"
	for listworld, count := range ListActiveWorld {
		if count > 0 {
			if listworld == "START" {
				world_packet += "add_floater|" + listworld + "|" + listworld + "|" + strconv.Itoa(count) + "|0.8|3529161471\n"

			} else {
				world_packet += "add_floater|" + listworld + "|" + listworld + "|" + strconv.Itoa(count) + "|0.5|3529161471\n"
			}
		} else {
			delete(ListActiveWorld, listworld)
		}
	}
	world_packet += "add_heading|Credits<CR>|\n"
	world_packet += "add_floater|KIPASGTS|KIPASGTS|0|0.5|2147418367\n"
	world_packet += "add_floater|EIKARNA|EIKARNA|0|0.5|2147418367\n"
	world_packet += "add_floater|AKBARDEV|AKBARDEV|0|0.5|2147418367\n"
	world_packet += "add_floater|TEAMNEVOLUTION|TEAMNEVOLUTION|0|0.5|2147418367\n"
	world_packet += "add_heading|Based On: https://github.com/eikarna/gotops<CR>|\n"

	variant := variant.NewVariant(0, -1) //delay netid
	variant.InsertString("OnRequestWorldSelectMenu")
	variant.InsertString(world_packet)
	variant.Send(peer)
}

func UpdateName(peer enet.Peer, name string) {
	if player.NotSafePlayer(peer) {
		return
	}
	pl := player.PInfo(peer)
	variant := variant.NewVariant(0, int(pl.NetID))
	variant.InsertString("OnNameChanged")
	variant.InsertString(name)
	if pl.CurrentWorld != "" {
		for _, currentPeer := range player.GetPeers(player.PlayerMap) {
			if player.PInfo(currentPeer).CurrentWorld == pl.CurrentWorld {
				variant.Send(currentPeer)
				break
			}
		}
	} else {
		variant.Send(peer)
	}
}

func TextOverlay(peer enet.Peer, text string) {
	if player.NotSafePlayer(peer) {
		return
	}
	variant := variant.NewVariant(0, -1)
	variant.InsertString("OnTextOverlay")
	variant.InsertString(text)
	variant.Send(peer)
}

func SetHasGrowID(peer enet.Peer) {
	if player.NotSafePlayer(peer) {
		return
	}
	pl := player.PInfo(peer)
	variant := variant.NewVariant(0, -1)
	variant.InsertString("SetHasGrowID")
	variant.InsertInt(1)
	variant.InsertString(pl.TankIDName)
	variant.InsertString(pl.TankIDPass)
	variant.Send(peer)
}

func UpdateWorld(peer enet.Peer, name string) {
	if player.NotSafePlayer(peer) {
		LogMsg(peer, "`4Invalid Player Data!``")
		return
	}
	//name := player.PInfo(peer).CurrentWorld
	wi := worlds.Worlds[name]
	if wi == nil {
		worlds.Worlds[name] = worlds.GenerateWorld(name, 100, 60)
		wi = worlds.Worlds[name]
		log.Error("Worlds with name: %s is not found in our database!", name)
	}
	nameLen := len(wi.Name)
	totalPacketLen := 78 + nameLen + len(wi.Tiles) + 24 + (8*len(wi.Tiles) + (0 * 16))
	worldPacket := make([]byte, totalPacketLen)
	worldPacket[0] = 4  //game message
	worldPacket[4] = 4  //world packet type
	worldPacket[16] = 8 //char state
	worldPacket[66] = byte(len(wi.Name))
	copy(worldPacket[68:], []byte(wi.Name))

	worldPacket[nameLen+68] = byte(wi.SizeX)
	worldPacket[nameLen+72] = byte(wi.SizeY)
	binary.LittleEndian.PutUint16(worldPacket[nameLen+76:], uint16(wi.TotalTiles))
	extraDataPos := 85 + nameLen

	for i := 0; i < int(wi.TotalTiles); i++ {
		// log.Info("Loaded Tiles: %v", world.Tiles[i])
		binary.LittleEndian.PutUint16(worldPacket[extraDataPos:], uint16(wi.Tiles[i].Fg))
		binary.LittleEndian.PutUint16(worldPacket[extraDataPos+2:], uint16(wi.Tiles[i].Bg))
		binary.LittleEndian.PutUint32(worldPacket[extraDataPos+4:], uint32(wi.Tiles[i].Flags))
		switch wi.Tiles[i].Fg {
		case 6:
			{
				worldPacket[extraDataPos+8] = 1 //block types
				binary.LittleEndian.PutUint16(worldPacket[extraDataPos+9:], uint16(len(wi.Tiles[i].Label)))
				copy(worldPacket[extraDataPos+11:], []byte(wi.Tiles[i].Label))

				player.PInfo(peer).SpawnX = uint32((i % int(wi.SizeX)) * 32)
				player.PInfo(peer).SpawnY = uint32((i / int(wi.SizeX)) * 32)
				extraDataPos += 4 + len(wi.Tiles[i].Label)
				break
			}
		case 7188:
			{
				worldPacket[extraDataPos+8] = 3
				worldPacket[extraDataPos+9] = 128
				binary.LittleEndian.PutUint32(worldPacket[extraDataPos+10:], player.PInfo(peer).UserID)
				worldPacket[extraDataPos+14] = 0
				binary.LittleEndian.PutUint32(worldPacket[extraDataPos+18:], 0)
				binary.LittleEndian.PutUint32(worldPacket[extraDataPos+22:], 0)
				extraDataPos += 26
				break
			}
		default:
			{
				break
			}
		}

		extraDataPos += 8
	}
	packet, err := enet.NewPacket(worldPacket, enet.PacketFlagReliable)
	if err != nil {
		panic(err)
	}
	peer.SendPacket(packet, 0)
	worldPacket = nil
}

func UpdateInventory(peer enet.Peer) {
	if player.NotSafePlayer(peer) {
		return
	}
	pl := player.PInfo(peer)
	if len(pl.Inventory) < 1 || pl.InventorySize < 1 {
		//NewInvent := pl.Inventory
		if pl.InventorySize == 0 {
			pl.InventorySize = 32
		}
		itemsToAdd := []player.ItemInfo{
			{ID: 18, Qty: 1},
			{ID: 32, Qty: 1},
			{ID: 7188, Qty: 3},
			{ID: 2, Qty: 200},
			{ID: 5480, Qty: 1},
		}
		for _, item := range itemsToAdd {
			pl.Inventory = append(pl.Inventory, item)
		}
	}
	packetLen := 66 + (pl.InventorySize * 4) + 4
	d_ := make([]byte, packetLen)
	binary.LittleEndian.PutUint16(d_[0:], 4)
	binary.LittleEndian.PutUint16(d_[4:], 9)
	binary.LittleEndian.PutUint16(d_[8:], 255)
	binary.LittleEndian.PutUint16(d_[16:], 8)
	binary.LittleEndian.PutUint16(d_[56:], 6+(pl.InventorySize*4)+4)
	binary.LittleEndian.PutUint16(d_[60:], 1)
	binary.LittleEndian.PutUint16(d_[61:], pl.InventorySize)
	binary.LittleEndian.PutUint16(d_[65:], pl.InventorySize)
	offset := 67
	for _, Inven := range pl.Inventory {
		if Inven.Qty > 0 {
			binary.LittleEndian.PutUint16(d_[offset:], uint16(Inven.ID))
			offset += 2
			binary.LittleEndian.PutUint16(d_[offset:], uint16(Inven.Qty))
			offset += 2
		}
	}
	//}
	// log.Info("SendInventory Byte: %b | String: %s", d_, d_)
	packet, err := enet.NewPacket(d_, enet.PacketFlagReliable)
	if err != nil {
		log.Error(err.Error())
	}
	peer.SendPacket(packet, 0)
	d_ = nil
}

/*func SendDoor(peer enet.Peer) {
	if player.GetPlayer(peer).CurrentWorld != "" {
		OnPlayerExitWorld(peer)
	}
	break
}*/

func ConsoleMsg(peer enet.Peer, delay int, a ...interface{}) {
	msg := fmt.Sprintf(a[0].(string), a[1:]...)
	variant := variant.NewVariant(delay, -1)
	variant.InsertString("OnConsoleMessage")
	variant.InsertString(msg)
	variant.Send(peer)
}

func TalkBubble(peer enet.Peer, netID uint32, delay int, isOverlay bool, a ...interface{}) {
	msg := fmt.Sprintf(a[0].(string), a[1:]...)
	variant := variant.NewVariant(delay, -1)
	variant.InsertString("OnTalkBubble")
	variant.InsertUnsignedInt(netID)
	variant.InsertString(msg)
	variant.InsertInt(utils.BoolToInt(isOverlay))
	variant.InsertInt(utils.BoolToInt(isOverlay))
	variant.Send(peer)
}

func OnSuperMain(peer enet.Peer, itemHash uint32) {

	variant := variant.NewVariant(0, -1)
	variant.InsertString("OnSuperMainStartAcceptLogonHrdxs47254722215a")
	variant.InsertUnsignedInt(itemHash)
	variant.InsertString("www.growtopia1.com")
	variant.InsertString("cache/")
	variant.InsertString("cc.cz.madkite.freedom org.aqua.gg idv.aqua.bulldog com.cih.gamecih2 com.cih.gamecih com.cih.game_cih cn.maocai.gamekiller com.gmd.speedtime org.dax.attack com.x0.strai.frep com.x0.strai.free org.cheatengine.cegui org.sbtools.gamehack com.skgames.traffikrider org.sbtoods.gamehaca com.skype.ralder org.cheatengine.cegui.xx.multi1458919170111 com.prohiro.macro me.autotouch.autotouch com.cygery.repetitouch.free com.cygery.repetitouch.pro com.proziro.zacro com.slash.gamebuster")
	variant.InsertString("proto=207|choosemusic=audio/mp3/about_theme.mp3|active_holiday=6|wing_week_day=0|ubi_week_day=2|server_tick=123665344|clash_active=0|drop_lavacheck_faster=1|isPayingUser=2|usingStoreNavigation=1|enableInventoryTab=1|bigBackpack=1|m_clientBits=0|eventButtons={\"EventButtonData\":[{\"Components\":[{\"Enabled\":true,\"Id\":\"Overlay\",\"Parameters\":\"target_child_entity_name:overlay_layer;var_name:alpha;target:0;interpolation:1;on_finish:1;duration_ms:1000;delayBeforeStartMS:1000\",\"Type\":\"InterpolateComponent\"}],\"DialogName\":\"openLnySparksPopup\",\"IsActive\":false,\"Name\":\"LnyButton\",\"Priority\":1,\"Text\":\"0/5\",\"TextOffset\":\"0.01,0.2\",\"Texture\":\"interface/large/event_button3.rttex\",\"TextureCoordinates\":\"0,2\"},{\"Components\":[{\"Enabled\":true,\"Parameters\":\"\",\"Type\":\"RenderDailyChallengeComponent\"}],\"DialogName\":\"dailychallengemenu\",\"IsActive\":false,\"Name\":\"DailyChallenge\",\"Priority\":2},{\"Components\":[{\"Enabled\":false,\"Id\":\"Overlay\",\"Parameters\":\"target_child_entity_name:overlay_layer;var_name:alpha;target:0;interpolation:1;on_finish:1;duration_ms:1000;delayBeforeStartMS:1000\",\"Type\":\"InterpolateComponent\"}],\"DialogName\":\"openStPatrickPiggyBank\",\"IsActive\":false,\"Name\":\"StPatrickPBButton\",\"Priority\":1,\"Text\":\"0/0\",\"TextOffset\":\"0.00,0.05\",\"Texture\":\"interface/large/event_button4.rttex\",\"TextureCoordinates\":\"0,0\"},{\"DialogName\":\"show_bingo_ui\",\"IsActive\":false,\"Name\":\"Bingo_Button\",\"Priority\":1,\"Texture\":\"interface/large/event_button4.rttex\"}]}")
	//p.Insert("654171113"); //tribute_data
	variant.Send(peer)
}

func LogMsg(peer enet.Peer, a ...interface{}) {
	msg := fmt.Sprintf(a[0].(string), a[1:]...)
	pkt.SendPacket(peer, 3, "action|log\nmsg|"+msg)
}

func PlayMsg(peer enet.Peer, delay int, name string) {
	msg := "action|play_sfx\nfile|" + name + "\ndelayMS|" + strconv.Itoa(delay)
	pkt.SendPacket(peer, 3, msg)
}

func OnSpawn(peer enet.Peer, netid int16, userid uint32, posX int32, posY int32, username string, country string, invis bool, mstate bool, smsate bool, local bool) {
	spawnAvatar := "spawn|avatar\n"
	spawnAvatar += "netID|" + strconv.Itoa(int(netid)) + "\n"
	spawnAvatar += "userID|" + strconv.Itoa(int(userid)) + "\n"
	spawnAvatar += "colrect|0|0|20|30\n"
	spawnAvatar += "posXY|" + strconv.Itoa(int(posX)) + "|" + strconv.Itoa(int(posY)) + "\n"
	spawnAvatar += "name|" + username + "\n"
	spawnAvatar += "country|" + country + "\n"
	spawnAvatar += "invis|" + utils.BoolToIntString(invis) + "\n"    //1 0
	spawnAvatar += "mstate|" + utils.BoolToIntString(mstate) + "\n"  //1 0
	spawnAvatar += "smstate|" + utils.BoolToIntString(smsate) + "\n" //1 0
	if local {
		spawnAvatar += "onlineID|\ntype|local\n"
	}

	variant := variant.NewVariant(0, -1)
	variant.InsertString("OnSpawn")
	variant.InsertString(spawnAvatar)
	variant.Send(peer)
}

func UpdateClothes(delay int, peer, otherPeer enet.Peer) {
	pData := player.PInfo(peer)
	variant := variant.NewVariant(delay, int(player.PInfo(peer).NetID))
	variant.InsertString("OnSetClothing")
	variant.InsertTripleFloat(pData.Clothes.Hair, pData.Clothes.Shirt, pData.Clothes.Pants)
	variant.InsertTripleFloat(pData.Clothes.Feet, pData.Clothes.Face, pData.Clothes.Hand)
	variant.InsertTripleFloat(pData.Clothes.Back, pData.Clothes.Mask, pData.Clothes.Necklace)
	variant.InsertInt(pData.SkinColor)
	variant.InsertTripleFloat(0, 0, 0)
	variant.Send(otherPeer)
}

func SetAccountHasSecured(peer enet.Peer) {
	variant := variant.NewVariant(0, -1)
	variant.InsertString("SetAccountHasSecured")
	variant.InsertInt(1)
	variant.Send(peer)
}

func SetRespawnPos(peer enet.Peer, pos int, delay int) {
	variant := variant.NewVariant(delay, int(player.PInfo(peer).NetID))
	variant.InsertString("SetRespawnPos")
	variant.InsertInt(pos)
	variant.Send(peer)
}

func OnSetFreezeState(peer enet.Peer, yes bool, delay int) {
	variant := variants.NewVariant(delay, int(player.PInfo(peer).NetID))
	variant.InsertString("OnSetFreezeState")
	if yes {
		variant.InsertInt(1)
	} else {
		variant.InsertInt(0)
	}
	variant.Send(peer)
}
