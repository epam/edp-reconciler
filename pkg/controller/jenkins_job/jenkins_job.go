package jenkins_job

import (
	"context"
	jenkinsApi "github.com/epam/edp-jenkins-operator/v2/pkg/apis/v2/v1alpha1"
	"github.com/epam/edp-reconciler/v2/pkg/controller/jenkins_job/service"
	"github.com/epam/edp-reconciler/v2/pkg/db"
	"github.com/go-logr/logr"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"time"
)

func NewReconcileJenkinsJob(client client.Client, scheme *runtime.Scheme, log logr.Logger) *ReconcileJenkinsJob {
	return &ReconcileJenkinsJob{
		client: client,
		scheme: scheme,
		jenkinsJob: service.JenkinsJobService{
			DB:     db.Instance,
			Client: client,
		},
		log: log.WithName("jenkins-job"),
	}
}

type ReconcileJenkinsJob struct {
	client     client.Client
	scheme     *runtime.Scheme
	jenkinsJob service.JenkinsJobService
	log        logr.Logger
}

func (r *ReconcileJenkinsJob) SetupWithManager(mgr ctrl.Manager) error {
	p := predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			oldObject := e.ObjectOld.(*jenkinsApi.JenkinsJob)
			newObject := e.ObjectNew.(*jenkinsApi.JenkinsJob)
			if oldObject.Status.Action != newObject.Status.Action ||
				oldObject.Status.Value != newObject.Status.Value {
				return true
			}
			return false
		},
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&jenkinsApi.JenkinsJob{}, builder.WithPredicates(p)).
		Complete(r)
}

func (r *ReconcileJenkinsJob) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	log := r.log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	log.V(2).Info("Reconciling JenkinsJob")
	i := &jenkinsApi.JenkinsJob{}
	if err := r.client.Get(ctx, request.NamespacedName, i); err != nil {
		if k8serrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	if err := r.jenkinsJob.UpdateActionLog(i); err != nil {
		return reconcile.Result{RequeueAfter: 5 * time.Second}, err
	}

	log.V(2).Info("Reconciling JenkinsJob has been finished successfully")
	return reconcile.Result{}, nil
}
