import { createRouter, createWebHistory } from "vue-router";
import SIA from "@/views/infoarea.vue";

const routes = [
  {
    path: "/",
    name: "Home",
    component: SIA,
  },

  {
    path: "/briefing",
    name: "Briefing",
    component: () => import("@/views/briefing.vue"),
  },
  {
    path: "/charts",
    name: "Charts",
    component: () => import("@/views/charts.vue"),
  },
  {
    path: "/charts/:airport",
    name: "AirportCharts",
    component: () => import("@/views/charts.airport.vue"),
  },
  {
    path: "/charts/:airport/:chart",
    name: "AirportChart",
    component: () => import("@/views/charts.chart.vue"),
  },
  {
    path: "/sops",
    name: "SOPs",
    component: () => import("@/views/sops.vue"),
  },
  {
    path: "/weather",
    name: "Weather",
    component: () => import("@/views/weather.vue"),
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;
