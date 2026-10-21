// Package locale provides the component library's maintained message catalog.
package locale

import (
	"embed"

	"github.com/egoist/mygo/i18n"
)

//go:embed en.json zh-Hans.json
var files embed.FS

// Catalog contains component messages. Hosts may override them with keys such as
// "mygoui.componentClose" in their own catalogs; only hosts register languages.
var Catalog *i18n.Catalog

func init() {
	var err error
	if Catalog, err = i18n.NewLibraryCatalog(files, "en", "mygoui"); err != nil {
		panic(err)
	}
}

// Messages provides component messages in a selected locale.
type Messages struct{ locale string }

// In selects a locale, falling back according to the library catalog's rules.
func In(locale string) Messages { return Messages{locale: locale} }

// ComponentLoading says "Loading"
// in the language the app shows.
func ComponentLoading() string { return In(i18n.Locale()).ComponentLoading() }

// ComponentLoading says "Loading"
// in the language of msgs.
func (msgs Messages) ComponentLoading() string {
	return Catalog.Format(msgs.locale, "componentLoading", nil)
}

// ComponentMorePeople says "{count, plural, one {# more person} other {# more people}}"
// in the language the app shows.
func ComponentMorePeople(count int) string { return In(i18n.Locale()).ComponentMorePeople(count) }

// ComponentMorePeople says "{count, plural, one {# more person} other {# more people}}"
// in the language of msgs.
func (msgs Messages) ComponentMorePeople(count int) string {
	return Catalog.Format(msgs.locale, "componentMorePeople", map[string]any{"count": count})
}

// ComponentClose says "Close"
// in the language the app shows.
func ComponentClose() string { return In(i18n.Locale()).ComponentClose() }

// ComponentClose says "Close"
// in the language of msgs.
func (msgs Messages) ComponentClose() string {
	return Catalog.Format(msgs.locale, "componentClose", nil)
}

// ComponentCancel says "Cancel"
// in the language the app shows.
func ComponentCancel() string { return In(i18n.Locale()).ComponentCancel() }

// ComponentCancel says "Cancel"
// in the language of msgs.
func (msgs Messages) ComponentCancel() string {
	return Catalog.Format(msgs.locale, "componentCancel", nil)
}

// ComponentContinue says "Continue"
// in the language the app shows.
func ComponentContinue() string { return In(i18n.Locale()).ComponentContinue() }

// ComponentContinue says "Continue"
// in the language of msgs.
func (msgs Messages) ComponentContinue() string {
	return Catalog.Format(msgs.locale, "componentContinue", nil)
}

// ComponentCloseToast says "Close toast"
// in the language the app shows.
func ComponentCloseToast() string { return In(i18n.Locale()).ComponentCloseToast() }

// ComponentCloseToast says "Close toast"
// in the language of msgs.
func (msgs Messages) ComponentCloseToast() string {
	return Catalog.Format(msgs.locale, "componentCloseToast", nil)
}

// ComponentBreadcrumb says "Breadcrumb"
// in the language the app shows.
func ComponentBreadcrumb() string { return In(i18n.Locale()).ComponentBreadcrumb() }

// ComponentBreadcrumb says "Breadcrumb"
// in the language of msgs.
func (msgs Messages) ComponentBreadcrumb() string {
	return Catalog.Format(msgs.locale, "componentBreadcrumb", nil)
}

// ComponentMoreAncestors says "More ancestors"
// in the language the app shows.
func ComponentMoreAncestors() string { return In(i18n.Locale()).ComponentMoreAncestors() }

// ComponentMoreAncestors says "More ancestors"
// in the language of msgs.
func (msgs Messages) ComponentMoreAncestors() string {
	return Catalog.Format(msgs.locale, "componentMoreAncestors", nil)
}

// ComponentCurrentPage says "Current page"
// in the language the app shows.
func ComponentCurrentPage() string { return In(i18n.Locale()).ComponentCurrentPage() }

// ComponentCurrentPage says "Current page"
// in the language of msgs.
func (msgs Messages) ComponentCurrentPage() string {
	return Catalog.Format(msgs.locale, "componentCurrentPage", nil)
}

// ComponentPagination says "Pagination"
// in the language the app shows.
func ComponentPagination() string { return In(i18n.Locale()).ComponentPagination() }

// ComponentPagination says "Pagination"
// in the language of msgs.
func (msgs Messages) ComponentPagination() string {
	return Catalog.Format(msgs.locale, "componentPagination", nil)
}

// ComponentPrevious says "Previous"
// in the language the app shows.
func ComponentPrevious() string { return In(i18n.Locale()).ComponentPrevious() }

// ComponentPrevious says "Previous"
// in the language of msgs.
func (msgs Messages) ComponentPrevious() string {
	return Catalog.Format(msgs.locale, "componentPrevious", nil)
}

// ComponentNext says "Next"
// in the language the app shows.
func ComponentNext() string { return In(i18n.Locale()).ComponentNext() }

// ComponentNext says "Next"
// in the language of msgs.
func (msgs Messages) ComponentNext() string {
	return Catalog.Format(msgs.locale, "componentNext", nil)
}

