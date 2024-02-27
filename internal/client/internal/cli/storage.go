package cli

type storage interface {
	Check() error
	Open(string) error
	GetServerAddr() (string, error)
	GetLoginAndPass() (string, string, error)
	SetServerAddr(string) error
	SetLoginAndPass(string, string) error
}
