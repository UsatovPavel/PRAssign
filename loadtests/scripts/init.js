import http from "k6/http";
import { check } from "k6";

export const options = {
  vus: 1,
  iterations: 1
};

import teams from "../data/teams.js";
import prRequests from "../data/pr-requests.js";

export default function () {
  teams.forEach(t => {
    const r = http.post("http://localhost:8080/team/add", JSON.stringify(t), {
      headers: { "Content-Type": "application/json" }
    });
    check(r, { ok: (res) => res.status === 200 });
  });

  prRequests.forEach(p => {
    const r = http.post("http://localhost:8080/pullRequest/create", JSON.stringify(p), {
      headers: { "Content-Type": "application/json" }
    });
    check(r, { ok: (res) => res.status === 200 || res.status === 409 });
  });
}
