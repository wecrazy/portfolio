// Package pagination provides reusable server-side pagination helpers for
// Fiber + GORM. It has no dependencies on internal packages and is safe to
// import from any handler or service in the project.
package pagination

import (
	"math"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Params holds parsed query parameters for a paginated request.
type Params struct {
	Page    int
	PerPage int
	Search  string
	SortBy  string
	SortDir string
}

// Result holds computed pagination metadata, ready to be passed to templates.
type Result struct {
	Page       int
	PerPage    int
	TotalItems int
	TotalPages int
	HasPrev    bool
	HasNext    bool
	Search     string
	SortBy     string
	SortDir    string
}

// ParseParams reads pagination query parameters from a Fiber context.
//
//   - defaultSort is the column used when the sort_by param is absent or invalid.
//   - allowedSorts is the whitelist of column names accepted for sort_by.
func ParseParams(c *fiber.Ctx, defaultSort string, allowedSorts []string) Params {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	if page < 1 {
		page = 1
	}

	perPage, _ := strconv.Atoi(c.Query("per_page", "5"))
	switch perPage {
	case 5, 10, 25, 50:
	default:
		perPage = 5
	}

	search := strings.TrimSpace(c.Query("search"))

	sortBy := c.Query("sort_by", defaultSort)
	allowed := false
	for _, s := range allowedSorts {
		if s == sortBy {
			allowed = true
			break
		}
	}
	if !allowed {
		sortBy = defaultSort
	}

	sortDir := strings.ToUpper(c.Query("sort_dir", "ASC"))
	if sortDir != "ASC" && sortDir != "DESC" {
		sortDir = "ASC"
	}

	return Params{
		Page:    page,
		PerPage: perPage,
		Search:  search,
		SortBy:  sortBy,
		SortDir: sortDir,
	}
}

// Paginate applies optional full-text search, counts total rows, then applies
// ordering and offset/limit to the query. It returns a scoped *gorm.DB ready
// for a .Find() call and a Result containing all metadata needed for the
// template pagination controls.
//
//   - model should be an empty pointer to the GORM model struct (e.g. &model.Post{}).
//   - searchCols is the list of column names searched with LIKE when params.Search is non-empty.
func Paginate(db *gorm.DB, model interface{}, params Params, searchCols []string) (*gorm.DB, Result) {
	query := db.Model(model)

	if params.Search != "" && len(searchCols) > 0 {
		like := "%" + params.Search + "%"
		cond := db.Where("1 = 0")
		for _, col := range searchCols {
			cond = cond.Or(col+" LIKE ?", like)
		}
		query = query.Where(cond)
	}

	var totalItems int64
	query.Count(&totalItems)

	totalPages := int(math.Ceil(float64(totalItems) / float64(params.PerPage)))
	if totalPages < 1 {
		totalPages = 1
	}
	if params.Page > totalPages {
		params.Page = totalPages
	}

	offset := (params.Page - 1) * params.PerPage
	query = query.Order(params.SortBy + " " + params.SortDir).
		Offset(offset).
		Limit(params.PerPage)

	return query, Result{
		Page:       params.Page,
		PerPage:    params.PerPage,
		TotalItems: int(totalItems),
		TotalPages: totalPages,
		HasPrev:    params.Page > 1,
		HasNext:    params.Page < totalPages,
		Search:     params.Search,
		SortBy:     params.SortBy,
		SortDir:    params.SortDir,
	}
}
