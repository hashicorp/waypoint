# Waypoint

Waypoint is a tool for deploying applications and services to any target. With Waypoint's predictable and uncomplicated workflow, you can deploy your code without spending days understanding how to interface with various endpoints. 

## Running Waypoint on Docker Desktop for Mac with Kubernetes

1. Install and Configure Docker Desktop for Kubernetes.
    1. Download and install Docker Desktop for Mac: https://docs.docker.com/docker-for-mac/install/
    1. Enable Kubernetes in Docker Desktop. It may take up to 5 minutes for Docker to apply your changes.Docker will restart and create a Kubernetes cluster.
    ![Image of Kubernetes Settings](./docs/images/d4m-k8s.png)
1. Configure Kubernetes to access the GitHub Docker registry
    1. Create a personal access token on GitHub: https://github.com/settings/tokens/new
    1. Configure Docker to authenticate to GitHub Packages: https://help.github.com/en/packages/using-github-packages-with-your-projects-ecosystem/configuring-docker-for-use-with-github-packages#authenticating-to-github-packages
    1. Add your personal access token to the Docker registry by running the following command in the Terminal:  
    ``kubectl create secret docker-registry github --docker-server=docker.pkg.github.com --docker-username=my-github-username --docker-password=my-github-personal-access-token --docker-email=my-github-email-address``
1. Download or build Waypoint.  
  Note: Be sure to [add Waypoint to your $PATH](#-Add-Waypoint-to-Your-Path).
    * To use the compiled version:
      1. Clone the Waypoint repository by using http. You will need to use configuration files and a test app from the Waypoint repository to complete this tutorial.
      1. Download the binary here: https://github.com/hashicorp/waypoint/releases/download/v0.1.0/waypoint-darwin-0.1.0.zip
      1. Extract the binary to a location outside of Downloads, such as ~/enlistments/waypoint/.
    * To build Waypoint:
      1. Clone the Waypoint repository by using http.
      1. In the Terminal, navigate to the **waypoint** folder and run `make bin`
1. Create a Waypoint URL Service account.
    1. Run `waypoint account register -email my-email-address -accept-eula` to create the token. 
    1. Copy the token output in the previous step and use it as the value for WAYPOINT\_URL\_TOKEN in your environment: `export WAYPOINT_URL_TOKEN=my-token-from-previous-step`
1. Create a Waypoint URL Service host.
    * Run `waypoint hostname register -l service=wpmini,env=dev` to create an auto-generated hostname.
    * To specify your own hostname, add the `-name my-hostname` option to the command.
1. Install the Waypoint server into your Kubernetes cluster.
    1. Run `waypoint install` to install the Waypoint server. The installation will download the server's Docker image.
    1. Once the installation is complete, verify that the server is running: `kubectl get pod -l app=waypoint-server`
    * The server's pod status should be **Running**.
    * If the pod status is **ImagePullBackOff**, your Docker registry access token is incorrect. Review step 2 to verify that you have properly set up your registry secret.
1. Add a Docker registry to your Waypoint install by running  
    `docker run -d -p 5000:5000 --restart=always --name registry registry:2`
1. In the Terminal, navigate to the **test-apps/wpmini** Waypoint directory.
1. Use Waypoint to build the wpmini app: `waypoint build`
1. Deploy the app with Waypoint: `waypoint deploy`
    * TODO: Make Ctrl-C cancel the deploy properly if waiting on the pods to start.
1. Release the app to make it available on your local server: `waypoint release`
1. Verify that the app is running: `curl localhost:8080`
   * The output should read, **Welcome to Waypoint!**
1. View the Waypoint logs to see requests for the app and the app's heartbeat: `waypoint logs`
1. Update the app's config to learn how it is constructed.
    1. Run `waypoint config set NAME=my-name`
      * NAME must be capitalized.
    1. Run `waypoint deploy`
    1. Run `waypoint release`
    1. Run `curl localhost:8080`
    * The new response will read: "Welcome your-selected-name, Welcome to Waypoint!"
1. Access the app with the Waypoint URL Service. 
    1. Use the hostname you registered in step 5 to run `curl https://my-registered-hostname.alpha.waypoint.run`

## Add Waypoint to Your Path
You will want to add Waypoint to your $PATH to use the `waypoint` command without always needing to include the full directory path.
1. Run `echo $PATH` to see what is in your current $PATH definition. If you do not have an entry for the Waypoint binary, add the following line to your shell profile:  
`export PATH=$PATH:my-directory-path-to-waypoint-binary`  
    For example: `export PATH=$PATH:~/enlistments/waypoint`
* To determine which shell you use, run `echo $SHELL`
* If using Bash, use `nano ~/.bash_profile` to edit your Bash profile.
* If using ZSH, use `nano ~/.zshrc` to edit your ZSH profile.
