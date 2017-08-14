package env

const (
	// SecretSSH is the identifier for storing the "id_rsa" and known_hosts file
	// secrets are stored.
	SecretSSH = "ssh"
	// SecretDockerCfg is the identifier for storing Docker configuration.
	SecretDockerCfg = "dockercfg"
)
