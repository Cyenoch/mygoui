// Package locales owns the gallery application's messages and languages.
package locales

import (
	"embed"
	"github.com/egoist/mygo/i18n"
	_ "github.com/egoist/mygo/locales/zh-Hans"
)

//go:embed *.json
var files embed.FS

var Catalog = func() *i18n.Catalog {
	catalog, err := i18n.NewCatalog(files, "en")
	if err != nil {
		panic(err)
	}
	i18n.Register(catalog)
	return catalog
}()

type Messages struct{ locale string }

func In(locale string) Messages { return Messages{locale: locale} }

// ComponentGalleryToolbar says "Gallery toolbar"
// in the language the app shows.
func ComponentGalleryToolbar() string { return In(i18n.Locale()).ComponentGalleryToolbar() }

// ComponentGalleryToolbar says "Gallery toolbar"
// in the language of msgs.
func (msgs Messages) ComponentGalleryToolbar() string {
	return Catalog.Format(msgs.locale, "componentGalleryToolbar", nil)
}

// ComponentGalleryBrand says "mygoui"
// in the language the app shows.
func ComponentGalleryBrand() string { return In(i18n.Locale()).ComponentGalleryBrand() }

// ComponentGalleryBrand says "mygoui"
// in the language of msgs.
func (msgs Messages) ComponentGalleryBrand() string {
	return Catalog.Format(msgs.locale, "componentGalleryBrand", nil)
}

// ComponentGallerySearch says "Search components"
// in the language the app shows.
func ComponentGallerySearch() string { return In(i18n.Locale()).ComponentGallerySearch() }

// ComponentGallerySearch says "Search components"
// in the language of msgs.
func (msgs Messages) ComponentGallerySearch() string {
	return Catalog.Format(msgs.locale, "componentGallerySearch", nil)
}

// ComponentGallerySearchPlaceholder says "Search components…"
// in the language the app shows.
func ComponentGallerySearchPlaceholder() string {
	return In(i18n.Locale()).ComponentGallerySearchPlaceholder()
}

// ComponentGallerySearchPlaceholder says "Search components…"
// in the language of msgs.
func (msgs Messages) ComponentGallerySearchPlaceholder() string {
	return Catalog.Format(msgs.locale, "componentGallerySearchPlaceholder", nil)
}

// ComponentGalleryDisplayOptions says "Display options"
// in the language the app shows.
func ComponentGalleryDisplayOptions() string {
	return In(i18n.Locale()).ComponentGalleryDisplayOptions()
}

// ComponentGalleryDisplayOptions says "Display options"
// in the language of msgs.
func (msgs Messages) ComponentGalleryDisplayOptions() string {
	return Catalog.Format(msgs.locale, "componentGalleryDisplayOptions", nil)
}

// ComponentGalleryBlueTheme says "App theme"
// in the language the app shows.
func ComponentGalleryBlueTheme() string { return In(i18n.Locale()).ComponentGalleryBlueTheme() }

// ComponentGalleryBlueTheme says "App theme"
// in the language of msgs.
func (msgs Messages) ComponentGalleryBlueTheme() string {
	return Catalog.Format(msgs.locale, "componentGalleryBlueTheme", nil)
}

// ComponentGalleryHighContrast says "High contrast"
// in the language the app shows.
func ComponentGalleryHighContrast() string { return In(i18n.Locale()).ComponentGalleryHighContrast() }

// ComponentGalleryHighContrast says "High contrast"
// in the language of msgs.
func (msgs Messages) ComponentGalleryHighContrast() string {
	return Catalog.Format(msgs.locale, "componentGalleryHighContrast", nil)
}

// ComponentGalleryLargeText says "Large text"
// in the language the app shows.
func ComponentGalleryLargeText() string { return In(i18n.Locale()).ComponentGalleryLargeText() }

// ComponentGalleryLargeText says "Large text"
// in the language of msgs.
func (msgs Messages) ComponentGalleryLargeText() string {
	return Catalog.Format(msgs.locale, "componentGalleryLargeText", nil)
}

