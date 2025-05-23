package protocol

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"io"
	"net"
	"strings"
	"wired/modules/globals"
	"wired/modules/logger"
)

func SerializeString(str string) []byte {
	buf := bytes.Buffer{}
	l := globals.VarInt(len(str))
	l.WriteTo(&buf)
	buf.Write([]byte(str))

	return buf.Bytes()
}

func DeserializeString(conn io.Reader) (string, error) {
	var l globals.VarInt
	_, err := l.ReadFrom(conn)
	if err != nil {
		return "", err
	}

	buf := make([]byte, l)
	_, err = io.ReadFull(conn, buf)
	if err != nil {
		return "", err
	}

	return string(buf), nil
}

func RemovePort(addr string) string {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		logger.Fatal("Error removing port from address: ", err)
	}

	return host
}

func IsIPv4(address string) bool {
	ip := net.ParseIP(address)
	return ip != nil && ip.To4() != nil
}

func IsIPv6(address string) bool {
	ip := net.ParseIP(address)
	return ip != nil && ip.To4() == nil && strings.Contains(address, ":")
}

func RandomUUID() string {
	var uuid [16]byte
	_, err := io.ReadFull(rand.Reader, uuid[:])
	if err != nil {
		logger.Printf("Error generating UUID: %v\n", err)
		return ""
	}

	uuid[6] = (uuid[6] & 0x0f) | 0x40 // Version 4
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // Variant is 10

	return hex.EncodeToString(uuid[:])
}
