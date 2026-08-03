package cmd

import "testing"

func TestRootCommandHasName(t *testing.T) {
	if RootCmd.Use != "dforge" {
		t.Fatalf("want use=dforge, got %q", RootCmd.Use)
	}
}

func TestRootHasGlobalFlags(t *testing.T) {
	if RootCmd.PersistentFlags().Lookup("yes") == nil {
		t.Fatal("missing --yes flag")
	}
	if RootCmd.PersistentFlags().Lookup("force") == nil {
		t.Fatal("missing --force flag")
	}
}
