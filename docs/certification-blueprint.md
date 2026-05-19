# LCP Certification Blueprint

This repository cannot self-certify an installation for EDRLab, but it can keep the evidence trail tidy before an official test run.

## Evidence to collect

1. Build SHA and deployment configuration snapshot
2. Public provider URI and certificate chain used for signing
3. Encrypted EPUB, PDF, and manifest samples
4. License create, download, status, extension, and revocation traces
5. Reader validation notes
6. Official EDRLab test output once run against production certificates

## Local report

```bash
sh scripts/demo-local.sh
sh scripts/certification-smoke.sh
```

The smoke script writes `certification-report.json` with machine-readable readiness checks.
