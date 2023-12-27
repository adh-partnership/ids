<template>
  <span class="clock" @click="showUTC = !showUTC">
    {{ time }}
  </span>
</template>

<script setup>
import { onMounted, onUnmounted, ref } from "vue";
import { useViewStore } from "@/store/viewstore";
import config from "@/../config.json";

const store = useViewStore();
const showUTC = ref(false);

const props = defineProps({
  timezone: {
    type: String,
    default: "UTC",
  },
});

const time = ref("");
let timer;

const updateClock = () => {
  let letter = "L";
  let timezone = props.timezone;
  if (props.timezone === "UTC" || showUTC.value) {
    letter = "Z";
    timezone = "UTC";
  }
  time.value = `${new Date()
    .toLocaleTimeString("en-UK", {
      timeZone: timezone,
      hour12: false,
    })
    .slice(-8)} ${letter}`;

  // Check if seconds are 0, if so, update store.heartbeat.minute...
  // the value doesn't matter.. we just need to change
  if (time.value.slice(-4) === "00 L") {
    store.heartbeat.minute = time.value;
  }
  store.heartbeat.second = time.value;
};

onMounted(() => {
  timer = setInterval(updateClock, 1000);
  updateClock();
});

onUnmounted(() => {
  clearInterval(timer);
});
</script>

<style scoped>
.clock {
  line-height: 1.5;
  padding: 0.5rem;
  font-size: 1.5rem;
  color: v-bind("config.colors.navbar.clock_foreground");
  background-color: v-bind("config.colors.navbar.clock_background");
  border-radius: 15px;
}
</style>
