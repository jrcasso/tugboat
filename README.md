# Tugboat

Tugboat is an orchestrator that manages the provisioning and deprovisioning of application services and dependencies needed to run microservices.

## Impetus

Often times an organization will reach a scale where it becomes necessary to automate the provisioning of resources for new microservices. This is often a laborious process to do by hand, and leaves room for human error, security holes, and configuration drift. Microservice applications typically require involve the provisioning of a code repository, databases, namespaces, secrets, security assurances, and possibly an initial set of skeleton files to make initial development easier for engineers.

Tugboat envdeavors to automate this process for a number of supported services.
