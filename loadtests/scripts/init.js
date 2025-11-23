import http from "k6/http";

export const options = {
  vus: 1,
  iterations: 1
};

import teams from "../data/teams.js";
import prRequests from "../data/pr-requests.js";

export default function () {
  teams.forEach(t => {
    http.post("http://localhost:8080/team/add", JSON.stringify(t), {
      headers: { "Content-Type": "application/json" }
    });
  });

  prRequests.forEach(p => {
    http.post("http://localhost:8080/pullRequest/create", JSON.stringify(p), {
      headers: { "Content-Type": "application/json" }
    });
  });
}