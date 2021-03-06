/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/pointer"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var clusterlog = logf.Log.WithName("cluster-resource")

func (r *Cluster) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
//

// +kubebuilder:webhook:path=/validate-scylla-scylladb-com-v1alpha1-cluster,failurePolicy=fail,groups=scylla.scylladb.com,resources=clusters,verbs=create;update,versions=v1alpha1,name=webhook.scylla.scylladb.com,mutating=false
var _ webhook.Validator = &Cluster{}

// +kubebuilder:webhook:path=/mutate-scylla-scylladb-com-v1alpha1-cluster,failurePolicy=fail,groups=scylla.scylladb.com,resources=clusters,verbs=create;update,versions=v1alpha1,name=webhook.scylla.scylladb.com,mutating=true
var _ webhook.Defaulter = &Cluster{}

func (c *Cluster) Default() {
	for i, repairTask := range c.Spec.Repairs {
		if repairTask.StartDate == nil {
			c.Spec.Repairs[i].StartDate = pointer.StringPtr("now")
		}
		if repairTask.Interval == nil {
			c.Spec.Repairs[i].Interval = pointer.StringPtr("0")
		}
		if repairTask.NumRetries == nil {
			c.Spec.Repairs[i].NumRetries = pointer.Int64Ptr(3)
		}
		if repairTask.SmallTableThreshold == nil {
			c.Spec.Repairs[i].SmallTableThreshold = pointer.StringPtr("1GiB")
		}
	}

	for i, backupTask := range c.Spec.Backups {
		if backupTask.StartDate == nil {
			c.Spec.Backups[i].StartDate = pointer.StringPtr("now")
		}
		if backupTask.Interval == nil {
			c.Spec.Backups[i].Interval = pointer.StringPtr("0")
		}
		if backupTask.NumRetries == nil {
			c.Spec.Backups[i].NumRetries = pointer.Int64Ptr(3)
		}
		if backupTask.Retention == nil {
			c.Spec.Backups[i].Retention = pointer.Int64Ptr(3)
		}
	}
}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *Cluster) ValidateCreate() error {
	clusterlog.Info("validate create", "name", r.Name)

	// First, check the values
	if err := checkValues(r); err != nil {
		return err
	}

	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *Cluster) ValidateUpdate(old runtime.Object) error {
	clusterlog.Info("validate update", "name", r.Name)

	// First, check the values
	if err := checkValues(r); err != nil {
		return err
	}

	if err := checkTransitions(old.(*Cluster), r); err != nil {
		return err
	}

	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *Cluster) ValidateDelete() error {
	clusterlog.Info("validate delete", "name", r.Name)
	// no validation during delete
	return nil
}
