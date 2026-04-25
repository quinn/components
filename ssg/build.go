package ssg

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Build writes every registered route and static mount to outDir.
//
// A route whose path has a file extension (for example "/feed.xml") is
// written verbatim. A route without an extension is treated as a directory
// and written as <route>/index.html so that the resulting tree works with
// any standard static host. The root route "/" is written as index.html.
//
// Static mounts are copied into outDir under their URL prefix. Routes are
// written first; static files therefore overlay routes if they collide,
// which lets callers drop in pre-rendered assets without modifying their
// route registrations.
func (s *Site) Build(outDir string) error {
	if outDir == "" {
		return fmt.Errorf("ssg: Build outDir must not be empty")
	}
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return fmt.Errorf("ssg: create output dir: %w", err)
	}

	for _, route := range s.Routes() {
		if err := writeRoute(outDir, route); err != nil {
			return err
		}
	}

	for _, mount := range s.staticsSnapshot() {
		if err := writeStatic(outDir, mount); err != nil {
			return err
		}
	}

	return nil
}

func writeRoute(outDir string, route Route) error {
	rel := routeOutputPath(route.Path)
	full := filepath.Join(outDir, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		return fmt.Errorf("ssg: mkdir %s: %w", filepath.Dir(full), err)
	}
	f, err := os.Create(full)
	if err != nil {
		return fmt.Errorf("ssg: create %s: %w", full, err)
	}
	if err := route.Renderer.Render(f); err != nil {
		_ = f.Close()
		return fmt.Errorf("ssg: render %s: %w", route.Path, err)
	}
	if err := f.Close(); err != nil {
		return fmt.Errorf("ssg: close %s: %w", full, err)
	}
	return nil
}

// routeOutputPath maps a URL path to a relative file path on disk.
func routeOutputPath(routePath string) string {
	if routePath == "" || routePath == "/" {
		return "index.html"
	}
	cleaned := strings.TrimPrefix(routePath, "/")
	cleaned = strings.TrimSuffix(cleaned, "/")
	if path.Ext(cleaned) != "" {
		return cleaned
	}
	return cleaned + "/index.html"
}

func writeStatic(outDir string, mount staticMount) error {
	destBase := outDir
	if mount.Prefix != "/" {
		destBase = filepath.Join(outDir, filepath.FromSlash(strings.TrimPrefix(mount.Prefix, "/")))
	}
	return fs.WalkDir(mount.FS, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		target := filepath.Join(destBase, filepath.FromSlash(p))
		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		return copyFSFile(mount.FS, p, target)
	})
}

func copyFSFile(fsys fs.FS, src, dst string) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	in, err := fsys.Open(src)
	if err != nil {
		return fmt.Errorf("ssg: open static %s: %w", src, err)
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("ssg: create %s: %w", dst, err)
	}
	if _, err := io.Copy(out, in); err != nil {
		_ = out.Close()
		return fmt.Errorf("ssg: copy %s: %w", src, err)
	}
	return out.Close()
}
