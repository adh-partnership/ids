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
      <div class="py-10 w-full flex items-center justify-center bg-red-950 bg-opacity-50" v-if="loggedIn && !connected">
        <div class="flex flex-col items-center justify-center">
          <div class="text-4xl font-bold text-red-500">DISCONNECTED FROM SIGNALR HUB</div>
          <div class="text-2xl font-bold text-red-500">Reconnecting...</div>
        </div>
      </div>
    </main>
    <footer class="fixed z-50 bg-neutral-700 text-white bottom-0 p-0 w-full">
      <div v-if="loggedIn">
        <button
          class="border-2 border-gray-500 bg-slate-800 hover:bg-gray-700 text-white font-bold py-2 px-4 rounded h-full w-[10rem] mr-1"
          @click="router.push('/')"
        >
          SIA
        </button>
        <button
          class="border-2 border-gray-500 bg-blue-800 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded h-full w-[10rem] mr-1"
          @click="router.push('/weather')"
        >
          WX
        </button>
        <button
          class="border-2 border-gray-500 bg-green-800 hover:bg-green-700 text-white font-bold py-2 px-4 rounded h-full w-[10rem] mr-1"
          @click="router.push('/sops')"
        >
          SOP
        </button>
        <button
          class="border-2 border-gray-500 bg-amber-600 hover:bg-amber-500 text-white font-bold py-2 px-4 rounded h-full w-[10rem] mr-1"
          @click="router.push('/pireps')"
        >
          PIREPs
        </button>
        <button
          class="border-2 border-gray-500 bg-purple-800 hover:bg-purple-700 text-white font-bold py-2 px-4 rounded h-full w-[10rem] mr-1"
          @click="router.push('/charts')"
        >
          CHARTS
        </button>
        <button
          class="border-2 border-gray-500 bg-rose-900 hover:bg-rose-800 text-white font-bold py-2 px-4 rounded h-full w-[10rem] mr-1"
          @click="router.push('/briefing')"
        >
          BRIEF
        </button>
      </div>
      <div v-else>
        <a :href="`${config.ids_api_base_url}/v1/oauth/login?redirect=${location}`">
          <button
            class="border-2 border-gray-500 bg-slate-800 hover:bg-gray-700 text-white font-bold py-2 px-4 rounded h-full w-[10rem] mr-1"
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
import { useRouter } from "vue-router";
import { storeToRefs } from "pinia";
import { useViewStore } from "@/store/viewstore";
import Clock from "./components/Clock.vue";
import config from "../config.json";
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
  config.views.forEach((v) => {
    v.facilities.forEach(async (f) => {
      if (store.sia[f] === undefined) {
        await store.addAirport(f);
      }
    });
  });

  // We now update via SignalR... so we don't need this?
  //store.updateMetars();

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
