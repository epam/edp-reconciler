package codebasebranch

import (
	"context"
	codebaseApi "github.com/epam/edp-codebase-operator/v2/pkg/apis/edp/v1alpha1"
	"github.com/epam/edp-reconciler/v2/pkg/controller/helper"
	"github.com/epam/edp-reconciler/v2/pkg/db"
	"github.com/epam/edp-reconciler/v2/pkg/model/codebasebranch"
	cbs "github.com/epam/edp-reconciler/v2/pkg/service/codebasebranch"
	"github.com/go-logr/logr"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"time"

	errWrap "github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const codebaseBranchReconcileFinalizerName = "codebasebranch.reconciler.finalizer.name"

func NewReconcileCodebaseBranch(client client.Client, scheme *runtime.Scheme, log logr.Logger) *ReconcileCodebaseBranch {
	return &ReconcileCodebaseBranch{
		client: client,
		scheme: scheme,
		branch: cbs.CodebaseBranchService{
			DB: db.Instance,
		},
		log: log.WithName("codebase-branch"),
	}
}

type ReconcileCodebaseBranch struct {
	client client.Client
	scheme *runtime.Scheme
	branch cbs.CodebaseBranchService
	log    logr.Logger
}

func (r *ReconcileCodebaseBranch) SetupWithManager(mgr ctrl.Manager) error {
	p := predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			oldObject := e.ObjectOld.(*codebaseApi.CodebaseBranch)
			newObject := e.ObjectNew.(*codebaseApi.CodebaseBranch)

			if oldObject.Status.Value != newObject.Status.Value ||
				oldObject.Status.Action != newObject.Status.Action {
				return true
			}

			if !reflect.DeepEqual(oldObject.Spec, newObject.Spec) {
				return true
			}

			if oldObject.Status.LastSuccessfulBuild != newObject.Status.LastSuccessfulBuild {
				return true
			}

			if oldObject.Status.Build != newObject.Status.Build {
				return true
			}

			if newObject.DeletionTimestamp != nil {
				return true
			}
			return false
		},
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&codebaseApi.CodebaseBranch{}, builder.WithPredicates(p)).
		Complete(r)
}

func (r *ReconcileCodebaseBranch) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	log := r.log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	log.Info("Reconciling CodebaseBranch")

	i := &codebaseApi.CodebaseBranch{}
	if err := r.client.Get(ctx, request.NamespacedName, i); err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	edpN, err := helper.GetEDPName(r.client, i.Namespace)
	if err != nil {
		return reconcile.Result{RequeueAfter: 2 * time.Second}, errWrap.Wrap(err, "couldn't get edp name")
	}

	if res, err := r.tryToDeleteCodebaseBranch(ctx, i, *edpN); err != nil || res != nil {
		return *res, err
	}

	app, err := codebasebranch.ConvertToCodebaseBranch(*i, *edpN)
	if err != nil {
		return reconcile.Result{RequeueAfter: 2 * time.Second}, errWrap.Wrap(err, "cannot convert to codebase branch dto")
	}
	if err := r.branch.PutCodebaseBranch(*app); err != nil {
		return reconcile.Result{RequeueAfter: 2 * time.Second}, errWrap.Wrap(err, "couldn't insert codebase branch")
	}
	log.Info("Reconciling has been finished successfully")
	return reconcile.Result{}, nil
}

func (r *ReconcileCodebaseBranch) tryToDeleteCodebaseBranch(ctx context.Context, cb *codebaseApi.CodebaseBranch, schema string) (*reconcile.Result, error) {
	if cb.GetDeletionTimestamp().IsZero() {
		if !helper.ContainsString(cb.ObjectMeta.Finalizers, codebaseBranchReconcileFinalizerName) {
			cb.ObjectMeta.Finalizers = append(cb.ObjectMeta.Finalizers, codebaseBranchReconcileFinalizerName)
			if err := r.client.Update(ctx, cb); err != nil {
				return &reconcile.Result{}, err
			}
		}
		return nil, nil
	}

	if err := r.branch.Delete(cb.Spec.CodebaseName, cb.Spec.BranchName, schema); err != nil {
		return &reconcile.Result{RequeueAfter: 2 * time.Second}, err
	}

	cb.ObjectMeta.Finalizers = helper.RemoveString(cb.ObjectMeta.Finalizers, codebaseBranchReconcileFinalizerName)
	if err := r.client.Update(ctx, cb); err != nil {
		return &reconcile.Result{RequeueAfter: 2 * time.Second}, err
	}
	return &reconcile.Result{}, nil
}
