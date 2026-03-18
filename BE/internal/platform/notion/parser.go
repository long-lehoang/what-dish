package notion

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

var durationRegex = regexp.MustCompile(`\((\d+)\s*phút\)`)

// Parsed types — notion package's own models to avoid importing recipe package.

type ParsedDish struct {
	ExternalID  string
	Name        string
	Slug        string
	Description *string
	ImageURL    *string
	PrepTime    *int
	CookTime    *int
	Servings    int
	Difficulty  string
	Status      string
	SourceURL   *string
}

type ParsedIngredient struct {
	Name      string
	Amount    *float64
	Unit      *string
	SortOrder int
}

type ParsedStep struct {
	StepNumber  int
	Title       *string
	Description string
	ImageURL    *string
	Duration    *int
	SortOrder   int
}

// ParsePageToDish converts a Notion page's properties into a ParsedDish.
func ParsePageToDish(page Page) ParsedDish {
	dish := ParsedDish{
		ExternalID: page.ID,
		SourceURL:  &page.URL,
		Status:     "PUBLISHED",
	}

	for key, prop := range page.Properties {
		switch key {
		case "Name":
			dish.Name = getTitle(prop)
		case "Slug":
			dish.Slug = getRichText(prop)
		case "Description":
			dish.Description = getRichTextPtr(prop)
		case "Difficulty":
			if prop.Select != nil {
				dish.Difficulty = prop.Select.Name
			}
		case "Prep Time":
			if prop.Number != nil {
				v := int(*prop.Number)
				dish.PrepTime = &v
			}
		case "Cook Time":
			if prop.Number != nil {
				v := int(*prop.Number)
				dish.CookTime = &v
			}
		case "Servings":
			if prop.Number != nil {
				dish.Servings = int(*prop.Number)
			}
		case "Cover":
			if url := getFileURL(prop); url != "" {
				dish.ImageURL = &url
			}
		}
	}

	// Get cover image from page cover if not set via properties
	if dish.ImageURL == nil && page.Cover != nil {
		if url := getCoverURL(page.Cover); url != "" {
			dish.ImageURL = &url
		}
	}

	// Generate slug from name if empty
	if dish.Slug == "" && dish.Name != "" {
		dish.Slug = slugify(dish.Name)
	}

	if dish.Servings == 0 {
		dish.Servings = 2
	}
	if dish.Difficulty == "" {
		dish.Difficulty = "EASY"
	}

	return dish
}

// ParseBlocksToContent parses Notion blocks into ingredients and steps.
func ParseBlocksToContent(blocks []Block) ([]ParsedIngredient, []ParsedStep) {
	var ingredients []ParsedIngredient
	var steps []ParsedStep

	section := "" // current section: "ingredients" or "steps"
	stepNum := 0

	for _, block := range blocks {
		switch block.Type {
		case "heading_2", "heading_3":
			heading := getBlockText(block)
			lower := strings.ToLower(heading)
			if strings.Contains(lower, "nguyên liệu") || strings.Contains(lower, "ingredient") {
				section = "ingredients"
			} else if strings.Contains(lower, "cách làm") || strings.Contains(lower, "hướng dẫn") || strings.Contains(lower, "step") || strings.Contains(lower, "bước") {
				section = "steps"
				stepNum = 0
			}

		case "bulleted_list_item":
			if section == "ingredients" {
				text := getListItemText(block)
				if text != "" {
					ing := parseIngredientText(text)
					ing.SortOrder = len(ingredients)
					ingredients = append(ingredients, ing)
				}
			}

		case "numbered_list_item":
			if section == "steps" || section == "" {
				text := getListItemText(block)
				if text != "" {
					stepNum++
					step := ParsedStep{
						StepNumber:  stepNum,
						Description: text,
						SortOrder:   stepNum,
					}
					if matches := durationRegex.FindStringSubmatch(text); len(matches) > 1 {
						if d, err := strconv.Atoi(matches[1]); err == nil {
							step.Duration = &d
						}
					}
					steps = append(steps, step)
				}
			}

		case "image":
			if section == "steps" && len(steps) > 0 {
				if url := getImageURL(block); url != "" {
					steps[len(steps)-1].ImageURL = &url
				}
			}

		case "callout":
			if len(steps) > 0 {
				tipText := getCalloutText(block)
				if tipText != "" {
					title := tipText
					steps[len(steps)-1].Title = &title
				}
			}
		}
	}

	return ingredients, steps
}

