import http from "k6/http";
import { check, sleep } from "k6";

export const options = {
  vus: 5,
  duration: "10m"
};

export default function () {
  const r = http.get("http://localhost:8080/users/getReview?user_id=u2");
  check(r, { ok: (res) => res.status === 200 });
  sleep(1);
}
