package m8s

import (
	"github.com/pkg/errors"
	"github.com/previousnext/m8s/server/k8s/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
)

const knownHosts = "fDF8dzJZU1RJSGNQN3pWSnJMZWJUOXZWZWJNekNjPXx1VmFQT1ZxMXorcisyem1vblhmejhzWkNxRGM9IHNzaC1yc2EgQUFBQUIzTnphQzF5YzJFQUFBQUJJd0FBQVFFQXEyQTdoUkdtZG5tOXRVRGJPOUlEU3dCSzZUYlFhK1BYWVBDUHk2cmJUclR0dzdQSGtjY0tycHAweVZocDVIZEVJY0tyNnBMbFZEQmZPTFg5UVVzeUNPVjB3emZqSUpObEdFWXNkbExKaXpIaGJuMm1VanZTQUhRcVpFVFlQODFlRnpMUU5uUEh0NEVWVlVoN1ZmREVTVTg0S2V6bUQ1UWxXcFhMbXZVMzEveU1mK1NlOHhoSFR2S1NDWklGSW1Xd29HNm1iVW9XZjluenBJb2FTakIrd2VxcVVVbXBhYWFzWFZhbDcySitVWDJCKzJSUFczUmNUMGVPelFncWxKTDNSS3JUSnZkc2pFM0pFQXZHcTNsR0hTWlh5MjhHM3NrdWEyU21WaS93NHlDRTZnYk9EcW5UV2xnNyt3QzYwNHlkR1hBOFZKaVM1YXA0M0pYaVVGRkFhUT09Cg=="

func installSecret(client *kubernetes.Clientset, namespace string) error {
	secret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "ssh",
			Namespace: namespace,
		},
		StringData: map[string]string{
			"known_hosts": knownHosts,
		},
	}

	_, err := utils.SecretCreate(client, secret)
	if err != nil {
		return errors.Wrap(err, "failed to install Secret")
	}

	return nil
}