// GetTagNames extracts multi-select tag names from a Notion page.
func GetTagNames(page Page) []string {
	prop, ok := page.Properties["Tags"]
	if !ok {
		return nil
	}
	names := make([]string, len(prop.MultiSelect))
	for i, s := range prop.MultiSelect {
		names[i] = s.Name
	}
	return names
}

// GetCategoryNames extracts category select values.
func GetCategoryNames(page Page) map[string]string {
	cats := make(map[string]string)

	if prop, ok := page.Properties["Category"]; ok && prop.Select != nil {
		cats["DISH_TYPE"] = prop.Select.Name
	}
	if prop, ok := page.Properties["Region"]; ok && prop.Select != nil {
		cats["REGION"] = prop.Select.Name
	}
	if prop, ok := page.Properties["Main Ingredient"]; ok && prop.Select != nil {
		cats["MAIN_INGREDIENT"] = prop.Select.Name
	}
	if prop, ok := page.Properties["Meal Type"]; ok && prop.Select != nil {
		cats["MEAL_TYPE"] = prop.Select.Name
	}

	return cats
}

func parseIngredientText(text string) ParsedIngredient {
	ing := ParsedIngredient{
		Name: text,
	}

	re := regexp.MustCompile(`^(\d+(?:\.\d+)?)\s*([a-zA-Z]+)?\s+(.+)$`)
	if matches := re.FindStringSubmatch(text); len(matches) > 3 {
		if amount, err := strconv.ParseFloat(matches[1], 64); err == nil {
			ing.Amount = &amount
			if matches[2] != "" {
				ing.Unit = &matches[2]
			}
			ing.Name = strings.TrimSpace(matches[3])
		}
	}

	return ing
}

// Helper functions

func getTitle(prop Property) string {
	if len(prop.Title) > 0 {
		return prop.Title[0].PlainText
	}
	return ""
}

func getRichText(prop Property) string {
	if len(prop.RichText) > 0 {
		var sb strings.Builder
		for _, rt := range prop.RichText {
			sb.WriteString(rt.PlainText)
		}
		return sb.String()
	}
	return ""
}

func getRichTextPtr(prop Property) *string {
	text := getRichText(prop)
	if text == "" {
		return nil
	}
	return &text
}

func getFileURL(prop Property) string {
	if len(prop.Files) == 0 {
		return ""
	}
	f := prop.Files[0]
	if f.File != nil {
		return f.File.URL
	}
	if f.External != nil {
		return f.External.URL
	}
	return ""
}

func getCoverURL(cover *FileCover) string {
	if cover == nil {
		return ""
	}
	if cover.File != nil {
		return cover.File.URL
	}
	if cover.External != nil {
		return cover.External.URL
	}
	return ""
}

func getBlockText(block Block) string {
	switch block.Type {
	case "heading_2":
		if block.Heading2 != nil && len(block.Heading2.RichText) > 0 {
			return block.Heading2.RichText[0].PlainText
		}
	case "heading_3":
		if block.Heading3 != nil && len(block.Heading3.RichText) > 0 {
			return block.Heading3.RichText[0].PlainText
		}
	}
	return ""
}

func getListItemText(block Block) string {
	switch block.Type {
	case "bulleted_list_item":
		if block.BulletedListItem != nil && len(block.BulletedListItem.RichText) > 0 {
			var sb strings.Builder
			for _, rt := range block.BulletedListItem.RichText {
				sb.WriteString(rt.PlainText)
			}
			return sb.String()
		}
	case "numbered_list_item":
		if block.NumberedListItem != nil && len(block.NumberedListItem.RichText) > 0 {
			var sb strings.Builder
			for _, rt := range block.NumberedListItem.RichText {
				sb.WriteString(rt.PlainText)
			}
			return sb.String()
		}
	}
	return ""
}

func getCalloutText(block Block) string {
	if block.Callout != nil && len(block.Callout.RichText) > 0 {
		var sb strings.Builder
		for _, rt := range block.Callout.RichText {
			sb.WriteString(rt.PlainText)
		}
		return sb.String()
	}
	return ""
}

func getImageURL(block Block) string {
	if block.Image == nil {
		return ""
	}
	if block.Image.File != nil {
		return block.Image.File.URL
	}
	if block.Image.External != nil {
		return block.Image.External.URL
	}
	return ""
}

func slugify(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	re := regexp.MustCompile(`[^a-z0-9\-]`)
	s = re.ReplaceAllString(s, "")
	re = regexp.MustCompile(`-+`)
	s = re.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	if s == "" {
		return uuid.New().String()[:8]
	}
	return s
}
