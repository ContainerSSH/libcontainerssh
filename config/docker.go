package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
	"gopkg.in/yaml.v3"
)

// DockerConfig is the base configuration structure of the Docker backend.
//
// swagger:model DockerConfig
type DockerConfig struct {
	// Connection configures how to connect to dockerd
	Connection DockerConnectionConfig `json:"connection" yaml:"connection"`
	// Execution drives how the container and the workload is executed
	Execution DockerExecutionConfig `json:"execution" yaml:"execution"`
	// Timeouts configures the various timeouts when interacting with dockerd.
	Timeouts DockerTimeoutConfig `json:"timeouts" yaml:"timeouts"`
}

// Validate validates the provided configuration and returns an error if invalid.
func (c DockerConfig) Validate() error {
	if err := c.Connection.Validate(); err != nil {
		return wrap(err, "connection")
	}
	if err := c.Execution.Validate(); err != nil {
		return wrap(err, "execution")
	}
	return nil
}

func (c DockerConnectionConfig) Validate() error {
	if c.Host == "" {
		return fmt.Errorf("missing host")
	}
	return nil
}

func parseRawDuration(rawValue interface{}, d *time.Duration) error {
	var err error
	switch value := rawValue.(type) {
	case nil:
		*d = time.Duration(0)
	case int32:
		*d = time.Duration(value)
	case int64:
		*d = time.Duration(value)
	case int:
		*d = time.Duration(value)
	case float32:
		*d = time.Duration(value)
	case float64:
		*d = time.Duration(value)
	case string:
		if *d, err = time.ParseDuration(value); err != nil {
			return err
		}
	default:
		return fmt.Errorf("invalid duration: %v", rawValue)
	}
	return nil
}

// DockerExecutionMode determines when a container is launched.
// DockerExecutionModeConnection launches one container per SSH connection (default), while DockerExecutionModeSession launches
// one container per SSH session.
//
// swagger:enum DockerExecutionMode
type DockerExecutionMode string

const (
	// DockerExecutionModeConnection launches one container per SSH connection.
	DockerExecutionModeConnection DockerExecutionMode = "connection"
	// DockerExecutionModeSession launches one container per SSH session (multiple containers per connection).
	DockerExecutionModeSession DockerExecutionMode = "session"
)

// Validate validates the execution config.
func (e DockerExecutionMode) Validate() error {
	switch e {
	case DockerExecutionModeConnection:
		fallthrough
	case DockerExecutionModeSession:
		return nil
	default:
		return fmt.Errorf("invalid execution mode: %s", e)
	}
}

