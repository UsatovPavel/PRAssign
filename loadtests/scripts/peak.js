import http from "k6/http";
import { check } from "k6";
import { BASE } from "./config/constants.js";

function getToken(username) {
  const res = http.post(`${BASE}/auth/token`,
    JSON.stringify({ username: username }),
    { headers: { "Content-Type": "application/json" } });
  try {
    return res.json().token || "";
  } catch (e) {
    return "";
  }
}

export const options = {
  vus: 100,
  duration: "10s"
};

export default function () {
  const token = getToken("integration-test");
  const r = http.get(`${BASE}/team/get?team_name=backend`, {
    headers: { token: token }
  });
  check(r, { ok: (res) => res.status === 200 || res.status === 404 });
}
