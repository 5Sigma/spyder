package endpoint

import (
	"fmt"
	"github.com/icrowley/fake"
	"regexp"
	"strings"
)

// fakeFuncs - A funtion map for faker keys and generators.
var fakeFuncs = map[string]func() string{
	"city":            fake.City,
	"color":           fake.Color,
	"country":         fake.Country,
	"email":           fake.EmailAddress,
	"femaleFirstName": fake.FemaleFirstName,
	"femaleFullName":  fake.FemaleFullName,
	"maleFirstName":   fake.MaleFirstName,
	"maleFullName":    fake.MaleFullName,
	"firstName":       fake.FirstName,
	"lastName":        fake.LastName,
	"fullName":        fake.FullName,
	"gender":          fake.Gender,
	"hexColor":        fake.HexColor,
	"ip":              fake.IPv4,
	"industry":        fake.Industry,
	"jobtitle":        fake.JobTitle,
	"month":           fake.Month,
	"monthNum":        func() string { return string(fake.MonthNum()) },
	"monthShort":      fake.MonthShort,
	"phone":           fake.Phone,
	"product":         fake.Product,
	"productName":     fake.ProductName,
	"street":          fake.Street,
	"streetAddress":   fake.StreetAddress,
	"weekDay":         fake.WeekDay,
	"zip":             fake.Zip,
}

// expandFakes - Expands faker templates into values and returns the expanded
// string
func expandFakes(str string) string {
	re := regexp.MustCompile(`\#\{([A-Za-z0-9]+)\}`)
	matches := re.FindAllStringSubmatch(str, -1)
	for _, match := range matches {
		if f, ok := fakeFuncs[match[1]]; ok {
			v := f()
			str = strings.Replace(str, fmt.Sprintf("#{%s}", match[1]), v, 1)
		}
	}
	return str
}
