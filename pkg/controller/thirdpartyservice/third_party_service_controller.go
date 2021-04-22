package thirdpartyservice

import (
	"context"
	codebaseApi "github.com/epam/edp-codebase-operator/v2/pkg/apis/edp/v1alpha1"
	"github.com/epam/edp-reconciler/v2/pkg/controller/helper"
	"github.com/epam/edp-reconciler/v2/pkg/db"
	dtoService "github.com/epam/edp-reconciler/v2/pkg/model/service"
	tps "github.com/epam/edp-reconciler/v2/pkg/service/thirdpartyservice"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func NewReconcileService(client client.Client, log logr.Logger) *ReconcileService {
	return &ReconcileService{
		client: client,
		tps: tps.ThirdPartyService{
			DB: db.Instance,
		},
		log: log.WithName("third-party-service"),
	}
}

type ReconcileService struct {
	client client.Client
	tps    tps.ThirdPartyService
	log    logr.Logger
}

func (r *ReconcileService) SetupWithManager(mgr ctrl.Manager) error {
	p := predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			return false
		},
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&codebaseApi.Service{}, builder.WithPredicates(p)).
		Complete(r)
}

func (r *ReconcileService) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	log := r.log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	log.Info("Reconciling ThirdPartyService CR")

	instance := &codebaseApi.Service{}
	if err := r.client.Get(ctx, request.NamespacedName, instance); err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	edpName, err := helper.GetEDPName(r.client, instance.Namespace)
	if err != nil {
		return reconcile.Result{}, err
	}

	dto := dtoService.ConvertToServiceDto(*instance, *edpName)
	if err := r.tps.PutService(dto); err != nil {
		return reconcile.Result{}, err
	}
	log.Info("Reconciling ThirdPartyService CR has been finished")
	return reconcile.Result{}, nil
}
