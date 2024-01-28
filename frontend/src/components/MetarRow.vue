<template>
  <tr>
    <td class="w-[2rem] border border-gray-500 text-center font-bold">
      <div class="w-4 h-4 rounded-full ml-auto mr-auto" :class="flightCategory()"></div>
    </td>
    <td class="w-[4rem] border border-gray-500 text-center font-bold">
      {{ props.airport }}
    </td>
    <td class="w-[10em] border border-gray-500 text-center text-yellow-400 font-bold">
      {{ wind }}
    </td>
    <td class="border border-gray-500">
      {{ props.metar || "No METAR Available" }}
    </td>
  </tr>
</template>

<script setup>
import { computed, ref } from "vue";
import { calcWindDir, parseMetar } from "@/utils/metar.js";

const props = defineProps({
  airport: {
    type: String,
    required: true,
  },
  metar: {
    type: String,
    required: true,
    default: 0,
  },
  mag_var: {
    type: Number,
    required: true,
    default: 0,
  },
});

const flightCategory = () => {
  if (props.metar === undefined || props.metar === "") {
    return "bg-gray-900";
  }

  const m = parseMetar(props.metar);
  if (m.flight_category === "LIFR") {
    return "bg-purple-500";
  }
  if (m.flight_category === "IFR") {
    return "bg-red-500";
  }
  if (m.flight_category === "MVFR") {
    return "bg-blue-500";
  }
  if (m.flight_category === "VFR") {
    return "bg-green-500";
  }
  return "bg-gray-900";
};

const wind = computed(() => {
  if (props.metar === undefined || props.metar === "") {
    return "??";
  }

  const m = parseMetar(props.metar);
  if (m.wind === undefined) {
    return "??";
  }

  if (m.wind.speed_kts < 3) {
    return "CALM";
  }

  let wind = `${calcWindDir(m.wind.degrees, props.mag_var).toString().padStart(3, "0")} @ ${m.wind.speed_kts}`;
  if (m.wind.gust_kts > m.wind.speed_kts + 6) {
    wind += `G${m.wind.gust_kts}`;
  }

  return wind;
});
</script>

<style lang="scss" scoped></style>
