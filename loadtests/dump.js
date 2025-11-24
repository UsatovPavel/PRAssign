//usage: k6 run loadtests/dump.js
// 0% - хорошо. Т.е. хорошо это 
// http_req_failed................: 0.00% 0 out of 1 
// http_req_failed................: 0.00% ✓ 0 ✗ 1
import http from "k6/http";
import { check } from "k6";

const BASE = "http://localhost:8080";

export default function () {
  const payload = JSON.stringify({ username: "admin" });
  const headers = { "Content-Type": "application/json" };

  const res = http.post(`${BASE}/auth/token`, payload, { headers });

  check(res, {
    "status is 200": (r) => r.status === 200,
    "body has token": (r) => r.json("token") !== undefined
  });

  console.log("status:", res.status);
  console.log("body:", res.body);
}
