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
	Object batchv1.Job
}

// NewJobResource creates and returns a new JobResource.
func NewJobResource(object metav1.Object) (*JobResource, error) {
	job := &batchv1.Job{}

	err := ToProper(job, object)
	if err != nil {
		return nil, err
	}

	return &JobResource{Object: *job}, nil
}

// IsReady checks to see if a job is ready.
func (job *JobResource) IsReady() (bool, error) {
	// if we have a name that is empty, we know we did not find the object
	if job.Object.Name == "" {
		return false, nil
	}

	// return immediately if the job is active or has no completion time
	if job.Object.Status.Active == 1 || job.Object.Status.CompletionTime == nil {
		return false, nil
	}

	// ensure the completion is actually successful
	if job.Object.Status.Succeeded != 1 {
		return false, fmt.Errorf("job %s was not successful", job.Object.Name)
	}

	return true, nil
}