// ComponentGalleryNavigation says "Component navigation"
// in the language the app shows.
func ComponentGalleryNavigation() string { return In(i18n.Locale()).ComponentGalleryNavigation() }

// ComponentGalleryNavigation says "Component navigation"
// in the language of msgs.
func (msgs Messages) ComponentGalleryNavigation() string {
	return Catalog.Format(msgs.locale, "componentGalleryNavigation", nil)
}

// ComponentGalleryUseLightTheme says "Use light theme"
// in the language the app shows.
func ComponentGalleryUseLightTheme() string { return In(i18n.Locale()).ComponentGalleryUseLightTheme() }

// ComponentGalleryUseLightTheme says "Use light theme"
// in the language of msgs.
func (msgs Messages) ComponentGalleryUseLightTheme() string {
	return Catalog.Format(msgs.locale, "componentGalleryUseLightTheme", nil)
}

// ComponentGalleryUseDarkTheme says "Use dark theme"
// in the language the app shows.
func ComponentGalleryUseDarkTheme() string { return In(i18n.Locale()).ComponentGalleryUseDarkTheme() }

// ComponentGalleryUseDarkTheme says "Use dark theme"
// in the language of msgs.
func (msgs Messages) ComponentGalleryUseDarkTheme() string {
	return Catalog.Format(msgs.locale, "componentGalleryUseDarkTheme", nil)
}

// ComponentGalleryPlannedName says "{name} · planned"
// in the language the app shows.
func ComponentGalleryPlannedName(name string) string {
	return In(i18n.Locale()).ComponentGalleryPlannedName(name)
}

// ComponentGalleryPlannedName says "{name} · planned"
// in the language of msgs.
func (msgs Messages) ComponentGalleryPlannedName(name string) string {
	return Catalog.Format(msgs.locale, "componentGalleryPlannedName", map[string]any{"name": name})
}

// ComponentGalleryNoMatches says "No matching components"
// in the language the app shows.
func ComponentGalleryNoMatches() string { return In(i18n.Locale()).ComponentGalleryNoMatches() }

// ComponentGalleryNoMatches says "No matching components"
// in the language of msgs.
func (msgs Messages) ComponentGalleryNoMatches() string {
	return Catalog.Format(msgs.locale, "componentGalleryNoMatches", nil)
}

// ComponentGalleryClearSearch says "Clear search"
// in the language the app shows.
func ComponentGalleryClearSearch() string { return In(i18n.Locale()).ComponentGalleryClearSearch() }

// ComponentGalleryClearSearch says "Clear search"
// in the language of msgs.
func (msgs Messages) ComponentGalleryClearSearch() string {
	return Catalog.Format(msgs.locale, "componentGalleryClearSearch", nil)
}

// ComponentGalleryPreview says "Component preview"
// in the language the app shows.
func ComponentGalleryPreview() string { return In(i18n.Locale()).ComponentGalleryPreview() }

// ComponentGalleryPreview says "Component preview"
// in the language of msgs.
func (msgs Messages) ComponentGalleryPreview() string {
	return Catalog.Format(msgs.locale, "componentGalleryPreview", nil)
}

// ComponentGalleryDocumentation says "Component documentation"
// in the language the app shows.
func ComponentGalleryDocumentation() string { return In(i18n.Locale()).ComponentGalleryDocumentation() }

// ComponentGalleryDocumentation says "Component documentation"
// in the language of msgs.
func (msgs Messages) ComponentGalleryDocumentation() string {
	return Catalog.Format(msgs.locale, "componentGalleryDocumentation", nil)
}

// ComponentGalleryPlanned says "Planned"
// in the language the app shows.
func ComponentGalleryPlanned() string { return In(i18n.Locale()).ComponentGalleryPlanned() }

