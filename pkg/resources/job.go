/*
	SPDX-License-Identifier: MIT
*/

package resources

import (
	"fmt"

	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	JobKind    = "Job"
	JobVersion = "batch/v1"
)

type JobResource struct {
	batchv1.Job
}

// NewJobResource creates and returns a new JobResource.
func NewJobResource(name, namespace string) *JobResource {
	return &JobResource{
		batchv1.Job{
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

// IsReady checks to see if a job is ready.
func (job *JobResource) IsReady() (bool, error) {
	// if we have a name that is empty, we know we did not find the object
	if job.Name == "" {
		return false, nil
	}

	// return immediately if the job is active or has no completion time
	if job.Status.Active == 1 || job.Status.CompletionTime == nil {
		return false, nil
	}

	// ensure the completion is actually successful
	if job.Status.Succeeded != 1 {
		return false, fmt.Errorf("job %s was not successful", job.GetName())
	}

	return true, nil
}
