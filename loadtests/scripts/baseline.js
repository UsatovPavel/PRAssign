import http from "k6/http";
import { check, sleep } from "k6";
import { BASE } from "./config/constants.js";

export const options = {
  vus: 5,
  duration: "10s"
};

export default function () {
  const r = http.get(`${BASE}/health`);
  check(r, { "healthy": (res) => res.status === 200 });
  sleep(1);
}
