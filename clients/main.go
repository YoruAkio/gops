package clients

import (
	"encoding/binary"
	"strings"

	log "github.com/codecat/go-libs/log"
	enet "github.com/eikarna/gotops"
	fn "github.com/eikarna/gotps/functions"
	items "github.com/eikarna/gotps/items"
	pkt "github.com/eikarna/gotps/packet"
	"github.com/eikarna/gotps/worlds"
)

func OnConnect(peer enet.Peer, host enet.Host, items *items.ItemInfo) {
	log.Info("New Client Connected %s", peer.GetAddress().String())
	pkt.SendPacket(peer, 1, "") //hello response
}

func OnDisConnect(peer enet.Peer, host enet.Host, items *items.ItemInfo) {
	log.Info("New Client Disconnected %s", peer.GetAddress().String())
}

func OnTextPacket(peer enet.Peer, host enet.Host, text string, items *items.ItemInfo) {
	if strings.Contains(text, "requestedName|") {
		fn.OnSuperMain(peer, items.GetItemHash())
	} else if len(text) > 6 && text[:6] == "action" {
		if strings.HasPrefix(text[7:], "enter_game") {
			fn.SendWorldMenu(peer)
		} else if strings.HasPrefix(text[7:], "join_request") {
			worldName := strings.ToUpper(strings.Split(text[25:], "\n")[0])
			fn.LogMsg(peer, "Sending you to world (%s) (%d)", worldName, len(worldName))
			OnEnterGameWorld(peer, host, worldName)
		}
	}

	log.Info("msg: %v", text)
}

var (
	SpawnX int
	SpawnY int
)

func OnEnterGameWorld(peer enet.Peer, host enet.Host, name string) {

	world, err := worlds.GetWorld(name)
	if err != nil {
		log.Error(err.Error())
	}
	nameLen := len(world.Name)
	totalPacketLen := 78 + nameLen + len(world.Tiles) + 24 + (8*len(world.Tiles) + (0 * 16))
	worldPacket := make([]byte, totalPacketLen)
	worldPacket[0] = 4  //game message
	worldPacket[4] = 4  //world packet type
	worldPacket[16] = 8 //char state
	worldPacket[66] = byte(len(world.Name))
	copy(worldPacket[68:], []byte(world.Name))

	worldPacket[nameLen+68] = byte(world.SizeX)
	worldPacket[nameLen+72] = byte(world.SizeY)
	binary.LittleEndian.PutUint16(worldPacket[nameLen+76:], uint16(world.TotalTiles))
	extraDataPos := 85 + nameLen

	for i := 0; i < int(world.TotalTiles); i++ {
		binary.LittleEndian.PutUint16(worldPacket[extraDataPos:], uint16(world.Tiles[i].Fg))
		binary.LittleEndian.PutUint16(worldPacket[extraDataPos+2:], uint16(world.Tiles[i].Bg))
		binary.LittleEndian.PutUint32(worldPacket[extraDataPos+4:], uint32(world.Tiles[i].Flags))

		switch world.Tiles[i].Fg {
		case 6:
			{
				worldPacket[extraDataPos+8] = 1 //block types
				binary.LittleEndian.PutUint16(worldPacket[extraDataPos+4:], uint16(len(world.Tiles[i].Label)))
				copy(worldPacket[extraDataPos+11:], []byte(world.Tiles[i].Label))
				SpawnX = (i % int(world.SizeX)) * 32
				SpawnY = (i / int(world.SizeX)) * 32
				extraDataPos += 4 + len(world.Tiles[i].Label)
				totalPacketLen += 4 + len(world.Tiles[i].Label)
				fn.LogMsg(peer, "x: %d, y: %d", SpawnX, SpawnY)
			}
		default:
			{
				break
			}
		}

		extraDataPos += 8
		totalPacketLen += 8
	}

	packet, err := enet.NewPacket(worldPacket, enet.PacketFlagReliable)
	if err != nil {
		panic(err)
	}
	peer.SendPacket(packet, 0)

	fn.OnSpawn(peer, 1, 1, int32(SpawnX), int32(SpawnY), "`6@Haikal_999", "id", false, true, true, true)
}
