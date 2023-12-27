<template>
  <iframe :src="url" width="100%" height="100%" frameborder="0" class="flex-1 mb-[40px]"></iframe>
</template>

<script setup>
import { computed } from "vue";
import config from "@/../config.json";
import { useViewStore } from "../store/viewstore";

const store = useViewStore();

const url = computed(() => {
  return config.views[store.view].sops;
});
</script>

<style lang="scss" scoped></style>
