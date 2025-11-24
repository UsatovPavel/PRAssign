import http from "k6/http";
import teams from "../data/teams.js";
import prRequests from "../data/pr-requests.js";
import { BASE } from "./config/constants.js";

export const options = {
  vus: 1,
  iterations: 1
};

export default function () {
  const adminRes = http.post(`${BASE}/auth/token`,
    JSON.stringify({ username: "integration-test" }),
    { headers: { "Content-Type": "application/json" } });
  const adminToken = (() => { try { return adminRes.json().token || ""; } catch (e) { return ""; } })();

  const headers = { "Content-Type": "application/json", token: adminToken };

  teams.forEach(t => {
    http.post(`${BASE}/team/add`, JSON.stringify(t), { headers });
  });

  prRequests.forEach(p => {
    http.post(`${BASE}/pullRequest/create`, JSON.stringify(p), { headers });
  });
}
