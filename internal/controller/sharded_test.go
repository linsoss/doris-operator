package controller

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	update_errors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

func TestAsResult(t *testing.T) {
	// Test case 1: No errors
	errSet := StCtrlErrSet{}
	result, err := errSet.AsResult()
	assert.NoError(t, err)
	assert.Equal(t, ctrl.Result{}, result)

	// Test case 2: Update conflict error
	statusErr := update_errors.StatusError{ErrStatus: metav1.Status{
		Reason: metav1.StatusReasonConflict,
		Code:   http.StatusConflict,
	}}
	errSet = StCtrlErrSet{Update: error(&statusErr)}
	result, err = errSet.AsResult()
	assert.NoError(t, err)
	assert.Equal(t, ctrl.Result{Requeue: true}, result)

	// Test case 3: Other Errors
	errSet = StCtrlErrSet{
		Rec:    errors.New("reconciliation error"),
		Sync:   errors.New("sync error"),
		Update: errors.New("sync error"),
	}
	result, err = errSet.AsResult()
	assert.Error(t, err)
	assert.Equal(t, ctrl.Result{Requeue: true}, result)
}
