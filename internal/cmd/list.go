package cmd

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/urfave/cli/v2"

	"github.com/h1s97x/gh-follow/internal/models"
	"github.com/h1s97x/gh-follow/internal/storage"
)

// List handles the follow list command
func List(c *cli.Context) error {
	// Get flags
	format := c.String("format")
	sortField := c.String("sort")
	order := c.String("order")
	limit := c.Int("limit")
	filter := c.String("filter")
	tagFilter := c.String("tag")
	dateFrom := c.String("date-from")
	dateTo := c.String("date-to")

	// Initialize storage
	st := storage.NewStorage(storage.DefaultStoragePath())

	// Load follow list
	list, err := st.Load()
	if err != nil {
		return fmt.Errorf("failed to load follow list: %w", err)
	}

	if len(list.Follows) == 0 {
		fmt.Println("Your follow list is empty")
		return nil
	}

	// Apply filters
	list.Follows = applyFilters(list.Follows, filter, tagFilter, dateFrom, dateTo)

	if len(list.Follows) == 0 {
		fmt.Println("No users match the specified filters")
		return nil
	}

	// Sort
	sortFollows(list.Follows, sortField, order)

	// Apply limit
	if limit > 0 && limit < len(list.Follows) {
		list.Follows = list.Follows[:limit]
	}

	// Output based on format
	switch format {
	case "json":
		return outputJSON(list)
	case "csv":
		return outputCSV(list)
	case "simple":
		return outputSimple(list)
	case "table":
		fallthrough
	default:
		return outputTable(list)
	}
}

// applyFilters applies all filters to the follow list
func applyFilters(follows []models.Follow, usernameFilter, tagFilter, dateFrom, dateTo string) []models.Follow {
	var result []models.Follow

	for _, f := range follows {
		// Username filter
		if usernameFilter != "" && !strings.Contains(strings.ToLower(f.Username), strings.ToLower(usernameFilter)) {
			continue
		}

		// Tag filter
		if tagFilter != "" && !containsTag(f.Tags, tagFilter) {
			continue
		}

		// Date from filter
		if dateFrom != "" {
			fromDate, err := time.Parse("2006-01-02", dateFrom)
			if err == nil && f.FollowedAt.Before(fromDate) {
				continue
			}
		}

		// Date to filter
		if dateTo != "" {
			toDate, err := time.Parse("2006-01-02", dateTo)
			if err == nil && f.FollowedAt.After(toDate.Add(24*time.Hour)) {
				continue
			}
		}

		result = append(result, f)
	}

	return result
}

// containsTag checks if a tag exists in the tag list (case-insensitive)
func containsTag(tags []string, target string) bool {
	target = strings.ToLower(target)
	for _, tag := range tags {
		if strings.ToLower(tag) == target {
			return true
		}
	}
	return false
}

// sortFollows sorts the follows slice based on field and order
func sortFollows(follows []models.Follow, sortField, order string) {
	sort.Slice(follows, func(i, j int) bool {
		var less bool
		switch sortField {
		case "name":
			less = strings.ToLower(follows[i].Username) < strings.ToLower(follows[j].Username)
		case "notes":
			less = follows[i].Notes < follows[j].Notes
		case "date":
			fallthrough
		default:
			less = follows[i].FollowedAt.Before(follows[j].FollowedAt)
		}

		if order == "desc" {
			return !less
		}
		return less
	})
}

// outputTable outputs the list in table format
func outputTable(list *models.FollowList) error {
	fmt.Printf("\n╔══════════════════════════════════════════════════════════════════════╗\n")
	fmt.Printf("║  Follow List (%d users)                                                \n", len(list.Follows))
	fmt.Printf("╠══════════════════════════════════════════════════════════════════════╣\n")
	fmt.Printf("║  %-25s %-18s %-20s ║\n", "Username", "Followed At", "Tags")
	fmt.Printf("╠══════════════════════════════════════════════════════════════════════╣\n")

	for _, f := range list.Follows {
		tags := strings.Join(f.Tags, ", ")
		if len(tags) > 20 {
			tags = tags[:17] + "..."
		}
		fmt.Printf("║  %-25s %-18s %-20s ║\n",
			f.Username,
			f.FollowedAt.Format("2006-01-02"),
			tags,
		)
	}

	fmt.Printf("╚══════════════════════════════════════════════════════════════════════╝\n")
	fmt.Printf("Last updated: %s\n", list.UpdatedAt.Format("2006-01-02 15:04:05"))

	return nil
}

// outputSimple outputs the list in simple format (just usernames)
func outputSimple(list *models.FollowList) error {
	for _, f := range list.Follows {
		fmt.Println(f.Username)
	}
	return nil
}

// outputJSON outputs the list in JSON format
func outputJSON(list *models.FollowList) error {
	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}

// outputCSV outputs the list in CSV format
func outputCSV(list *models.FollowList) error {
	fmt.Println("username,followed_at,notes,tags")
	for _, f := range list.Follows {
		tags := strings.Join(f.Tags, ";")
		// Escape commas and quotes in notes
		notes := f.Notes
		if strings.Contains(notes, ",") || strings.Contains(notes, "\"") {
			notes = "\"" + strings.ReplaceAll(notes, "\"", "\"\"") + "\""
		}
		fmt.Printf("%s,%s,%s,%s\n",
			f.Username,
			f.FollowedAt.Format(time.RFC3339),
			notes,
			tags,
		)
	}
	return nil
}
