# [k8sbuilder](https://go.kubebuilder.io/introduction.html)

```shell
curl -L -o kubebuilder https://go.kubebuilder.io/dl/latest/$(go env GOOS)/$(go env GOARCH)
chmod +x kubebuilder && sudo mv kubebuilder /usr/local/bin/

mkdir demo
cd demo
go mod init github.com/futugyousuzu/k8sbuilder
kubebuilder init --domain vishel.io --repo github.com/futugyousuzu/k8sbuilder
kubebuilder edit --multigroup=true

kubebuilder create api --group webapp --kind Guestbook --version v1
kubebuilder create api --group webapp --kind Welcome --version v1

kubebuilder create api --group batch  --kind CronJob --version v1
kubebuilder create api --group batch  --kind CronJob --version v2

kubebuilder create webhook --group batch --version v1 --kind CronJob --defaulting --programmatic-validation 
kubebuilder create webhook --group batch --version v2 --kind CronJob --conversion

kubebuilder create api --group config --version v2 --kind ProjectConfig --resource=true --controller=false --make=false

kubebuilder create api --group apps --kind SimpleDeployment --version v1
kubebuilder create api --group apps --kind ConfigDeployment --version v1

make docker-build IMG=cronjobs:latest
// minikube
minikube image load cronjobs
minikube image ls --format table
make deploy IMG=cronjobs

kubectl apply -f config/samples/batch_v1_cronjob.yaml
kubectl apply -f config/samples/batch_v2_cronjob.yaml
kubectl get cronjob.batch.vishel.io -o yaml
kubectl get job
```

update DeepCopy, DeepCopyInto, and DeepCopyObject

```shell
make generate
```

after etid CRDs

```shell
make manifests
```

Install the CRDs into the cluster

```shell
make install
make install deploy
```

Install Instances of CR

```shell
kubectl apply -f config/crd/patches/  ??
kubectl apply -f config/samples
```

Check

```shell
kubectl get crd
kubectl get welcome
```

Run your controller

```shell
make run
make run ENABLE_WEBHOOKS=false

make uninstall
make undeploy
```

## Description

// TODO(user): An in-depth paragraph about your project and overview of use

## Getting Started

You’ll need a Kubernetes cluster to run against. You can use [KIND](https://sigs.k8s.io/kind) to
get a local cluster for testing, or run against a remote cluster.

**Note:** Your controller will automatically use the current context in your kubeconfig file
(i.e. whatever cluster `kubectl cluster-info` shows).

## Running on the cluster

1. Install Instances of Custom Resources:

    ```sh
    kubectl apply -f config/samples/
    ```

2. Build and push your image to the location specified by `IMG`:

    ```sh
    make docker-build docker-push IMG=<some-registry>/k8sbuilder:tag
    ```

3. Deploy the controller to the cluster with the image specified by `IMG`:

    ```sh
    make deploy IMG=<some-registry>/k8sbuilder:tag
    ```

## Uninstall CRDs

To delete the CRDs from the cluster:

```sh
make uninstall
```

## Undeploy controller

UnDeploy the controller from the cluster:

```sh
make undeploy
```

## Contributing

// TODO(user): Add detailed information on how you would like others to contribute to this project

## How it works

This project aims to follow the Kubernetes [Operator pattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/).

It uses [Controllers](https://kubernetes.io/docs/concepts/architecture/controller/),
which provide a reconcile function responsible for synchronizing resources until the desired
state is reached on the cluster.

## Test It Out

1. Install the CRDs into the cluster:

    ```sh
    make install
    ```

2. Run your controller (this will run in the foreground, so switch to a new terminal
if you want to leave it running):

    ```sh
    make run
    ```

**NOTE:** You can also run this in one step by running: `make install run`

## Modifying the API definitions

If you are editing the API definitions, generate the manifests such as CRs or CRDs using:

```sh
make manifests
```

**NOTE:** Run `make --help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)
