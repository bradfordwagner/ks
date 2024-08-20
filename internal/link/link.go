package link

import "os"

// ForceLink creates a symbolic link from source to target, removing the target if it already exists
func ForceLink(source, target string) (err error) {
	// Remove the target if it exists
	if err = os.Remove(target); err != nil && !os.IsNotExist(err) {
		return err
	}

	// Create the symbolic link
	if err = os.Symlink(source, target); err != nil {
		return err
	}

	return nil
}
