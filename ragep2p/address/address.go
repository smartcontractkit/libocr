package address

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/libocr/ragep2p/types"
	"net"
	"strconv"
)

func AnnounceAddrs(a types.Address) ([]types.Address, error) {
	host, port, err := splitAddress(a)
	if err != nil {
		return nil, errors.Wrap(err, "invalid address")
	}
	ip := net.ParseIP(host)
	if ip != nil && ip.IsUnspecified() {
		allIPv4s, allIPv6s, err := getAddressesForAllInterfaces(port)
		if err != nil {
			return nil, err
		} else if ip.To4() != nil {
			return allIPv4s, nil
		} else {
			return allIPv6s, nil
		}
	}
	// host is fully specified already
	return []types.Address{a}, nil
}

func IsValid(a types.Address) bool {
	_, _, err := splitAddress(a)
	return err == nil
}

// splitAddress splits an address into host and port. A third error argument is also returned.
func splitAddress(a types.Address) (string, uint16, error) {
	host, portString, err := net.SplitHostPort(string(a))
	if err != nil {
		return "", 0, err
	}
	ip := net.ParseIP(host)
	if ip == nil {
		return "", 0, fmt.Errorf("host [%s] is not an IP", host)
	}
	port, err := strconv.Atoi(portString)
	return host, uint16(port), err
}

func getAddressesForAllInterfaces(port uint16) (ipv4s []types.Address, ipv6s []types.Address, err error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return
	}

	for _, addr := range addrs {
		// By default addr contains the subnet, so decompose.
		ipNet, ok := addr.(*net.IPNet)
		if !ok {
			// Since we are working with interface addresses here, we should never reach this.
			// However, if for some reason we do, it is safe to say that the address would not be a useful
			// address to advertise, so it's ok to ignore.
			continue
		}
		ip := ipNet.IP
		// would be nice to use a switch here but IPv4 and IPv6 share a type in the standard library
		addr := types.Address(net.JoinHostPort(ip.String(), fmt.Sprintf("%d", port)))
		if ip.To4() != nil {
			ipv4s = append(ipv4s, addr)
		} else {
			ipv6s = append(ipv6s, addr)
		}
	}
	return
}
