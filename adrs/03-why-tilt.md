# Context
This is a monolith application and it will be for a long time. But we use third-party services, like a database to store our data. So, to make it easier to test/start and develop the application, we need to improve our development environemnt.

# Decision
We firstly started using docker compose, which works percectly, but is not so easy to view your logs and to auto-reload your containers when you change something in your Golang code.

So, we chose (Tilt)[https://tilt.dev/] to improve our dev environment.

This way, we just need to run `tilt up` and it will spawn all our services, also, Tilt provides a web interface at http://localhost:10350/ where you can see your logs etc.

Tilt can run services locally, using docker compose or through a kubernetes cluster. I've chose a kubernetes cluster, because I want later to integrate with some observability tools which mainly operates on kubernetes.

I also want to later deploy this in production using a kubernetes cluster, so I will use most of the kubernetes configs I've created to run this service locally.

# Consequences
## Tilt pros
- Auto-reload Golang projects when something changes
- Easy to setup kuberentes services
- Easy to replicate
- Easy to start the environment
- We can run different scripts to do more things, like cloning a repository
- Inteface to see the logs and to interact with the services

## Tilt cons
Now we need to have a kubernetes clusters running :(

# Status
Approved