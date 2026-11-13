package form

import (
	"slices"
	"testing"

	"github.com/Cyenoch/mygoui/locale"
	"github.com/egoist/mygo/ui"
)

func TestUploadReportsAcceptedDroppedPaths(t *testing.T) {
	var got [][]string
	o := UploadOptions{
		Title:    "Drop files",
		Accept:   []string{".png", "jpg"},
		Multiple: true,
		OnFiles:  func(paths []string) { got = append(got, paths) },
	}
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() { Upload(c, o) })
	}, 360, 240)
	zone, ok := tt.Find("Drop files")
	if !ok {
		t.Fatal("the drop zone is not rendered")
	}
	x, y := zone.X+zone.W/2, zone.Y+zone.H/2

	if !tt.DropFiles(x, y, "/tmp/a.PNG", "/tmp/notes.txt", "/tmp/b.jpg") {
		t.Fatal("the zone did not take the drop")
	}
	if len(got) != 1 || !slices.Equal(got[0], []string{"/tmp/a.PNG", "/tmp/b.jpg"}) {
		t.Fatalf("OnFiles got %v, want one report of the png and jpg paths", got)
	}

	// Without Multiple only the first accepted file is reported.
	got, o.Multiple = nil, false
	tt.Frame()
	tt.DropFiles(x, y, "/tmp/notes.txt", "/tmp/a.png", "/tmp/b.jpg")
	if len(got) != 1 || !slices.Equal(got[0], []string{"/tmp/a.png"}) {
		t.Fatalf("a single upload got %v, want only /tmp/a.png", got)
	}

	// A disabled zone reports nothing.
	got, o.Disabled = nil, true
	tt.Frame()
	tt.DropFiles(x, y, "/tmp/a.png")
	if len(got) != 0 {
		t.Fatalf("a disabled zone reported %v", got)
	}
}

func TestUploadListRemovesAndChooses(t *testing.T) {
	removed, chosen := -1, 0
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() {
			Upload(c, UploadOptions{
				Title:       "Drop files",
				ChooseLabel: "Choose files",
				OnChoose:    func() { chosen++ },
				Files: []UploadFile{
					{Name: "a.png", Detail: "12 KB", Status: UploadDone},
					{Name: "b.png", Status: UploadUploading, Progress: 0.4},
				},
				OnRemove: func(i int) { removed = i },
			})
		})
	}, 360, 360)
	if err := tt.Click("Remove b.png"); err != nil {
		t.Fatal(err)
	}
	if removed != 1 {
		t.Fatalf("removing b.png reported index %d, want 1", removed)
	}
	if err := tt.Click("Choose files"); err != nil {
		t.Fatal(err)
	}
	if chosen != 1 {
		t.Fatalf("the choose action ran %d times, want 1", chosen)
	}
}

func TestUploadFilesAreAListWithStatus(t *testing.T) {
	tt := formTextTester(t, func(c *ui.Context) {
		Upload(c, UploadOptions{Title: "Files", Files: []UploadFile{
			{Name: "a.txt", Status: UploadDone},
			{Name: "b.txt", Status: UploadFailed},
		}})
	})
	if info, ok := tt.AccessByRole(ui.RoleListItem, "b.txt"); !ok || info.PosInSet != 2 || info.SetSize != 2 {
		t.Fatalf("row b.txt: found %v, %d of %d", ok, info.PosInSet, info.SetSize)
	}
	for _, status := range []string{locale.M().Uploaded, locale.M().UploadFailed} {
		if _, ok := tt.AccessByRole(ui.RoleImage, status); !ok {
			t.Fatalf("no status %q", status)
		}
	}
}
