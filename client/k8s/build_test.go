package k8s

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/previousnext/compose"
	"github.com/previousnext/skpr/k8stest/hash"
	"github.com/stretchr/testify/assert"

	"github.com/previousnext/m8s/client/types"
	"github.com/previousnext/m8s/config"
)

func TestBuild(t *testing.T) {
	client, err := New(types.ClientParams{
		Master: "http://localhost:8080",
	})
	assert.Nil(t, err)

	var buffer bytes.Buffer

	retention, err := time.ParseDuration("30m")
	assert.Nil(t, err)

	err = client.Build(&buffer, types.BuildParams{
		Name:   "test",
		Domain: "sdflkjsdf",
		Annotations: map[string]string{
			"nick": "rocks",
		},
		Repository: "https://github.com/example/repo.git",
		Revision:   "abc123",
		Config: config.Config{
			Namespace: hash.Generate(10),
			Retention: retention,
			Build: config.Build{
				Container: "",
				Steps: []string{
					"echo 1",
					"echo 2",
				},
			},
			Cache: config.Cache{
				Paths: []string{
					"/root/.composer",
				},
			},
			Secrets: config.Secrets{
				DockerCfg: "test-dockercfg",
				SSH:       "test-ssh",
			},
		},
		DockerCompose: compose.DockerCompose{
			Services: map[string]compose.Service{
				"app": {
					Image: "test/app:0.0.1",
					Volumes: []string{
						".:/data",
					},
					Entrypoint: []string{
						"echo 1",
						"echo 2",
					},
					Ports: []string{
						"80:80",
					},
					Environment: []string{
						"FOO=bar",
					},
					CapAdd: []string{
						"ALLOFIT",
					},
					Tmpfs: []string{
						"/tmp",
					},
					Deploy: compose.ServiceDeploy{
						Resources: compose.ServiceDeployResources{
							Limits: compose.ServiceDeployResource{
								CPUs:   "10m",
								Memory: "128m",
							},
							Reservations: compose.ServiceDeployResource{
								CPUs:   "10m",
								Memory: "128m",
							},
						},
					},
				},
			},
		},
	})
	assert.Nil(t, err)

	fmt.Println(buffer.String())

	assert.True(t, strings.Contains(buffer.String(), "Creating: PersistentVolumeClaim: /root/.composer"))
	assert.True(t, strings.Contains(buffer.String(), "Creating: Service"))
	assert.True(t, strings.Contains(buffer.String(), "Creating: Ingress"))
	assert.True(t, strings.Contains(buffer.String(), "Creating: Pod"))
}
