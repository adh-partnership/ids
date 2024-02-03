<template>
  <div v-if="chartLength > 0" class="relative flex flex-col min-w-0 break-words w-full mb-6 shadow-lg rounded">
    <div class="px-4 py-5 flex-auto">
      <div class="flex w-full justify-left mb-2">
        <div class="flex items-center">
          <button
            class="border-2 border-gray-500 bg-purple-800 hover:bg-purple-700 text-white font-bold py-2 px-4 rounded h-full w-full"
            @click="$router.push({ name: 'Charts' })"
          >
            <i class="fa-solid fa-backward text-white"></i>
            Back
          </button>
        </div>
      </div>
      <div>
        <h2 class="text-2xl font-bold text-white">{{ airport }}</h2>
        <div
          v-if="charts.filter((chart) => chart.chart_code === 'STAR').length > 0"
          class="grid grid-cols-3 gap-x-8 gap-y-4 py-2"
        >
          <div v-for="(chart, index) of charts.filter((c) => c.chart_code === 'STAR')" :key="index">
            <button
              class="border-2 border-gray-500 bg-blue-800 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded h-full w-full mr-1"
              @click.stop="$router.push({ name: 'AirportChart', params: { airport, chart: encode(chart.chart_name) } })"
            >
              <i class="fa-solid fa-plane"></i> {{ chart.chart_name }}
            </button>
          </div>
        </div>
        <div
          v-if="charts.filter((chart) => chart.chart_code === 'DP').length > 0"
          class="grid grid-cols-3 gap-x-8 gap-y-4 py-2"
        >
          <div v-for="(chart, index) of charts.filter((c) => c.chart_code === 'DP')" :key="index">
            <button
              class="border-2 border-gray-500 bg-green-900 hover:bg-green-800 text-white font-bold py-2 px-4 rounded h-full w-full mr-1"
              @click.stop="$router.push({ name: 'AirportChart', params: { airport, chart: encode(chart.chart_name) } })"
            >
              <i class="fa-solid fa-plane-departure"></i> {{ chart.chart_name }}
            </button>
          </div>
        </div>
        <div
          v-if="charts.filter((chart) => chart.chart_code === 'IAP').length > 0"
          class="grid grid-cols-3 gap-x-8 gap-y-4 py-2"
        >
          <div v-for="(chart, index) of charts.filter((c) => c.chart_code === 'IAP')" :key="index">
            <button
              class="border-2 border-gray-500 bg-teal-800 hover:bg-teal-700 text-white font-bold py-2 px-4 rounded h-full w-full mr-1"
              @click.stop="$router.push({ name: 'AirportChart', params: { airport, chart: encode(chart.chart_name) } })"
            >
              <i class="fa-solid fa-plane-arrival"></i> {{ chart.chart_name }}
            </button>
          </div>
        </div>
        <div
          v-if="charts.filter((chart) => chart.chart_code === 'OTHER').length > 0"
          class="grid grid-cols-3 gap-x-8 gap-y-4 py-2"
        >
          <div v-for="(chart, index) of charts.filter((c) => c.chart_code === 'OTHER')" :key="index">
            <button
              class="border-2 border-gray-500 bg-gray-800 hover:bg-gray-700 text-white font-bold py-2 px-4 rounded h-full w-full mr-1"
              @click.stop="$router.push({ name: 'AirportChart', params: { airport, chart: encode(chart.chart_name) } })"
            >
              <i class="fa-solid fa-paper-plane"></i> {{ chart.chart_name }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, ref } from "vue";
import { useRoute } from "vue-router";
import { Base64 } from "js-base64";
import { useViewStore } from "../store/viewstore";

const store = useViewStore();
const route = useRoute();
const airport = ref(route.params.airport);

const chartLength = computed(() => {
  return Object.keys(store.charts).length;
});

const charts = computed(() => {
  return store.charts[airport.value].sort((a, b) => {
    // sort by chart_code, then chart_name
    if (a.chart_code < b.chart_code) return -1;
    if (a.chart_code > b.chart_code) return 1;
    if (a.chart_name < b.chart_name) return -1;
    if (a.chart_name > b.chart_name) return 1;
    return 0;
  });
});

const encode = (string) => {
  return Base64.encode(string, true);
};
</script>

<style lang="scss" scoped></style>
