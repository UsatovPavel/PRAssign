import http from "k6/http";
import { sleep } from "k6";

export const options = {
  vus: 1,
  iterations: 1,
  thresholds: {
    'http_req_failed': ['rate<1'], // чтобы k6 не считался провальным, мы сами бросим ошибку при таймауте
  },
};

export default function () {
  const maxAttempts = 60;
  for (let i = 0; i < maxAttempts; i++) {
    let res = http.get("http://localhost:8080/health");
    if (res.status === 200) {
      return;
    }
    sleep(1);
  }
  // если не поднялся за timeout — завершить с ошибкой
  throw new Error("app /health did not become healthy within timeout");
}