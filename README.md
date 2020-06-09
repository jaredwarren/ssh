# DB
SSH connection wrapper. One step ssh (tcp) connection.

The main purpose is to provide support for github.com/jaredwarren/db

TODO: add support for password, right now ssh key is required

## Usage
### connect
```go
c :=  &ssh.Conn {
        Host: "123.123.123.123",
        Port: 22,
        User: "root",
        Key: "~/.ssh/id_rsa",
    }
sshc := c.Connect()
``` 