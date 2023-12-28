<template>
  <div v-if="chartLength > 0" class="relative flex flex-col min-w-0 break-words w-full mb-6 shadow-lg rounded">
    <div class="px-4 py-5 flex-auto">
      <div class="flex w-full justify-left mb-2">
        <div class="flex items-center">
          <label class="block text-gray-100 font-bold pr-4">Filter:</label>
          <input
            v-model="filter"
            type="text"
            class="w-[10em] px-4 py-2 text-gray-100 bg-gray-600 rounded focus:bg-gray-600 focus:outline-none uppercase"
            ref="filterInput"
            autofocus
          />
        </div>
      </div>
      <div class="grid grid-cols-8 gap-4">
        <div v-for="apt in filteredCharts" :key="apt">
          <button
            class="border-2 border-gray-500 bg-gray-800 hover:bg-gray-700 text-white font-bold py-2 px-4 rounded h-full w-full mr-1"
            @click.stop="$router.push({ name: 'AirportCharts', params: { airport: apt } })"
          >
            {{ apt }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, nextTick, onMounted, onUpdated, ref } from "vue";
import { useViewStore } from "../store/viewstore";

const store = useViewStore();
const filter = ref("");
const filterInput = ref(null);

onMounted(() => {
  nextTick(() => {
    filterInput.value.focus();
  });
});

const chartLength = computed(() => {
  return Object.keys(store.charts).length;
});

const filteredCharts = computed(() => {
  return Object.keys(store.charts).filter((apt) => {
    return filter != "" && apt.toLowerCase().includes(filter.value.toLowerCase());
  });
});
</script>

<style lang="scss" scoped></style>
