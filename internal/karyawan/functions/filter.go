package functions

import (
	"net/url"
	"regexp"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Filter struct {
	Name         string
	Position     string
	SalaryMin    *int
	SalaryMax    *int
	BirthDateOp  string
	BirthDate    *time.Time
	BirthDateEnd *time.Time
	JoinedOp     string
	Joined       *time.Time
	JoinedEnd    *time.Time
}

func (f Filter) ToBSON() bson.M {
	filter := bson.M{}

	if f.Name != "" {
		filter["Name"] = bson.M{"$regex": regexp.QuoteMeta(f.Name), "$options": "i"}
	}

	if f.Position != "" {
		filter["Position"] = bson.M{"$regex": regexp.QuoteMeta(f.Position), "$options": "i"}
	}

	if f.SalaryMin != nil || f.SalaryMax != nil {
		salary := bson.M{}
		if f.SalaryMin != nil {
			salary["$gte"] = *f.SalaryMin
		}
		if f.SalaryMax != nil {
			salary["$lte"] = *f.SalaryMax
		}
		filter["Salary"] = salary
	}

	if dateRange := dateRangeBSON(f.BirthDateOp, f.BirthDate, f.BirthDateEnd); dateRange != nil {
		filter["BirthDate"] = dateRange
	}

	if dateRange := dateRangeBSON(f.JoinedOp, f.Joined, f.JoinedEnd); dateRange != nil {
		filter["Joined"] = dateRange
	}

	return filter
}

func dateRangeBSON(op string, start *time.Time, end *time.Time) bson.M {
	switch op {
	case "on":
		if start == nil {
			return nil
		}
		return bson.M{"$eq": *start}
	case "before":
		if start == nil {
			return nil
		}
		return bson.M{"$lt": *start}
	case "after":
		if start == nil {
			return nil
		}
		return bson.M{"$gt": *start}
	case "between":
		if start == nil || end == nil {
			return nil
		}
		return bson.M{"$gte": *start, "$lte": *end}
	default:
		return nil
	}
}

// QueryString rebuilds the query string for the currently active filters, so
// pagination links can carry them forward.
func (f Filter) QueryString() string {
	values := url.Values{}

	if f.Name != "" {
		values.Set("name", f.Name)
	}
	if f.Position != "" {
		values.Set("position", f.Position)
	}
	if f.SalaryMin != nil {
		values.Set("salaryMin", strconv.Itoa(*f.SalaryMin))
	}
	if f.SalaryMax != nil {
		values.Set("salaryMax", strconv.Itoa(*f.SalaryMax))
	}
	if f.BirthDateOp != "" {
		values.Set("birthDateOp", f.BirthDateOp)
	}
	if f.BirthDate != nil {
		values.Set("birthDate", f.BirthDate.Format("2006-01-02"))
	}
	if f.BirthDateEnd != nil {
		values.Set("birthDateEnd", f.BirthDateEnd.Format("2006-01-02"))
	}
	if f.JoinedOp != "" {
		values.Set("joinedOp", f.JoinedOp)
	}
	if f.Joined != nil {
		values.Set("joined", f.Joined.Format("2006-01-02"))
	}
	if f.JoinedEnd != nil {
		values.Set("joinedEnd", f.JoinedEnd.Format("2006-01-02"))
	}

	return values.Encode()
}

func (f Filter) SalaryMinString() string {
	if f.SalaryMin == nil {
		return ""
	}
	return strconv.Itoa(*f.SalaryMin)
}

func (f Filter) SalaryMaxString() string {
	if f.SalaryMax == nil {
		return ""
	}
	return strconv.Itoa(*f.SalaryMax)
}

func (f Filter) BirthDateString() string {
	return dateString(f.BirthDate)
}

func (f Filter) BirthDateEndString() string {
	return dateString(f.BirthDateEnd)
}

func (f Filter) JoinedString() string {
	return dateString(f.Joined)
}

func (f Filter) JoinedEndString() string {
	return dateString(f.JoinedEnd)
}

func dateString(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("2006-01-02")
}
