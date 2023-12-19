/*
Copyright 2023 Flant JSC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"bufio"
	"errors"
	"fmt"
	"golang.org/x/sys/unix"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/coreos/go-iptables/iptables"
	"github.com/vishvananda/netlink"
)

type iptablesRule struct {
	Table string
	Chain string
	Rule  []string
	Pos   int
}

func main() {
	network, ok := os.LookupEnv("TUNNEL_NETWORK")
	if !ok {
		log.Fatal("The TUNNEL_NETWORK environment variable does not exist.")
	}

	protocol, ok := os.LookupEnv("OPENVPN_PROTO")
	if !ok {
		log.Fatal("The TUNNEL_NETWORK environment variable does not exist.")
	}

	var mgmtport string
	var routeTable int

	switch protocol {
	case "tcp":
		mgmtport = "8989"
		routeTable = 10
	case "udp":
		mgmtport = "9090"
		routeTable = 11
	default:
		log.Fatalf("OPENVPN_PROTO env value must be tcp or udp: %s", protocol)
	}

	iptablesMgr, err := iptables.NewWithProtocol(iptables.ProtocolIPv4)
	if err != nil {
		log.Fatal(err)
	}

	nat := iptablesRule{
		Table: "nat",
		Chain: "POSTROUTING",
		Rule:  strings.Fields(fmt.Sprintf("-s %s ! -d %s -j MASQUERADE", network, network)),
		Pos:   1,
	}

	manglePreroutingSetMark := iptablesRule{
		Table: "mangle",
		Chain: "PREROUTING",
		Rule:  strings.Fields(fmt.Sprintf("-i tun-%s -j CONNMARK --set-mark %d", protocol, routeTable)),
		Pos:   1,
	}

	manglePreroutingRestoreMark := iptablesRule{
		Table: "mangle",
		Chain: "PREROUTING",
		Rule:  strings.Fields("! -i tun+ -j CONNMARK --restore-mark"),
		Pos:   2,
	}
	mangleOutputRestoreMark := iptablesRule{
		Table: "mangle",
		Chain: "OUTPUT",
		Rule:  strings.Fields("-j CONNMARK --restore-mark"),
		Pos:   1,
	}

	err = insertUnique(iptablesMgr, nat.Table, nat.Chain, nat.Rule, nat.Pos)
	if err != nil {
		log.Fatal(err)
	}

	err = insertUnique(iptablesMgr, manglePreroutingSetMark.Table, manglePreroutingSetMark.Chain, manglePreroutingSetMark.Rule, manglePreroutingSetMark.Pos)
	if err != nil {
		log.Fatal(err)
	}

	err = insertUnique(iptablesMgr, manglePreroutingRestoreMark.Table, manglePreroutingRestoreMark.Chain, manglePreroutingRestoreMark.Rule, manglePreroutingRestoreMark.Pos)
	if err != nil {
		log.Fatal(err)
	}

	err = insertUnique(iptablesMgr, mangleOutputRestoreMark.Table, mangleOutputRestoreMark.Chain, mangleOutputRestoreMark.Rule, mangleOutputRestoreMark.Pos)
	if err != nil {
		log.Fatal(err)
	}

	rule := netlink.NewRule()
	rule.Table = routeTable
	rule.Priority = 1
	rule.Mark = routeTable
	err = netlink.RuleAdd(rule)

	if err != nil {
		if !os.IsExist(err) {
			log.Fatal(err)
		}
	}

	err = mknodDevNetTun()
	if err != nil {
		log.Fatal(err)
	}

	err = netLinkCreateTuntap(fmt.Sprintf("tun-%s", protocol), 1400)
	if err != nil {
		log.Fatal(err)
	}

	routeAdd(network, fmt.Sprintf("tun-%s", protocol), routeTable)

	requiredFiles := []string{
		"/etc/openvpn/certs/pki/ca.crt",
		"/etc/openvpn/certs/pki/private/server.key",
		"/etc/openvpn/certs/pki/issued/server.crt",
		"/etc/openvpn/certs/pki/ta.key",
		"/etc/openvpn/certs/pki/dh.pem",
		"/etc/openvpn/certs/pki/crl.pem",
	}
	for _, path := range requiredFiles {
		waitingForFile(path)
	}

	var args []string
	args = append(args, "--config")
	args = append(args, "/etc/openvpn/openvpn.conf")
	args = append(args, "--proto")
	args = append(args, protocol)
	args = append(args, "--management")
	args = append(args, "127.0.0.1")
	args = append(args, mgmtport)
	args = append(args, "--dev")
	args = append(args, fmt.Sprintf("tun-%s", protocol))
	log.Println(args)

	cmd := exec.Command("/usr/sbin/openvpn", args...)
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			log.Println(scanner.Text())
		}
	}()
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}
}

func insertUnique(iptablesMgr *iptables.IPTables, table, chain string, rule []string, pos int) error {
	ok, err := iptablesMgr.Exists(table, chain, rule...)
	if err != nil {
		return err
	}
	if !ok {
		err := iptablesMgr.Insert(table, chain, pos, rule...)
		if err != nil {
			return err
		}
	}

	return nil
}

func mknodDevNetTun() error {
	_, err := os.Stat("/dev")
	if os.IsNotExist(err) {
		err := unix.Mount("dev", "/dev", "devtmpfs", unix.MS_NOSUID|unix.MS_NOEXEC|unix.MS_RELATIME, "size=10m,nr_inodes=248418,mode=755")
		if err != nil {
			return fmt.Errorf("error mounting %s to %s: %v", "dev", "/dev", err)
		}
	}

	err = unix.Mkdir("/dev/net", 0753)
	if err != nil {
		if !os.IsExist(err) {
			return fmt.Errorf("error create dir /dev/net: %v", err)
		}
	}

	command := fmt.Sprintf("/bin/mknod")
	var args []string
	args = append(args, "/dev/net/tun")
	args = append(args, "c")
	args = append(args, "10")
	args = append(args, "200")
	cmd := exec.Command(command, args...)
	_, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error create /dev/net/tun: %v", err)
	}

	return nil
}

func netLinkCreateTuntap(name string, mtu int) error {
	linkAttrs := netlink.NewLinkAttrs()
	linkAttrs.Name = name
	linkAttrs.MTU = mtu

	l := &netlink.Tuntap{
		LinkAttrs: linkAttrs,
		Mode:      unix.IFF_TUN,
	}
	err := netlink.LinkAdd(l)
	if err != nil {
		return err
	}

	link, _ := netlink.LinkByName(linkAttrs.Name)
	err = netlink.LinkSetUp(link)
	if err != nil {
		return err
	}
	return nil
}

func routeAdd(dstNet string, linkName string, table int) {
	dstIPNet, err := parseIPNet(dstNet)
	if err != nil {
		log.Fatal("error parse IPNet: ", err)
	}

	link, _ := netlink.LinkByName(linkName)
	if err != nil {
		log.Fatal("error parse IPNet: ", err)
	}
	route := netlink.Route{Dst: dstIPNet, Table: table, LinkIndex: link.Attrs().Index}
	err = netlink.RouteAdd(&route)
	if err != nil {
		log.Fatalf("add route: %s error: %s\n", route.String(), err.Error())
	}
}

func parseIPNet(address string) (*net.IPNet, error) {
	ip := net.ParseIP(address)
	if ip != nil {
		return &net.IPNet{IP: ip, Mask: net.CIDRMask(32, 8*net.IPv4len)}, nil
	}
	_, IPNet, err := net.ParseCIDR(address)
	if err != nil {
		return nil, err
	}
	return IPNet, nil
}

func waitingForFile(path string) {
	for {
		_, err := os.Stat(path)
		if errors.Is(err, os.ErrNotExist) {
			log.Println("waiting for", path)
			time.Sleep(2 * time.Second)
		} else {
			break
		}
	}
}
