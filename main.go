package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
)

const (
	MAX_OSPATH     = 260
	DEMO_HEADER_ID = "HL2DEMO"
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
type demofile struct {
	header demoheader
}

func (d *demofile) PrintInfo() {
	fmt.Printf("Map: %s\n", d.header.Mapname)
	fmt.Printf("Ticks: %d\n", d.header.Playback_ticks)
	fmt.Printf("Game Directory: %s\n", d.header.Gamedirectory)
	fmt.Printf("Client name: %s\n", d.header.Clientname)
	fmt.Printf("Playback time: %f seconds\n", d.header.Playback_time)
	fmt.Printf("Server Name: %s\n", d.header.Servername)
}

func (d *demofile) Open(path string) {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	d.header = demoheader{}
	// Read the header from the .dem
	buf := make([]byte, 4096)
	_, err = f.Read(buf)
	if err != nil {
		panic(err)
	}
	reader := bytes.NewReader(buf)
	err = binary.Read(reader, binary.LittleEndian, &d.header)
	if err != nil {
		panic(err)
	}
	if string(d.header.Demofilestamp[:7]) != DEMO_HEADER_ID {
		log.Fatal("Invalid demo header, are you sure this is a .dem?\n")
	}
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
}