// DockerExecutionConfig contains the configuration of what container to run in Docker.
//
// swagger:model DockerExecutionConfig
type DockerExecutionConfig struct {
	// DockerLaunchConfig contains the Docker-specific launch configuration.
	DockerLaunchConfig `json:",inline,omitempty" yaml:",inline,omitempty"`

	// Mode influences how commands are executed.
	//
	// - If "connection" is chosen (default) a new container is launched per connection. In this mode
	//   sessions are executed using the "docker exec" functionality and the main container console runs a script that
	//   waits for a termination signal.
	// - If "session" is chosen a new container is launched per session, leading to potentially multiple
	//   containers per connection. In this mode the program is launched directly as the main process of the container.
	//   When configuring this mode you should explicitly configure the "cmd" option to an empty list if you want the
	//   default command in the container to launch.
	//
	// default: connection
	Mode DockerExecutionMode `json:"mode,omitempty" yaml:"mode,omitempty" default:"connection"`

	// IdleCommand is the command that runs as the first process in the container in DockerExecutionModeConnection. Ignored in DockerExecutionModeSession.
	//
	// default: ["/usr/bin/containerssh-agent", "wait-signal", "--signal", "INT", "--signal", "TERM"]
	IdleCommand []string `json:"idleCommand,omitempty" yaml:"idleCommand,omitempty" comment:"Run this command to wait for container exit" default:"[\"/usr/bin/containerssh-agent\", \"wait-signal\", \"--signal\", \"INT\", \"--signal\", \"TERM\"]"`
	// ShellCommand is the command used for launching shells when the container is in DockerExecutionModeConnection. Ignored in DockerExecutionModeSession.
	//
	// default: ["/bin/bash"]
	ShellCommand []string `json:"shellCommand,omitempty" yaml:"shellCommand,omitempty" comment:"Run this command as a default shell." default:"[\"/bin/bash\"]"`
	// AgentPath contains the path to the ContainerSSH Guest Agent.
	//
	// default: /usr/bin/containerssh-agent
	AgentPath string `json:"agentPath,omitempty" yaml:"agentPath,omitempty" default:"/usr/bin/containerssh-agent"`
	// DisableAgent enables using the ContainerSSH Guest Agent.
	//
	// default: false
	DisableAgent bool `json:"disableAgent,omitempty" yaml:"disableAgent,omitempty"`
	// Subsystems contains a map of subsystem names and their corresponding binaries in the container.
	//
	// default: {"sftp":"/usr/lib/openssh/sftp-server"}
	Subsystems map[string]string `json:"subsystems,omitempty" yaml:"subsystems,omitempty" comment:"Subsystem names and binaries map." default:"{\"sftp\":\"/usr/lib/openssh/sftp-server\"}"`

	// ImagePullPolicy controls when to pull container images.
	//
	// default: IfNotPresent
	ImagePullPolicy DockerImagePullPolicy `json:"imagePullPolicy,omitempty" yaml:"imagePullPolicy,omitempty" comment:"Image pull policy" default:"IfNotPresent"`

	// ExposeAuthMetadataAsEnv lets you expose the authentication metadata (e.g. GITHUB_TOKEN) as an environment variable
	// in the container. In contrast to the environment variables set in the SSH connection these environment variables
	// are available to all processes in the container, including the idle command.
	//
	// default: false
	ExposeAuthMetadataAsEnv bool `json:"exposeAuthMetadataAsEnv,omitempty" yaml:"exposeAuthMetadataAsEnv,omitempty"`
}

type tmpDockerExecutionConfig struct {
	Container interface{} `json:"container" yaml:"container"`
	Host interface{} `json:"host" yaml:"host"`
	Network interface{} `json:"network" yaml:"network"`
	Platform interface{} `json:"platform" yaml:"platform"`
	ContainerName interface{} `json:"containername" yaml:"containername"`
	Mode DockerExecutionMode `json:"mode,omitempty" yaml:"mode,omitempty" default:"connection"`
	IdleCommand []string `json:"idleCommand,omitempty" yaml:"idleCommand,omitempty" comment:"Run this command to wait for container exit" default:"[\"/usr/bin/containerssh-agent\", \"wait-signal\", \"--signal\", \"INT\", \"--signal\", \"TERM\"]"`
	ShellCommand []string `json:"shellCommand,omitempty" yaml:"shellCommand,omitempty" comment:"Run this command as a default shell." default:"[\"/bin/bash\"]"`
	AgentPath string `json:"agentPath,omitempty" yaml:"agentPath,omitempty" default:"/usr/bin/containerssh-agent"`
	DisableAgent bool `json:"disableAgent,omitempty" yaml:"disableAgent,omitempty"`
	Subsystems map[string]string `json:"subsystems,omitempty" yaml:"subsystems,omitempty" comment:"Subsystem names and binaries map." default:"{\"sftp\":\"/usr/lib/openssh/sftp-server\"}"`
	ImagePullPolicy DockerImagePullPolicy `json:"imagePullPolicy,omitempty" yaml:"imagePullPolicy,omitempty" comment:"Image pull policy" default:"IfNotPresent"`
	ExposeAuthMetadataAsEnv bool `json:"exposeAuthMetadataAsEnv,omitempty" yaml:"exposeAuthMetadataAsEnv,omitempty"`
}


