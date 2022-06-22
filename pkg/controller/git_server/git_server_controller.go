package git_server

import (
	"context"

	codebaseApi "github.com/epam/edp-codebase-operator/v2/pkg/apis/edp/v1"
	"github.com/go-logr/logr"
	errWrap "github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/epam/edp-reconciler/v2/pkg/controller/helper"
	"github.com/epam/edp-reconciler/v2/pkg/db"
	"github.com/epam/edp-reconciler/v2/pkg/model/gitserver"
	"github.com/epam/edp-reconciler/v2/pkg/service/git"
	"github.com/epam/edp-reconciler/v2/pkg/service/infrastructure"
)

func NewReconcileGitServer(client client.Client, log logr.Logger) *ReconcileGitServer {
	return &ReconcileGitServer{
		client: client,
		git: git.GitServerService{
			DB: db.Instance,
		},
		infraDb: infrastructure.InfrastructureDbService{
			DB: db.Instance,
		},
		log: log.WithName("git-server"),
	}
}

type ReconcileGitServer struct {
	client  client.Client
	git     git.GitServerService
	infraDb infrastructure.InfrastructureDbService
	log     logr.Logger
}

func (r *ReconcileGitServer) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&codebaseApi.GitServer{}).
		Complete(r)
}

func (r *ReconcileGitServer) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	log := r.log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	log.Info("Reconciling GitServer")

	instance := &codebaseApi.GitServer{}
	if err := r.client.Get(ctx, request.NamespacedName, instance); err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}
	log.WithValues("GitServer", instance)
	edpN, err := helper.GetEDPName(r.client, instance.Namespace)
	if err != nil {
		return reconcile.Result{}, err
	}
	gitServer, err := gitserver.ConvertToGitServer(*instance, *edpN)
	if err != nil {
		return reconcile.Result{}, err
	}

	exists, err := r.infraDb.DoesSchemaExist(gitServer.Tenant)
	if err != nil {
		return reconcile.Result{}, errWrap.Wrap(err, "an error has occurred while checking schema in BD")
	}
	log.Info("Check schema: ", "schema", gitServer.Tenant, "exists", exists)

	if exists {
		err := r.git.PutGitServer(*gitServer)
		if err != nil {
			return reconcile.Result{}, err
		}

	}

	return reconcile.Result{}, nil
}
