package perfserver

import (
	"context"

	perfApi "github.com/epam/edp-perf-operator/v2/pkg/apis/edp/v1"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/epam/edp-reconciler/v2/pkg/controller/helper"
	"github.com/epam/edp-reconciler/v2/pkg/db"
	perfServerModel "github.com/epam/edp-reconciler/v2/pkg/model/perfserver"
	"github.com/epam/edp-reconciler/v2/pkg/service/perfserver"
)

func NewReconcilePerfServer(client client.Client, log logr.Logger) *ReconcilePerfServer {
	return &ReconcilePerfServer{
		client: client,
		perfService: perfserver.PerfServerService{
			DB: db.Instance,
		},
		log: log.WithName("perf-server"),
	}
}

type ReconcilePerfServer struct {
	client      client.Client
	perfService perfserver.PerfServerService
	log         logr.Logger
}

func (r *ReconcilePerfServer) SetupWithManager(mgr ctrl.Manager) error {
	p := predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			oldObject := e.ObjectOld.(*perfApi.PerfServer)
			newObject := e.ObjectNew.(*perfApi.PerfServer)
			if oldObject.Spec != newObject.Spec {
				return true
			}
			if oldObject.Status.Available != newObject.Status.Available {
				return true
			}
			return false
		},
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&perfApi.PerfServer{}, builder.WithPredicates(p)).
		Complete(r)
}

func (r *ReconcilePerfServer) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	log := r.log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	log.Info("Reconciling PerfServer")

	i := &perfApi.PerfServer{}
	if err := r.client.Get(ctx, request.NamespacedName, i); err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	schema, err := helper.GetEDPName(r.client, i.Namespace)
	if err != nil {
		return reconcile.Result{}, err
	}

	if err := r.perfService.PutPerfServer(perfServerModel.ConvertPerfServerToDto(*i), *schema); err != nil {
		return reconcile.Result{}, err
	}

	log.Info("PerfServer reconciling has been finished successfully")
	return reconcile.Result{}, nil
}
