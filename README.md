# ADH Partnership IDS

This is still very much a work in progress. This README will be expanded as development progresses.

## Introduction

This is a rewrite of the [ZAN IDS](https://github.com/vpaza/ids) to be donated to the [ADH Partnership](https://github.com/adh-partnership) with some caching and efficiencies written in, and designed to be configurable and usable by multiple facilities. This is still very much a work in progress.

While the ADH Partnership is the target, considerations are being made in the design that will allow it to be implemented in most environments, though some concessions will need to be made. For example, members of the partnership will be able to restrict the IDS to members on the roster only whereas non-members will not be able to do so due to the differences in the API results.

There may be some consideration of VATUSA API integration once the new API is released.

## Requirements

This project will also rely on the [chart-parser](https://github.com/adh-partnership/chart-parser) to populate the charts table. This is a separate project and will need to be run via cron or a kubernetes CronJob.

## Configuration

### Backend Configuration

Please consult the [backend README](backend/README.md) for configuration information.

### Frontend Configuration

Please consult the [frontend README](frontend/README.md) for configuration information.

### Shared

There are some items that are shared between backend and frontend. This configuration is shared so that we can reduce
the amount of duplication.

Please consult the [shared README](shared/README.md) for configuration information.

## Local Development

To run locally for development purposes, you will need:

- [Golang](https://go.dev) v1.21
- [NodeJS](https://nodejs.org/en/) LTS
- [Yarn](https://yarnpkg.com/)

You will also need DNS entries pointing at localhost, for the purpose of cookies. Browsers will typically not store,
send, or receive cookies when the domain is `localhost`.

**NOTE** Commands below assume your current working directory is the project root (where this README)
resides.

### Authentication

If you are using VATSIM Connect, you should use the development credentials. More information can be found [here](https://vatsim.dev/services/connect/sandbox).

If you are using the ADH Stack, you will need to add the OAuth2 client credentials to the database. Please consult the ADH Partnership members if you do not know how to do this already.

### Backend

1. Create a config.yaml file in the root project directory. This is already part of the `.gitignore`. For more information,
  please consult the [backend README](backend/README.md).
2. If you need to seed the database, make sure shared/airports.json is correct, and run:

  ```shell
  go run backend/cmd/api/main.go --log-level debug update
  ```

3. To start the backend, run:
  
    ```shell
    go run backend/cmd/api/main.go --log-level debug server
    ```

**Note**: This assumes that airports.json is in shared/ and the config file is in the current directory.

### Frontend

1. Make sure `frontend/config.json` is populated with the desired values. See the 
  [frontend README](frontend/README.md) for more information.
2. Make sure all modules are installed by running `yarn`.
3. Run `yarn dev` to start the dev environment.

## Subdivision Setup

If you are a member of the ADH Partnership, getting started is fairly easy.

1. Open a PR to add your facility's frontend configuration, and adding it to the workflow [here](.github/workflows/build.yaml).
2. Setup your cluster as appropriate. You can look at ZAN's GitOps repo for more inforation [here](https://github.com/vpaza/gitops/tree/main/overlays/prod/ids).
   1. You will need to setup the configmap to load the backend's configuration. Override any sensitive information (such as session's block and hash keys, database credentials, etc)
      utilizing secrets. See the [backend README](backend/README.md) for more information.

If you are not apart of the ADH Partnership, you'll need to do a bit more work as we do not provide building and deployment for non-members.

1. Typically you'll need to clone the repo.
2. You'll need to provide the backend with a YAML configuration. See the [backend README](backend/README.md) for more information. This can be done at runtime.
3. Update frontend/config.json with your desired config. This *must* be done prior to building as the config is embedded into the frontend.
4. Use scripts/build.sh to build the container image. If you don't want to use Docker, you'll need to look at the script and Dockerfile to see the steps needed
   to build the frontend and backend.
5. Deploy to your infrastructure as appropriate for it. The ADH Partnership utilizes Kubernetes, so you can look at step #2 above for an example. Any other environments
   is up to your webmaster to deploy. You will need a reverse proxy to handle the backend, the frontend is static JavaScript and can be served by any web server.

If you have questions, please find Daniel Hawton (@TheFrozenDev on Discord, in the ZAN or VATUSA servers) and ask. I will not provide detailed step-by-step for all environments
but can answer more general questions and provide guidance.

## License

This project is licensed under the Apache 2.0 license. See the [LICENSE](LICENSE) file for more information. Some packages may have alternate licenses that apply to only code in that directory and inward, which will be noted in the LICENSE file in that directory.
