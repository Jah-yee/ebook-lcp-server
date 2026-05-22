package observability

import "testing"

func TestMetricsCountersAndSnapshot(t *testing.T) {
	IncWebhookOK()
	IncWebhookFailed()
	IncS3StoreOK()
	IncS3StoreFailed()
	IncS3OpenOK()
	IncS3OpenFailed()
	IncS3SignedURLOK()
	IncS3SignedURLFail()
	IncLicensesOK()
	IncLicensesFailed()
	IncAuthFailed()

	snapshot := Current()
	if snapshot.WebhookOK < 1 || snapshot.WebhookFailed < 1 {
		t.Fatalf("unexpected webhook counters: %+v", snapshot)
	}
	if snapshot.S3StoreOK < 1 || snapshot.S3StoreFailed < 1 {
		t.Fatalf("unexpected s3 store counters: %+v", snapshot)
	}
	if snapshot.S3OpenOK < 1 || snapshot.S3OpenFailed < 1 {
		t.Fatalf("unexpected s3 open counters: %+v", snapshot)
	}
	if snapshot.S3SignedURLOK < 1 || snapshot.S3SignedURLFail < 1 {
		t.Fatalf("unexpected s3 signed url counters: %+v", snapshot)
	}
	if snapshot.LicensesOK < 1 || snapshot.LicensesFailed < 1 || snapshot.AuthFailed < 1 {
		t.Fatalf("unexpected auth/license counters: %+v", snapshot)
	}
}
