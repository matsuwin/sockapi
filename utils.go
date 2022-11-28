package sockapi

import (
	"fmt"
	"strconv"
	"strings"
	"syscall"
)

const (
	ServerMode = iota
	ClientMode
)

func (s *SocketAddr) Close() {
	_ = syscall.Close(s.fd)
}

func (s *SocketAddr) Listen(mode int) (err error) {
	switch mode {
	default:
		err = fmt.Errorf("socket mode: Server|Client")
	case ServerMode:
		err = syscall.Bind(s.fd, s.Src)
		if err != nil {
			err = socketError(err, "bind")
		} else {
			err = syscall.Listen(s.fd, 1024) // backlog 1024
			if err != nil {
				err = socketError(err, "listen")
			}
		}
	case ClientMode:
		err = syscall.Connect(s.fd, s.Src)
	}
	return
}

func (s *SocketAddr) Accept() (*SocketAddr, syscall.Sockaddr, error) {
	nfd, sa, err := syscall.Accept(s.fd)
	if err != nil {
		err = socketError(err, "accept")
	}
	return &SocketAddr{fd: nfd}, sa, err
}

func (s *SocketAddr) Write(data []byte) (n int, err error) {
	n, err = syscall.Write(s.fd, data)
	if err != nil {
		err = socketError(err, "write")
	}
	return
}

func (s *SocketAddr) Read(length int) ([]byte, error) {
	buf := make([]byte, length)
	if _, err := syscall.Read(s.fd, buf); err != nil {
		return nil, socketError(err, "read")
	}
	return buf, nil
}

func socketError(err error, desc string) error {
	if err == nil {
		return fmt.Errorf("socket %s", desc)
	}
	return fmt.Errorf("socket %s %s", desc, err)
}

func ipv4AddrSplit(fp string) (ipv4 [4]byte, port int, err error) {
	addr := strings.Split(fp, ":")
	if len(addr) != 2 {
		err = socketError(err, "bad address")
		return
	}
	if addr[0] != "" {
		ipv := strings.Split(addr[0], ".")
		if len(ipv) != 4 {
			err = socketError(err, "bad address")
		}
		for i := 0; i < len(ipv); i++ {
			value, _ := strconv.Atoi(addr[i])
			ipv4[i] = byte(value)
		}
	}
	port, err = strconv.Atoi(addr[1])
	if err != nil {
		err = socketError(err, "bad address")
		return
	}
	return
}

func SocketAddrDo(proto, fp string, mode int) (sa *SocketAddr, err error) {
	switch proto {
	default:
		err = socketError(err, "bad proto")
		return
	case "unix":
		sa, err = sa.SourceUNIX(fp)
	case "tcp":
		sa, err = sa.SourceINET(fp)
	}
	if err != nil {
		return
	}
	err = sa.Listen(mode)
	if err != nil {
		err = socketError(err, fp)
	}
	return
}
