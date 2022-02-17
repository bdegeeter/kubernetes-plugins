// +build mage

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	. "get.porter.sh/plugin/kubernetes/mage"
	"github.com/carolynvs/magex/mgx"
	"github.com/carolynvs/magex/pkg"
	"github.com/carolynvs/magex/pkg/gopath"
	"github.com/carolynvs/magex/shx"
	"github.com/magefile/mage/mg"
	"github.com/pkg/errors"

	// mage:import
	. "get.porter.sh/porter/mage/tests"
)

const (
	// Version of KIND to install if not already present
	kindVersion = "v0.10.0"

	// Name of the KIND cluster used for testing
	kindClusterName = "porter"

	// Namespace where you can do manual testing
	testNamespace = "test"

	// Relative location of the KUBECONFIG for the test cluster
	kubeconfig = "kind.config"

	// Namespace of the porter plugin
	pluginNamespace = "porter-plugin-test-ns"

	// Container name of the local registry
	registryContainer = "registry"
)

// Build a command that stops the build on if the command fails
var must = shx.CommandBuilder{StopOnError: true}

// We are migrating to mage, but for now keep using make as the main build script interface.

// Publish the cross-compiled binaries.
// func Publish(plugin string, version string, permalink string) {
// 	releases.PreparePluginForPublish(plugin, version, permalink)
// 	releases.PublishPlugin(plugin, version, permalink)
// 	releases.PublishPluginFeed(plugin, version)
// }

// Add GOPATH/bin to the path on the GitHub Actions agent
// TODO: Add to magex
func addGopathBinOnGithubActions() error {
	githubPath := os.Getenv("GITHUB_PATH")
	if githubPath == "" {
		return nil
	}

	log.Println("Adding GOPATH/bin to the PATH for the GitHub Actions Agent")
	gopathBin := gopath.GetGopathBin()
	return ioutil.WriteFile(githubPath, []byte(gopathBin), 0644)
}

// Ensure mage is installed.
func EnsureMage() error {
	addGopathBinOnGithubActions()
	return pkg.EnsureMage("v1.11.0")
}

func Fmt() {
	must.RunV("go", "fmt", "./...")
}

func Vet() {
	must.RunV("go", "vet", "./...")
}

// Run integration tests against the test cluster.
func TestIntegration() {
	mg.Deps(UseTestEnvironment)
	// mg.Deps(UseTestEnvironment, EnsureDeployed)

	must.RunV("go", "test", "-tags=integration", "./...", "-coverprofile=coverage-integration.out")
}

// Build the operator and deploy it to the test cluster using
// func DeployOperator() {
// 	mg.Deps(UseTestEnvironment, EnsureTestCluster, EnsureKubectl)

// meta := releases.LoadMetadata()
// PublishLocalPorterAgent()
// PublishBundle()

// porter("credentials", "apply", "hack/creds.yaml", "-n=operator", "--debug", "--debug-plugins").Must().RunV()
// bundleRef := Env.BundlePrefix + meta.Version
// porter("install", "operator", "-r", bundleRef, "-c=kind", "--force", "-n=operator").Must().RunV()
// }

// get the config of the current kind cluster, if available
func getClusterConfig() (kubeconfig string, ok bool) {
	contents, err := shx.OutputE("kind", "get", "kubeconfig", "--name", kindClusterName)
	return contents, err == nil
}

// setup environment to use the current kind cluster, if available
func useCluster() bool {
	contents, ok := getClusterConfig()
	if ok {
		log.Println("Reusing existing kind cluster")

		userKubeConfig, _ := filepath.Abs(os.Getenv("KUBECONFIG"))
		currentKubeConfig := filepath.Join(pwd(), kubeconfig)
		if userKubeConfig != currentKubeConfig {
			fmt.Printf("ATTENTION! You should set your KUBECONFIG to match the cluster used by this project\n\n\texport KUBECONFIG=%s\n\n", currentKubeConfig)
		}
		os.Setenv("KUBECONFIG", currentKubeConfig)

		err := ioutil.WriteFile(kubeconfig, []byte(contents), 0644)
		mgx.Must(errors.Wrapf(err, "error writing %s", kubeconfig))

		setClusterNamespace(pluginNamespace)
		return true
	}

	return false
}

func setClusterNamespace(name string) {
	shx.RunE("kubectl", "config", "set-context", "--current", "--namespace", name)
}

// Check if the operator is deployed to the test cluster.
// func EnsureDeployed() {
// 	if !isDeployed() {
// 		DeployOperator()
// 	}
// }

func isDeployed() bool {
	if useCluster() {
		if err := kubectl("rollout", "status", "deployment", "porter-operator-controller-manager", "--namespace", pluginNamespace).Must(false).Run(); err != nil {
			log.Println("the operator is not installed")
			return false
		}
		if err := kubectl("rollout", "status", "deployment", "mongodb", "--namespace", pluginNamespace).Must(false).Run(); err != nil {
			log.Println("the database is not installed")
			return false
		}
		log.Println("the operator is installed and ready to use")
		return true
	}
	log.Println("could not connect to the test cluster")
	return false
}

// Remove the test cluster and registry.
func Clean() {
	//mg.Deps(DeleteTestCluster, StopDockerRegistry)
	mg.Deps(DeleteTestCluster)
	os.RemoveAll("bin")
}

func pwd() string {
	wd, _ := os.Getwd()
	return wd
}

func kubectl(args ...string) shx.PreparedCommand {
	kubeconfig := fmt.Sprintf("KUBECONFIG=%s", os.Getenv("KUBECONFIG"))
	return must.Command("kubectl", args...).Env(kubeconfig)
}

// Run porter using the local storage, not the in-cluster storage
func porter(args ...string) shx.PreparedCommand {
	return shx.Command("porter").Args(args...).
		Env("PORTER_DEFAULT_STORAGE=", "PORTER_DEFAULT_STORAGE_PLUGIN=mongodb-docker")
}