// UnmarshalJSON takes a JSON byte array and unmarshalls it into a structure.
func (c *DockerExecutionConfig) UnmarshalJSON(b []byte) error {
	decoder := json.NewDecoder(bytes.NewReader(b))
	tmp := &tmpDockerExecutionConfig{}
	if err := decoder.Decode(tmp); err != nil {
		return err
	}
	launch := DockerLaunchConfig{}
	decoder = json.NewDecoder(bytes.NewReader(b))
	if err := decoder.Decode(&launch); err != nil {
		return err
	}
	c.fill(launch, tmp)
	return nil
}

// UnmarshalYAML takes a YAML byte array and unmarshalls it into a structure.
func (c *DockerExecutionConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	tmp := &tmpDockerExecutionConfig{}
	if err := unmarshal(tmp); err != nil {
		return err
	}
	launch := DockerLaunchConfig{}
	if err := unmarshal(&launch); err != nil {
		return err
	}
	c.fill(launch, tmp)
	return nil
}

func (c *DockerExecutionConfig) fill(launch DockerLaunchConfig, tmp *tmpDockerExecutionConfig) {
	c.DockerLaunchConfig = launch
	c.Mode = tmp.Mode
	c.IdleCommand = tmp.IdleCommand
	c.ShellCommand = tmp.ShellCommand
	c.AgentPath = tmp.AgentPath
	c.DisableAgent = tmp.DisableAgent
	c.Subsystems = tmp.Subsystems
	c.ImagePullPolicy = tmp.ImagePullPolicy
	c.ExposeAuthMetadataAsEnv = tmp.ExposeAuthMetadataAsEnv
}

// Validate validates the docker config structure.
func (c DockerExecutionConfig) Validate() error {
	if c.Mode == DockerExecutionModeConnection && len(c.IdleCommand) == 0 {
		return newError("idleCommand", "idle command required for execution mode \"connection\"")
	}
	if c.Mode == DockerExecutionModeConnection && len(c.ShellCommand) == 0 {
		return newError("shellCommand", "shell command required for execution mode \"connection\"")
	}
	switch c.Mode {
	case DockerExecutionModeSession:
		if c.DockerLaunchConfig.HostConfig != nil && !c.DockerLaunchConfig.HostConfig.RestartPolicy.IsNone() {
			return wrap(
				newError(
					"restartPolicy",
					"unsupported restart policy for execution mode \"session\": %s (session containers may not restart)",
					c.DockerLaunchConfig.HostConfig.RestartPolicy.Name,
				),
				"hostConfig",
			)
		}
	}
	if err := c.ImagePullPolicy.Validate(); err != nil {
		return wrap(err, "imagePullPolicy")
	}
	if err := c.DockerLaunchConfig.Validate(); err != nil {
		return err
	}
	if err := c.Mode.Validate(); err != nil {
		return wrap(err, "mode")
	}
	return nil
}

// DockerImagePullPolicy drives how and when images are pulled. The values are closely aligned with the Kubernetes image pull
// policy.
//
// - ImagePullPolicyAlways means that the container image will be pulled on every connection.
// - ImagePullPolicyIfNotPresent means the image will be pulled if the image is not present locally, an empty tag, or
//	 the "latest" tag was specified.
// - ImagePullPolicyNever means that the image will be never pulled, and if the image is not available locally the
//	 connection will fail.
//
// swagger:enum DockerImagePullPolicy
type DockerImagePullPolicy string

const (
	// ImagePullPolicyAlways means that the container image will be pulled on every connection.
	ImagePullPolicyAlways DockerImagePullPolicy = "Always"
	// ImagePullPolicyIfNotPresent means the image will be pulled if the image is not present locally, an empty tag, or
	// the "latest" tag was specified.
	ImagePullPolicyIfNotPresent DockerImagePullPolicy = "IfNotPresent"
	// ImagePullPolicyNever means that the image will be never pulled, and if the image is not available locally the
	// connection will fail.
	ImagePullPolicyNever DockerImagePullPolicy = "Never"
)

// Validate checks if the given image pull policy is valid.
func (p DockerImagePullPolicy) Validate() error {
	switch p {
	case ImagePullPolicyAlways:
		fallthrough
	case ImagePullPolicyIfNotPresent:
		fallthrough
	case ImagePullPolicyNever:
		return nil
	default:
		return fmt.Errorf("invalid image pull policy: %s", p)
	}
}

