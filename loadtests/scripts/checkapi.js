import http from "k6/http";
import { check } from "k6";
import { BASE } from "./config/constants.js";

export const options = {
  vus: 1,
  iterations: 1
};

export default function () {
  const adminTokenRes = http.post(`${BASE}/auth/token`,
    JSON.stringify({ username: "integration-test" }),
    { headers: { "Content-Type": "application/json" } });
  const adminToken = (() => { try { return adminTokenRes.json().token || ""; } catch (e) { return ""; } })();

  const userTokenRes = http.post(`${BASE}/auth/token`,
    JSON.stringify({ username: "u1" }),
    { headers: { "Content-Type": "application/json" } });
  const userToken = (() => { try { return userTokenRes.json().token || ""; } catch (e) { return ""; } })();

  let res = http.get(`${BASE}/team/get?team_name=backend`, {
    headers: { token: adminToken }
  });
  check(res, { "team get OK": (r) => r.status === 200 || r.status === 404 });

  res = http.get(`${BASE}/users/getReview?user_id=u1`, {
    headers: { token: userToken }
  });
  check(res, { "users getReview OK": (r) => r.status === 200 || r.status === 404 });
}
