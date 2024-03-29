<template>
  <div class="flex flex-col min-h-screen">
    <header class="fixed bg-neutral-700 text-white top-0 p-4 w-full shadow-md mb-[10px] h-[83px]">
      <div class="py-1">
        <div class="grid grid-cols-4 items-center justify-between">
          <div class="col-span-1 title">
            <span class="facility">{{ config.subdivision.id }}</span> <span class="idscolor">IDS</span>
          </div>
          <div class="col-span-2">
            <div v-if="loggedIn" class="flex items-center">
              <label class="block text-gray-100 font-bold pr-4"> Views: </label>
              <select
                id="view"
                v-model="view"
                class="block w-full bg-neutral-800 text-white py-3 px-4 pr-8 rounded leading-tight focus:outline-none focus:border-gray-500"
                @change="changeView()"
              >
                <option v-for="(v, i) in config.views" :key="i" :value="i">{{ v.name }}</option>
              </select>
            </div>
          </div>
          <div class="col-span-1 text-right ml-auto">
            <Clock :timezone="config.timezone.name" />
          </div>
        </div>
      </div>
    </header>
    <main class="pt-[95px] flex-1 overflow-y-auto p-5 pb-[44px] flex flex-col">
      <router-view v-if="loggedIn && connected" />
      <div v-if="loggedIn && !connected" class="py-10 w-full flex items-center justify-center bg-red-950 bg-opacity-50">
        <div class="flex flex-col items-center justify-center">
          <div class="text-4xl font-bold text-red-500">DISCONNECTED FROM SIGNALR HUB</div>
          <div class="text-2xl font-bold text-red-500">Reconnecting...</div>
        </div>
      </div>
    </main>
    <footer class="fixed z-50 bg-neutral-700 text-white bottom-0 p-0 w-full">
      <div v-if="loggedIn">
        <button
          class="border-2 border-gray-500 font-bold py-2 px-4 rounded h-full w-[10rem] mr-1"
          :class="`${config.colors.buttons.SIA.background} hover:${config.colors.buttons.SIA.hover} ${config.colors.buttons.SIA.foreground}`"
          @click="router.push('/')"
        >
          SIA
        </button>
        <button
          class="border-2 border-gray-500 font-bold py-2 px-4 rounded h-full w-[10rem] mr-1"
          :class="`${config.colors.buttons.WX.background} hover:${config.colors.buttons.WX.hover} ${config.colors.buttons.WX.foreground}`"
          @click="router.push('/weather')"
        >
          WX
        </button>
        <button
          class="border-2 border-gray-500 font-bold py-2 px-4 rounded h-full w-[10rem] mr-1"
          :class="`${config.colors.buttons.SOP.background} hover:${config.colors.buttons.SOP.hover} ${config.colors.buttons.SOP.foreground}`"
          @click="router.push('/sops')"
        >
          SOP
        </button>
        <button
          class="border-2 border-gray-500 font-bold py-2 px-4 rounded h-full w-[10rem] mr-1"
          :class="`${config.colors.buttons.PIREPS.background} hover:${config.colors.buttons.PIREPS.hover} ${config.colors.buttons.PIREPS.foreground}`"
          @click="router.push('/pireps')"
        >
          PIREPs
        </button>
        <button
          class="border-2 border-gray-500 font-bold py-2 px-4 rounded h-full w-[10rem] mr-1"
          :class="`${config.colors.buttons.CHARTS.background} hover:${config.colors.buttons.CHARTS.hover} ${config.colors.buttons.CHARTS.foreground}`"
          @click="router.push('/charts')"
        >
          CHARTS
        </button>
        <button
          class="border-2 border-gray-500 font-bold py-2 px-4 rounded h-full w-[10rem] mr-1"
          :class="`${config.colors.buttons.BRIEF.background} hover:${config.colors.buttons.BRIEF.hover} ${config.colors.buttons.BRIEF.foreground}`"
          @click="router.push('/briefing')"
        >
          BRIEF
        </button>
      </div>
      <div v-else>
        <a :href="`${config.ids_api_base_url}/v1/auth/login?redirect=${location}`">
          <button
            class="border-2 border-gray-500 font-bold py-2 px-4 rounded h-full w-[10rem] mr-1"
            :class="`${config.colors.buttons.Login.background} hover:${config.colors.buttons.Login.hover} ${config.colors.buttons.Login.foreground}`"
          >
            Login
          </button>
        </a>
      </div>
    </footer>
  </div>
</template>

<script setup>
import { onMounted, ref, watch } from "vue";
import { storeToRefs } from "pinia";
import { useRouter } from "vue-router";
import { useViewStore } from "@/store/viewstore.js";
import config from "../config.json";
import Clock from "./components/Clock.vue";
import signalrConnection from "@/main.js";

const store = useViewStore();
const router = useRouter();
const { view, loggedIn } = storeToRefs(store);
const location = ref(window.location.href);
const connected = ref(false);

const changeView = () => {
  store.view = view.value;
};

store.getAuthed();

onMounted(() => {
  // Create a timer to cleanup PIREPs every minute, this will drop any PIREPs older than one hour
  setInterval(() => {
    store.cleanupPIREPs();
  }, 60000);

  config.views.forEach((v) => {
    v.facilities.forEach(async (f) => {
      if (store.sia[f] === undefined) {
        await store.addAirport(f);
      }
    });
  });

  // We now update via SignalR... so we don't need this?
  // store.updateMetars();

  // We aren't connected... yet... but will when logged in is verified.
  // So this *should* work here, if not move to the watch below
  signalrConnection.on("airports", store.signalRAirports);
  signalrConnection.on("charts", store.signalRCharts);
  signalrConnection.on("airportUpdate", store.signalRAirportUpdate);
  signalrConnection.on("pirepUpdate", store.signalRPIREPUpdate);
  signalrConnection.onclose(() => {
    console.log("SignalR connection closed");
    connected.value = false;
  });
  signalrConnection.onreconnecting(() => {
    console.log("SignalR connection reconnecting");
    connected.value = false;
  });
  signalrConnection.onreconnected(() => {
    console.log("SignalR connection reconnected");
    connected.value = true;
  });
  signalrConnection.on("connected", () => {
    console.log("SignalR connection connected");
    connected.value = true;
  });
});

watch(
  () => store.loggedIn,
  async () => {
    console.log(`Logged in? ${store.loggedIn}`);
    if (store.loggedIn) {
      await signalrConnection.start();
      console.log("SignalR connection started");
      connected.value = true;
    }
  }
);
</script>

<style scoped>
.title {
  font-size: 1.5rem;
  font-weight: 700;
  color: #ffb612;
}
.facility {
  color: v-bind("config.colors.navbar.facility");
}
.idscolor {
  color: v-bind("config.colors.navbar.ids");
}
</style>
