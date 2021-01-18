package Zvpn

func NewProtocol(pro Protocol, o Obfuscate) Protocol {
	return ProtocolWrapper{pro, o}
}

type ProtocolWrapper struct {
	pro Protocol
	o   Obfuscate
}

func (p ProtocolWrapper) Bind(addr string) (listener Listener, err error) {
	listener, err = p.pro.Bind(addr)
	return ListerWrapper{listener, p.o}, err
}
func (p ProtocolWrapper) Dial(addr string) (conn Conn, err error) {
	conn, err = p.pro.Dial(addr)
	if err != nil {
		return nil, err
	}
	e, d, err := p.o.ClientHandShake(conn)
	if err != nil {
		return nil, err
	}
	return ConnWrapper{conn, e, d}, err
}

type ListerWrapper struct {
	listener Listener
	o        Obfuscate
}

func (l ListerWrapper) Accept() (Conn, error) {
	conn, err := l.listener.Accept()
	if err != nil {
		return nil, err
	}
	e, d, err := l.o.ServerHandShake(conn)
	if err != nil {
		return nil, err
	}
	return ConnWrapper{conn, e, d}, nil
}
func (l ListerWrapper) Close() error { return l.listener.Close() }

type ConnWrapper struct {
	conn Conn
	e    Encrypter
	d    Decrypter
}

func (c ConnWrapper) Read(buf []byte) (retLen int, retErr error) {
	return c.d.Read(c.conn, buf)
}
func (c ConnWrapper) Write(buf []byte) (retLen int, retErr error) {
	return c.e.Write(c.conn, buf)
}
func (c ConnWrapper) Close() error { return c.conn.Close() }
