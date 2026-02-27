// Package fileutil provides general-purpose filesystem and MIME-type helpers.
// It has no dependencies outside the Go standard library and is safe to import
// from any package in the project.
package fileutil

import (
	"io"
	"mime"
	"os"
	"strings"
)

// CopyFile copies the file at src to dst, creating dst (or truncating it) as
// needed. The destination directory must already exist.
func CopyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

// AllowedExts returns a set of lowercase file extensions that correspond to the
// given MIME types (e.g. ["image/jpeg", "image/png"]). When the standard library
// cannot resolve a MIME type, the subtype is used directly as the extension
// (e.g. "image/webp" -> ".webp").
func AllowedExts(mimeTypes []string) map[string]bool {
	exts := make(map[string]bool)
	for _, mt := range mimeTypes {
		list, err := mime.ExtensionsByType(mt)
		if err != nil || len(list) == 0 {
			// Fallback: derive extension from the MIME subtype.
			if parts := strings.SplitN(mt, "/", 2); len(parts) == 2 {
				exts["."+parts[1]] = true
			}
			continue
		}
		for _, e := range list {
			exts[strings.ToLower(e)] = true
		}
	}
	return exts
}

// MimeByExt returns the MIME type for the given lowercase file extension
// (e.g. ".jpg"). Falls back to "application/octet-stream" when the extension
// is not recognised by the standard library.
func MimeByExt(ext string) string {
	if mt := mime.TypeByExtension(ext); mt != "" {
		return mt
	}
	return "application/octet-stream"
}

// Exists reports whether the path exists on disk and is a regular file (not a
// directory).
func Exists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}
