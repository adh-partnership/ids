import { createRouter, createWebHistory } from "vue-router";
import SIA from "@/views/infoarea.vue";

const routes = [
  {
    path: "/",
    name: "Home",
    component: SIA,
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
  {
    path: "/briefing",
    name: "Briefing",
    component: () => import("@/views/briefing.vue"),
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;
