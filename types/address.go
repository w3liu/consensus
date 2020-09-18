package types

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
)

type Address struct {
	Id   string `json:"id"`
	Ip   net.IP `json:"ip"`
	Port uint16 `json:"port"`
}

func NewAddress(addr string) (*Address, error) {
	arr := strings.Split(addr, "@")
	if len(arr) != 2 {
		return nil, errors.New("invalid address")
	}
	id, hostport := arr[0], arr[1]
	host, portstr, err := net.SplitHostPort(hostport)
	if err != nil {
		return nil, err
	}
	if len(host) == 0 {
		return nil, errors.New("host is empty")
	}
	ip := net.ParseIP(host)
	if ip == nil {
		ips, err := net.LookupIP(host)
		if err != nil {
			return nil, err
		}
		ip = ips[0]
	}
	port, err := strconv.ParseUint(portstr, 10, 16)
	if err != nil {
		return nil, err
	}

	address := &Address{
		Id:   id,
		Ip:   ip,
		Port: uint16(port),
	}

	return address, nil
}

func (a *Address) ToIpPortString() string {
	return fmt.Sprintf("%s:%d", a.Ip, a.Port)
}
