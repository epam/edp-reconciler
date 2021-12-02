package job_provisioning

import (
	"context"
	"github.com/epam/edp-reconciler/v2/pkg/db"
	"github.com/go-logr/logr"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sort"
	"time"

	"github.com/epam/edp-reconciler/v2/pkg/controller/helper"
	jp "github.com/epam/edp-reconciler/v2/pkg/service/job-provisioning"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	jenkinsApi "github.com/epam/edp-jenkins-operator/v2/pkg/apis/v2/v1alpha1"
	errWrap "github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func NewReconcileJobProvision(client client.Client, log logr.Logger) *ReconcileJobProvision {
	return &ReconcileJobProvision{
		client: client,
		jobProvision: jp.JobProvisionService{
			DB: db.Instance,
		},
		log: log.WithName("job-provision"),
	}
}

type ReconcileJobProvision struct {
	client       client.Client
	jobProvision jp.JobProvisionService
	log          logr.Logger
}

func (r *ReconcileJobProvision) SetupWithManager(mgr ctrl.Manager) error {
	p := predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			old := e.ObjectOld.(*jenkinsApi.Jenkins).Status.JobProvisions
			new := e.ObjectNew.(*jenkinsApi.Jenkins).Status.JobProvisions

			sort.Slice(old, func(i, j int) bool {
				return old[i].Name < old[j].Name
			})
			sort.Slice(new, func(i, j int) bool {
				return new[i].Name < new[j].Name
			})

			return !reflect.DeepEqual(old, new)
		},
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&jenkinsApi.Jenkins{}, builder.WithPredicates(p)).
		Complete(r)
}

func (r *ReconcileJobProvision) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	log := r.log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	log.Info("Reconciling Jenkins CR to handle Job Provisios")

	instance := &jenkinsApi.Jenkins{}
	if err := r.client.Get(ctx, request.NamespacedName, instance); err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	jp := instance.Status.JobProvisions
	edpN, err := helper.GetEDPName(r.client, instance.Namespace)
	if err != nil {
		return reconcile.Result{}, err
	}
	err = r.jobProvision.PutJobProvisions(jp, *edpN)
	if err != nil {
		return reconcile.Result{RequeueAfter: time.Second * 120},
			errWrap.Wrapf(err, "an error has occurred while adding {%v} job provisions into DB", jp)
	}

	return reconcile.Result{}, nil
}
