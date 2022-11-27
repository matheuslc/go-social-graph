# Nice to see you here!
Feel free to contribute with features, fixes, new ideas, architecture improvements, documentation etc.

## API First
We follow an API first approach, where we first define the contracts for our REST API an later implement the logic to meet the contract requirements. We use the Open API spec version 3.0 to defines our REST API and the [oapi-codegen](https://github.com/deepmap/oapi-codegen) library to generate our server (echo) interface.

### Steps
To add a new REST endpoint

* Define your endpoint spec at openapi/openapi.yml.
* Define the types you endpoint use.
* Run go:generate to generante the code for your new endpoint
* Imeplement the handler and the business logic

## Running the API
We have some ways to run this project. The most recommended way is using [Tilt](https://tilt.dev/). Tilt let's us define our dev environment using code. We can run services on localhost, docker compose or using kubernetes.

For this project, we are using Kubernetes, so you must have a k8s cluster running on your machine. I recommend [Rancher Desktop](https://www.rancher.com/).

Make sure to have installed Golang ([asdf](https://asdf-vm.com/) recommended), [Helm](https://helm.sh/docs/intro/install/), Nodejs ([asdf](https://asdf-vm.com/), and Docker.

Now, it's the easiest part. Just run `tilt up --namespace=internal`. Tilt will provide you an interface where you can start some services and test the API.

### Alternatives
You can directly run the cmd/main.go file, just make sure to fill all environment variables. Also, make sure to have all dependencies running, like neo4j.