// ComponentBack says "Back"
// in the language the app shows.
func ComponentBack() string { return In(i18n.Locale()).ComponentBack() }

// ComponentBack says "Back"
// in the language of msgs.
func (msgs Messages) ComponentBack() string {
	return Catalog.Format(msgs.locale, "componentBack", nil)
}

// ComponentNoPages says "No pages"
// in the language the app shows.
func ComponentNoPages() string { return In(i18n.Locale()).ComponentNoPages() }

// ComponentNoPages says "No pages"
// in the language of msgs.
func (msgs Messages) ComponentNoPages() string {
	return Catalog.Format(msgs.locale, "componentNoPages", nil)
}

// ComponentPageOf says "Page {page} of {pages}"
// in the language the app shows.
func ComponentPageOf(page int, pages int) string {
	return In(i18n.Locale()).ComponentPageOf(page, pages)
}

// ComponentPageOf says "Page {page} of {pages}"
// in the language of msgs.
func (msgs Messages) ComponentPageOf(page int, pages int) string {
	return Catalog.Format(msgs.locale, "componentPageOf", map[string]any{"page": page, "pages": pages})
}

// ComponentMorePages says "More pages"
// in the language the app shows.
func ComponentMorePages() string { return In(i18n.Locale()).ComponentMorePages() }

// ComponentMorePages says "More pages"
// in the language of msgs.
func (msgs Messages) ComponentMorePages() string {
	return Catalog.Format(msgs.locale, "componentMorePages", nil)
}

// ComponentGoToPage says "Go to page {page}"
// in the language the app shows.
func ComponentGoToPage(page int) string { return In(i18n.Locale()).ComponentGoToPage(page) }

// ComponentGoToPage says "Go to page {page}"
// in the language of msgs.
func (msgs Messages) ComponentGoToPage(page int) string {
	return Catalog.Format(msgs.locale, "componentGoToPage", map[string]any{"page": page})
}

// ComponentScrollChoicesUp says "Scroll choices up"
// in the language the app shows.
func ComponentScrollChoicesUp() string { return In(i18n.Locale()).ComponentScrollChoicesUp() }

// ComponentScrollChoicesUp says "Scroll choices up"
// in the language of msgs.
func (msgs Messages) ComponentScrollChoicesUp() string {
	return Catalog.Format(msgs.locale, "componentScrollChoicesUp", nil)
}

// ComponentScrollChoicesDown says "Scroll choices down"
// in the language the app shows.
func ComponentScrollChoicesDown() string { return In(i18n.Locale()).ComponentScrollChoicesDown() }

// ComponentScrollChoicesDown says "Scroll choices down"
// in the language of msgs.
func (msgs Messages) ComponentScrollChoicesDown() string {
	return Catalog.Format(msgs.locale, "componentScrollChoicesDown", nil)
}

// ComponentCommandSearchSurface says "Command search surface"
// in the language the app shows.
func ComponentCommandSearchSurface() string { return In(i18n.Locale()).ComponentCommandSearchSurface() }

// ComponentCommandSearchSurface says "Command search surface"
// in the language of msgs.
func (msgs Messages) ComponentCommandSearchSurface() string {
	return Catalog.Format(msgs.locale, "componentCommandSearchSurface", nil)
}

// ComponentNoCommandsFound says "No commands found"
// in the language the app shows.
func ComponentNoCommandsFound() string { return In(i18n.Locale()).ComponentNoCommandsFound() }

// ComponentNoCommandsFound says "No commands found"
// in the language of msgs.
func (msgs Messages) ComponentNoCommandsFound() string {
	return Catalog.Format(msgs.locale, "componentNoCommandsFound", nil)
}

// ComponentSidebar says "Sidebar"
// in the language the app shows.
func ComponentSidebar() string { return In(i18n.Locale()).ComponentSidebar() }

// ComponentSidebar says "Sidebar"
// in the language of msgs.
func (msgs Messages) ComponentSidebar() string {
	return Catalog.Format(msgs.locale, "componentSidebar", nil)
}

// ComponentToggleNamed says "Toggle {name}"
// in the language the app shows.
func ComponentToggleNamed(name string) string { return In(i18n.Locale()).ComponentToggleNamed(name) }

// ComponentToggleNamed says "Toggle {name}"
// in the language of msgs.
func (msgs Messages) ComponentToggleNamed(name string) string {
	return Catalog.Format(msgs.locale, "componentToggleNamed", map[string]any{"name": name})
}

// ComponentCloseNamed says "Close {name}"
// in the language the app shows.
func ComponentCloseNamed(name string) string { return In(i18n.Locale()).ComponentCloseNamed(name) }

// ComponentCloseNamed says "Close {name}"
// in the language of msgs.
func (msgs Messages) ComponentCloseNamed(name string) string {
	return Catalog.Format(msgs.locale, "componentCloseNamed", map[string]any{"name": name})
}

// ComponentFilterRowsPlaceholder says "Filter rows…"
// in the language the app shows.
func ComponentFilterRowsPlaceholder() string {
	return In(i18n.Locale()).ComponentFilterRowsPlaceholder()
}

