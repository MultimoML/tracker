# Tracker

Microservice returns all new products. 

Available endpoints:
- [`/live`](https://multimo.ml/tracker/live): Liveliness check
- [`/ready`](https://multimo.ml/tracker/ready): Readiness check
- [`/v1/all`](https://multimo.ml/tracker/v1/all): Returns a list of all new products

Branches:
- [`main`](https://github.com/MultimoML/tracker/tree/main): Contains latest development version
- [`prod`](https://github.com/MultimoML/tracker/tree/prod): Contains stable, tagged releases

## Setup/installation

Prerequisites:
- [Go](https://go.dev/)
- [Docker](https://www.docker.com/)

Example usage:
- See all available options: `make help`
- Run microservice in a container: `make run`
- Release a new version: `make release ver=x.y.z`

All work should be done on `main`, `prod` should never be checked out or manually edited.
When releasing, the changes are merged into `prod` and both branches are pushed.
A GitHub Action workflow will then build and publish the image to GHCR, and deploy it to Kubernetes.

## License

Multimo is licensed under the [GNU AGPLv3 license](LICENSE).
