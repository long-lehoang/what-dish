package notion

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ---------------------------------------------------------------------------
// Helpers to build Notion objects for tests
// ---------------------------------------------------------------------------

func ptrFloat(v float64) *float64 { return &v }
func ptrString(v string) *string  { return &v }

func richText(text string) []RichText {
	return []RichText{{PlainText: text}}
}

func titleProp(text string) Property {
	return Property{Type: "title", Title: richText(text)}
}

func richTextProp(text string) Property {
	return Property{Type: "rich_text", RichText: richText(text)}
}

func selectProp(name string) Property {
	return Property{Type: "select", Select: &SelectValue{Name: name}}
}

func numberProp(v float64) Property {
	return Property{Type: "number", Number: &v}
}

func multiSelectProp(names ...string) Property {
	vals := make([]SelectValue, len(names))
	for i, n := range names {
		vals[i] = SelectValue{Name: n}
	}
	return Property{Type: "multi_select", MultiSelect: vals}
}

func fileProp(url string) Property {
	return Property{
		Type: "files",
		Files: []FileObject{
			{Type: "external", External: &FileURL{URL: url}},
		},
	}
}

func heading2Block(text string) Block {
	return Block{
		Type:     "heading_2",
		Heading2: &HeadingBlock{RichText: richText(text)},
	}
}

func heading3Block(text string) Block {
	return Block{
		Type:     "heading_3",
		Heading3: &HeadingBlock{RichText: richText(text)},
	}
}

func bulletBlock(text string) Block {
	return Block{
		Type:             "bulleted_list_item",
		BulletedListItem: &ListItemBlock{RichText: richText(text)},
	}
}

func numberedBlock(text string) Block {
	return Block{
		Type:             "numbered_list_item",
		NumberedListItem: &ListItemBlock{RichText: richText(text)},
	}
}

func imageBlock(url string) Block {
	return Block{
		Type:  "image",
		Image: &ImageBlock{Type: "external", External: &FileURL{URL: url}},
	}
}

func calloutBlock(text string) Block {
	return Block{
		Type:    "callout",
		Callout: &CalloutBlock{RichText: richText(text)},
	}
}

// ---------------------------------------------------------------------------
// ParsePageToDish
// ---------------------------------------------------------------------------

func TestParsePageToDish_FullProperties(t *testing.T) {
	page := Page{
		ID:  "page-123",
		URL: "https://notion.so/page-123",
		Properties: map[string]Property{
			"Name":        titleProp("Pho Bo"),
			"Slug":        richTextProp("pho-bo"),
			"Description": richTextProp("Vietnamese beef noodle soup"),
			"Difficulty":  selectProp("MEDIUM"),
			"Prep Time":   numberProp(15),
			"Cook Time":   numberProp(60),
			"Servings":    numberProp(4),
			"Cover":       fileProp("https://img.example.com/pho.jpg"),
		},
	}

	dish := ParsePageToDish(page)

	assert.Equal(t, "page-123", dish.ExternalID)
	assert.Equal(t, "Pho Bo", dish.Name)
	assert.Equal(t, "pho-bo", dish.Slug)
	assert.NotNil(t, dish.Description)
	assert.Equal(t, "Vietnamese beef noodle soup", *dish.Description)
	assert.Equal(t, "MEDIUM", dish.Difficulty)
	assert.NotNil(t, dish.PrepTime)
	assert.Equal(t, 15, *dish.PrepTime)
	assert.NotNil(t, dish.CookTime)
	assert.Equal(t, 60, *dish.CookTime)
	assert.Equal(t, 4, dish.Servings)
	assert.NotNil(t, dish.ImageURL)
	assert.Equal(t, "https://img.example.com/pho.jpg", *dish.ImageURL)
	assert.NotNil(t, dish.SourceURL)
	assert.Equal(t, "https://notion.so/page-123", *dish.SourceURL)
	assert.Equal(t, "PUBLISHED", dish.Status)
}

func TestParsePageToDish_MinimalProperties(t *testing.T) {
	page := Page{
		ID:  "page-456",
		URL: "https://notion.so/page-456",
		Properties: map[string]Property{
			"Name": titleProp("Banh Mi"),
		},
	}

	dish := ParsePageToDish(page)

	assert.Equal(t, "page-456", dish.ExternalID)
	assert.Equal(t, "Banh Mi", dish.Name)
	assert.Equal(t, "banh-mi", dish.Slug) // auto-generated from name
	assert.Nil(t, dish.Description)
	assert.Nil(t, dish.PrepTime)
	assert.Nil(t, dish.CookTime)
}

func TestParsePageToDish_DefaultServings(t *testing.T) {
	page := Page{
		ID:         "page-789",
		URL:        "https://notion.so/page-789",
		Properties: map[string]Property{"Name": titleProp("Com Tam")},
	}

	dish := ParsePageToDish(page)
	assert.Equal(t, 2, dish.Servings, "default servings should be 2")
}

func TestParsePageToDish_DefaultDifficulty(t *testing.T) {
	page := Page{
		ID:         "page-abc",
		URL:        "https://notion.so/page-abc",
		Properties: map[string]Property{"Name": titleProp("Bun Cha")},
	}

	dish := ParsePageToDish(page)
	assert.Equal(t, "EASY", dish.Difficulty, "default difficulty should be EASY")
}

