package main

import (
	"bitbucket.org/polyu-named-data-network/ndn/packet"
	"bitbucket.org/polyu-named-data-network/ndn/packet/contentname"
	"bitbucket.org/polyu-named-data-network/ndn/packet/datatype"
	"bitbucket.org/polyu-named-data-network/ndn/packet/packettype"
	"bitbucket.org/polyu-named-data-network/ndn/packet/returncode"
	"bitbucket.org/polyu-named-data-network/ndn/utils"
	"encoding/json"
	"fmt"
	"github.com/aabbcc1241/goutils/log"
	"net"
)

func init() {
	log.Init(true, true, true, log.ShortCommFlag)
}
func main() {
	fmt.Println("NDN application demo - ping-pong client start")

	/* connect to proxy */
	conn, err := net.Dial("tcp", "127.0.0.1:8123")
	if err != nil {
		log.Error.Println("failed to connect to proxy", err)
		panic(1)
	}

	encoder := json.NewEncoder(conn)
	decoder := json.NewDecoder(conn)

	/* request data */
	log.Info.Println("preparing request packet")
	p1 := packet.InterestPacket_s{
		ContentName: contentname.ContentName_s{
			Name:           "ping",
			ContentType:    contentname.ExactMatch,
			ContentParam:   nil,
			AcceptDataType: []datatype.Base{datatype.RAW},
		},
		SeqNum:             1,
		AllowCache:         true,
		PublisherPublicKey: utils.ZeroKey,
	}
	p2, err := p1.ToGenericPacket()
	if err != nil {
		log.Error.Println("failed to marshal interest data", err)
		panic(2)
	}
	if err := encoder.Encode(p2); err != nil {
		log.Error.Println("failed to sent packet", err)
		panic(3)
	}
	log.Info.Println("sent request packet")

	/* wait for data */
	log.Info.Println("waiting for data")
	var in_packet packet.GenericPacket_s
	if err := decoder.Decode(&in_packet); err != nil {
		log.Error.Println("failed decode income packet", err)
		panic(4)
	}
	if in_packet.PacketType == packettype.DataPacket_c {
		var p packet.DataPacket_s
		err := json.Unmarshal(in_packet.Payload, &p)
		if err != nil {
			log.Error.Println("failed to parse data packet")
			panic(5)
		}
		log.Info.Println("received data:", string(p.ContentData))
	} else if in_packet.PacketType == packettype.InterestReturnPacket_c {
		var p packet.InterestReturnPacket_s
		if err := json.Unmarshal(in_packet.Payload, &p); err != nil {
			log.Error.Println("failed to parse interest return packet")
			panic(6)
		}
		reason := "replease check returncode.go"
		if p.ReturnCode == returncode.NoRoute {
			reason = "NoRoute"
		}
		log.Info.Println("interest return, resultcode:", p.ReturnCode, "reason:", reason)
	} else {
		log.Error.Println("unexpected packet", in_packet)
	}

	fmt.Println("NDN application demo - ping-pong client end")
}
