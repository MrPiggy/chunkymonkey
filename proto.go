package main

import (
	"net"
	"os"
	"fmt"
)

const (
	protocolVersion = 2

	// Packet type IDs
	packetIDLogin = 0x1
	packetIDHandshake = 0x2
	packetIDSpawnPosition = 0x6
)

func ReadByte(conn net.Conn) (b byte, err os.Error) {
	bs := make([]byte, 1)
	_, err = conn.Read(bs)
	return bs[0], err
}

func WriteByte(conn net.Conn, b byte) (err os.Error) {
	_, err = conn.Write([]byte{b})
	return
}

func ReadShort(conn net.Conn) (i int, err os.Error) {
	bs := make([]byte, 2)
	_, err = conn.Read(bs)
	return int(uint16(bs[0]) << 8 | uint16(bs[1])), err
}

func WriteShort(conn net.Conn, i int) (err os.Error) {
	_, err = conn.Write([]byte{byte(i >> 8), byte(i)})
	return
}

func ReadInt32(conn net.Conn) (i int32, err os.Error) {
	bs := make([]byte, 4)
	_, err = conn.Read(bs)
	return int32(uint32(bs[0]) << 24 | uint32(bs[1]) << 16 | uint32(bs[2]) << 8 | uint32(bs[3])), err
}

func WriteInt32(conn net.Conn, i int32) (err os.Error) {
	_, err = conn.Write([]byte{byte(i >> 24), byte(i >> 16), byte(i >> 8), byte(i)})
	return
}

func ReadString(conn net.Conn) (s string, err os.Error) {
	n, e := ReadShort(conn)
	if e != nil {
		return "", e
	}

	bs := make([]byte, n)
	_, err = conn.Read(bs)
	return string(bs), err
}

func WriteString(conn net.Conn, s string) (err os.Error) {
	bs := []byte(s)

	err = WriteShort(conn, len(bs))
	if err != nil {
		return err
	}

	_, err = conn.Write(bs)
	return err
}

func ReadHandshake(conn net.Conn) (username string, err os.Error) {
	packetID, e := ReadByte(conn)
	if e != nil {
		return "", e
	}
	if packetID != packetIDHandshake {
		panic(fmt.Sprintf("ReadHandshake: invalid packet ID %#x", packetID))
	}

	return ReadString(conn)
}

func WriteHandshake(conn net.Conn, reply string) (err os.Error) {
	_, err = conn.Write([]byte{packetIDHandshake})
	if err != nil {
		return err
	}

	return WriteString(conn, reply)
}

func ReadLogin(conn net.Conn) (username, password string, err os.Error) {
	packetID, e := ReadByte(conn)
	if e != nil {
		return "", "", e
	}
	if packetID != packetIDLogin {
		panic(fmt.Sprintf("ReadLogin: invalid packet ID %#x", packetID))
	}

	version, e2 := ReadInt32(conn)
	if e2 != nil {
		return "", "", e2
	}
	if version != protocolVersion {
		panic(fmt.Sprintf("ReadLogin: unsupported protocol version %#x", version))
	}

	username, e3 := ReadString(conn)
	if e3 != nil {
		return "", "", e3
	}

	password, e4 := ReadString(conn)
	if e4 != nil {
		return "", "", e4
	}

	return username, password, nil
}

func WriteLogin(conn net.Conn) (err os.Error) {
	_, err = conn.Write([]byte{packetIDLogin, 0, 0, 0, 0, 0, 0, 0, 0})
	return err
}

func WriteSpawnPosition(conn net.Conn, position *XYZ) (err os.Error) {
	err = WriteByte(conn, packetIDSpawnPosition)
	if err != nil {
		return
	}

	err = WriteInt32(conn, position.x)
	if err != nil {
		return
	}

	err = WriteInt32(conn, position.y)
	if err != nil {
		return
	}

	err = WriteInt32(conn, position.z)
	return
}