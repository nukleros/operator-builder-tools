// SPDX-License-Identifier: MIT

package phases

import (
	"strings"
)

// optimisticLockErrorMsg is the error we see from the k8s api when optimistic locking occurs.
const optimisticLockErrorMsg = "the object has been modified; please apply your changes to the latest version and try again"

// IsOptimisticLockError checks to see if the error is a locking error.
func IsOptimisticLockError(err error) bool {
	return strings.Contains(err.Error(), optimisticLockErrorMsg)
}
