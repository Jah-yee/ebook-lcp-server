import http from "k6/http";
import { check } from "k6";

export const options = {
  vus: 100,
  duration: "1m",
  thresholds: {
    http_req_duration: ["p(95)<200"],
    http_req_failed: ["rate<0.01"]
  }
};

const baseURL = __ENV.BASE_URL || "http://localhost:8080";
const token = __ENV.JWT || "";

export default function () {
  const response = http.get(`${baseURL}/api/v1/lcp/status`, {
    headers: {
      Authorization: `Bearer ${token}`
    }
  });

  check(response, {
    "status is 200": (res) => res.status === 200
  });
}
