package main

import (
  "bitbucket.org/polyu-named-data-network/ndn/packet"
  "bitbucket.org/polyu-named-data-network/ndn/packet/contentname"
  "encoding/json"
  "fmt"
  "github.com/aabbcc1241/goutils/log"
  "net"
  "strconv"
  "sync"
)

func init() {
  log.Init(true, true, true, log.ShortCommFlag)
}
func main() {
  fmt.Println("NDN application demo - ping-pong client start")

  wg := sync.WaitGroup{}

  /* establish data connection */
  dataConn, err := net.Dial("tcp", "127.0.0.1:8124")
  if err != nil {
    fmt.Println("failed to connect to proxy data service", err)
    panic(3)
  }
  fmt.Println("started data socket", dataConn.LocalAddr().String())
  _, dataPort_string, _ := net.SplitHostPort(dataConn.LocalAddr().String())
  dataPort, _ := strconv.Atoi(dataPort_string)
  wg.Add(1)
  go func() {
    defer wg.Done()
    var in_packet packet.DataPacket_s
    log.Debug.Println("data connection", dataConn.LocalAddr())
    /* wait for data packet (response)*/
    fmt.Println("wait for data packet")
    json.NewDecoder(dataConn).Decode(&in_packet)
    fmt.Println("received data packet", in_packet)
    fmt.Println("data content (string)", string(in_packet.ContentData))
  }()

  /* establish interest return connection */
  //interestReturnConn,err:=net.Dial("tcp","127.0.0.1:8123")

  /* establish interest connection */
  interestConn, err := net.Dial("tcp", "127.0.0.1:8123")
  if err != nil {
    fmt.Println("failed to connect to proxy interest service", err)
    panic(1)
  }

  fmt.Println("preparing interest packet")
  out_packet := packet.InterestPacket_s{
    ContentName: contentname.ContentName_s{
      Name: "ping",
      Type: contentname.ExactMatch,
    },
    SeqNum:     1,
    AllowCache: false,
    DataPort:   dataPort,
  }
  err = json.NewEncoder(interestConn).Encode(out_packet)
  if err != nil {
    fmt.Println("failed to encode interest packet", err)
    panic(2)
  }
  fmt.Println("sent interest packet")

  /* prepare interest packet (request) */
  wg.Add(1)
  go func() {
    defer wg.Done()
    var in_packet packet.InterestReturnPacket_s
    fmt.Println("wait for interestReturn packet")
    json.NewDecoder(interestConn).Decode(&in_packet)
    fmt.Println("received interestReturn packet", in_packet)
  }()

  /* wait for interest return (NAK) */

  wg.Wait()
  fmt.Println("NDN application demo - ping-pong client end")
}
