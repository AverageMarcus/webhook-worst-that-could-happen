# Webhooks - What's the worst that could happen?

This repo contains proof of concepts for various misconfigured and malicious Kubernetes webhook configurations.

It was designed to aid in putting together a talk I gave at [Kubernetes Community Days UK](https://community.cncf.io/events/details/cncf-kcd-uk-presents-kubernetes-community-days-uk-2023/) and [Cloud Native Rejekts](https://cloud-native.rejekts.io/) in 2023.

[Accompanying Slides](https://speaking.marcusnoble.co.uk/4YvpTx/webhooks-whats-the-worst-that-could-happen)

> ⚠️ Don't run these scenarios in a cluster you care about!

## Building & Running

1. Build the container image:
    ```sh
    make docker-build
    ```
2. Publish the container image:
    ```sh
    make docker-publish
    ```
3. Create a new Kind cluster. (To ensure all the needed configuration is set I use this script: [kind-create-cluster](https://github.com/AverageMarcus/dotfiles/blob/master/home/.bin/kind-create-cluster))
4. Once the Kind cluster is ready, deploy the chosen scenario:
    ```sh
    make deploy-[scenario]
    ```
    e.g. for the `fork-bomb` scenario
    ```sh
    make deploy-fork-bomb
    ```
5. Perform whatever action is required to trigger the scenario. In most cases this is just creating a new pod.