// ComponentFilterRowsPlaceholder says "Filter rows…"
// in the language of msgs.
func (msgs Messages) ComponentFilterRowsPlaceholder() string {
	return Catalog.Format(msgs.locale, "componentFilterRowsPlaceholder", nil)
}

// ComponentFilterRows says "Filter rows"
// in the language the app shows.
func ComponentFilterRows() string { return In(i18n.Locale()).ComponentFilterRows() }

// ComponentFilterRows says "Filter rows"
// in the language of msgs.
func (msgs Messages) ComponentFilterRows() string {
	return Catalog.Format(msgs.locale, "componentFilterRows", nil)
}

// ComponentColumns says "Columns"
// in the language the app shows.
func ComponentColumns() string { return In(i18n.Locale()).ComponentColumns() }

// ComponentColumns says "Columns"
// in the language of msgs.
func (msgs Messages) ComponentColumns() string {
	return Catalog.Format(msgs.locale, "componentColumns", nil)
}

// ComponentRowsSelected says "{rows, plural, one {# row} other {# rows}} · {selected} selected"
// in the language the app shows.
func ComponentRowsSelected(rows int, selected int) string {
	return In(i18n.Locale()).ComponentRowsSelected(rows, selected)
}

// ComponentRowsSelected says "{rows, plural, one {# row} other {# rows}} · {selected} selected"
// in the language of msgs.
func (msgs Messages) ComponentRowsSelected(rows int, selected int) string {
	return Catalog.Format(msgs.locale, "componentRowsSelected", map[string]any{"rows": rows, "selected": selected})
}

// ComponentNoMatchingRows says "No matching rows"
// in the language the app shows.
func ComponentNoMatchingRows() string { return In(i18n.Locale()).ComponentNoMatchingRows() }

// ComponentNoMatchingRows says "No matching rows"
// in the language of msgs.
func (msgs Messages) ComponentNoMatchingRows() string {
	return Catalog.Format(msgs.locale, "componentNoMatchingRows", nil)
}

// ComponentSlideOf says "Slide {slide} of {slides}"
// in the language the app shows.
func ComponentSlideOf(slide int, slides int) string {
	return In(i18n.Locale()).ComponentSlideOf(slide, slides)
}

// ComponentSlideOf says "Slide {slide} of {slides}"
// in the language of msgs.
func (msgs Messages) ComponentSlideOf(slide int, slides int) string {
	return Catalog.Format(msgs.locale, "componentSlideOf", map[string]any{"slide": slide, "slides": slides})
}

// ComponentPositionOf says "{position} of {count}"
// in the language the app shows.
func ComponentPositionOf(position int, count int) string {
	return In(i18n.Locale()).ComponentPositionOf(position, count)
}

// ComponentPositionOf says "{position} of {count}"
// in the language of msgs.
func (msgs Messages) ComponentPositionOf(position int, count int) string {
	return Catalog.Format(msgs.locale, "componentPositionOf", map[string]any{"position": position, "count": count})
}

// ComponentPreviousSlide says "Previous slide"
// in the language the app shows.
func ComponentPreviousSlide() string { return In(i18n.Locale()).ComponentPreviousSlide() }

// ComponentPreviousSlide says "Previous slide"
// in the language of msgs.
func (msgs Messages) ComponentPreviousSlide() string {
	return Catalog.Format(msgs.locale, "componentPreviousSlide", nil)
}

// ComponentNextSlide says "Next slide"
// in the language the app shows.
func ComponentNextSlide() string { return In(i18n.Locale()).ComponentNextSlide() }

// ComponentNextSlide says "Next slide"
// in the language of msgs.
func (msgs Messages) ComponentNextSlide() string {
	return Catalog.Format(msgs.locale, "componentNextSlide", nil)
}

// ComponentNoSlides says "No slides"
// in the language the app shows.
func ComponentNoSlides() string { return In(i18n.Locale()).ComponentNoSlides() }

// ComponentNoSlides says "No slides"
// in the language of msgs.
func (msgs Messages) ComponentNoSlides() string {
	return Catalog.Format(msgs.locale, "componentNoSlides", nil)
}

// ComponentNotifications says "Notifications"
// in the language the app shows.
func ComponentNotifications() string { return In(i18n.Locale()).ComponentNotifications() }

// ComponentNotifications says "Notifications"
// in the language of msgs.
func (msgs Messages) ComponentNotifications() string {
	return Catalog.Format(msgs.locale, "componentNotifications", nil)
}

// ComponentToastNamed says "Toast: {id}"
// in the language the app shows.
func ComponentToastNamed(id string) string { return In(i18n.Locale()).ComponentToastNamed(id) }

// ComponentToastNamed says "Toast: {id}"
// in the language of msgs.
func (msgs Messages) ComponentToastNamed(id string) string {
	return Catalog.Format(msgs.locale, "componentToastNamed", map[string]any{"id": id})
}

// ComponentInfo says "info"
// in the language the app shows.
func ComponentInfo() string { return In(i18n.Locale()).ComponentInfo() }

// ComponentInfo says "info"
// in the language of msgs.
func (msgs Messages) ComponentInfo() string {
	return Catalog.Format(msgs.locale, "componentInfo", nil)
}

