package glance

import (
	"fmt"
	"html/template"
)

var bookmarksWidgetTemplate = mustParseTemplate("bookmarks.html", "widget-base.html")

type bookmarksWidget struct {
	widgetBase `yaml:",inline"`
	cachedHTML template.HTML `yaml:"-"`
	Groups     []struct {
		Title     string         `yaml:"title"`
		Color     *hslColorField `yaml:"color"`
		SameTab   bool           `yaml:"same-tab"`
		HideArrow bool           `yaml:"hide-arrow"`
		Target    string         `yaml:"target"`
		Links     []struct {
			Title       string          `yaml:"title"`
			URL         string          `yaml:"url"`
			Description string          `yaml:"description"`
			Icon        customIconField `yaml:"icon"`
			// we need a pointer to bool to know whether a value was provided,
			// however there's no way to dereference a pointer in a template so
			// {{ if not .SameTab }} would return true for any non-nil pointer
			// which leaves us with no way of checking if the value is true or
			// false, hence the duplicated fields below
			SameTabRaw   *bool  `yaml:"same-tab"`
			SameTab      bool   `yaml:"-"`
			HideArrowRaw *bool  `yaml:"hide-arrow"`
			HideArrow    bool   `yaml:"-"`
			Target       string `yaml:"target"`
			// Icon-only links shown on the link's row, right-aligned. They
			// open the way the link does (same-tab, target).
			SecondaryLinks []struct {
				Title string          `yaml:"title"`
				URL   string          `yaml:"url"`
				Icon  customIconField `yaml:"icon"`
			} `yaml:"secondary-links"`
		} `yaml:"links"`
	} `yaml:"groups"`

	TargetOverride   string        `yaml:"-"`
	cachedNewTabHTML template.HTML `yaml:"-"`
}

func (widget *bookmarksWidget) initialize() error {
	widget.withTitle("Bookmarks").withError(nil)

	for g := range widget.Groups {
		group := &widget.Groups[g]
		for l := range group.Links {
			link := &group.Links[l]
			if link.SameTabRaw == nil {
				link.SameTab = group.SameTab
			} else {
				link.SameTab = *link.SameTabRaw
			}

			if link.HideArrowRaw == nil {
				link.HideArrow = group.HideArrow
			} else {
				link.HideArrow = *link.HideArrowRaw
			}

			for _, secondary := range link.SecondaryLinks {
				if secondary.Icon.URL == "" {
					return fmt.Errorf("bookmarks: secondary link %q of %q has no icon", secondary.Title, link.Title)
				}
			}

			if link.Target == "" {
				if group.Target != "" {
					link.Target = group.Target
				} else {
					if link.SameTab {
						link.Target = ""
					} else {
						link.Target = "_blank"
					}
				}
			}
		}
	}

	widget.cachedHTML = widget.renderTemplate(widget, bookmarksWidgetTemplate)

	widget.TargetOverride = "_top"
	widget.cachedNewTabHTML = widget.renderTemplate(widget, bookmarksWidgetTemplate)
	widget.TargetOverride = ""

	return nil
}

func (widget *bookmarksWidget) renderNewTab() template.HTML {
	return widget.cachedNewTabHTML
}

func (widget *bookmarksWidget) Render() template.HTML {
	return widget.cachedHTML
}
