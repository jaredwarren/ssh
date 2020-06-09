package ssh

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

// Conn ...
type Conn struct {
	// ssh config
	Host string `json:"host,omitempty"`
	Port int    `json:"port,omitempty"`
	User string `json:"user,omitempty"`
	Pass string `json:"pass,omitempty"`
	Key  string `json:"key,omitempty"`
	// connections
	conn   net.Conn
	sshcon *ssh.Client
}

// Connect ...
func (c *Conn) Connect() error {
	// start ssh
	var agentClient agent.Agent
	// Establish a connection to the local ssh-agent
	conn, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK"))
	if err != nil {
		return err
	}
	c.conn = conn
	// Create a new instance of the ssh agent
	agentClient = agent.NewClient(conn)

	// TODO: figure out how to fix implement password or key

	key, err := ioutil.ReadFile(c.Key)
	if err != nil {
		return err
	}
	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return err
	}

	// The client configuration with configuration option to use the ssh-agent
	sshConfig := &ssh.ClientConfig{
		User: c.User,
		Auth: []ssh.AuthMethod{
			// Use the PublicKeys method for remote authentication.
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	// When the agentClient connection succeeded, add them as AuthMethod
	if agentClient != nil {
		sshConfig.Auth = append(sshConfig.Auth, ssh.PublicKeysCallback(agentClient.Signers))
	}

	// Connect to the SSH Server
	sshcon, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", c.Host, c.Port), sshConfig)
	if err != nil {
		return err
	}
	c.sshcon = sshcon

	return nil
}

// Close all connections
func (c *Conn) Close() {
	c.sshcon.Close()
	c.conn.Close()
}

// DialFunc ...
func (c *Conn) DialFunc(addr string) (net.Conn, error) {
	return c.sshcon.Dial("tcp", addr)
}