// ComponentSuccess says "success"
// in the language the app shows.
func ComponentSuccess() string { return In(i18n.Locale()).ComponentSuccess() }

// ComponentSuccess says "success"
// in the language of msgs.
func (msgs Messages) ComponentSuccess() string {
	return Catalog.Format(msgs.locale, "componentSuccess", nil)
}

// ComponentWarning says "warning"
// in the language the app shows.
func ComponentWarning() string { return In(i18n.Locale()).ComponentWarning() }

// ComponentWarning says "warning"
// in the language of msgs.
func (msgs Messages) ComponentWarning() string {
	return Catalog.Format(msgs.locale, "componentWarning", nil)
}

// ComponentError says "error"
// in the language the app shows.
func ComponentError() string { return In(i18n.Locale()).ComponentError() }

// ComponentError says "error"
// in the language of msgs.
func (msgs Messages) ComponentError() string {
	return Catalog.Format(msgs.locale, "componentError", nil)
}

// ComponentMenu says "Menu"
// in the language the app shows.
func ComponentMenu() string { return In(i18n.Locale()).ComponentMenu() }

// ComponentMenu says "Menu"
// in the language of msgs.
func (msgs Messages) ComponentMenu() string {
	return Catalog.Format(msgs.locale, "componentMenu", nil)
}

// ComponentInvalidValue says "Invalid value"
// in the language the app shows.
func ComponentInvalidValue() string { return In(i18n.Locale()).ComponentInvalidValue() }

// ComponentInvalidValue says "Invalid value"
// in the language of msgs.
func (msgs Messages) ComponentInvalidValue() string {
	return Catalog.Format(msgs.locale, "componentInvalidValue", nil)
}

// ComponentReadOnly says "Read only"
// in the language the app shows.
func ComponentReadOnly() string { return In(i18n.Locale()).ComponentReadOnly() }

// ComponentReadOnly says "Read only"
// in the language of msgs.
func (msgs Messages) ComponentReadOnly() string {
	return Catalog.Format(msgs.locale, "componentReadOnly", nil)
}

// ComponentRemoveNamed says "Remove {name}"
// in the language the app shows.
func ComponentRemoveNamed(name string) string { return In(i18n.Locale()).ComponentRemoveNamed(name) }

// ComponentRemoveNamed says "Remove {name}"
// in the language of msgs.
func (msgs Messages) ComponentRemoveNamed(name string) string {
	return Catalog.Format(msgs.locale, "componentRemoveNamed", map[string]any{"name": name})
}

// ComponentNoMatches says "No matches"
// in the language the app shows.
func ComponentNoMatches() string { return In(i18n.Locale()).ComponentNoMatches() }

// ComponentNoMatches says "No matches"
// in the language of msgs.
func (msgs Messages) ComponentNoMatches() string {
	return Catalog.Format(msgs.locale, "componentNoMatches", nil)
}

// ComponentNavigationMenu says "Navigation menu"
// in the language the app shows.
func ComponentNavigationMenu() string { return In(i18n.Locale()).ComponentNavigationMenu() }

// ComponentNavigationMenu says "Navigation menu"
// in the language of msgs.
func (msgs Messages) ComponentNavigationMenu() string {
	return Catalog.Format(msgs.locale, "componentNavigationMenu", nil)
}

// ComponentCodeBlock says "Code block"
// in the language the app shows.
func ComponentCodeBlock() string { return In(i18n.Locale()).ComponentCodeBlock() }

// ComponentCodeBlock says "Code block"
// in the language of msgs.
func (msgs Messages) ComponentCodeBlock() string {
	return Catalog.Format(msgs.locale, "componentCodeBlock", nil)
}

// ComponentCodeContent says "Code content"
// in the language the app shows.
func ComponentCodeContent() string { return In(i18n.Locale()).ComponentCodeContent() }

// ComponentCodeContent says "Code content"
// in the language of msgs.
func (msgs Messages) ComponentCodeContent() string {
	return Catalog.Format(msgs.locale, "componentCodeContent", nil)
}

// ComponentAddedLine says "Added line"
// in the language the app shows.
func ComponentAddedLine() string { return In(i18n.Locale()).ComponentAddedLine() }

// ComponentAddedLine says "Added line"
// in the language of msgs.
func (msgs Messages) ComponentAddedLine() string {
	return Catalog.Format(msgs.locale, "componentAddedLine", nil)
}

// ComponentRemovedLine says "Removed line"
// in the language the app shows.
func ComponentRemovedLine() string { return In(i18n.Locale()).ComponentRemovedLine() }

// ComponentRemovedLine says "Removed line"
// in the language of msgs.
func (msgs Messages) ComponentRemovedLine() string {
	return Catalog.Format(msgs.locale, "componentRemovedLine", nil)
}

// ComponentCopyCode says "Copy code"
// in the language the app shows.
func ComponentCopyCode() string { return In(i18n.Locale()).ComponentCopyCode() }

// ComponentCopyCode says "Copy code"
// in the language of msgs.
func (msgs Messages) ComponentCopyCode() string {
	return Catalog.Format(msgs.locale, "componentCopyCode", nil)
}

