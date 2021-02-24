package meilisearch

import (
	"testing"
)

func contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

func TestClientDumps_CreateAndGetStatus(t *testing.T) {
	resp, err := client.Dumps().Create()

	if err != nil {
		t.Fatal(err)
	}

	if resp.Status != "in_progress" {
		t.Fatal("response create dump does not have the 'in_progress' status")
	}
	var dumpUID = resp.UID
	resp, err = client.Dumps().GetStatus(dumpUID)
	if err != nil {
		t.Fatal(err)
	}
	if resp.UID != dumpUID {
		t.Fatal("response get dump status does not have the same UID")
	}

	var possibleStatuses = []string{"in_progress", "failed", "done"}
	if !contains(possibleStatuses, resp.Status) {
		t.Fatalf("response get dump status must be from %q", possibleStatuses)
	}
}
