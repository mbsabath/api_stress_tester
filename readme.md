# Stress Tester Utility

Command line utility to generate traffic against an endpoint at a predefined rate.
Also records the response time of each request and response status code.

Intended for use stress testing specific API endpoints.

# Usage

```bash
stress-tester -url <url> -rate <rate> -duration <duration> -method <method> -body <request body> -results <result file path>
```