func TestParsePageToDish_MissingFields(t *testing.T) {
	page := Page{
		ID:         "page-empty",
		URL:        "https://notion.so/page-empty",
		Properties: map[string]Property{},
	}

	dish := ParsePageToDish(page)

	assert.Equal(t, "page-empty", dish.ExternalID)
	assert.Equal(t, "", dish.Name)
	assert.Equal(t, 2, dish.Servings)
	assert.Equal(t, "EASY", dish.Difficulty)
	assert.Equal(t, "PUBLISHED", dish.Status)
}

func TestParsePageToDish_CoverFromPageCover(t *testing.T) {
	page := Page{
		ID:         "page-cover",
		URL:        "https://notion.so/page-cover",
		Properties: map[string]Property{"Name": titleProp("Goi Cuon")},
		Cover: &FileCover{
			Type:     "external",
			External: &FileURL{URL: "https://img.example.com/cover.jpg"},
		},
	}

	dish := ParsePageToDish(page)
	assert.NotNil(t, dish.ImageURL)
	assert.Equal(t, "https://img.example.com/cover.jpg", *dish.ImageURL)
}

func TestParsePageToDish_PropertyCoverTakesPrecedenceOverPageCover(t *testing.T) {
	page := Page{
		ID:  "page-dual-cover",
		URL: "https://notion.so/page-dual-cover",
		Properties: map[string]Property{
			"Name":  titleProp("Bun Bo Hue"),
			"Cover": fileProp("https://img.example.com/property-cover.jpg"),
		},
		Cover: &FileCover{
			Type:     "external",
			External: &FileURL{URL: "https://img.example.com/page-cover.jpg"},
		},
	}

	dish := ParsePageToDish(page)
	assert.NotNil(t, dish.ImageURL)
	assert.Equal(t, "https://img.example.com/property-cover.jpg", *dish.ImageURL)
}

// ---------------------------------------------------------------------------
// ParseBlocksToContent
// ---------------------------------------------------------------------------

func TestParseBlocksToContent_IngredientsAndSteps(t *testing.T) {
	blocks := []Block{
		heading2Block("Nguyên liệu"),
		bulletBlock("200g thịt bò"),
		bulletBlock("1 củ hành tây"),
		heading2Block("Cách làm"),
		numberedBlock("Sơ chế nguyên liệu"),
		numberedBlock("Xào thịt bò (5 phút)"),
	}

	ingredients, steps := ParseBlocksToContent(blocks)

	assert.Len(t, ingredients, 2)
	assert.Len(t, steps, 2)

	assert.Equal(t, 1, steps[0].StepNumber)
	assert.Equal(t, "Sơ chế nguyên liệu", steps[0].Description)

	assert.Equal(t, 2, steps[1].StepNumber)
	assert.Contains(t, steps[1].Description, "5 phút")
}

func TestParseBlocksToContent_DurationExtraction(t *testing.T) {
	blocks := []Block{
		heading2Block("Cách làm"),
		numberedBlock("Ninh xương (30 phút)"),
		numberedBlock("Nêm gia vị"),
	}

	_, steps := ParseBlocksToContent(blocks)

	assert.Len(t, steps, 2)
	assert.NotNil(t, steps[0].Duration)
	assert.Equal(t, 30, *steps[0].Duration)
	assert.Nil(t, steps[1].Duration)
}

func TestParseBlocksToContent_ImageAttachedToStep(t *testing.T) {
	blocks := []Block{
		heading2Block("Cách làm"),
		numberedBlock("Trộn đều"),
		imageBlock("https://img.example.com/step1.jpg"),
	}

	_, steps := ParseBlocksToContent(blocks)

	assert.Len(t, steps, 1)
	assert.NotNil(t, steps[0].ImageURL)
	assert.Equal(t, "https://img.example.com/step1.jpg", *steps[0].ImageURL)
}

func TestParseBlocksToContent_CalloutBecomesTip(t *testing.T) {
	blocks := []Block{
		heading2Block("Cách làm"),
		numberedBlock("Nấu canh"),
		calloutBlock("Mẹo: thêm chút đường để dậy vị"),
	}

	_, steps := ParseBlocksToContent(blocks)

	assert.Len(t, steps, 1)
	assert.NotNil(t, steps[0].Title)
	assert.Contains(t, *steps[0].Title, "Mẹo")
}

func TestParseBlocksToContent_Heading3AlsoWorks(t *testing.T) {
	blocks := []Block{
		heading3Block("Ingredients"),
		bulletBlock("100ml soy sauce"),
		heading3Block("Steps"),
		numberedBlock("Mix everything"),
	}

	ingredients, steps := ParseBlocksToContent(blocks)

	assert.Len(t, ingredients, 1)
	assert.Len(t, steps, 1)
}

func TestParseBlocksToContent_EmptyBlocks(t *testing.T) {
	ingredients, steps := ParseBlocksToContent([]Block{})
	assert.Empty(t, ingredients)
	assert.Empty(t, steps)
}

