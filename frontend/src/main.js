import { createApp } from "vue";
import { createPinia } from "pinia";
import App from "./App.vue";
import router from "./router";
import config from "@/../config.json";

import "./assets/css/style.css";
import { VueSignalR } from "./signalr";
import { HttpTransportType, HubConnectionBuilder } from "@microsoft/signalr";

const pinia = createPinia();

const app = createApp(App);
app.config.devtools = true;

const signalrConnection = new HubConnectionBuilder()
  .withUrl(`${config.ids_api_base_url}/signalr`, {
    skipNegotiation: true,
    transport: HttpTransportType.WebSockets,
  })
  .build();

export default signalrConnection;

app
  .use(router)
  .use(pinia)
  .use(VueSignalR, { connection: signalrConnection })
  .mount("#app");

