package karyawan

import (
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/arie741/go-karyawan/internal/karyawan/functions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	repo *Repository
}

type PageInfo struct {
	Number int
	Offset int
	Active bool
}

func NewHandler(repo *Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) List(c *gin.Context) {
	offset := 0
	limit := 10
	if c.Query("offset") != "" {
		offsetValue, err := strconv.Atoi(c.Query("offset"))
		if err != nil {
			panic(err)
		}
		offset = offsetValue
	}

	if c.Query("limit") != "" {
		limitValue, err := strconv.Atoi(c.Query("limit"))
		if err != nil {
			panic(err)
		}
		limit = limitValue
	}

	filter := functions.Filter{
		Name:         c.Query("name"),
		Position:     c.Query("position"),
		SalaryMin:    parseIntQuery(c, "salaryMin"),
		SalaryMax:    parseIntQuery(c, "salaryMax"),
		BirthDateOp:  c.Query("birthDateOp"),
		BirthDate:    parseDateQuery(c, "birthDate"),
		BirthDateEnd: parseDateQuery(c, "birthDateEnd"),
		JoinedOp:     c.Query("joinedOp"),
		Joined:       parseDateQuery(c, "joined"),
		JoinedEnd:    parseDateQuery(c, "joinedEnd"),
	}

	karyawans, err := h.repo.GetAll(offset, limit, filter.ToBSON())
	if err != nil {
		panic(err)
	}

	totalDocuments, err := h.repo.CountAll(filter.ToBSON())
	if err != nil {
		panic(err)
	}

	totalPages := int(math.Ceil(float64(totalDocuments) / float64(limit)))
	currentPage := offset/limit + 1

	var pages []PageInfo
	for i := 0; i < totalPages; i++ {
		pages = append(pages, PageInfo{
			Number: i + 1,
			Offset: i * limit,
			Active: i+1 == currentPage,
		})
	}

	c.HTML(http.StatusOK, "pages/index.html", gin.H{
		"karyawans":   karyawans,
		"pages":       pages,
		"limit":       limit,
		"offset":      offset,
		"filter":      filter,
		"filterQuery": filter.QueryString(),
	})
}

func parseIntQuery(c *gin.Context, key string) *int {
	value := c.Query(key)
	if value == "" {
		return nil
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		panic(err)
	}

	return &parsed
}

func parseDateQuery(c *gin.Context, key string) *time.Time {
	value := c.Query(key)
	if value == "" {
		return nil
	}

	parsed, err := time.Parse("2006-01-02", value)
	if err != nil {
		panic(err)
	}

	return &parsed
}

func (h *Handler) FindById(c *gin.Context) {
	id := c.Param("id")

	karyawan, err := h.repo.FindById(id)
	if err != nil {
		panic(err)
	}

	c.HTML(http.StatusOK, "pages/karyawan.html", gin.H{
		"karyawan": karyawan,
	})
}

func (h *Handler) Edit(c *gin.Context) {
	birthDate, birthDateErr := time.Parse("2006-01-02", c.PostForm("birthDate"))
	if birthDateErr != nil {
		panic(birthDateErr)
	}

	salary, salaryErr := strconv.Atoi(c.PostForm("salary"))
	if salaryErr != nil {
		panic(salaryErr)
	}

	joined, joinedErr := time.Parse("2006-01-02", c.PostForm("joined"))
	if joinedErr != nil {
		panic(joinedErr)
	}

	karyawan := Karyawan{
		Id:        c.Param("id"),
		Name:      c.PostForm("name"),
		BirthDate: birthDate,
		Salary:    salary,
		Position:  c.PostForm("position"),
		Joined:    joined,
	}

	updatedKaryawan, err := h.repo.Update(&karyawan)
	if err != nil {
		panic(err)
	}

	c.HTML(http.StatusOK, "karyawan-info-element", gin.H{
		"karyawan": updatedKaryawan,
	})
}

func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")

	err := h.repo.Delete(id)
	if err != nil {
		panic(err)
	}

	c.HTML(http.StatusOK, "karyawan-info-element", gin.H{
		"deleted": true,
	})
}

func (h *Handler) Add(c *gin.Context) {
	birthDate, birthDateErr := time.Parse("2006-01-02", c.PostForm("birthDate"))
	if birthDateErr != nil {
		panic(birthDateErr)
	}

	salary, salaryErr := strconv.Atoi(c.PostForm("salary"))
	if salaryErr != nil {
		panic(salaryErr)
	}

	joined, joinedErr := time.Parse("2006-01-02", c.PostForm("joined"))
	if joinedErr != nil {
		panic(joinedErr)
	}

	karyawan := Karyawan{
		Id:        uuid.New().String(),
		Name:      c.PostForm("name"),
		BirthDate: birthDate,
		Salary:    salary,
		Position:  c.PostForm("position"),
		Joined:    joined,
	}

	_, err := h.repo.Add(&karyawan)
	if err != nil {
		panic(err)
	}

	h.List(c)
}

func (h *Handler) Modal(c *gin.Context) {
	c.HTML(http.StatusOK, "add-karyawan-modal", gin.H{})
}
