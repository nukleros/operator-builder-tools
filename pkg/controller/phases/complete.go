// SPDX-License-Identifier: MIT

package phases

import "github.com/nukleros/operator-builder-tools/pkg/controller/workload"

// CompletePhase executes the completion of a reconciliation loop.
func CompletePhase(r workload.Reconciler, req *workload.Request) (bool, error) {
	req.Workload.SetReadyStatus(true)
	req.Log.Info("successfully reconciled")

	return true, nil
}