// ComponentCopied says "Copied"
// in the language the app shows.
func ComponentCopied() string { return In(i18n.Locale()).ComponentCopied() }

// ComponentCopied says "Copied"
// in the language of msgs.
func (msgs Messages) ComponentCopied() string {
	return Catalog.Format(msgs.locale, "componentCopied", nil)
}

// ComponentWrapLines says "Wrap lines"
// in the language the app shows.
func ComponentWrapLines() string { return In(i18n.Locale()).ComponentWrapLines() }

// ComponentWrapLines says "Wrap lines"
// in the language of msgs.
func (msgs Messages) ComponentWrapLines() string {
	return Catalog.Format(msgs.locale, "componentWrapLines", nil)
}

// ComponentDisableLineWrap says "Disable line wrap"
// in the language the app shows.
func ComponentDisableLineWrap() string { return In(i18n.Locale()).ComponentDisableLineWrap() }

// ComponentDisableLineWrap says "Disable line wrap"
// in the language of msgs.
func (msgs Messages) ComponentDisableLineWrap() string {
	return Catalog.Format(msgs.locale, "componentDisableLineWrap", nil)
}

// ComponentDownload says "Download"
// in the language the app shows.
func ComponentDownload() string { return In(i18n.Locale()).ComponentDownload() }

// ComponentDownload says "Download"
// in the language of msgs.
func (msgs Messages) ComponentDownload() string {
	return Catalog.Format(msgs.locale, "componentDownload", nil)
}

// ComponentLoadEarlierMessages says "Load earlier messages"
// in the language the app shows.
func ComponentLoadEarlierMessages() string { return In(i18n.Locale()).ComponentLoadEarlierMessages() }

// ComponentLoadEarlierMessages says "Load earlier messages"
// in the language of msgs.
func (msgs Messages) ComponentLoadEarlierMessages() string {
	return Catalog.Format(msgs.locale, "componentLoadEarlierMessages", nil)
}

// ComponentLoadingEarlierMessages says "Loading earlier messages…"
// in the language the app shows.
func ComponentLoadingEarlierMessages() string {
	return In(i18n.Locale()).ComponentLoadingEarlierMessages()
}

// ComponentLoadingEarlierMessages says "Loading earlier messages…"
// in the language of msgs.
func (msgs Messages) ComponentLoadingEarlierMessages() string {
	return Catalog.Format(msgs.locale, "componentLoadingEarlierMessages", nil)
}

// ComponentMessages says "Messages"
// in the language the app shows.
func ComponentMessages() string { return In(i18n.Locale()).ComponentMessages() }

// ComponentMessages says "Messages"
// in the language of msgs.
func (msgs Messages) ComponentMessages() string {
	return Catalog.Format(msgs.locale, "componentMessages", nil)
}

// ComponentJumpLatestMessage says "Jump to latest message"
// in the language the app shows.
func ComponentJumpLatestMessage() string { return In(i18n.Locale()).ComponentJumpLatestMessage() }

// ComponentJumpLatestMessage says "Jump to latest message"
// in the language of msgs.
func (msgs Messages) ComponentJumpLatestMessage() string {
	return Catalog.Format(msgs.locale, "componentJumpLatestMessage", nil)
}

// ComponentReactionCount says "{label} {count}"
// in the language the app shows.
func ComponentReactionCount(label string, count int) string {
	return In(i18n.Locale()).ComponentReactionCount(label, count)
}

// ComponentReactionCount says "{label} {count}"
// in the language of msgs.
func (msgs Messages) ComponentReactionCount(label string, count int) string {
	return Catalog.Format(msgs.locale, "componentReactionCount", map[string]any{"label": label, "count": count})
}

// ComponentReactionAccessibleCount says "{label}, {count}"
// in the language the app shows.
func ComponentReactionAccessibleCount(label string, count int) string {
	return In(i18n.Locale()).ComponentReactionAccessibleCount(label, count)
}

// ComponentReactionAccessibleCount says "{label}, {count}"
// in the language of msgs.
func (msgs Messages) ComponentReactionAccessibleCount(label string, count int) string {
	return Catalog.Format(msgs.locale, "componentReactionAccessibleCount", map[string]any{"label": label, "count": count})
}

// ComponentOpenAttachment says "Open attachment {title}"
// in the language the app shows.
func ComponentOpenAttachment(title string) string {
	return In(i18n.Locale()).ComponentOpenAttachment(title)
}

// ComponentOpenAttachment says "Open attachment {title}"
// in the language of msgs.
func (msgs Messages) ComponentOpenAttachment(title string) string {
	return Catalog.Format(msgs.locale, "componentOpenAttachment", map[string]any{"title": title})
}

// ComponentRetryAttachment says "Retry attachment {title}"
// in the language the app shows.
func ComponentRetryAttachment(title string) string {
	return In(i18n.Locale()).ComponentRetryAttachment(title)
}

// ComponentRetryAttachment says "Retry attachment {title}"
// in the language of msgs.
func (msgs Messages) ComponentRetryAttachment(title string) string {
	return Catalog.Format(msgs.locale, "componentRetryAttachment", map[string]any{"title": title})
}

