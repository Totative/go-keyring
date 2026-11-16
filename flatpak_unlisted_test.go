//go:build linux

package keyring

import (
	"context"
	"errors"
	"os"
	"testing"
)

func TestFlatpakUnlistedService(t *testing.T) {
	t.Setenv("FLATPAK_ID", "org.example.App")

	newStore := func(t *testing.T) (*Store, string) {
		t.Helper()
		dir := t.TempDir()
		s, err := New(Options{
			Service:        "example",
			SealInfo:       "example secrets v1",
			Folder:         "example-secrets",
			SharedServices: []string{"listed"},
		})
		if err != nil {
			t.Fatal(err)
		}
		s.storeDir = func() (string, error) { return dir, nil }
		s.portalSecret = func(context.Context) ([]byte, error) {
			return []byte("0123456789abcdef0123456789abcdef"), nil
		}
		s.get = func(string, string) (string, error) { return "", ErrNotFound }
		s.set = func(string, string, string) error { return errors.New("no host") }
		s.del = func(string, string) error { return ErrNotFound }
		return s, dir
	}

	t.Run("unlisted refused with no file", func(t *testing.T) {
		s, dir := newStore(t)
		if err := s.SetIn("other", "item", "secret"); !errors.Is(err, ErrServiceNotShared) {
			t.Fatalf("SetIn: %v", err)
		}
		if _, err := s.GetIn("other", "item"); !errors.Is(err, ErrServiceNotShared) {
			t.Fatalf("GetIn: %v", err)
		}
		if err := s.DeleteIn("other", "item"); !errors.Is(err, ErrServiceNotShared) {
			t.Fatalf("DeleteIn: %v", err)
		}
		ents, err := os.ReadDir(dir)
		if err != nil {
			t.Fatal(err)
		}
		if len(ents) != 0 {
			t.Fatalf("wrote %v", ents)
		}
	})

	t.Run("listed still works", func(t *testing.T) {
		s, _ := newStore(t)
		if err := s.SetIn("listed", "item", "shared-secret"); err != nil {
			t.Fatal(err)
		}
		v, err := s.GetIn("listed", "item")
		if err != nil || v != "shared-secret" {
			t.Fatalf("GetIn listed: %q %v", v, err)
		}
	})

	t.Run("own dotted name unaffected", func(t *testing.T) {
		s, _ := newStore(t)
		if err := s.Set("other.item", "own-secret"); err != nil {
			t.Fatal(err)
		}
		v, err := s.Get("other.item")
		if err != nil || v != "own-secret" {
			t.Fatalf("Get own: %q %v", v, err)
		}
		if _, err := s.GetIn("other", "item"); !errors.Is(err, ErrServiceNotShared) {
			t.Fatalf("GetIn: %v", err)
		}
		if err := s.DeleteIn("other", "item"); !errors.Is(err, ErrServiceNotShared) {
			t.Fatalf("DeleteIn: %v", err)
		}
		v, err = s.Get("other.item")
		if err != nil || v != "own-secret" {
			t.Fatalf("own item after DeleteIn: %q %v", v, err)
		}
	})
}
