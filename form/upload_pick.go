package form

import (
	"strings"

	"github.com/egoist/mygo"
)

// PickFilesOptions configures PickFiles. Accept and Multiple mean what they
// do in UploadOptions; Parent attaches the dialog to a window.
type PickFilesOptions struct {
	Parent   *mygo.Window
	Accept   []string
	Multiple bool
}

// PickFiles shows the system's open dialog and returns the chosen paths, nil
// when canceled. It blocks until the dialog is dismissed, so applications
// call it from an UploadOptions.OnChoose worker and publish the result with
// Window.Update. Platforms without a native dialog, as the browser, return
// MyGo's error; leave OnChoose unset there.
func PickFiles(o PickFilesOptions) ([]string, error) {
	var filters []mygo.FileFilter
	if len(o.Accept) > 0 {
		exts := make([]string, len(o.Accept))
		for i, a := range o.Accept {
			exts[i] = strings.TrimPrefix(a, ".")
		}
		filters = []mygo.FileFilter{{Name: "*." + strings.Join(exts, ", *."), Extensions: exts}}
	}
	return mygo.Dialog.Open(mygo.OpenDialogOptions{
		Parent: o.Parent, Filters: filters, Multiple: o.Multiple,
	})
}