// ComponentGalleryPlanned says "Planned"
// in the language of msgs.
func (msgs Messages) ComponentGalleryPlanned() string {
	return Catalog.Format(msgs.locale, "componentGalleryPlanned", nil)
}

// ComponentGalleryImplementedUnverified says "Implemented / unverified"
// in the language the app shows.
func ComponentGalleryImplementedUnverified() string {
	return In(i18n.Locale()).ComponentGalleryImplementedUnverified()
}

// ComponentGalleryImplementedUnverified says "Implemented / unverified"
// in the language of msgs.
func (msgs Messages) ComponentGalleryImplementedUnverified() string {
	return Catalog.Format(msgs.locale, "componentGalleryImplementedUnverified", nil)
}

// ComponentGalleryInteractivePreview says "Interactive preview"
// in the language the app shows.
func ComponentGalleryInteractivePreview() string {
	return In(i18n.Locale()).ComponentGalleryInteractivePreview()
}

// ComponentGalleryInteractivePreview says "Interactive preview"
// in the language of msgs.
func (msgs Messages) ComponentGalleryInteractivePreview() string {
	return Catalog.Format(msgs.locale, "componentGalleryInteractivePreview", nil)
}

// ComponentGalleryViewCode says "View Code"
// in the language the app shows.
func ComponentGalleryViewCode() string { return In(i18n.Locale()).ComponentGalleryViewCode() }

// ComponentGalleryViewCode says "View Code"
// in the language of msgs.
func (msgs Messages) ComponentGalleryViewCode() string {
	return Catalog.Format(msgs.locale, "componentGalleryViewCode", nil)
}

// ComponentGalleryHideCode says "Hide Code"
// in the language the app shows.
func ComponentGalleryHideCode() string { return In(i18n.Locale()).ComponentGalleryHideCode() }

// ComponentGalleryHideCode says "Hide Code"
// in the language of msgs.
func (msgs Messages) ComponentGalleryHideCode() string {
	return Catalog.Format(msgs.locale, "componentGalleryHideCode", nil)
}

// ComponentGalleryActualSource says "Actual preview source"
// in the language the app shows.
func ComponentGalleryActualSource() string { return In(i18n.Locale()).ComponentGalleryActualSource() }

// ComponentGalleryActualSource says "Actual preview source"
// in the language of msgs.
func (msgs Messages) ComponentGalleryActualSource() string {
	return Catalog.Format(msgs.locale, "componentGalleryActualSource", nil)
}

// ComponentGalleryImplementationPlanned says "Native implementation planned"
// in the language the app shows.
func ComponentGalleryImplementationPlanned() string {
	return In(i18n.Locale()).ComponentGalleryImplementationPlanned()
}

// ComponentGalleryImplementationPlanned says "Native implementation planned"
// in the language of msgs.
func (msgs Messages) ComponentGalleryImplementationPlanned() string {
	return Catalog.Format(msgs.locale, "componentGalleryImplementationPlanned", nil)
}

// ComponentGalleryPreviewUnavailable says "A preview will appear when this capability is implemented. The catalog records the native plan and remaining acceptan…"
// in the language the app shows.
func ComponentGalleryPreviewUnavailable() string {
	return In(i18n.Locale()).ComponentGalleryPreviewUnavailable()
}

// ComponentGalleryPreviewUnavailable says "A preview will appear when this capability is implemented. The catalog records the native plan and remaining acceptan…"
// in the language of msgs.
func (msgs Messages) ComponentGalleryPreviewUnavailable() string {
	return Catalog.Format(msgs.locale, "componentGalleryPreviewUnavailable", nil)
}

// ComponentGalleryVariants says "{name} variants and states"
// in the language the app shows.
func ComponentGalleryVariants(name string) string {
	return In(i18n.Locale()).ComponentGalleryVariants(name)
}

// ComponentGalleryVariants says "{name} variants and states"
// in the language of msgs.
func (msgs Messages) ComponentGalleryVariants(name string) string {
	return Catalog.Format(msgs.locale, "componentGalleryVariants", map[string]any{"name": name})
}

