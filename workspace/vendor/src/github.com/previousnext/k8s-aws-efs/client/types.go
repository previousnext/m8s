package client

import (
	"encoding/json"
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type LifeCycleState string

const (
	LifeCycleStateReady    LifeCycleState = "Ready"
	LifeCycleStateNotReady LifeCycleState = "Not Ready"
	LifeCycleStateUnknown  LifeCycleState = "Unknown"
)

type PerformanceMode string

const (
	PerformanceModeGeneralPurpose PerformanceMode = "generalPurpose"
	PerformanceModeMaxIo          PerformanceMode = "maxIO"
)

type EfsSpec struct {
	Region        string          `json:"region"`
	Performance   PerformanceMode `json:"performance"`
	Subnets       []string        `json:"subnets"`
	SecurityGroup string          `json:"securityGroup"`
}

type EfsStatus struct {
	LastUpdate     time.Time      `json:"lastUpdate"`
	LifeCycleState LifeCycleState `json:"lifeCycleState"`
	ID             string         `json:"id"`
}

type Efs struct {
	metav1.TypeMeta `json:",inline"`
	Metadata        metav1.ObjectMeta `json:"metadata"`

	Spec   EfsSpec   `json:"spec"`
	Status EfsStatus `json:"status"`
}

// This creates a url based on Spec.Region and Status.ID.
// http://docs.aws.amazon.com/efs/latest/ug/mounting-fs.html
func (e *Efs) Endpoint() (string, error) {
	if e.Status.ID == "" {
		return "", fmt.Errorf("Cannot find filesystem ID")
	}

	if e.Spec.Region == "" {
		return "", fmt.Errorf("Cannot find filesystem region")
	}

	return fmt.Sprintf("%s.efs.%s.amazonaws.com", e.Status.ID, e.Spec.Region), nil
}

// Changing the format of this out WILL result in new filesystems being created.
func (e *Efs) CreationToken() string {
	return fmt.Sprintf("%s-%s", e.Metadata.Namespace, e.Metadata.Name)
}

// This is used for setting default values for an EFS.
// The main default we want to enforce is "General Purpose" for our mounts.
func (e *Efs) Defaults() {
	if e.Spec.Performance == "" {
		e.Spec.Performance = PerformanceModeGeneralPurpose
	}
}

// Ensures the user has provided the required values.
func (e *Efs) Validate() error {
	if e.Spec.Region == "" {
		return fmt.Errorf("Region was not provided")
	}
	return nil
}

type EfsList struct {
	metav1.TypeMeta `json:",inline"`
	Metadata        metav1.ListMeta `json:"metadata"`

	Items []Efs `json:"items"`
}

// Required to satisfy Object interface
func (e *Efs) GetObjectKind() schema.ObjectKind {
	return &e.TypeMeta
}

// Required to satisfy ObjectMetaAccessor interface
func (e *Efs) GetObjectMeta() metav1.Object {
	return &e.Metadata
}

// Required to satisfy Object interface
func (el *EfsList) GetObjectKind() schema.ObjectKind {
	return &el.TypeMeta
}

// Required to satisfy ListMetaAccessor interface
func (el *EfsList) GetListMeta() metav1.List {
	return &el.Metadata
}

// The code below is used only to work around a known problem with third-party
// resources and ugorji. If/when these issues are resolved, the code below
// should no longer be required.

type EfsListCopy EfsList
type EfsCopy Efs

func (e *Efs) UnmarshalJSON(data []byte) error {
	tmp := EfsCopy{}
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}
	tmp2 := Efs(tmp)
	*e = tmp2
	return nil
}

func (el *EfsList) UnmarshalJSON(data []byte) error {
	tmp := EfsListCopy{}
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}
	tmp2 := EfsList(tmp)
	*el = tmp2
	return nil
}