// ComponentRemoveAttachment says "Remove attachment {title}"
// in the language the app shows.
func ComponentRemoveAttachment(title string) string {
	return In(i18n.Locale()).ComponentRemoveAttachment(title)
}

// ComponentRemoveAttachment says "Remove attachment {title}"
// in the language of msgs.
func (msgs Messages) ComponentRemoveAttachment(title string) string {
	return Catalog.Format(msgs.locale, "componentRemoveAttachment", map[string]any{"title": title})
}

// ComponentAttachmentUploadProgress says "Attachment upload progress"
// in the language the app shows.
func ComponentAttachmentUploadProgress() string {
	return In(i18n.Locale()).ComponentAttachmentUploadProgress()
}

// ComponentAttachmentUploadProgress says "Attachment upload progress"
// in the language of msgs.
func (msgs Messages) ComponentAttachmentUploadProgress() string {
	return Catalog.Format(msgs.locale, "componentAttachmentUploadProgress", nil)
}

// ComponentNotUploaded says "Not uploaded"
// in the language the app shows.
func ComponentNotUploaded() string { return In(i18n.Locale()).ComponentNotUploaded() }

// ComponentNotUploaded says "Not uploaded"
// in the language of msgs.
func (msgs Messages) ComponentNotUploaded() string {
	return Catalog.Format(msgs.locale, "componentNotUploaded", nil)
}

// ComponentUploading says "Uploading"
// in the language the app shows.
func ComponentUploading() string { return In(i18n.Locale()).ComponentUploading() }

// ComponentUploading says "Uploading"
// in the language of msgs.
func (msgs Messages) ComponentUploading() string {
	return Catalog.Format(msgs.locale, "componentUploading", nil)
}

// ComponentProcessing says "Processing"
// in the language the app shows.
func ComponentProcessing() string { return In(i18n.Locale()).ComponentProcessing() }

// ComponentProcessing says "Processing"
// in the language of msgs.
func (msgs Messages) ComponentProcessing() string {
	return Catalog.Format(msgs.locale, "componentProcessing", nil)
}

// ComponentUploadFailed says "Upload failed"
// in the language the app shows.
func ComponentUploadFailed() string { return In(i18n.Locale()).ComponentUploadFailed() }

// ComponentUploadFailed says "Upload failed"
// in the language of msgs.
func (msgs Messages) ComponentUploadFailed() string {
	return Catalog.Format(msgs.locale, "componentUploadFailed", nil)
}

// ComponentReady says "Ready"
// in the language the app shows.
func ComponentReady() string { return In(i18n.Locale()).ComponentReady() }

// ComponentReady says "Ready"
// in the language of msgs.
func (msgs Messages) ComponentReady() string {
	return Catalog.Format(msgs.locale, "componentReady", nil)
}

// ComponentUnit says "Unit: {unit}"
// in the language the app shows.
func ComponentUnit(unit string) string { return In(i18n.Locale()).ComponentUnit(unit) }

// ComponentUnit says "Unit: {unit}"
// in the language of msgs.
func (msgs Messages) ComponentUnit(unit string) string {
	return Catalog.Format(msgs.locale, "componentUnit", map[string]any{"unit": unit})
}

// ComponentRangeValues says "{low} to {high}"
// in the language the app shows.
func ComponentRangeValues(low string, high string) string {
	return In(i18n.Locale()).ComponentRangeValues(low, high)
}

// ComponentRangeValues says "{low} to {high}"
// in the language of msgs.
func (msgs Messages) ComponentRangeValues(low string, high string) string {
	return Catalog.Format(msgs.locale, "componentRangeValues", map[string]any{"low": low, "high": high})
}

// ComponentCalendar says "Calendar"
// in the language the app shows.
func ComponentCalendar() string { return In(i18n.Locale()).ComponentCalendar() }

// ComponentCalendar says "Calendar"
// in the language of msgs.
func (msgs Messages) ComponentCalendar() string {
	return Catalog.Format(msgs.locale, "componentCalendar", nil)
}

// ComponentPreviousMonth says "Previous month"
// in the language the app shows.
func ComponentPreviousMonth() string { return In(i18n.Locale()).ComponentPreviousMonth() }

// ComponentPreviousMonth says "Previous month"
// in the language of msgs.
func (msgs Messages) ComponentPreviousMonth() string {
	return Catalog.Format(msgs.locale, "componentPreviousMonth", nil)
}

// ComponentNextMonth says "Next month"
// in the language the app shows.
func ComponentNextMonth() string { return In(i18n.Locale()).ComponentNextMonth() }

// ComponentNextMonth says "Next month"
// in the language of msgs.
func (msgs Messages) ComponentNextMonth() string {
	return Catalog.Format(msgs.locale, "componentNextMonth", nil)
}

// ComponentChooseMonth says "Choose month"
// in the language the app shows.
func ComponentChooseMonth() string { return In(i18n.Locale()).ComponentChooseMonth() }

// ComponentChooseMonth says "Choose month"
// in the language of msgs.
func (msgs Messages) ComponentChooseMonth() string {
	return Catalog.Format(msgs.locale, "componentChooseMonth", nil)
}

