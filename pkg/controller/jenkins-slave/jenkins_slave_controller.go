package jenkins_slave

import (
	"context"
	"github.com/epam/edp-reconciler/v2/pkg/controller/helper"
	"github.com/epam/edp-reconciler/v2/pkg/db"
	"github.com/epam/edp-reconciler/v2/pkg/service/jenkins-slave"
	"github.com/go-logr/logr"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sort"
	"time"

	jenkinsApi "github.com/epam/edp-jenkins-operator/v2/pkg/apis/v2/v1alpha1"
	errWrap "github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func NewReconcileJenkinsSlave(client client.Client, log logr.Logger) *ReconcileJenkinsSlave {
	return &ReconcileJenkinsSlave{
		client: client,
		jenkinsSlave: jenkins_slave.JenkinsSlaveService{
			DB: db.Instance,
		},
		log: log.WithName("jenkins-slave"),
	}
}

type ReconcileJenkinsSlave struct {
	client       client.Client
	jenkinsSlave jenkins_slave.JenkinsSlaveService
	log          logr.Logger
}

func (r *ReconcileJenkinsSlave) SetupWithManager(mgr ctrl.Manager) error {
	p := predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			old := e.ObjectOld.(*jenkinsApi.Jenkins).Status.Slaves
			new := e.ObjectNew.(*jenkinsApi.Jenkins).Status.Slaves

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

func (r *ReconcileJenkinsSlave) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	log := r.log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	log.Info("Reconciling Jenkins")

	jenkins := &jenkinsApi.Jenkins{}
	if err := r.client.Get(ctx, request.NamespacedName, jenkins); err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}
	log.WithValues("Jenkins", jenkins)

	edpN, err := helper.GetEDPName(r.client, jenkins.Namespace)
	if err != nil {
		return reconcile.Result{}, err
	}

	if err := r.jenkinsSlave.CreateSlavesOrDoNothing(jenkins.Status.Slaves, *edpN); err != nil {
		return reconcile.Result{RequeueAfter: time.Second * 120},
			errWrap.Wrapf(err, "an error has occurred while adding {%v} slaves into DB", jenkins.Status.Slaves)
	}

	return reconcile.Result{}, nil
}
