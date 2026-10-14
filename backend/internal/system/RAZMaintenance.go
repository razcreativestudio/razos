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

// RestartHost mengeksekusi reboot aman pada perangkat host Linux STB
func RestartHost() error {
	cmd := exec.Command("reboot")
	if err := cmd.Run(); err != nil {
		// Fallback perintah shutdown -r now jika reboot direct ditolak
		fallback := exec.Command("shutdown", "-r", "now")
		return fallback.Run()
	}
	return nil
}

// ShutdownHost mematikan daya perangkat STB host secara aman
func ShutdownHost() error {
	cmd := exec.Command("poweroff")
	if err := cmd.Run(); err != nil {
		// Fallback perintah shutdown -h now jika poweroff direct ditolak
		fallback := exec.Command("shutdown", "-h", "now")
		return fallback.Run()
	}
	return nil
}