// ComponentChooseYear says "Choose year"
// in the language the app shows.
func ComponentChooseYear() string { return In(i18n.Locale()).ComponentChooseYear() }

// ComponentChooseYear says "Choose year"
// in the language of msgs.
func (msgs Messages) ComponentChooseYear() string {
	return Catalog.Format(msgs.locale, "componentChooseYear", nil)
}

// ComponentEnterISODate says "Enter a valid date as YYYY-MM-DD"
// in the language the app shows.
func ComponentEnterISODate() string { return In(i18n.Locale()).ComponentEnterISODate() }

// ComponentEnterISODate says "Enter a valid date as YYYY-MM-DD"
// in the language of msgs.
func (msgs Messages) ComponentEnterISODate() string {
	return Catalog.Format(msgs.locale, "componentEnterISODate", nil)
}

// ComponentDateUnavailable says "This date is unavailable"
// in the language the app shows.
func ComponentDateUnavailable() string { return In(i18n.Locale()).ComponentDateUnavailable() }

// ComponentDateUnavailable says "This date is unavailable"
// in the language of msgs.
func (msgs Messages) ComponentDateUnavailable() string {
	return Catalog.Format(msgs.locale, "componentDateUnavailable", nil)
}

// ComponentChooseAvailableRange says "Choose an ordered range containing only available dates"
// in the language the app shows.
func ComponentChooseAvailableRange() string { return In(i18n.Locale()).ComponentChooseAvailableRange() }

// ComponentChooseAvailableRange says "Choose an ordered range containing only available dates"
// in the language of msgs.
func (msgs Messages) ComponentChooseAvailableRange() string {
	return Catalog.Format(msgs.locale, "componentChooseAvailableRange", nil)
}

// ComponentDate says "Date"
// in the language the app shows.
func ComponentDate() string { return In(i18n.Locale()).ComponentDate() }

// ComponentDate says "Date"
// in the language of msgs.
func (msgs Messages) ComponentDate() string {
	return Catalog.Format(msgs.locale, "componentDate", nil)
}

// ComponentDateRangeStart says "{label} start"
// in the language the app shows.
func ComponentDateRangeStart(label string) string {
	return In(i18n.Locale()).ComponentDateRangeStart(label)
}

// ComponentDateRangeStart says "{label} start"
// in the language of msgs.
func (msgs Messages) ComponentDateRangeStart(label string) string {
	return Catalog.Format(msgs.locale, "componentDateRangeStart", map[string]any{"label": label})
}

// ComponentDateRangeEnd says "{label} end"
// in the language the app shows.
func ComponentDateRangeEnd(label string) string {
	return In(i18n.Locale()).ComponentDateRangeEnd(label)
}

// ComponentDateRangeEnd says "{label} end"
// in the language of msgs.
func (msgs Messages) ComponentDateRangeEnd(label string) string {
	return Catalog.Format(msgs.locale, "componentDateRangeEnd", map[string]any{"label": label})
}

// ComponentOpenCalendar says "Open calendar"
// in the language the app shows.
func ComponentOpenCalendar() string { return In(i18n.Locale()).ComponentOpenCalendar() }

// ComponentOpenCalendar says "Open calendar"
// in the language of msgs.
func (msgs Messages) ComponentOpenCalendar() string {
	return Catalog.Format(msgs.locale, "componentOpenCalendar", nil)
}

// ComponentClearDate says "Clear date"
// in the language the app shows.
func ComponentClearDate() string { return In(i18n.Locale()).ComponentClearDate() }

// ComponentClearDate says "Clear date"
// in the language of msgs.
func (msgs Messages) ComponentClearDate() string {
	return Catalog.Format(msgs.locale, "componentClearDate", nil)
}

// ComponentQuestionAnswerRequired says "This question needs an answer."
// in the language the app shows.
func ComponentQuestionAnswerRequired() string {
	return In(i18n.Locale()).ComponentQuestionAnswerRequired()
}

// ComponentQuestionAnswerRequired says "This question needs an answer."
// in the language of msgs.
func (msgs Messages) ComponentQuestionAnswerRequired() string {
	return Catalog.Format(msgs.locale, "componentQuestionAnswerRequired", nil)
}

// ComponentNoQuestions says "No questions."
// in the language the app shows.
func ComponentNoQuestions() string { return In(i18n.Locale()).ComponentNoQuestions() }

// ComponentNoQuestions says "No questions."
// in the language of msgs.
func (msgs Messages) ComponentNoQuestions() string {
	return Catalog.Format(msgs.locale, "componentNoQuestions", nil)
}

// ComponentQuestionOf says "Question {question} of {questions}"
// in the language the app shows.
func ComponentQuestionOf(question int, questions int) string {
	return In(i18n.Locale()).ComponentQuestionOf(question, questions)
}

// ComponentQuestionOf says "Question {question} of {questions}"
// in the language of msgs.
func (msgs Messages) ComponentQuestionOf(question int, questions int) string {
	return Catalog.Format(msgs.locale, "componentQuestionOf", map[string]any{"question": question, "questions": questions})
}

