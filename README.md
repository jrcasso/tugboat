# Tugboat

Tugboat is an orchestrator that manages the provisioning and deprovisioning of application services and dependencies needed to run in microservice architecures.

## Impetus

Often times an organization will reach a scale where it becomes necessary to automate the provisioning of resources for new microservices. This is often a laborious process to do by hand, and leaves room for human error, security holes, and configuration drift. Microservice applications typically require involve the provisioning of a code repository, databases, namespaces, secrets, security assurances, and possibly an initial set of skeleton files to make initial development easier for engineers.

Tugboat envdeavors to automate this process for a number of supported services.

## Design Principles

- The Tugboat approach is that of convention over configuration; a consistent set of sensible defaults across providers is used.
- Conventions need not impact flexibility - defaults may be overridden for each provider.

## Currently Supported Providers

The initial provisioning feature-set is small while providers are being flushed out.

| Service | Description |
|-|-|
| GitHub | Tugboat will create a GitHub repository for the new microservice if it does not already exist. |
| Kubernetes | Tugboat will create a Kubernetes namespace for the new microservice if it does not already exist. |

# Development Setup

This is likely the easiest local development setup process you'll have ever encountered if you use VS Code.

Ensure you have the following prerequisites satisfied:
 - Docker for Desktop
 - VS Code Extensions: Remote Containers
   - Download and install Microsoft's VS Code extension for developing in [Remote Containers](vscode:extension/ms-vscode-remote.remote-containers)

>Note: This is a VS Code Remote Containers development project: all development is done within a container to reduce initial time-to-develop. Getting this project up and running on your machine can be as simple as pulling down the repository, running the Docker daemon the host machine, opening the project in VS Code, and clicking twice.


## Directions

- Clone the repository

```sh
git clone git@github.com:jrcasso/tugboat
```

- Open the repository in VS Code
```sh
code tugboat
```
- Ensure you have the necessary environment variables to run this application set in your local environment; they will be used transferred to the remote container environment. You can see which are required in the `.devcontainer/devcontainer.json` file under the `remoteEnv` key.
- In the bottom-left corner of the VS Code window, click the highlighted "><" button (or navigate to the Remote Containers extension).
- From the dropdown, select "Remote Containers: Reopen in Container"

_**That's it!**_

>Note: When you enter the remote VS Code environment for the first time, a pop-up will appear in the corner indicating that you should install some of the go binaries required for tooling (e.g. gopls, dlv-dap - a Go debugging server). Click "Install", and these binaries will be persisted for all future spin-ups via a docker-compose volume for the environment.

## Development Details

VS Code will begin to build an image that is specified in `.devcontainer/`; it will be the container image that you develop in. When it's done, it'll automatically throw your entire VS Code interface/environment inside that container where you may begin deveopment. The current configuration will also mount your Docker engine socket into this container, so that Docker commands may be issued from within to manage containers on the host. Utilitarian tools like git and all the things needed to run a Go program are in that environment. It's still a container, so all of the idempotency and innate destructivity of containers are in fact *features* of this development strategy. If everyone develops in the same way, the time-to-develop becomes incredibly small.

Additional tooling that might be needed can be done so during container runtime; however, if it is something that should stick around for every other developer too (i.e. they might also run into this same issue), please modify the `.devcontainer/Dockerfile` and open a pull request.
