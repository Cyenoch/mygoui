package locale_test

import (
	"testing"
	"testing/fstest"

	"github.com/egoist/mygo/i18n"

	"github.com/Cyenoch/mygoui/locale"
)

func TestHostOverrideUsesLibraryNamespaceAndKeepsOtherTranslations(t *testing.T) {
	host, err := i18n.NewCatalog(fstest.MapFS{
		"en.json":      &fstest.MapFile{Data: []byte(`{"mygoui.componentClose":"Dismiss panel"}`)},
		"zh-Hans.json": &fstest.MapFile{Data: []byte(`{"mygoui.componentClose":"收起面板"}`)},
	}, "en")
	if err != nil {
		t.Fatal(err)
	}
	i18n.Register(host)
	if got := locale.In("en").ComponentClose(); got != "Dismiss panel" {
		t.Fatalf("host English override = %q", got)
	}
	if got := locale.In("zh-Hans").ComponentClose(); got != "收起面板" {
		t.Fatalf("host Chinese override = %q", got)
	}
	if got := locale.In("en").ComponentLoading(); got != "Loading" {
		t.Fatalf("unmodified English message = %q", got)
	}
	if got := locale.In("zh-Hans").ComponentLoading(); got == "Loading" || got == "" {
		t.Fatalf("library Chinese fallback was lost: %q", got)
	}
}