// ComponentQuestionnaireProgress says "Questionnaire progress"
// in the language the app shows.
func ComponentQuestionnaireProgress() string {
	return In(i18n.Locale()).ComponentQuestionnaireProgress()
}

// ComponentQuestionnaireProgress says "Questionnaire progress"
// in the language of msgs.
func (msgs Messages) ComponentQuestionnaireProgress() string {
	return Catalog.Format(msgs.locale, "componentQuestionnaireProgress", nil)
}

// ComponentAnotherAnswer says "Another answer"
// in the language the app shows.
func ComponentAnotherAnswer() string { return In(i18n.Locale()).ComponentAnotherAnswer() }

// ComponentAnotherAnswer says "Another answer"
// in the language of msgs.
func (msgs Messages) ComponentAnotherAnswer() string {
	return Catalog.Format(msgs.locale, "componentAnotherAnswer", nil)
}

// ComponentSkip says "Skip"
// in the language the app shows.
func ComponentSkip() string { return In(i18n.Locale()).ComponentSkip() }

// ComponentSkip says "Skip"
// in the language of msgs.
func (msgs Messages) ComponentSkip() string {
	return Catalog.Format(msgs.locale, "componentSkip", nil)
}

// ComponentSubmitAnswers says "Submit answers"
// in the language the app shows.
func ComponentSubmitAnswers() string { return In(i18n.Locale()).ComponentSubmitAnswers() }

// ComponentSubmitAnswers says "Submit answers"
// in the language of msgs.
func (msgs Messages) ComponentSubmitAnswers() string {
	return Catalog.Format(msgs.locale, "componentSubmitAnswers", nil)
}

// ComponentPolarNonnegative says "Polar charts require nonnegative values; negative values are omitted."
// in the language the app shows.
func ComponentPolarNonnegative() string { return In(i18n.Locale()).ComponentPolarNonnegative() }

// ComponentPolarNonnegative says "Polar charts require nonnegative values; negative values are omitted."
// in the language of msgs.
func (msgs Messages) ComponentPolarNonnegative() string {
	return Catalog.Format(msgs.locale, "componentPolarNonnegative", nil)
}

// ComponentRadarMinimumCategories says "Radar charts require at least three categories"
// in the language the app shows.
func ComponentRadarMinimumCategories() string {
	return In(i18n.Locale()).ComponentRadarMinimumCategories()
}

// ComponentRadarMinimumCategories says "Radar charts require at least three categories"
// in the language of msgs.
func (msgs Messages) ComponentRadarMinimumCategories() string {
	return Catalog.Format(msgs.locale, "componentRadarMinimumCategories", nil)
}

// ComponentNoVisibleSeries says "No visible series"
// in the language the app shows.
func ComponentNoVisibleSeries() string { return In(i18n.Locale()).ComponentNoVisibleSeries() }

// ComponentNoVisibleSeries says "No visible series"
// in the language of msgs.
func (msgs Messages) ComponentNoVisibleSeries() string {
	return Catalog.Format(msgs.locale, "componentNoVisibleSeries", nil)
}

// ComponentNoData says "No data"
// in the language the app shows.
func ComponentNoData() string { return In(i18n.Locale()).ComponentNoData() }

// ComponentNoData says "No data"
// in the language of msgs.
func (msgs Messages) ComponentNoData() string {
	return Catalog.Format(msgs.locale, "componentNoData", nil)
}

// ComponentAllValuesZero says "All values are zero"
// in the language the app shows.
func ComponentAllValuesZero() string { return In(i18n.Locale()).ComponentAllValuesZero() }

// ComponentAllValuesZero says "All values are zero"
// in the language of msgs.
func (msgs Messages) ComponentAllValuesZero() string {
	return Catalog.Format(msgs.locale, "componentAllValuesZero", nil)
}

// ComponentSeriesMissing says "{series}: missing"
// in the language the app shows.
func ComponentSeriesMissing(series string) string {
	return In(i18n.Locale()).ComponentSeriesMissing(series)
}

// ComponentSeriesMissing says "{series}: missing"
// in the language of msgs.
func (msgs Messages) ComponentSeriesMissing(series string) string {
	return Catalog.Format(msgs.locale, "componentSeriesMissing", map[string]any{"series": series})
}

// ComponentMissing says "missing"
// in the language the app shows.
func ComponentMissing() string { return In(i18n.Locale()).ComponentMissing() }

// ComponentMissing says "missing"
// in the language of msgs.
func (msgs Messages) ComponentMissing() string {
	return Catalog.Format(msgs.locale, "componentMissing", nil)
}

// ComponentChartValue says "{category} · {series}: {value}"
// in the language the app shows.
func ComponentChartValue(category string, series string, value string) string {
	return In(i18n.Locale()).ComponentChartValue(category, series, value)
}

// ComponentChartValue says "{category} · {series}: {value}"
// in the language of msgs.
func (msgs Messages) ComponentChartValue(category string, series string, value string) string {
	return Catalog.Format(
		msgs.locale,
		"componentChartValue",
		map[string]any{"category": category, "series": series, "value": value},
	)
}