// DockerTimeoutConfig drives the various timeouts in the Docker backend.
//
// swagger:model DockerTimeoutConfig
type DockerTimeoutConfig struct {
	// ContainerStart is the maximum time starting a container may take. It may be configured as an integer in
	// nanoseconds or as a time formatting string.
	//
	// required: false
	// example: 60s
	// swagger:type string
	ContainerStart time.Duration `json:"containerStart" yaml:"containerStart" default:"60s"`
	// ContainerStop is the maximum time to wait for a container to stop.
	// This should always be set higher than the Docker StopTimeout. It may be configured as an integer in
	// nanoseconds or as a time formatting string.
	//
	// required: true
	// example: 60s
	// swagger:type string
	ContainerStop time.Duration `json:"containerStop" yaml:"containerStop" default:"60s"`
	// CommandStart sets the maximum time starting a command may take. It may be configured as an integer in
	// nanoseconds or as a time formatting string.
	//
	// required: true
	// example: 60s
	// swagger:type string
	CommandStart time.Duration `json:"commandStart" yaml:"commandStart" default:"60s"`
	// Signal sets the maximum time sending a signal may take. It may be configured as an integer in
	// nanoseconds or as a time formatting string.
	//
	// required: true
	// example: 60s
	// swagger:type string
	Signal time.Duration `json:"signal" yaml:"signal" default:"60s"`
	// Signal sets the maximum time setting the window size may take. It may be configured as an integer in
	// nanoseconds or as a time formatting string.
	//
	// required: true
	// example: 60s
	// swagger:type string
	Window time.Duration `json:"window" yaml:"window" default:"60s"`
	// HTTP is the timeout for the HTTP calls themselves.
	//
	// required: true
	// example: 60s
	// swagger:type string
	HTTP time.Duration `json:"http" yaml:"http" default:"15s"`
}

type dockerTmpTimeoutConfig struct {
	// ContainerStart is the maximum time starting a container may take.
	ContainerStart interface{} `json:"containerStart" yaml:"containerStart" default:"60s"`
	// ContainerStop is the maximum time to wait for a container to stop. This should always be set higher than the Docker StopTimeout.
	ContainerStop interface{} `json:"containerStop" yaml:"containerStop" default:"60s"`
	// CommandStart sets the maximum time starting a command may take.
	CommandStart interface{} `json:"commandStart" yaml:"commandStart" default:"60s"`
	// Signal sets the maximum time sending a signal may take.
	Signal interface{} `json:"signal" yaml:"signal" default:"60s"`
	// Signal sets the maximum time setting the window size may take.
	Window interface{} `json:"window" yaml:"window" default:"60s"`
	// HTTP
	HTTP interface{} `json:"http" yaml:"http" default:"15s"`
}

// UnmarshalJSON takes a JSON byte array and unmarshalls it into a structure.
func (t *DockerTimeoutConfig) UnmarshalJSON(b []byte) error {
	decoder := json.NewDecoder(bytes.NewReader(b))
	tmp := &dockerTmpTimeoutConfig{}
	if err := decoder.Decode(tmp); err != nil {
		return err
	}

	return t.unmarshalTmp(tmp)
}

// UnmarshalYAML takes a YAML byte array and unmarshalls it into a structure.
func (t *DockerTimeoutConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	tmp := &dockerTmpTimeoutConfig{}
	if err := unmarshal(tmp); err != nil {
		return err
	}

	return t.unmarshalTmp(tmp)
}

func (t *DockerTimeoutConfig) unmarshalTmp(tmp *dockerTmpTimeoutConfig) error {
	if err := parseRawDuration(tmp.ContainerStart, &t.ContainerStart); err != nil {
		return err
	}
	if err := parseRawDuration(tmp.ContainerStop, &t.ContainerStop); err != nil {
		return err
	}
	if err := parseRawDuration(tmp.CommandStart, &t.CommandStart); err != nil {
		return err
	}
	if err := parseRawDuration(tmp.Signal, &t.Signal); err != nil {
		return err
	}
	if err := parseRawDuration(tmp.Window, &t.Window); err != nil {
		return err
	}
	if err := parseRawDuration(tmp.HTTP, &t.HTTP); err != nil {
		return err
	}
	return nil
}

