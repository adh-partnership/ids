<template>
  <div class="w-full">
    <ul class="flex mb-0 list-none flex-wrap pt-3 pb-4 flex-row">
      <li class="-mb-px mr-2 last:mr-0 flex-auto text-center">
        <a
          class="cursor-pointer text-xs font-bold uppercase px-5 py-3 shadow-lg rounded block leading-normal"
          :class="{
            'text-white bg-gray-900': openTab !== 1,
            'text-white bg-blue-900': openTab === 1,
          }"
          @click="toggleTabs(1)"
        >
          METAR
        </a>
      </li>
      <li class="-mb-px mr-2 last:mr-0 flex-auto text-center">
        <a
          class="cursor-pointer text-xs font-bold uppercase px-5 py-3 shadow-lg rounded block leading-normal"
          :class="{
            'text-white bg-gray-900': openTab !== 2,
            'text-white bg-blue-900': openTab === 2,
          }"
          @click="toggleTabs(2)"
        >
          PRE-DUTY BRIEF
        </a>
      </li>
      <li class="-mb-px mr-2 last:mr-0 flex-auto text-center">
        <a
          class="cursor-pointer text-xs font-bold uppercase px-5 py-3 shadow-lg rounded block leading-normal"
          :class="{
            'text-white bg-gray-900': openTab !== 3,
            'text-white bg-blue-900': openTab === 3,
          }"
          @click="toggleTabs(3)"
        >
          SATELITE
        </a>
      </li>
      <li class="-mb-px mr-2 last:mr-0 flex-auto text-center">
        <a
          class="cursor-pointer text-xs font-bold uppercase px-5 py-3 shadow-lg rounded block leading-normal"
          :class="{
            'text-white bg-gray-900': openTab !== 4,
            'text-white bg-blue-900': openTab === 4,
          }"
          @click="toggleTabs(4)"
        >
          REQ PIREP
        </a>
      </li>
      <li class="-mb-px mr-2 last:mr-0 flex-auto text-center">
        <a
          class="cursor-pointer text-xs font-bold uppercase px-5 py-3 shadow-lg rounded block leading-normal"
          :class="{
            'text-white bg-gray-900': openTab !== 5,
            'text-white bg-blue-900': openTab === 5,
          }"
          @click="toggleTabs(5)"
        >
          ICING+CONV
        </a>
      </li>
      <li class="-mb-px mr-2 last:mr-0 flex-auto text-center">
        <a
          class="cursor-pointer text-xs font-bold uppercase px-5 py-3 shadow-lg rounded block leading-normal"
          :class="{
            'text-white bg-gray-900': openTab !== 6,
            'text-white bg-blue-900': openTab === 6,
          }"
          @click="toggleTabs(6)"
        >
          PROG + SIG WX CHART
        </a>
      </li>
    </ul>
    <div class="relative flex flex-col min-w-0 break-words w-full mb-6 shadow-lg rounded">
      <div class="px-4 py-5 flex-auto">
        <div class="tab-content tab-space">
          <div :class="{ hidden: openTab !== 1, block: openTab === 1 }">
            <div class="flex w-full justify-left mb-2">
              <div class="flex items-center">
                <label class="block text-gray-100 font-bold pr-4"> Filter: </label>
                <input
                  v-model="filter"
                  type="text"
                  class="w-[10em] px-4 py-2 text-gray-100 bg-gray-600 rounded focus:bg-gray-600 focus:outline-none uppercase"
                />
              </div>
            </div>
            <table class="table-fixed w-full text-white border-collapse">
              <MetarRow
                v-for="metar in metarFiltered()"
                :key="metar"
                :mag_var="sia[metar].mag_var"
                :airport="metar"
                :metar="sia[metar].metar"
              ></MetarRow>
            </table>
          </div>
          <div :class="{ hidden: openTab !== 2, block: openTab === 2 }" class="align-center">
            <center>
              <video controls="" width="956" height="717">
                <source :src="config.weather['PRE-DUTY BRIEF'].videoUrl" type="video/mp4" />
              </video>
            </center>
          </div>
          <div :class="{ hidden: openTab !== 3, block: openTab === 3 }" class="text-white grid grid-cols-3 gap-4">
            <div v-for="(img, i) in config.weather['SATELITE'].images" :key="i">
              <center>
                <img :src="img.src" :alt="img.title" />
                <p class="text-xs">{{ img.title }}</p>
              </center>
            </div>
          </div>
          <div :class="{ hidden: openTab !== 4, block: openTab === 4 }">
            <div v-for="(img, i) in config.weather['REQ PIREP'].images" :key="i">
              <center>
                <img :src="img.src" :alt="img.title" />
                <p class="text-xs">{{ img.title }}</p>
              </center>
            </div>
          </div>
          <div :class="{ hidden: openTab !== 5, block: openTab === 5 }" class="grid grid-cols-4 gap-4">
            <div v-for="(img, i) in config.weather['ICING+CONV'].images" :key="i">
              <center>
                <img :src="img.src" alt="icing graphic" />
              </center>
            </div>
          </div>
          <div :class="{ hidden: openTab !== 6, block: openTab === 6 }" class="grid grid-cols-2 gap-4">
            <div v-for="(img, i) in config.weather['PROG+SIG WX CHART'].images" :key="i">
              <center>
                <img :src="img.src" alt="prog+sig wx chart" />
              </center>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, ref } from "vue";
import { storeToRefs } from "pinia";
import { useViewStore } from "../store/viewstore";
import config from "@/../config.json";
import MetarRow from "@/components/MetarRow.vue";

const store = useViewStore();
const { sia } = storeToRefs(store);
const openTab = ref(1);
const filter = ref("");

const metarFiltered = () => {
  // Return keys of metars filtered by filter
  return Object.keys(sia.value)
    .filter((key) => {
      return key.toUpperCase().includes(filter.value.toUpperCase()) || sia.value[key].metar.includes(filter.value);
    })
    .sort();
};

const toggleTabs = (tab) => {
  openTab.value = tab;
};
</script>

<style lang="scss" scoped></style>
