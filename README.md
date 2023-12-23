# ADH Partnership IDS

This is still very much a work in progress. This README will be expanded as development progresses.

## Introduction

This is a rewrite of the [ZAN IDS](https://github.com/vpaza/ids) to be donated to the Partnership with some caching and efficiencies written in, and designed to be configurable and usable by multiple facilities. This is still very much a work in progress.

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

## License

This project is licensed under the Apache 2.0 license. See the [LICENSE](LICENSE) file for more information.