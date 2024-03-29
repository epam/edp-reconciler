package perfdatasourcesonar

import (
	"context"
	"time"

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
	"github.com/epam/edp-reconciler/v2/pkg/service/perfdatasource"
	"github.com/epam/edp-reconciler/v2/pkg/util/cluster"
)

const (
	codebaseKind                          = "Codebase"
	sonarDataSourceReconcileFinalizerName = "sonar.data.source.reconciler.finalizer.name"
)

func NewReconcilePerfDataSourceSonar(client client.Client, log logr.Logger) *ReconcilePerfDataSourceSonar {
	return &ReconcilePerfDataSourceSonar{
		client: client,
		dsService: perfdatasource.PerfDataSourceService{
			DB: db.Instance,
		},
		log: log.WithName("perf-data-source-sonar"),
	}
}

type ReconcilePerfDataSourceSonar struct {
	client    client.Client
	dsService perfdatasource.PerfDataSourceService
	log       logr.Logger
}

func (r *ReconcilePerfDataSourceSonar) SetupWithManager(mgr ctrl.Manager) error {
	p := predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			return e.ObjectNew.(*perfApi.PerfDataSourceSonar).DeletionTimestamp != nil
		},
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&perfApi.PerfDataSourceSonar{}, builder.WithPredicates(p)).
		Complete(r)
}

func (r *ReconcilePerfDataSourceSonar) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	log := r.log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	log.Info("Reconciling PerfDataSourceSonar")

	i := &perfApi.PerfDataSourceSonar{}
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

	result, err := r.tryToDeleteCodebasePerfDataSourceSonar(ctx, i, *schema)
	if err != nil || result != nil {
		return *result, err
	}

	log.Info("PerfDataSourceSonar reconciling has been finished successfully")
	return reconcile.Result{}, nil
}

func (r *ReconcilePerfDataSourceSonar) tryToDeleteCodebasePerfDataSourceSonar(ctx context.Context,
	ds *perfApi.PerfDataSourceSonar, schema string) (*reconcile.Result, error) {
	if ds.GetDeletionTimestamp().IsZero() {
		if !helper.ContainsString(ds.ObjectMeta.Finalizers, sonarDataSourceReconcileFinalizerName) {
			ds.ObjectMeta.Finalizers = append(ds.ObjectMeta.Finalizers, sonarDataSourceReconcileFinalizerName)
			if err := r.client.Update(ctx, ds); err != nil {
				return &reconcile.Result{}, err
			}
		}
		return nil, nil
	}

	ow := cluster.GetOwnerReference(codebaseKind, ds.GetOwnerReferences())
	if ow == nil {
		r.log.Info("sonar data source doesn't contain Codebase owner reference", "data source", ds.Name)
		return &reconcile.Result{RequeueAfter: 30 * time.Second}, nil
	}

	if err := r.dsService.RemoveCodebaseDataSource(ow.Name, ds.Spec.Type, schema); err != nil {
		return &reconcile.Result{}, err
	}

	ds.ObjectMeta.Finalizers = helper.RemoveString(ds.ObjectMeta.Finalizers, sonarDataSourceReconcileFinalizerName)
	if err := r.client.Update(ctx, ds); err != nil {
		return &reconcile.Result{}, err
	}
	return &reconcile.Result{}, nil
}