func TestParseBlocksToContent_IngredientParsing(t *testing.T) {
	// parseIngredientText is unexported; test it indirectly.
	blocks := []Block{
		heading2Block("Nguyên liệu"),
		bulletBlock("200g thịt bò"),
		bulletBlock("muối"),
		bulletBlock("1.5kg gạo"),
	}

	ingredients, _ := ParseBlocksToContent(blocks)

	assert.Len(t, ingredients, 3)

	// "200g thịt bò" -> amount=200, unit=g, name=thịt bò
	assert.NotNil(t, ingredients[0].Amount)
	assert.Equal(t, 200.0, *ingredients[0].Amount)
	assert.NotNil(t, ingredients[0].Unit)
	assert.Equal(t, "g", *ingredients[0].Unit)
	assert.Equal(t, "thịt bò", ingredients[0].Name)

	// "muối" -> no amount/unit parsed, name is the whole text
	assert.Nil(t, ingredients[1].Amount)
	assert.Nil(t, ingredients[1].Unit)
	assert.Equal(t, "muối", ingredients[1].Name)

	// "1.5kg gạo" -> amount=1.5, unit=kg, name=gạo
	assert.NotNil(t, ingredients[2].Amount)
	assert.Equal(t, 1.5, *ingredients[2].Amount)
	assert.NotNil(t, ingredients[2].Unit)
	assert.Equal(t, "kg", *ingredients[2].Unit)
	assert.Equal(t, "gạo", ingredients[2].Name)
}

func TestParseBlocksToContent_SortOrder(t *testing.T) {
	blocks := []Block{
		heading2Block("Nguyên liệu"),
		bulletBlock("a"),
		bulletBlock("b"),
		bulletBlock("c"),
		heading2Block("Cách làm"),
		numberedBlock("step 1"),
		numberedBlock("step 2"),
	}

	ingredients, steps := ParseBlocksToContent(blocks)

	for i, ing := range ingredients {
		assert.Equal(t, i, ing.SortOrder)
	}
	for i, s := range steps {
		assert.Equal(t, i+1, s.SortOrder)
		assert.Equal(t, i+1, s.StepNumber)
	}
}

func TestParseBlocksToContent_NumberedListOutsideStepsSection(t *testing.T) {
	// Numbered list items before any heading should still be captured as steps
	// because the code uses section == "" as a fallback for numbered_list_item.
	blocks := []Block{
		numberedBlock("Do this first"),
		numberedBlock("Then do this"),
	}

	_, steps := ParseBlocksToContent(blocks)
	assert.Len(t, steps, 2)
}

// ---------------------------------------------------------------------------
// GetTagNames
// ---------------------------------------------------------------------------

func TestGetTagNames_WithTags(t *testing.T) {
	page := Page{
		Properties: map[string]Property{
			"Tags": multiSelectProp("quick", "healthy", "spicy"),
		},
	}

	tags := GetTagNames(page)
	assert.Equal(t, []string{"quick", "healthy", "spicy"}, tags)
}

func TestGetTagNames_NoTagsProperty(t *testing.T) {
	page := Page{Properties: map[string]Property{}}
	tags := GetTagNames(page)
	assert.Nil(t, tags)
}

func TestGetTagNames_EmptyMultiSelect(t *testing.T) {
	page := Page{
		Properties: map[string]Property{
			"Tags": {Type: "multi_select", MultiSelect: []SelectValue{}},
		},
	}

	tags := GetTagNames(page)
	assert.Empty(t, tags)
}

// ---------------------------------------------------------------------------
// GetCategoryNames
// ---------------------------------------------------------------------------

func TestGetCategoryNames_AllCategories(t *testing.T) {
	page := Page{
		Properties: map[string]Property{
			"Category":        selectProp("Canh"),
			"Region":          selectProp("Bắc"),
			"Main Ingredient": selectProp("Gà"),
			"Meal Type":       selectProp("Trưa"),
		},
	}

	cats := GetCategoryNames(page)

	assert.Equal(t, "Canh", cats["DISH_TYPE"])
	assert.Equal(t, "Bắc", cats["REGION"])
	assert.Equal(t, "Gà", cats["MAIN_INGREDIENT"])
	assert.Equal(t, "Trưa", cats["MEAL_TYPE"])
	assert.Len(t, cats, 4)
}

func TestGetCategoryNames_PartialCategories(t *testing.T) {
	page := Page{
		Properties: map[string]Property{
			"Category": selectProp("Xào"),
		},
	}

	cats := GetCategoryNames(page)
	assert.Equal(t, "Xào", cats["DISH_TYPE"])
	assert.Len(t, cats, 1)
}

func TestGetCategoryNames_NoCategories(t *testing.T) {
	page := Page{Properties: map[string]Property{}}
	cats := GetCategoryNames(page)
	assert.Empty(t, cats)
}

func TestGetCategoryNames_NilSelect(t *testing.T) {
	page := Page{
		Properties: map[string]Property{
			"Category": {Type: "select", Select: nil},
		},
	}

	cats := GetCategoryNames(page)
	assert.Empty(t, cats)
}
