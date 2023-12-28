import { defineStore } from "pinia";
import { API } from "@/utils/api";
import config from "@/../config.json";

export const useViewStore = defineStore("view", {
  state: () => ({
    view: config.views.findIndex((v) => v.default),
    sia: {},
    metars: {},
    charts: {},
    pireps: [],
    fetching: [],
    metarFetching: false,
    timers: {},
    loggedIn: false,
    initialized: false,
    heartbeat: {
      second: 0,
      minute: 0,
    },
    metarTimer: null,
  }),
  actions: {
    // Airports SignalR Invocation... which is only sent on connection
    signalRAirports(airports) {
      airports.forEach((airport) => {
        this.sia[airport.id] = {
          atis: airport.atis,
          atis_time: airport.atis_time,
          arrival_atis: airport.arrival_atis,
          arrival_atis_time: airport.arrival_atis_time,
          departure_runways: airport.departure_runways,
          arrival_runways: airport.arrival_runways,
          metar: airport.metar,
          mag_var: airport.mag_var,
          first: true,
        }
      });
    },
    signalRAirportUpdate(oldAirport, newAirport) {
      this.sia[newAirport.id] = {
        atis: newAirport.atis,
        atis_time: newAirport.atis_time,
        arrival_atis: newAirport.arrival_atis,
        arrival_atis_time: newAirport.arrival_atis_time,
        departure_runways: newAirport.departure_runways,
        arrival_runways: newAirport.arrival_runways,
        metar: newAirport.metar,
        mag_var: newAirport.mag_var,
        first: false,
      }
    },
    signalRCharts(charts) {
      console.log(`signalRCharts:`, charts)
      this.charts = charts;
    },
    signalRPIREPUpdate(pirep) {
      this.pireps.push(pirep);
    },
    async addAirport(airport) {
      if (this.sia[airport] !== undefined) return;
      this.sia[airport] = {
        faa_id: "",
        icao_id: "",
        atis: "",
        arrival_atis: "",
        atis_time: new Date(),
        arrival_atis_time: new Date(),
        departure_runways: "",
        arrival_runways: "",
        metar: "",
      };
    },
    // We could send this via SignalR, but let's do it this way to also refresh the auth cookie
    // we won't update the data on our side, that will be when we get the update via SignalR.
    async patchSIA(airport, field, value) {
      if (field === "ATIS") {
        field = "atis";
      }
      if (field === "arrival_ATIS") {
        field = "arrival_atis";
      }

      try {
        const data = {};
        data[field] = value;
        await API.patch(`/v1/airports/${airport}`, data);
      } catch (err) {
        console.error(err);
      }
    },
    async getAuthed() {
      console.log(`in getAuthed... stack:`, new Error().stack);
      try {
        await API.get("/v1/auth/check");
        this.loggedIn = true;
        return true;
      } catch (err) {
        console.error(err);
        this.loggedIn = false;
        return false;
      }
    },
  },
});
