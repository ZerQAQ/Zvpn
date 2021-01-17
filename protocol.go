package Zvpn

type Conn interface {
	//write to the connect
	//l must equal len(buf), otherwise error is not nil
	Write(buf []byte) (l int, err error)
	//read from the connect and return how many byte is read
	//it's that l smaller than len(buf)
	Read(buf []byte) (l int, err error)
	Close() error
}

type Listener interface {
	//this is a blocking function
	Accept() (Conn, error)
	Close() error
}

type Protocol interface {
	Bind(string) (Listener, error)
	Dial(string) (Conn, error)
}

//var TCP Protocol = &TCPProtocol{}
