import axios from "axios";
import config from "@/../config.json";

export const API = axios.create({
  baseURL: config.ids_api_base_url,
  withCredentials: true,
});
