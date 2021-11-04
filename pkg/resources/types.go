/*
	SPDX-License-Identifier: MIT
*/

package resources

// resourceChecker is an interface which allows checking of a resource to see
// if it is in a ready state.
type resourceChecker interface {
	IsReady() (bool, error)
}
