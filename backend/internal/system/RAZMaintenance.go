/**
 * RAZOS System Maintenance
 * Updated: 2026-10-08
 * Deskripsi: Modul pemeliharaan memori RAM dan ruang disk dengan eksekusi aman langsung ke kernel.
 */

package system

import (
	"os/exec"
)

// CleanRAM menjalankan sync buffer filesystem dan memerintahkan kernel membersihkan pagecache, dentries, dan inodes
func CleanRAM() error {
	cmd := exec.Command("sh", "-c", "sync && echo 3 > /proc/sys/vm/drop_caches")
	return cmd.Run()
}

// CleanStorage menjalankan pembersihan container & image Docker yatim piatu
func CleanStorage() (string, error) {
	cmd := exec.Command("docker", "system", "prune", "-f")
	out, err := cmd.CombinedOutput()
	return string(out), err
}
