package git_server

import (
	"context"

	"github.com/epam/edp-codebase-operator/v2/pkg/apis/edp/v1alpha1"
	"github.com/epmd-edp/reconciler/v2/pkg/controller/helper"
	"github.com/epmd-edp/reconciler/v2/pkg/db"
	jiramodel "github.com/epmd-edp/reconciler/v2/pkg/model/jira-server"
	jiraserver "github.com/epmd-edp/reconciler/v2/pkg/service/jira-server"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_jira_server")

func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileJiraServer{
		client:  mgr.GetClient(),
		service: jiraserver.JiraServerService{DB: db.Instance},
	}
}

func add(mgr manager.Manager, r reconcile.Reconciler) error {
	c, err := controller.New("jira-server-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	p := predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			oldObject := e.ObjectOld.(*v1alpha1.JiraServer)
			newObject := e.ObjectNew.(*v1alpha1.JiraServer)
			if oldObject.Spec != newObject.Spec {
				return true
			}
			if oldObject.Status.Available != newObject.Status.Available {
				return true
			}
			return false
		},
	}

	if err = c.Watch(&source.Kind{Type: &v1alpha1.JiraServer{}}, &handler.EnqueueRequestForObject{}, p); err != nil {
		return err
	}
	return nil
}

var _ reconcile.Reconciler = &ReconcileJiraServer{}

type ReconcileJiraServer struct {
	client  client.Client
	service jiraserver.JiraServerService
}

func (r *ReconcileJiraServer) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	rl := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	rl.V(2).Info("Reconciling JiraServer")

	i := &v1alpha1.JiraServer{}
	if err := r.client.Get(context.TODO(), request.NamespacedName, i); err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	tenant, err := helper.GetEDPName(r.client, i.Namespace)
	if err != nil {
		return reconcile.Result{}, err
	}

	if err := r.service.PutJiraServer(jiramodel.ConvertSpecToJira(*i, *tenant)); err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}
