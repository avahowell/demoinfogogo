/*
Created on June 13, 2015
Author: me@johnathanhowell.com
License: MIT
*/
package main

import (
	"encoding/binary"
	"fmt"
	"github.com/golang/protobuf/proto"
	"log"
	"os"
)

const (
	MAX_OSPATH              = 260
	PACKET_OFFSET           = 160
	DEMO_HEADER_ID          = "HL2DEMO"
	DEM_SIGNON              = 1
	DEM_PACKET              = 2
	DEM_SYNCTICK            = 3
	DEM_CONSOLECMD          = 4
	DEM_USERCMD             = 5
	DEM_DATATABLES          = 6
	DEM_STOP                = 7
	DEM_CUSTOMDATA          = 8
	DEM_STRINGTABLES        = 9
	DEM_LASTCOMMAND         = 9
	MSG_SERVER_INFO         = 8
	MSG_DATA_TABLE          = 0
	MSG_CREATE_STRING_TABLE = 12
	MSG_UPDATE_STRING_TABLE = 13
	MSG_USER_MESSAGE        = 23
	MSG_GAME_EVENT          = 25
	MSG_PACKET_ENTITIES     = 26
	MSG_GAME_EVENTS_LIST    = 30
)

type demoheader struct {
	Demofilestamp   [8]byte
	Demoprotocol    int32
	Networkprotocol int32
	Servername      [MAX_OSPATH]byte
	Clientname      [MAX_OSPATH]byte
	Mapname         [MAX_OSPATH]byte
	Gamedirectory   [MAX_OSPATH]byte
	Playback_time   float32
	Playback_ticks  int32
	Playback_frames int32
	Signonlength    int32
}
type democmdheader struct {
	Cmd        byte
	Tick       int32
	Playerslot byte
}
type demofile struct {
	header demoheader
	tick   int32
	frame  int32
	stream *demostream
}

func (d *demofile) PrintInfo() {
	fmt.Printf("Map: %s\n", d.header.Mapname)
	fmt.Printf("Ticks: %d\n", d.header.Playback_ticks)
	fmt.Printf("Game Directory: %s\n", d.header.Gamedirectory)
	fmt.Printf("Client name: %s\n", d.header.Clientname)
	fmt.Printf("Playback time: %f seconds\n", d.header.Playback_time)
	fmt.Printf("Server Name: %s\n", d.header.Servername)
	fmt.Printf("Signon Length: %d\n", d.header.Signonlength)
	fmt.Printf("Frames: %d\n", d.header.Playback_frames)
	fmt.Printf("Ticks: %d\n", d.header.Playback_ticks)
}
func (d *demofile) readCommandHeader() democmdheader {
	return democmdheader{Cmd: d.stream.GetByte(),
		Tick:       d.stream.GetInt(),
		Playerslot: d.stream.GetByte()}
}
func (d *demofile) parseProtobufMessage() {
	messagetype := d.stream.GetVarInt()
	fmt.Printf("messagetype: %d\n", messagetype)
	messagelen := d.stream.GetVarInt()
	switch messagetype {
	case MSG_SERVER_INFO:
		msg := CSVCMsg_ServerInfo{}
		buf := make([]byte, messagelen)
		_, err := d.stream.Read(buf)
		if err != nil {
			panic(err)
		}
		err = proto.Unmarshal(buf, &msg)
		if err != nil {
			panic(err)
		}
		fmt.Println(msg)

		//	default:
		//	skiplen := d.stream.GetVarInt()
		//d.f.Seek(skiplen, 1)
	}
}
func (d *demofile) readPacket() {
	d.stream.Skip(160)
	blocksize := d.stream.GetInt()
	fmt.Printf("CHUNK SIZE: %d\n", blocksize)
	d.parseProtobufMessage()
}
func (d *demofile) GetFrame() {
	cmdheader := d.readCommandHeader()
	switch cmdheader.Cmd {
	case DEM_SIGNON, DEM_PACKET:
		fmt.Println("Got packet")
		d.readPacket()
	case DEM_SYNCTICK:
		// handle synctick
	case DEM_CONSOLECMD:
		fmt.Println("Got consolecmd")
		// handle consolecommand
	case DEM_USERCMD:
		fmt.Println("got usercmd")
		// handle usercommand
	case DEM_DATATABLES:
		fmt.Println("got datatables")
		// handle datatables
	case DEM_STOP:
		fmt.Println("GOT STOP")
		return
		// handle stop
	case DEM_CUSTOMDATA:
		fmt.Println("Got customdata")
		// handle customdata
	case DEM_STRINGTABLES:
		fmt.Println("Got stringtables")
		// handle stringtables
	}
	d.frame++
}
func (d *demofile) Open(path string) {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	d.stream = NewDemoStream(f)
	d.header = demoheader{}
	err = binary.Read(d.stream, binary.LittleEndian, &d.header)
	if err != nil {
		panic(err)
	}
	if string(d.header.Demofilestamp[:7]) != DEMO_HEADER_ID {
		log.Fatal("Invalid demo header, are you sure this is a .dem?\n")
	}
	d.tick = 0
	d.frame = 0
}
func usage() {
	fmt.Printf("Usage: %s [demo.dem]\n", os.Args[0])
	os.Exit(2)
}
func main() {
	if len(os.Args) != 2 {
		usage()
	}
	d := demofile{}
	d.Open(os.Args[1])
	d.PrintInfo()
	d.GetFrame()
}
