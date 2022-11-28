//go:build unix

package sockapi

import "syscall"

type SocketAddr struct {
	fd  int
	Src syscall.Sockaddr
}

func (s *SocketAddr) FD() int { return s.fd }

func (s *SocketAddr) SourceUNIX(fp string) (*SocketAddr, error) {
	fd, err := syscall.Socket(syscall.AF_UNIX, syscall.SOCK_STREAM, 0)
	if err != nil {
		return nil, socketError(err, "fd")
	}
	return &SocketAddr{
		fd: fd,
		Src: &syscall.SockaddrUnix{
			Name: fp,
		},
	}, nil
}

func (s *SocketAddr) SourceINET(fp string) (_ *SocketAddr, err error) {
	var ipv4 [4]byte
	var port int
	var fd int
	ipv4, port, err = ipv4AddrSplit(fp)
	if err != nil {
		return
	}
	fd, err = syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		err = socketError(err, "fd")
		return
	}
	return &SocketAddr{
		fd: fd,
		Src: &syscall.SockaddrInet4{
			Port: port,
			Addr: ipv4,
		},
	}, nil
}
