package grouper_test

import (
	"testing"

	"github.com/user/envdiff/internal/grouper"
)

func baseEnv() map[string]string {
	return map[string]string{
		"DB_HOST":     "localhost",
		"DB_PORT":     "5432",
		"DB_NAME":     "mydb",
		"APP_DEBUG":   "true",
		"APP_PORT":    "8080",
		"STANDALONE":  "yes",
		"LOG_LEVEL":   "info",
	}
}

func TestGroupEnv_BasicPrefixes(t *testing.T) {
	res := grouper.GroupEnv(baseEnv(), grouper.DefaultOptions())

	if len(res.Groups) != 3 {
		t.Fatalf("expected 3 groups, got %d", len(res.Groups))
	}
	if len(res.Ungrouped) != 1 || res.Ungrouped[0] != "STANDALONE" {
		t.Fatalf("expected STANDALONE in ungrouped, got %v", res.Ungrouped)
	}
}

func TestGroupEnv_GroupContents(t *testing.T) {
	res := grouper.GroupEnv(baseEnv(), grouper.DefaultOptions())

	var db *grouper.Group
	for _, g := range res.Groups {
		if g.Prefix == "DB" {
			db = g
			break
		}
	}
	if db == nil {
		t.Fatal("expected DB group")
	}
	if len(db.Keys) != 3 {
		t.Errorf("expected 3 keys in DB group, got %d", len(db.Keys))
	}
	if db.Env["DB_HOST"] != "localhost" {
		t.Errorf("unexpected value for DB_HOST: %s", db.Env["DB_HOST"])
	}
}

func TestGroupEnv_SortedOutput(t *testing.T) {
	res := grouper.GroupEnv(baseEnv(), grouper.DefaultOptions())

	for i := 1; i < len(res.Groups); i++ {
		if res.Groups[i].Prefix < res.Groups[i-1].Prefix {
			t.Errorf("groups not sorted: %s before %s", res.Groups[i-1].Prefix, res.Groups[i].Prefix)
		}
	}
}

func TestGroupEnv_EmptyEnv(t *testing.T) {
	res := grouper.GroupEnv(map[string]string{}, grouper.DefaultOptions())
	if len(res.Groups) != 0 || len(res.Ungrouped) != 0 {
		t.Error("expected empty result for empty env")
	}
}

func TestGroupEnv_CustomDelimiter(t *testing.T) {
	env := map[string]string{
		"APP.DEBUG": "true",
		"APP.PORT":  "8080",
		"NODEL":     "val",
	}
	opts := grouper.Options{Delimiter: ".", MaxDepth: 1}
	res := grouper.GroupEnv(env, opts)

	if len(res.Groups) != 1 || res.Groups[0].Prefix != "APP" {
		t.Errorf("expected APP group, got %v", res.Groups)
	}
	if len(res.Ungrouped) != 1 {
		t.Errorf("expected 1 ungrouped, got %v", res.Ungrouped)
	}
}

func TestGroupEnv_MaxDepth2(t *testing.T) {
	env := map[string]string{
		"AWS_S3_BUCKET": "my-bucket",
		"AWS_S3_REGION": "us-east-1",
		"AWS_RDS_HOST":  "db.example.com",
	}
	opts := grouper.Options{Delimiter: "_", MaxDepth: 2}
	res := grouper.GroupEnv(env, opts)

	if len(res.Groups) != 2 {
		t.Fatalf("expected 2 groups at depth 2, got %d", len(res.Groups))
	}
	if res.Groups[0].Prefix != "AWS_RDS" && res.Groups[1].Prefix != "AWS_RDS" {
		t.Errorf("expected AWS_RDS group, got %v", res.Groups)
	}
}
