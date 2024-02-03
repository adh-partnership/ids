<template>
  <div class="grid grid-cols-4 gap-4">
    <div class="col-span-2 lg:col-span-3">
      <h2 class="text-2xl text-white pb-10">PIREPs</h2>
      <ul class="list-none text-white">
        <li v-for="pirep in pireps" :key="pirep.raw">
          {{ pirep.raw }}
        </li>
        <li v-if="pireps.length === 0">No PIREPS</li>
      </ul>
    </div>
    <div class="colspan-2 lg:col-span-1 text-white">
      <form>
        <table>
          <thead>
            <tr>
              <th colspan="2" class="text-2xl">PIREP Form</th>
            </tr>
            <tr>
              <th colspan="2">Pilot Weather Report</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="formError !== ''">
              <td class="text-white bg-red-800 font-bold" colspan="2">Error in form: {{ formError }}</td>
            </tr>
            <tr>
              <td>1.</td>
              <td>
                <input
                  id="type-ua"
                  v-model="newPirep.ty"
                  value="UA"
                  type="radio"
                  class="w-4 h-4 text-blue-600 bg-gray-100 border-gray-300 focus:ring-blue-500 dark:focus:ring-blue-600 dark:ring-offset-gray-700 dark:focus:ring-offset-gray-700 focus:ring-2 dark:bg-gray-600 dark:border-gray-500"
                />
                <label for="type-ua" class="w-full py-3 ms-2 font-medium text-gray-900 dark:text-gray-300 pr-4"
                  >UA</label
                >
                <input
                  id="type-uua"
                  v-model="newPirep.ty"
                  value="UUA"
                  type="radio"
                  class="w-4 h-4 text-blue-600 bg-gray-100 border-gray-300 focus:ring-blue-500 dark:focus:ring-blue-600 dark:ring-offset-gray-700 dark:focus:ring-offset-gray-700 focus:ring-2 dark:bg-gray-600 dark:border-gray-500"
                />
                <label for="type-uua" class="w-full py-3 ms-2 font-medium text-gray-900 dark:text-gray-300">UUA</label>
              </td>
            </tr>
            <tr>
              <td>2. /OV <i class="fa-solid fa-arrow-right"></i></td>
              <td>
                <input
                  id="location"
                  v-model="newPirep.ov"
                  type="text"
                  class="w-full px-4 text-sm py-2 text-gray-100 bg-gray-600 rounded focus:bg-gray-600 focus:outline-none uppercase"
                />
              </td>
            </tr>
            <tr>
              <td>3. /TM <i class="fa-solid fa-arrow-right"></i></td>
              <td>
                <input
                  id="time"
                  v-model="newPirep.tm"
                  type="text"
                  class="w-full px-4 text-sm py-2 text-gray-100 bg-gray-600 rounded focus:bg-gray-600 focus:outline-none uppercase"
                />
              </td>
            </tr>
            <tr>
              <td>4. /FL</td>
              <td>
                <input
                  id="flight-level"
                  v-model="newPirep.fl"
                  type="text"
                  class="w-full px-4 text-sm py-2 text-gray-100 bg-gray-600 rounded focus:bg-gray-600 focus:outline-none uppercase"
                />
              </td>
            </tr>
            <tr>
              <td class="pb-2">5. /TP <i class="fa-solid fa-arrow-right"></i></td>
              <td class="pb-2">
                <input
                  id="aircraft-type"
                  v-model="newPirep.tp"
                  type="text"
                  class="w-full px-4 text-sm py-2 text-gray-100 bg-gray-600 rounded focus:bg-gray-600 focus:outline-none uppercase"
                />
              </td>
            </tr>
            <tr>
              <td colspan="2" class="border-y-2 border-gray-600">Items 1 through 5 are mandatory for all PIREPs.</td>
            </tr>
            <tr>
              <td class="pt-2">6. /SK <i class="fa-solid fa-arrow-right"></i></td>
              <td class="pt-2">
                <input
                  id="skycondition"
                  v-model="newPirep.sk"
                  type="text"
                  class="w-full px-4 text-sm py-2 text-gray-100 bg-gray-600 rounded focus:bg-gray-600 focus:outline-none uppercase"
                />
              </td>
            </tr>
            <tr>
              <td>7. /WX <i class="fa-solid fa-arrow-right"></i></td>
              <td>
                <input
                  id="weather"
                  v-model="newPirep.wx"
                  type="text"
                  class="w-full px-4 text-sm py-2 text-gray-100 bg-gray-600 rounded focus:bg-gray-600 focus:outline-none uppercase"
                />
              </td>
            </tr>
            <tr>
              <td>8. /TA <i class="fa-solid fa-arrow-right"></i></td>
              <td>
                <input
                  id="temperature"
                  v-model="newPirep.ta"
                  type="text"
                  class="w-full px-4 text-sm py-2 text-gray-100 bg-gray-600 rounded focus:bg-gray-600 focus:outline-none uppercase"
                />
              </td>
            </tr>
            <tr>
              <td>9. /WV <i class="fa-solid fa-arrow-right"></i></td>
              <td>
                <input
                  id="wind-velocity"
                  v-model="newPirep.wv"
                  type="text"
                  class="w-full px-4 text-sm py-2 text-gray-100 bg-gray-600 rounded focus:bg-gray-600 focus:outline-none uppercase"
                />
              </td>
            </tr>
            <tr>
              <td>10. /TB <i class="fa-solid fa-arrow-right"></i></td>
              <td>
                <input
                  id="turbulence"
                  v-model="newPirep.tb"
                  type="text"
                  class="w-full px-4 text-sm py-2 text-gray-100 bg-gray-600 rounded focus:bg-gray-600 focus:outline-none uppercase"
                />
              </td>
            </tr>
            <tr>
              <td>11. /IC <i class="fa-solid fa-arrow-right"></i></td>
              <td>
                <input
                  id="icing"
                  v-model="newPirep.ic"
                  type="text"
                  class="w-full px-4 text-sm py-2 text-gray-100 bg-gray-600 rounded focus:bg-gray-600 focus:outline-none uppercase"
                />
              </td>
            </tr>
            <tr>
              <td>12. /RM <i class="fa-solid fa-arrow-right"></i></td>
              <td>
                <input
                  id="remarks"
                  v-model="newPirep.rm"
                  type="text"
                  class="w-full px-4 text-sm py-2 text-gray-100 bg-gray-600 rounded focus:bg-gray-600 focus:outline-none"
                />
              </td>
            </tr>
            <tr>
              <td colspan="2">
                <a
                  href="https://www.faa.gov/documentLibrary/media/Order/7110.10DD_Basic_dtd_10-5-23.pdf"
                  target="_blank"
                  class="text-blue-500 hover:text-blue-700 underline"
                  >Ref: FAAO 7110.10DD 8-1</a
                >
              </td>
            </tr>
            <tr>
              <td colspan="2">
                <button
                  class="w-full px-4 py-2 mt-4 text-sm font-medium text-white uppercase bg-blue-600 rounded-lg shadow hover:bg-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 focus:ring-offset-blue-200"
                  @click.prevent="addPirep"
                >
                  Submit PIREP
                </button>
              </td>
            </tr>
            <tr>
              <td colspan="2">
                <button
                  class="w-full px-4 py-2 mt-4 text-sm font-medium text-white uppercase bg-red-600 rounded-lg shadow hover:bg-red-500 focus:outline-none focus:ring-2 focus:ring-red-500 focus:ring-offset-2 focus:ring-offset-red-200"
                  @click.prevent="resetPirep"
                >
                  Cancel
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </form>
    </div>
  </div>
