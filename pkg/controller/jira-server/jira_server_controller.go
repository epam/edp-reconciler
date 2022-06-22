package jiraserver

import (
	"context"

	codebaseApi "github.com/epam/edp-codebase-operator/v2/pkg/apis/edp/v1"
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
	jiramodel "github.com/epam/edp-reconciler/v2/pkg/model/jira-server"
	jiraserver "github.com/epam/edp-reconciler/v2/pkg/service/jira-server"
)

func NewReconcileJiraServer(client client.Client, log logr.Logger) *ReconcileJiraServer {
	return &ReconcileJiraServer{
		client: client,
		jiraServer: jiraserver.JiraServerService{
			DB: db.Instance,
		},
		log: log.WithName("jira-server"),
	}
}

type ReconcileJiraServer struct {
	client     client.Client
	jiraServer jiraserver.JiraServerService
	log        logr.Logger
}

func (r *ReconcileJiraServer) SetupWithManager(mgr ctrl.Manager) error {
	p := predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			oldObject := e.ObjectOld.(*codebaseApi.JiraServer)
			newObject := e.ObjectNew.(*codebaseApi.JiraServer)
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
		For(&codebaseApi.JiraServer{}, builder.WithPredicates(p)).
		Complete(r)
}

func (r *ReconcileJiraServer) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	log := r.log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	log.V(2).Info("Reconciling JiraServer")

	i := &codebaseApi.JiraServer{}
	if err := r.client.Get(ctx, request.NamespacedName, i); err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	tenant, err := helper.GetEDPName(r.client, i.Namespace)
	if err != nil {
		return reconcile.Result{}, err
	}

	if err := r.jiraServer.PutJiraServer(jiramodel.ConvertSpecToJira(*i, *tenant)); err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}
