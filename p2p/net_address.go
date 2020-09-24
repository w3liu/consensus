package p2p

import (
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
)

type NetAddress struct {
	ID   ID     `json:"id"`
	IP   net.IP `json:"ip"`
	Port uint16 `json:"port"`
}

func IDAddressString(id ID, protocolHostPort string) string {
	hostPort := removeProtocolIfDefined(protocolHostPort)
	return fmt.Sprintf("%s@%s", id, hostPort)
}

func removeProtocolIfDefined(addr string) string {
	if strings.Contains(addr, "://") {
		return strings.Split(addr, "://")[1]
	}
	return addr

}

func NewNetAddressString(addr string) (*NetAddress, error) {
	addrWithoutProtocol := removeProtocolIfDefined(addr)
	spl := strings.Split(addrWithoutProtocol, "@")
	if len(spl) != 2 {
		return nil, ErrNetAddressNoID{addr}
	}

	// get ID
	if err := validateID(ID(spl[0])); err != nil {
		return nil, ErrNetAddressInvalid{addrWithoutProtocol, err}
	}
	var id ID
	id, addrWithoutProtocol = ID(spl[0]), spl[1]

	// get host and port
	host, portStr, err := net.SplitHostPort(addrWithoutProtocol)
	if err != nil {
		return nil, ErrNetAddressInvalid{addrWithoutProtocol, err}
	}
	if len(host) == 0 {
		return nil, ErrNetAddressInvalid{
			addrWithoutProtocol,
			errors.New("host is empty")}
	}

	ip := net.ParseIP(host)
	if ip == nil {
		ips, err := net.LookupIP(host)
		if err != nil {
			return nil, ErrNetAddressLookup{host, err}
		}
		ip = ips[0]
	}

	port, err := strconv.ParseUint(portStr, 10, 16)
	if err != nil {
		return nil, ErrNetAddressInvalid{portStr, err}
	}

	na := NewNetAddressIPPort(ip, uint16(port))
	na.ID = id
	return na, nil
}

func NewNetAddressIPPort(ip net.IP, port uint16) *NetAddress {
	return &NetAddress{
		IP:   ip,
		Port: port,
	}
}

func validateID(id ID) error {
	if len(id) == 0 {
		return errors.New("no ID")
	}
	idBytes, err := hex.DecodeString(string(id))
	if err != nil {
		return err
	}
	if len(idBytes) != IDByteLength {
		return fmt.Errorf("invalid hex length - got %d, expected %d", len(idBytes), IDByteLength)
	}
	return nil
}
