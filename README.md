# Waypoint

Waypoint is a tool for deploying applications and services to any target. With Waypoint's predictable and uncomplicated workflow, you can deploy your code without spending days understanding how to interface with various endpoints.

## Running Waypoint on Docker Desktop for Mac with Kubernetes

1. Install and configure Docker Desktop for Kubernetes.
    1. Download and install Docker Desktop for Mac: https://docs.docker.com/docker-for-mac/install/
    1. Enable Kubernetes in Docker Desktop. It may take up to 5 minutes for Docker to apply your changes. Docker will restart and create a Kubernetes cluster.
    ![Image of Kubernetes Settings](./docs/images/d4m-k8s.png)
1. Configure Kubernetes to access the GitHub Docker registry.
    1. Create a personal access token on GitHub: https://github.com/settings/tokens/new
    1. Configure Docker to authenticate to GitHub Packages: https://help.github.com/en/packages/using-github-packages-with-your-projects-ecosystem/configuring-docker-for-use-with-github-packages#authenticating-to-github-packages
    1. Add your personal access token to the Docker:
    ``kubectl create secret docker-registry github --docker-server=docker.pkg.github.com --docker-username=my-github-username --docker-password=my-github-personal-access-token --docker-email=my-github-email-address``
1. Download or build Waypoint.
  Note: Be sure to [add Waypoint to your $PATH](#add-waypoint-to-your-path).
    * To use the compiled version:
      1. Clone the Waypoint repository by using http. You will need to use configuration files and a test app from the Waypoint repository to complete this tutorial.
      1. Download the binary here: https://github.com/hashicorp/waypoint/releases/download/v0.1.0/waypoint-darwin-0.1.0.zip
      1. Extract the binary to a location outside of Downloads, such as ~/enlistments/waypoint/.
    * To build Waypoint:
      1. Clone the Waypoint repository by using http.
      1. In the Terminal, navigate to the **waypoint** folder and run `make bin`
1. Install the Waypoint server into your Kubernetes cluster.
    1. Install the Waypoint server:
    `waypoint install`
    The installation will download the server's Docker image.
    1. Once the installation is complete, verify that the server is running:
    `kubectl get pod -l app=waypoint-server`
    * The server's pod status should be **Running**.
    * If the pod status is **ImagePullBackOff**, your Docker registry access token is incorrect. Review step 2 to verify that you have properly set up your registry secret.
1. In the Terminal, navigate to the **test-apps/wpmini** Waypoint directory.
1. Initialize the wpmini app with Waypoint:
    `waypoint init`
1. Use Waypoint to build the wpmini app:
    `waypoint build`
1. Deploy the app with Waypoint:
    `waypoint deploy`
    * TODO: Make Ctrl-C cancel the deploy properly if waiting on the pods to start.
1. Verify that the app is running:
    `curl localhost:8080`
   * The output should read, **Welcome to Waypoint!**
1. View the Waypoint logs to see requests for the app and the app's heartbeat:
    `waypoint logs`
1. Update the app's config to learn how it is constructed.
    1. Run `waypoint config set NAME=my-name`
    Note: NAME must be capitalized.
    1. Run `waypoint deploy`
    1. Run `curl localhost:8080`
    * The new response will read: "Hello your-selected-name, Welcome to Waypoint!"
1. Access the app with the Waypoint URL Service.
    1. Register a hostname by executing `waypoint hostname register` (no args or flags)
	2. Access the domain that is output in the previous step
	`curl https://<domain>`

## Add Waypoint to Your Path

You will want to add Waypoint to your $PATH to use the `waypoint` command
without always needing to include the full directory path.

1. Run `echo $PATH` to see what is in your current $PATH definition. If you do not have an entry for the Waypoint binary, add the following line to your shell profile:
`export PATH=$PATH:my-directory-path-to-waypoint-binary`
    For example: `export PATH=$PATH:~/enlistments/waypoint`
* To determine which shell you use, run `echo $SHELL`
* If using Bash, use `nano ~/.bash_profile` to edit your Bash profile.
* If using ZSH, use `nano ~/.zshrc` to edit your ZSH profile.

## Removing Waypoint Server

Waypoint Server creates several resources in Kubernetes and Docker that should be removed to either reinstall Waypoint or to completely remove it from a system.

### Waypoint Server in Kubernetes

`waypoint install` for Kubernetes creates a StatefulSet, Service and PersistentVolumeClaim. These resources should be removed when Waypoint Server is no longer needed. These are some example `kubectl` commands that should clean up after a Waypoint Server installation.

```
kubectl delete statefulset waypoint-server
kubectl get pv #note that pv name associated with $NAMESPACE/data-waypoint-server-0
kubectl delete pv pvc-9063b5f7-dc86-4fa7-b120-bf084bbbbf93 #this id will be different for you
kubectl delete svc waypoint
```

### Waypoint Server in Docker

`waypoint install` for Docker creates a container and a volume. These resources should be removed when Waypoint Server is no longer needed. These are some example `docker` commands that should clean up after a Waypoint Server installation.

```
docker stop waypoint-server
docker rm waypoint-server
docker volume prune -f
```
