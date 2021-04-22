package main

import (
	"flag"
	cdPipeApi "github.com/epam/edp-cd-pipeline-operator/v2/pkg/apis/edp/v1alpha1"
	codebaseApi "github.com/epam/edp-codebase-operator/v2/pkg/apis/edp/v1alpha1"
	edpCompApi "github.com/epam/edp-component-operator/pkg/apis/v1/v1alpha1"
	"github.com/epam/edp-gerrit-operator/v2/pkg/controller/helper"
	jenkinsApi "github.com/epam/edp-jenkins-operator/v2/pkg/apis/v2/v1alpha1"
	perfApi "github.com/epam/edp-perf-operator/v2/pkg/apis/edp/v1alpha1"
	reconcilerApi "github.com/epam/edp-reconciler/v2/pkg/apis/edp/v1alpha1"
	"github.com/epam/edp-reconciler/v2/pkg/controller/cdpipeline"
	"github.com/epam/edp-reconciler/v2/pkg/controller/codebase"
	"github.com/epam/edp-reconciler/v2/pkg/controller/codebasebranch"
	edpComponent "github.com/epam/edp-reconciler/v2/pkg/controller/edp-component"
	gitServer "github.com/epam/edp-reconciler/v2/pkg/controller/git_server"
	jenkinsSlave "github.com/epam/edp-reconciler/v2/pkg/controller/jenkins-slave"
	jenkinsJob "github.com/epam/edp-reconciler/v2/pkg/controller/jenkins_job"
	"github.com/epam/edp-reconciler/v2/pkg/controller/jira-server"
	job_provisioning "github.com/epam/edp-reconciler/v2/pkg/controller/job-provisioning"
	"github.com/epam/edp-reconciler/v2/pkg/controller/perfdatasourcejenkins"
	"github.com/epam/edp-reconciler/v2/pkg/controller/perfdatasourcesonar"
	perfserverCtrl "github.com/epam/edp-reconciler/v2/pkg/controller/perfserver"
	"github.com/epam/edp-reconciler/v2/pkg/controller/stage"
	"github.com/epam/edp-reconciler/v2/pkg/controller/thirdpartyservice"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/client-go/rest"
	"os"
	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	//+kubebuilder:scaffold:imports
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

const reconcilerOperatorLock = "edp-reconciler-operator-lock"

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(reconcilerApi.AddToScheme(scheme))

	utilruntime.Must(cdPipeApi.AddToScheme(scheme))

	utilruntime.Must(codebaseApi.AddToScheme(scheme))

	utilruntime.Must(edpCompApi.AddToScheme(scheme))

	utilruntime.Must(jenkinsApi.AddToScheme(scheme))

	utilruntime.Must(perfApi.AddToScheme(scheme))
}

func main() {
	var (
		metricsAddr          string
		enableLeaderElection bool
		probeAddr            string
	)

	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", helper.RunningInCluster(),
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")

	mode, err := helper.GetDebugMode()
	if err != nil {
		setupLog.Error(err, "unable to get debug mode value")
		os.Exit(1)
	}

	opts := zap.Options{
		Development: mode,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	ns, err := helper.GetWatchNamespace()
	if err != nil {
		setupLog.Error(err, "unable to get watch namespace")
		os.Exit(1)
	}

	cfg := ctrl.GetConfigOrDie()
	mgr, err := ctrl.NewManager(cfg, ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     metricsAddr,
		HealthProbeBindAddress: probeAddr,
		Port:                   9443,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       reconcilerOperatorLock,
		MapperProvider: func(c *rest.Config) (meta.RESTMapper, error) {
			return apiutil.NewDynamicRESTMapper(cfg)
		},
		Namespace: ns,
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	ctrlLog := ctrl.Log.WithName("controllers")

	pipelineCtrl, err := cdpipeline.NewReconcileCDPipeline(mgr.GetClient(), mgr.GetScheme(), ctrlLog)
	if err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "cd-pipeline")
		os.Exit(1)
	}

	if err := pipelineCtrl.SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "cd-pipeline")
		os.Exit(1)
	}

	codebaseCtrl := codebase.NewReconcileCodebase(mgr.GetClient(), mgr.GetScheme(), ctrlLog)
	if err := codebaseCtrl.SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "codebase")
		os.Exit(1)
	}

	branchCtrl := codebasebranch.NewReconcileCodebaseBranch(mgr.GetClient(), mgr.GetScheme(), ctrlLog)
	if err := branchCtrl.SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "codebase-branch")
		os.Exit(1)
	}

	componentCtrl := edpComponent.NewEDPComponent(mgr.GetClient(), ctrlLog)
	if err := componentCtrl.SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "edp-component")
		os.Exit(1)
	}

	gitServerCtrl := gitServer.NewReconcileGitServer(mgr.GetClient(), ctrlLog)
	if err := gitServerCtrl.SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "git-server")
		os.Exit(1)
	}

	jenkinsSlaveCtrl := jenkinsSlave.NewReconcileJenkinsSlave(mgr.GetClient(), ctrlLog)
	if err := jenkinsSlaveCtrl.SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "jenkins-slave")
		os.Exit(1)
	}

	jenkinsJobCtrl := jenkinsJob.NewReconcileJenkinsJob(mgr.GetClient(), mgr.GetScheme(), ctrlLog)
	if err := jenkinsJobCtrl.SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "jenkins-job")
		os.Exit(1)
	}

	jiraServerCtrl := jiraserver.NewReconcileJiraServer(mgr.GetClient(), ctrlLog)
	if err := jiraServerCtrl.SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "jira-server")
		os.Exit(1)
	}

	jobProvisionCtrl := job_provisioning.NewReconcileJobProvision(mgr.GetClient(), ctrlLog)
	if err := jobProvisionCtrl.SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "job-provision")
		os.Exit(1)
	}

	pdsjCtrl := perfdatasourcejenkins.NewReconcilePerfDataSourceJenkins(mgr.GetClient(), ctrlLog)
	if err := pdsjCtrl.SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "perf-data-source-jenkins")
		os.Exit(1)
	}

	pdssCtrl := perfdatasourcesonar.NewReconcilePerfDataSourceSonar(mgr.GetClient(), ctrlLog)
	if err := pdssCtrl.SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "perf-data-source-sonar")
		os.Exit(1)
	}

	psCtrl := perfserverCtrl.NewReconcilePerfServer(mgr.GetClient(), ctrlLog)
	if err := psCtrl.SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "perf-server")
		os.Exit(1)
	}

	stageCtrl, err := stage.NewReconcileStage(mgr.GetClient(), mgr.GetScheme(), ctrlLog)
	if err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "cd-stage")
		os.Exit(1)
	}

	if err := stageCtrl.SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "cd-stage")
		os.Exit(1)
	}

	serviceCtrl := thirdpartyservice.NewReconcileService(mgr.GetClient(), ctrlLog)
	if err := serviceCtrl.SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "third-party-service")
		os.Exit(1)
	}

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}

	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
