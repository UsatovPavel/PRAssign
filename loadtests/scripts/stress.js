import http from "k6/http";
import { check } from "k6";

export const options = {
  stages: [
    { duration: "20s", target: 10 },
    { duration: "20s", target: 30 },
    { duration: "20s", target: 60 },
    { duration: "20s", target: 0 }
  ]
};

export default function () {
  const r = http.get("http://localhost:8080/users/getReview?user_id=u1");
  check(r, { ok: (res) => res.status === 200 });
}
