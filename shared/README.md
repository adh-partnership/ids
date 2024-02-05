# Shared Configs

These are JSON configurations that are shared/used by both the frontend and backend. This file is not yet used
by the frontend but may in the future.

Currently this file is only run by the backend's update command. Non-ZAN subdivisions do not need to commit their 
airports.json file.

## airports.json

This file contains an array of strings that represent all the airports covered by this IDS. It does not
dictate views, TRACONs, or ARTCCs, but just straight airports that the IDS will attempt to fetch weather for.
This file should use FAA IDs as the backend update command will use this to populate the database utilizing NASR
data. This data will populate the ICAO IDs in the database to allow either to be used.