// DockerLaunchConfig contains the container configuration for the Docker client version 20.
type DockerLaunchConfig struct {
	// ContainerConfig contains container-specific configuration options.
	//
	// default: {"image":"containerssh/containerssh-guest-image"}
	ContainerConfig *container.Config `json:"container" yaml:"container" comment:"DockerConfig configuration." default:"{\"image\":\"containerssh/containerssh-guest-image\"}"`
	// HostConfig contains the host-specific configuration options.
	HostConfig *container.HostConfig `json:"host" yaml:"host" comment:"Host configuration"`
	// NetworkConfig contains the network settings.
	NetworkConfig *network.NetworkingConfig `json:"network" yaml:"network" comment:"Network configuration"`
	// Platform contains the platform specification.
	Platform *specs.Platform `json:"platform" yaml:"platform" comment:"Platform specification"`
	// ContainerName is the name of the container to launch. It is recommended to leave this empty, otherwise
	// ContainerSSH may not be able to start the container if a container with the same name already exists.
	ContainerName string `json:"containername" yaml:"containername" comment:"Name for the container to be launched"`
}

type dockerTmpLaunchConfig struct {
	// ContainerConfig contains container-specific configuration options.
	ContainerConfig *container.Config `json:"container" yaml:"container"`
	// HostConfig contains the host-specific configuration options.
	HostConfig *container.HostConfig `json:"host" yaml:"host"`
	// NetworkConfig contains the network settings.
	NetworkConfig *network.NetworkingConfig `json:"network" yaml:"network"`
	// Platform contains the platform specification.
	Platform *specs.Platform `json:"platform" yaml:"platform"`
	// ContainerName is the name of the container to launch. It is recommended to leave this empty, otherwise
	// ContainerSSH may not be able to start the container if a container with the same name already exists.
	ContainerName string `json:"containername" yaml:"containername"`
}

// UnmarshalJSON implements the special unmarshalling of the DockerLaunchConfig that ignores unknown fields.
// This is needed because Docker treats removing fields as backwards-compatible.
// See https://github.com/moby/moby/pull/39158#issuecomment-489704731
func (l *DockerLaunchConfig) UnmarshalJSON(b []byte) error {
	decoder := json.NewDecoder(bytes.NewReader(b))
	tmp := &dockerTmpLaunchConfig{}
	if err := decoder.Decode(tmp); err != nil {
		return err
	}
	l.ContainerConfig = tmp.ContainerConfig
	l.HostConfig = tmp.HostConfig
	l.NetworkConfig = tmp.NetworkConfig
	l.Platform = tmp.Platform
	l.ContainerName = tmp.ContainerName
	return nil
}

// UnmarshalYAML implements the special unmarshalling of the DockerLaunchConfig that ignores unknown fields.
// This is needed because Docker treats removing fields as backwards-compatible.
// See https://github.com/moby/moby/pull/39158#issuecomment-489704731
func (l *DockerLaunchConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	lc := &map[string]interface{}{}
	if err := unmarshal(lc); err != nil {
		return err
	}
	substructure, err := yaml.Marshal(lc)
	if err != nil {
		return err
	}
	tmp := &dockerTmpLaunchConfig{}
	if err = yaml.Unmarshal(substructure, tmp); err != nil {
		return err
	}
	l.ContainerConfig = tmp.ContainerConfig
	l.HostConfig = tmp.HostConfig
	l.NetworkConfig = tmp.NetworkConfig
	l.Platform = tmp.Platform
	l.ContainerName = tmp.ContainerName
	return nil
}

// Validate validates the launch configuration.
func (l *DockerLaunchConfig) Validate() error {
	if l.ContainerConfig == nil {
		return newError("container", "no container config provided")
	}
	if l.ContainerConfig.Image == "" {
		return wrap(newError("image", "no image name provided"), "container")
	}
	return nil
}
