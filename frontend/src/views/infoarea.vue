<template>
  <table class="table-fixed w-full text-white border-collapse">
    <AirportRow v-for="airport in viewFacility.facilities" :key="airport" :airport="airport" />
  </table>
</template>

<script setup>
import { ref, watch } from "vue";
import { storeToRefs } from "pinia";
import config from "@/../config.json";
import AirportRow from "@/components/AirportRow.vue";
import { useViewStore } from "../store/viewstore";

const store = useViewStore();
const { view } = storeToRefs(store);

watch(view, () => {
  viewFacility.value = config.views[view.value];
});
const viewFacility = ref(config.views[view.value]);
</script>

<style lang="scss" scoped></style>