// ComponentGalleryReference says "Reference: {url}"
// in the language the app shows.
func ComponentGalleryReference(url string) string {
	return In(i18n.Locale()).ComponentGalleryReference(url)
}

// ComponentGalleryReference says "Reference: {url}"
// in the language of msgs.
func (msgs Messages) ComponentGalleryReference(url string) string {
	return Catalog.Format(msgs.locale, "componentGalleryReference", map[string]any{"url": url})
}

// ComponentGalleryOnThisPage says "On this page"
// in the language the app shows.
func ComponentGalleryOnThisPage() string { return In(i18n.Locale()).ComponentGalleryOnThisPage() }

// ComponentGalleryOnThisPage says "On this page"
// in the language of msgs.
func (msgs Messages) ComponentGalleryOnThisPage() string {
	return Catalog.Format(msgs.locale, "componentGalleryOnThisPage", nil)
}

// ComponentGalleryJumpSection says "Jump to {section}"
// in the language the app shows.
func ComponentGalleryJumpSection(section string) string {
	return In(i18n.Locale()).ComponentGalleryJumpSection(section)
}

// ComponentGalleryJumpSection says "Jump to {section}"
// in the language of msgs.
func (msgs Messages) ComponentGalleryJumpSection(section string) string {
	return Catalog.Format(msgs.locale, "componentGalleryJumpSection", map[string]any{"section": section})
}

// ComponentGalleryWindowTitle says "mygoui Gallery"
// in the language the app shows.
func ComponentGalleryWindowTitle() string { return In(i18n.Locale()).ComponentGalleryWindowTitle() }

// ComponentGalleryWindowTitle says "mygoui Gallery"
// in the language of msgs.
func (msgs Messages) ComponentGalleryWindowTitle() string {
	return Catalog.Format(msgs.locale, "componentGalleryWindowTitle", nil)
}

// ComponentGallerySection says "{section, select, Preview {Preview} Usage {Usage} Examples {Examples} Behavior_and_API {Behavior and API} other {{sec…"
// in the language the app shows.
func ComponentGallerySection(section string) string {
	return In(i18n.Locale()).ComponentGallerySection(section)
}

// ComponentGallerySection says "{section, select, Preview {Preview} Usage {Usage} Examples {Examples} Behavior_and_API {Behavior and API} other {{sec…"
// in the language of msgs.
func (msgs Messages) ComponentGallerySection(section string) string {
	return Catalog.Format(msgs.locale, "componentGallerySection", map[string]any{"section": section})
}

// ComponentGalleryFamily says "{family, select, Navigation {Navigation} Presentation {Presentation} Overlays {Overlays} Conversation {Conversation} …"
// in the language the app shows.
func ComponentGalleryFamily(family string) string {
	return In(i18n.Locale()).ComponentGalleryFamily(family)
}

// ComponentGalleryFamily says "{family, select, Navigation {Navigation} Presentation {Presentation} Overlays {Overlays} Conversation {Conversation} …"
// in the language of msgs.
func (msgs Messages) ComponentGalleryFamily(family string) string {
	return Catalog.Format(msgs.locale, "componentGalleryFamily", map[string]any{"family": family})
}

// ComponentGalleryName says "{id, select, accordion {Accordion} alert {Alert} alert_dialog {Alert Dialog} aspect_ratio {Aspect Ratio} attachment {…"
// in the language the app shows.
func ComponentGalleryName(id string) string { return In(i18n.Locale()).ComponentGalleryName(id) }

// ComponentGalleryName says "{id, select, accordion {Accordion} alert {Alert} alert_dialog {Alert Dialog} aspect_ratio {Aspect Ratio} attachment {…"
// in the language of msgs.
func (msgs Messages) ComponentGalleryName(id string) string {
	return Catalog.Format(msgs.locale, "componentGalleryName", map[string]any{"id": id})
}
