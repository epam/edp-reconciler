package stage

import (
	"context"
	"github.com/epmd-edp/reconciler/v2/pkg/db"
	"github.com/epmd-edp/reconciler/v2/pkg/model"
	"github.com/epmd-edp/reconciler/v2/pkg/platform"
	"github.com/epmd-edp/reconciler/v2/pkg/service"
	"log"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	edpv1alpha1 "github.com/epmd-edp/reconciler/v2/pkg/apis/edp/v1alpha1"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// Add creates a new Stage Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	clientSet, err := platform.CreateOpenshiftClients()
	if err != nil {
		panic(err)
	}

	return &ReconcileStage{
		client: mgr.GetClient(),
		scheme: mgr.GetScheme(),
		service: service.StageService{
			DB:        db.Instance,
			ClientSet: *clientSet,
		},
	}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("stage-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	pred := predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			oldObject := e.ObjectOld.(*edpv1alpha1.Stage)
			newObject := e.ObjectNew.(*edpv1alpha1.Stage)

			if oldObject.Status.Value != newObject.Status.Value {
				return true
			}

			if !reflect.DeepEqual(oldObject.Spec, newObject.Spec) {
				return true
			}

			return false
		},
	}

	// Watch for changes to primary resource Stage
	err = c.Watch(&source.Kind{Type: &edpv1alpha1.Stage{}}, &handler.EnqueueRequestForObject{}, pred)
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileStage{}

// ReconcileStage reconciles a Stage object
type ReconcileStage struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client  client.Client
	scheme  *runtime.Scheme
	service service.StageService
}

// Reconcile reads that state of the cluster for a Stage object and makes changes based on the state read
// and what is in the Stage.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileStage) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	log.Println("Reconciling Stage")

	// Fetch the Stage instance
	instance := &edpv1alpha1.Stage{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	log.Printf("Stage: %v", instance)

	stage, _ := model.ConvertToStage(*instance)

	_ = r.service.PutStage(*stage)

	log.Printf("Reconciling Stage %v/%v has been finished", request.Namespace, request.Name)
	return reconcile.Result{Requeue: false}, nil
}
