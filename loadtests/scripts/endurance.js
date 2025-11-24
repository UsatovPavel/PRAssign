import http from "k6/http";
import { check, sleep } from "k6";
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
  vus: 5,
  duration: "10m"
};

export default function () {
  const token = getToken("u2");
  const r = http.get(`${BASE}/users/getReview?user_id=u2`, {
    headers: { token: token }
  });
  check(r, { ok: (res) => res.status === 200 || res.status === 404 });
  sleep(1);
}
