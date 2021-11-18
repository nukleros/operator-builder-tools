// SPDX-License-Identifier: MIT

package status

import "time"

// PhaseState defines the current state of the phase.
// +kubebuilder:validation:Enum=Complete;Reconciling;Failed;Pending
type PhaseState string

const (
	PhaseStatePending     PhaseState = "Pending"
	PhaseStateReconciling PhaseState = "Reconciling"
	PhaseStateFailed      PhaseState = "Failed"
	PhaseStateComplete    PhaseState = "Complete"
)

// PhaseCondition describes an event that has occurred during a phase
// of the controller reconciliation loop.
type PhaseCondition struct {
	State PhaseState `json:"state"`

	// Phase defines the phase in which the condition was set.
	Phase string `json:"phase"`

	// Message defines a helpful message from the phase.
	Message string `json:"message"`

	// LastModified defines the time in which this component was updated.
	LastModified string `json:"lastModified"`
}

// GetSuccessCondition defines the success condition for the phase.
func GetSuccessCondition(name string) PhaseCondition {
	return PhaseCondition{
		Phase:        name,
		LastModified: time.Now().UTC().String(),
		State:        PhaseStateComplete,
		Message:      "Successfully Completed Phase",
	}
}

// GetPendingCondition defines the pending condition for the phase.
func GetPendingCondition(name string) PhaseCondition {
	return PhaseCondition{
		Phase:        name,
		LastModified: time.Now().UTC().String(),
		State:        PhaseStatePending,
		Message:      "Pending Execution of Phase",
	}
}

// GetFailCondition defines the fail condition for the phase.
func GetFailCondition(name string, err error) PhaseCondition {
	return PhaseCondition{
		Phase:        name,
		LastModified: time.Now().UTC().String(),
		State:        PhaseStateFailed,
		Message:      "Failed Phase with Error; " + err.Error(),
	}
}
