package Zvpn

import (
	"errors"
	"github.com/ZerQAQ/Zvpn/lib"
	"io"
	"net"
	"strings"
)

func errIgnoreChecker(err error) bool {
	return err != nil && (err == io.EOF || strings.HasSuffix(err.Error(), "closed network connection"))
}

func BuildPipe(src Conn, des Conn) error {
	defer src.Close()
	defer des.Close()

	ch := make(chan error)
	go ConnCopy(src, des, ch)
	go ConnCopy(des, src, ch)
	err := <-ch
	return errors.New("in Pipe: " + err.Error())
}

func ConnCopy(src Conn, des Conn, errCh chan error) {
	buf := make([]byte, 1<<17)
	for {
		rLen, err := src.Read(buf)
		if errIgnoreChecker(err) || lib.ErrHandler("err when read", err) {
			errCh <- err
			return
		}
		_, err = des.Write(buf[:rLen])
		if errIgnoreChecker(err) || lib.ErrHandler("err when write", err) {
			errCh <- err
			return
		}
	}
}

type WallCrosserImply struct {
	protocol   Protocol
	proxy      Proxy
	localAddr  string
	serverAddr string
}

func (w WallCrosserImply) ClientConnHandler(conn Conn) {
	//connect the server
	remoteConn, err := w.protocol.Dial(w.serverAddr)
	if lib.ErrHandler("client", err) {
		return
	}

	//proxy protocol handshake
	remoteConn, err = w.proxy.ClientHandshake(conn, remoteConn, w.protocol)
	if lib.ErrHandler("client handshake", err) {
		return
	}

	//proxy conn establish
	err = BuildPipe(conn, remoteConn)
	lib.ErrHandler("in clientConnHandler", BuildPipe(conn, remoteConn))
}

func (w *WallCrosserImply) StartClient(loc, serv string) {
	w.localAddr = loc
	w.serverAddr = serv
	l, err := net.Listen("tcp4", w.localAddr)
	if err != nil {
		panic(err)
	}
	defer func() {
		err := l.Close()
		lib.ErrHandler("listener", err)
	}()
	println("client start on", w.localAddr)

	for {
		conn, err := l.Accept()
		if lib.ErrHandler("err when accept", err) {
			continue
		}
		go w.ClientConnHandler(conn)
	}
}

func (w WallCrosserImply) ServerConnHandler(conn Conn) {
	//build the conn the client want
	remoteConn, err := w.proxy.ServerHandshake(conn, w.protocol)
	if lib.ErrHandler("server proxy handshake", err) {
		return
	}

	//proxy conn establish
	lib.ErrHandler("in serverConnHandler", BuildPipe(conn, remoteConn))
}

func (w *WallCrosserImply) StartServer(addr string) {
	w.serverAddr = addr
	l, err := w.protocol.Bind(w.serverAddr)
	if err != nil {
		panic(err)
	}
	defer func() {
		err := l.Close()
		lib.ErrHandler("listener", err)
	}()
	println("server start on", w.serverAddr)

	for {
		conn, err := l.Accept()
		if lib.ErrHandler("err when accept", err) {
			continue
		}
		go w.ServerConnHandler(conn)
	}
}
