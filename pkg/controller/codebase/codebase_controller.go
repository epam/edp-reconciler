package codebase

import (
	"context"
	"github.com/epam/edp-reconciler/v2/pkg/controller/helper"
	"github.com/epam/edp-reconciler/v2/pkg/db"
	"github.com/epam/edp-reconciler/v2/pkg/model/codebase"
	"github.com/epam/edp-reconciler/v2/pkg/service"
	"github.com/epam/edp-reconciler/v2/pkg/service/codebaseperfdatasource"
	"github.com/epam/edp-reconciler/v2/pkg/service/perfdatasource"
	"github.com/epam/edp-reconciler/v2/pkg/service/perfserver"
	"github.com/go-logr/logr"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"time"

	codebaseApi "github.com/epam/edp-codebase-operator/v2/pkg/apis/edp/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const codebaseReconcileFinalizerName = "codebase.reconciler.finalizer.name"

func NewReconcileCodebase(client client.Client, scheme *runtime.Scheme, log logr.Logger) *ReconcileCodebase {
	return &ReconcileCodebase{
		client: client,
		scheme: scheme,
		codebase: service.CodebaseService{
			DB: db.Instance,
			DataSourceService: perfdatasource.PerfDataSourceService{
				DB: db.Instance,
			},
			PerfService: perfserver.PerfServerService{
				DB: db.Instance,
			},
			CodebaseDsService: codebaseperfdatasource.CodebasePerfDataSourceService{
				DB: db.Instance,
			},
		},
		log: log.WithName("codebase"),
	}
}

type ReconcileCodebase struct {
	client   client.Client
	scheme   *runtime.Scheme
	codebase service.CodebaseService
	log      logr.Logger
}

func (r *ReconcileCodebase) SetupWithManager(mgr ctrl.Manager) error {
	p := predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			oldObject := e.ObjectOld.(*codebaseApi.Codebase)
			newObject := e.ObjectNew.(*codebaseApi.Codebase)

			if oldObject.Status.Value != newObject.Status.Value ||
				oldObject.Status.Action != newObject.Status.Action {
				return true
			}

			if !reflect.DeepEqual(oldObject.Spec, newObject.Spec) {
				return true
			}

			if newObject.DeletionTimestamp != nil {
				return true
			}
			return false
		},
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&codebaseApi.Codebase{}, builder.WithPredicates(p)).
		Complete(r)
}

func (r *ReconcileCodebase) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	log := r.log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	log.Info("Reconciling Codebase")

	i := &codebaseApi.Codebase{}
	if err := r.client.Get(ctx, request.NamespacedName, i); err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}
	log.Info("Codebase has been retrieved", "codebase", i)

	edpN, err := helper.GetEDPName(r.client, i.Namespace)
	if err != nil {
		log.Error(err, "cannot get edp name")
		return reconcile.Result{RequeueAfter: 2 * time.Second}, nil
	}

	result, err := r.tryToDeleteCodebase(ctx, i, *edpN)
	if err != nil || result != nil {
		return *result, err
	}

	c, err := codebase.Convert(*i, *edpN)
	if err != nil {
		log.Error(err, "cannot convert codebase to dto")
		return reconcile.Result{RequeueAfter: 2 * time.Second}, nil
	}

	if err = r.codebase.PutCodebase(*c); err != nil {
		log.Error(err, "cannot put codebase", "name", c.Name)
		return reconcile.Result{RequeueAfter: 2 * time.Second}, nil
	}

	log.Info("Reconciling has been finished successfully")
	return reconcile.Result{}, nil
}

func (r *ReconcileCodebase) tryToDeleteCodebase(ctx context.Context, i *codebaseApi.Codebase, schema string) (*reconcile.Result, error) {
	if i.GetDeletionTimestamp().IsZero() {
		if !helper.ContainsString(i.ObjectMeta.Finalizers, codebaseReconcileFinalizerName) {
			i.ObjectMeta.Finalizers = append(i.ObjectMeta.Finalizers, codebaseReconcileFinalizerName)
			if err := r.client.Update(ctx, i); err != nil {
				return &reconcile.Result{}, err
			}
		}
		return nil, nil
	}
	if err := r.codebase.Delete(i.Spec.Perf, i.Name, schema); err != nil {
		return &reconcile.Result{}, err
	}

	i.ObjectMeta.Finalizers = helper.RemoveString(i.ObjectMeta.Finalizers, codebaseReconcileFinalizerName)
	if err := r.client.Update(ctx, i); err != nil {
		return &reconcile.Result{}, err
	}
	return &reconcile.Result{}, nil
}
