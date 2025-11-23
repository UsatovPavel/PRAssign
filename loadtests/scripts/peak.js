import http from "k6/http";
import { check } from "k6";

export const options = {
  vus: 100,
  duration: "10s"
};

export default function () {
  const r = http.get("http://localhost:8080/pullRequest/list");
  check(r, { ok: (res) => res.status === 200 });
}
