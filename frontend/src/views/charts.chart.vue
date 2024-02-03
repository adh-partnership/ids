<template>
  <div v-if="chartLength > 0" class="relative flex flex-col min-w-0 break-words w-full mb-6 shadow-lg rounded">
    <div class="px-4 py-5 flex-auto">
      <div class="flex w-full justify-left mb-2">
        <div class="flex items-center">
          <button
            class="border-2 border-gray-500 bg-purple-800 hover:bg-purple-700 text-white font-bold py-2 px-4 rounded h-full w-full"
            @click="$router.push({ name: 'AirportCharts', params: { airport } })"
          >
            <i class="fa-solid fa-backward text-white"></i>
            Back
          </button>
        </div>
      </div>
      <div>
        <h2 class="text-2xl font-bold text-white">{{ airport }} - {{ chart.chart_name }} ({{ chart.chart_code }})</h2>
        <iframe
          :src="chart.chart_url"
          width="100%"
          height="100%"
          frameborder="0"
          class="flex-1 h-[80vh] mb-[40px]"
        ></iframe>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import { Base64 } from "js-base64";
import { useViewStore } from "../store/viewstore";

const store = useViewStore();
const route = useRoute();
const airport = ref(route.params.airport);
const chartname = ref(Base64.decode(route.params.chart));
const chart = ref({});

const chartLength = computed(() => {
  return Object.keys(store.charts).length;
});

onMounted(() => {
  if (chartname.value != "") {
    if (store.charts[airport.value] === undefined) {
      useRouter().push({ name: "Charts" });
      return;
    }
    const c = store.charts[airport.value].find((c) => c.chart_name === chartname.value);
    if (c) {
      chart.value = c;
    }
  }
});
</script>

<style lang="scss" scoped></style>