</template>

<script setup>
import { storeToRefs } from "pinia";
import { ref } from "vue";
import { useViewStore } from "../store/viewstore";

const store = useViewStore();
const { pireps } = storeToRefs(store);
const newPirep = ref({
  ty: "UA",
  ov: "",
  tm: "",
  fl: "",
  tp: "",
  sk: "",
  wx: "",
  ta: "",
  wv: "",
  tb: "",
  ic: "",
  rm: "",
});
const formError = ref("");

function validate() {
  if (newPirep.value.ov === "") {
    document.getElementById("location").focus();
    formError.value = "Location is required";
    return false;
  }
  if (newPirep.value.tm === "") {
    document.getElementById("time").focus();
    formError.value = "Time is required";
    return false;
  }
  if (newPirep.value.tp === "") {
    document.getElementById("aircraft-type").focus();
    formError.value = "Aircraft type is required";
    return false;
  }
  if (newPirep.value.fl === "") {
    document.getElementById("flight-level").focus();
    formError.value = "Flight level is required";
    return false;
  }

  return true;
}

async function addPirep() {
  // Validation
  if (!validate()) {
    return;
  }

  formError.value = "";

  await store.submitPirep(newPirep.value);
  resetPirep();
}

function resetPirep() {
  formError.value = "";
  newPirep.value = {
    ty: "UA",
    ov: "",
    tm: "",
    fl: "",
    tp: "",
    sk: "",
    wx: "",
    ta: "",
    wv: "",
    tb: "",
    ic: "",
    rm: "",
  };
}
</script>

<style scoped></style>
