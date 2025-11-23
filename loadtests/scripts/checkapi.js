import http from "k6/http";
import { check, sleep } from "k6";

export const options = {
  vus: 1,
  iterations: 1
};

export default function () {
    const BASE_URL = "http://app:8080";
    const r1 = http.get(`${BASE_URL}/team/get`);
  check(r1, { t: (res) => res.status === 200 });

  const r2 = http.get("http://localhost:8080/users/getReview?user_id=u1");
  check(r2, { u: (res) => res.status === 200 });

  sleep(1);
}
