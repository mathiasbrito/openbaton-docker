package proto

const (
	// DefaultAddress is the default address for the client to connect to.
	DefaultAddress = "localhost:60000"

	// DefaultListenProtocol is the default protocol the server uses for Listen.
	DefaultListenProtocol = "tcp"

	// DefaultListenAddress is the default address the server uses for Listen.
	DefaultListenAddress = ":60000"

	// LoginMethod is the signature of the login method. Check this string 
	// carefully.
	LoginMethod = "/vim_pop.Pop/Login"
)
