package edp_component

import (
	"context"
	"github.com/epam/edp-reconciler/v2/pkg/controller/helper"
	"github.com/epam/edp-reconciler/v2/pkg/db"
	"github.com/epam/edp-reconciler/v2/pkg/model"
	ec "github.com/epam/edp-reconciler/v2/pkg/service/edp-component"
	"github.com/go-logr/logr"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"time"

	edpCompApi "github.com/epam/edp-component-operator/pkg/apis/v1/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func NewEDPComponent(client client.Client, log logr.Logger) *EDPComponent {
	return &EDPComponent{
		client: client,
		component: ec.EDPComponentService{
			DB: db.Instance,
		},
		log: log.WithName("edp-component"),
	}
}

type EDPComponent struct {
	client    client.Client
	component ec.EDPComponentService
	log       logr.Logger
}

func (r *EDPComponent) SetupWithManager(mgr ctrl.Manager) error {
	p := predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			old := e.ObjectOld.(*edpCompApi.EDPComponent).Spec
			new := e.ObjectNew.(*edpCompApi.EDPComponent).Spec

			if reflect.DeepEqual(old, new) {
				return false
			}
			return true
		},
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&edpCompApi.EDPComponent{}, builder.WithPredicates(p)).
		Complete(r)
}

func (r *EDPComponent) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	log := r.log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	log.Info("Reconciling EDPComponent CR")

	i := &edpCompApi.EDPComponent{}
	if err := r.client.Get(ctx, request.NamespacedName, i); err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	c, err := model.ConvertToEDPComponent(*i)
	if err != nil {
		return reconcile.Result{}, err
	}
	log.Info("start reconciling for component", "type", c.Type, "url", c.Url)
	edpN, err := helper.GetEDPName(r.client, i.Namespace)
	if err != nil {
		return reconcile.Result{}, err
	}
	err = r.component.PutEDPComponent(*c, *edpN)
	if err != nil {
		return reconcile.Result{RequeueAfter: time.Second * 120}, err
	}

	return reconcile.Result{}, nil
}
