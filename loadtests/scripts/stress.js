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
  stages: [
    { duration: "20s", target: 10 },
    { duration: "20s", target: 30 },
    { duration: "20s", target: 60 },
    { duration: "20s", target: 0 }
  ]
};

export default function () {
  const token = getToken("u1");
  const r = http.get(`${BASE}/users/getReview?user_id=u1`, {
    headers: { token: token }
  });
  check(r, { ok: (res) => res.status === 200 || res.status === 404 });
}
