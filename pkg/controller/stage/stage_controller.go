package stage

import (
	"context"
	cdPipeApi "github.com/epam/edp-cd-pipeline-operator/v2/pkg/apis/edp/v1alpha1"
	"github.com/epam/edp-reconciler/v2/pkg/controller/helper"
	"github.com/epam/edp-reconciler/v2/pkg/db"
	"github.com/epam/edp-reconciler/v2/pkg/model/stage"
	"github.com/epam/edp-reconciler/v2/pkg/platform"
	stageService "github.com/epam/edp-reconciler/v2/pkg/service/stage"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"time"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const stageReconcileFinalizerName = "stage.reconciler.finalizer.name"

func NewReconcileStage(client client.Client, scheme *runtime.Scheme, log logr.Logger) (*ReconcileStage, error) {
	cs, err := platform.CreateOpenshiftClients()
	if err != nil {
		return nil, errors.Wrap(err, "unable to create openshift clients")
	}

	return &ReconcileStage{
		client: client,
		scheme: scheme,
		service: stageService.StageService{
			DB:        db.Instance,
			ClientSet: *cs,
		},
		log: log.WithName("cd-stage"),
	}, nil
}

type ReconcileStage struct {
	client  client.Client
	scheme  *runtime.Scheme
	service stageService.StageService
	log     logr.Logger
}

func (r *ReconcileStage) SetupWithManager(mgr ctrl.Manager) error {
	p := predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			oldObject := e.ObjectOld.(*cdPipeApi.Stage)
			newObject := e.ObjectNew.(*cdPipeApi.Stage)

			if oldObject.Status.Value != newObject.Status.Value {
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
		For(&cdPipeApi.Stage{}, builder.WithPredicates(p)).
		Complete(r)
}

func (r *ReconcileStage) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	log := r.log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	log.V(2).Info("Reconciling Stage")
	i := &cdPipeApi.Stage{}
	if err := r.client.Get(ctx, request.NamespacedName, i); err != nil {
		if k8serrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	edpN, err := helper.GetEDPName(r.client, i.Namespace)
	if err != nil {
		return reconcile.Result{RequeueAfter: 2 * time.Second}, errors.Wrap(err, "cannot get edp name")
	}

	if res, err := r.tryToDeleteCDStage(ctx, i, *edpN); err != nil || res != nil {
		return *res, err
	}

	st, err := stage.ConvertToStage(*i, *edpN)
	if err != nil {
		return reconcile.Result{RequeueAfter: 2 * time.Second}, errors.Wrap(err, "couldn't convert to stage dto")
	}

	if err = r.service.PutStage(*st); err != nil {
		return reconcile.Result{RequeueAfter: 2 * time.Second}, errors.Wrap(err, "couldn't put stage")
	}
	log.V(2).Info("Reconciling has been finished successfully")
	return reconcile.Result{}, nil
}

func (r ReconcileStage) tryToDeleteCDStage(ctx context.Context, i *cdPipeApi.Stage, schema string) (*reconcile.Result, error) {
	if i.GetDeletionTimestamp().IsZero() {
		if !helper.ContainsString(i.ObjectMeta.Finalizers, stageReconcileFinalizerName) {
			i.ObjectMeta.Finalizers = append(i.ObjectMeta.Finalizers, stageReconcileFinalizerName)
			if err := r.client.Update(ctx, i); err != nil {
				return &reconcile.Result{}, err
			}
		}
		return nil, nil
	}

	if err := r.service.DeleteCDStage(i.Spec.CdPipeline, i.Spec.Name, schema); err != nil {
		return &reconcile.Result{RequeueAfter: 2 * time.Second}, err
	}

	i.ObjectMeta.Finalizers = helper.RemoveString(i.ObjectMeta.Finalizers, stageReconcileFinalizerName)
	if err := r.client.Update(ctx, i); err != nil {
		return &reconcile.Result{RequeueAfter: 2 * time.Second}, err
	}
	return &reconcile.Result{}, nil
}
