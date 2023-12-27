import { defineStore } from "pinia";
import { API } from "@/utils/api";
import config from "@/../config.json";

export const useViewStore = defineStore("view", {
  state: () => ({
    view: config.views.findIndex((v) => v.default),
    sia: {},
    metars: {},
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
          first: false,
        }
      });
    },
    signalRAirportUpdate(oldAirport, newAirport) {
      console.log(`signalRAirportUpdate: oldAirport:`, oldAirport, `newAirport:`, newAirport)
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
      /*      if (this.timers[airport] === undefined) {
              await this.updateSIA(airport);
              this.sia[airport].first = true;
              setTimeout(() => {
                this.sia[airport].first = false;
              }, 1000);
              this.timers[airport] = setInterval(() => {
                this.updateSIA(airport);
              }, 15000);
            }*/
    },
    async updateSIA(airport) {
      /*      if (this.fetching.includes(airport)) return;
            this.fetching.push(airport);
            try {
              const response = await API.get(`/v1/sia/${airport}`);
              this.sia[airport] = response.data;
            } catch (error) {
              console.error(error);
              // Try again in 15 seconds
              setTimeout(() => {
                this.updateSIA(airport);
              }, 15000);
            } finally {
              this.fetching.splice(this.fetching.indexOf(airport), 1);
            }*/
    },
    async updateMetars() {
      /*      try {
              const response = await API.get("/v1/weather/metar/all");
              Object.keys(response.data).forEach((key) => {
                if (response.data[key] !== "") {
                  this.metars[key] = response.data[key];
                }
              });
            } catch (error) {
              console.error(error);
            }
            if (this.metarTimer === null) {
              this.metarTimer = setInterval(() => {
                this.updateMetars();
              }, 60000);
            }*/
    },
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
        console.log(`got authed`);
        this.loggedIn = true;
        console.log(`set loggedIn to true`);
        return true;
      } catch (err) {
        console.error(err);
        console.log(`got error`);
        this.loggedIn = false;
        return false;
      }
    },
  },
});
