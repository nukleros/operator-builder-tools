/*
	SPDX-License-Identifier: MIT
*/

package resources

import (
	"fmt"

	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	JobKind    = "Job"
	JobVersion = "batch/v1"
)

type JobResource struct {
	parent *batchv1.Job
}

// NewJobResource creates and returns a new JobResource.
func NewJobResource(name, namespace string) *JobResource {
	return &JobResource{
		parent: &batchv1.Job{
			TypeMeta: metav1.TypeMeta{
				Kind:       JobKind,
				APIVersion: JobVersion,
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
			},
		},
	}
}

// GetParent returns the parent attribute of the resource.
func (job *JobResource) GetParent() client.Object {
	return job.parent
}

// IsReady checks to see if a job is ready.
func (job *JobResource) IsReady(resource *Resource) (bool, error) {
	// if we have a name that is empty, we know we did not find the object
	if job.parent.Name == "" {
		return false, nil
	}

	// return immediately if the job is active or has no completion time
	if job.parent.Status.Active == 1 || job.parent.Status.CompletionTime == nil {
		return false, nil
	}

	// ensure the completion is actually successful
	if job.parent.Status.Succeeded != 1 {
		return false, fmt.Errorf("job %s was not successful", job.parent.GetName())
	}

	return true, nil
}
