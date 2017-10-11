package env

import (
	"testing"

	"github.com/previousnext/m8s/cmd/metadata"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/pkg/api/v1"
)

func TestService(t *testing.T) {
	want := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "test",
			Name:      "pr1",
			Annotations: map[string]string{
				metadata.AnnotationBitbucketRepoOwner: "nick",
			},
		},
		Spec: v1.ServiceSpec{
			ClusterIP: "None",
			Ports: []v1.ServicePort{
				{
					Name: "http",
					Port: 80,
				},
				{
					Name: "mailhog",
					Port: 8025,
				},
				{
					Name: "solr",
					Port: 8983,
				},
			},
			Selector: map[string]string{
				"env": "pr1",
			},
		},
	}

	annotations, err := metadata.Annotations([]string{"BITBUCKET_REPO_OWNER=nick"})
	assert.Nil(t, err)

	assert.Equal(t, want, Service("test", "pr1", annotations))
}
