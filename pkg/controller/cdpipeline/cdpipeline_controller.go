package cdpipeline

import (
	"context"
	cdPipeApi "github.com/epam/edp-cd-pipeline-operator/v2/pkg/apis/edp/v1alpha1"
	"github.com/epam/edp-reconciler/v2/pkg/controller/helper"
	"github.com/epam/edp-reconciler/v2/pkg/db"
	"github.com/epam/edp-reconciler/v2/pkg/model/cdpipeline"
	"github.com/epam/edp-reconciler/v2/pkg/platform"
	"github.com/epam/edp-reconciler/v2/pkg/service/cd-pipeline"
	tps "github.com/epam/edp-reconciler/v2/pkg/service/thirdpartyservice"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"time"

	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const cdPipelineReconcileFinalizerName = "cdpipeline.reconciler.finalizer.name"

func NewReconcileCDPipeline(client client.Client, scheme *runtime.Scheme, log logr.Logger) (*ReconcileCDPipeline, error) {
	cs, err := platform.CreateOpenshiftClients()
	if err != nil {
		return nil, errors.Wrap(err, "unable to create openshift clients")
	}

	return &ReconcileCDPipeline{
		client: client,
		scheme: scheme,
		pipe: cd_pipeline.CdPipelineService{
			DB:        db.Instance,
			ClientSet: *cs,
			ThirdPartyService: tps.ThirdPartyService{
				DB: db.Instance,
			},
		},
		log: log.WithName("cd-pipeline"),
	}, nil
}

type ReconcileCDPipeline struct {
	client client.Client
	scheme *runtime.Scheme
	pipe   cd_pipeline.CdPipelineService
	log    logr.Logger
}

func (r *ReconcileCDPipeline) SetupWithManager(mgr ctrl.Manager) error {
	p := predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			oldObject := e.ObjectOld.(*cdPipeApi.CDPipeline)
			newObject := e.ObjectNew.(*cdPipeApi.CDPipeline)

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
		For(&cdPipeApi.CDPipeline{}, builder.WithPredicates(p)).
		Complete(r)
}

func (r *ReconcileCDPipeline) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	log := r.log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	log.Info("Reconciling CDPipeline")

	instance := &cdPipeApi.CDPipeline{}
	if err := r.client.Get(ctx, request.NamespacedName, instance); err != nil {
		if k8sErrors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	log.Info("CD pipeline has been retrieved", "cd pipeline", instance)

	edpN, err := helper.GetEDPName(r.client, instance.Namespace)
	if err != nil {
		log.Error(err, "cannot get edp name")
		return reconcile.Result{RequeueAfter: 2 * time.Second}, nil
	}

	if res, err := r.tryToDeleteCDPipeline(ctx, instance, *edpN); err != nil || res != nil {
		return *res, err
	}

	cdp, err := cdpipeline.ConvertToCDPipeline(*instance, *edpN)
	if err != nil {
		log.Error(err, "cannot convert to cd pipeline dto")
		return reconcile.Result{RequeueAfter: 2 * time.Second}, nil
	}
	err = r.pipe.PutCDPipeline(*cdp)
	if err != nil {
		log.Error(err, "cannot put cd pipeline")
		return reconcile.Result{RequeueAfter: 2 * time.Second}, nil
	}

	log.Info("Reconciling has been finished successfully")
	return reconcile.Result{}, nil
}

func (r *ReconcileCDPipeline) tryToDeleteCDPipeline(ctx context.Context, p *cdPipeApi.CDPipeline, schema string) (*reconcile.Result, error) {
	if p.GetDeletionTimestamp().IsZero() {
		if !helper.ContainsString(p.ObjectMeta.Finalizers, cdPipelineReconcileFinalizerName) {
			p.ObjectMeta.Finalizers = append(p.ObjectMeta.Finalizers, cdPipelineReconcileFinalizerName)
			if err := r.client.Update(ctx, p); err != nil {
				return &reconcile.Result{}, err
			}
		}
		return nil, nil
	}

	if err := r.pipe.DeleteCDPipeline(p.Name, schema); err != nil {
		return &reconcile.Result{RequeueAfter: 2 * time.Second}, err
	}

	p.ObjectMeta.Finalizers = helper.RemoveString(p.ObjectMeta.Finalizers, cdPipelineReconcileFinalizerName)
	if err := r.client.Update(ctx, p); err != nil {
		return &reconcile.Result{RequeueAfter: 2 * time.Second}, err
	}
	return &reconcile.Result{}, nil
}
