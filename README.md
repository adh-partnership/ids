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

Soon (TM).

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

## License

This project is licensed under the Apache 2.0 license. See the [LICENSE](LICENSE) file for more information. Some packages may have alternate licenses that apply to only code in that directory and inward, which will be noted in the LICENSE file in that directory.
