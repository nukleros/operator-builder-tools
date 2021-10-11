/*
	SPDX-License-Identifier: MIT
*/

package resources

import (
	appsv1 "k8s.io/api/apps/v1"
)

const (
	DaemonSetKind = "DaemonSet"
)

// DaemonSetIsReady checks to see if a daemonset is ready.
func DaemonSetIsReady(
	resource *Resource,
) (bool, error) {
	var daemonSet appsv1.DaemonSet
	if err := GetObject(resource, &daemonSet, true); err != nil {
		return false, err
	}

	// ensure the desired number is scheduled and ready
	if daemonSet.Status.DesiredNumberScheduled == daemonSet.Status.NumberReady {
		if daemonSet.Status.NumberReady > 0 && daemonSet.Status.NumberUnavailable < 1 {
			return true, nil
		}
	}

	return false, nil
}
