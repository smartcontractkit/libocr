package ragedisco

import (
	"fmt"
	"github.com/pkg/errors"
	ragetypes "github.com/smartcontractkit/libocr/ragep2p/types"
	"net"
	"strconv"
)

func announceAddrs(a ragetypes.Address) ([]ragetypes.Address, error) {
	host, port, err := splitAddress(a)
	if err != nil {
		return nil, errors.Wrap(err, "invalid address")
	}
	ip := net.ParseIP(host)
	if ip == nil {
		return nil, fmt.Errorf("host was not an IP")
	}
	if ip.IsUnspecified() {
		allIPv4s, allIPv6s, err := getAddressesForAllInterfaces(port)
		if err != nil {
			return nil, err
		} else if ip.To4() != nil {
			return allIPv4s, nil
		} else {
			return allIPv6s, nil
		}
	}
	return []ragetypes.Address{a}, nil
}

func combinedAnnounceAddrs(as []string) (combined []ragetypes.Address, err error) {
	dedup := make(map[ragetypes.Address]struct{})
	for _, addr := range as {
		announceAddresses, err := announceAddrs(ragetypes.Address(addr))
		if err != nil {
			return nil, err
		}
		for _, annAddr := range announceAddresses {
			dedup[annAddr] = struct{}{}
		}
	}
	for addr := range dedup {
		combined = append(combined, addr)
	}
	return
}

// isValidForAnnouncement checks that the provided address is in the form ip:port.
// Hostnames or domain names are not allowed.
func isValidForAnnouncement(a ragetypes.Address) bool {
	host, _, err := splitAddress(a)
	if err != nil {
		return false
	}
	ip := net.ParseIP(host)
	return ip != nil
}

// splitAddress splits an address into host and port. A third error argument is also returned.
func splitAddress(a ragetypes.Address) (string, uint16, error) {
	host, portString, err := net.SplitHostPort(string(a))
	if err != nil {
		return "", 0, err
	}
	port, err := strconv.Atoi(portString)
	return host, uint16(port), err
}

func getAddressesForAllInterfaces(port uint16) (ipv4s []ragetypes.Address, ipv6s []ragetypes.Address, err error) {
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
		addr := ragetypes.Address(net.JoinHostPort(ip.String(), fmt.Sprintf("%d", port)))
		if ip.To4() != nil {
			ipv4s = append(ipv4s, addr)
		} else {
			ipv6s = append(ipv6s, addr)
		}
	}
	return
}
