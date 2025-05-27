package adminfeatures

import (
	"fixify_backend/middleware"
	errors "fixify_backend/model/error"
	"fixify_backend/model/response"
	"fixify_backend/model/users"
	"time"

	"github.com/gofiber/fiber/v2"
)

func PercentageRepairman(c *fiber.Ctx) error {
	db := middleware.DBConn
	now := time.Now()
	location := now.Location()

	// Determine the start of this week (Monday)
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7 // Sunday
	}
	weekStart := time.Date(now.Year(), now.Month(), now.Day()-weekday+1, 0, 0, 0, 0, location)
	weekEnd := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, int(time.Second-time.Nanosecond), location)

	// Determine the start and end of last week
	lastWeekStart := weekStart.AddDate(0, 0, -7)
	lastWeekEnd := weekStart.Add(-time.Nanosecond)

	// Fetch this week's repairmen
	var thisWeekRepairmen []users.User
	err := db.Model(&users.User{}).
		Where("type = ? AND createdat BETWEEN ? AND ?", "Repairman", weekStart, weekEnd).
		Find(&thisWeekRepairmen).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Failed to fetch this week's repairmen",
			Data: errors.ErrorModel{
				Message:   "Database error",
				IsSuccess: false,
				Error:     err.Error(),
			},
		})
	}

	// Fetch last week's repairmen
	var lastWeekRepairmen []users.User
	err = db.Model(&users.User{}).
		Where("type = ? AND createdat BETWEEN ? AND ?", "Repairman", lastWeekStart, lastWeekEnd).
		Find(&lastWeekRepairmen).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Failed to fetch last week's repairmen",
			Data: errors.ErrorModel{
				Message:   "Database error",
				IsSuccess: false,
				Error:     err.Error(),
			},
		})
	}

	// Calculate percentage change
	lastWeekCount := len(lastWeekRepairmen)
	thisWeekCount := len(thisWeekRepairmen)

	var percentageChange float64
	if lastWeekCount == 0 {
		if thisWeekCount > 0 {
			percentageChange = 100 // From 0 to something = 100% increase (or "new" growth)
		} else {
			percentageChange = 0 // 0 to 0 = no change
		}
	} else {
		percentageChange = (float64(thisWeekCount-lastWeekCount) / float64(lastWeekCount)) * 100
	}

	// Collect creation dates for this week's repairmen
	var timestamps []string
	for _, r := range thisWeekRepairmen {
		timestamps = append(timestamps, r.Created_at.Format("2006-01-02"))
	}

	return c.JSON(response.ResponseModel{
		RetCode: "200",
		Message: "Success",
		Data: fiber.Map{
			"role":              "repairman",
			"week_start":        weekStart.Format("2006-01-02"),
			"week_end":          weekEnd.Format("2006-01-02"),
			"this_week_count":   thisWeekCount,
			"last_week_count":   lastWeekCount,
			"percentage_change": percentageChange,
			"creation_dates":    timestamps,
		},
	})
}

func PercentageClient(c *fiber.Ctx) error {
	db := middleware.DBConn
	now := time.Now()
	location := now.Location()

	// Determine the start of this week (Monday)
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7 // Sunday
	}
	weekStart := time.Date(now.Year(), now.Month(), now.Day()-weekday+1, 0, 0, 0, 0, location)
	weekEnd := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, int(time.Second-time.Nanosecond), location)

	// Determine the start and end of last week
	lastWeekStart := weekStart.AddDate(0, 0, -7)
	lastWeekEnd := weekStart.Add(-time.Nanosecond)

	// Fetch this week's repairmen
	var thisWeekClient []users.User
	err := db.Model(&users.User{}).
		Where("type = ? AND createdat BETWEEN ? AND ?", "Client", weekStart, weekEnd).
		Find(&thisWeekClient).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Failed to fetch this week's repairmen",
			Data: errors.ErrorModel{
				Message:   "Database error",
				IsSuccess: false,
				Error:     err.Error(),
			},
		})
	}

	// Fetch last week's repairmen
	var lastWeekClient []users.User
	err = db.Model(&users.User{}).
		Where("type = ? AND createdat BETWEEN ? AND ?", "Client", lastWeekStart, lastWeekEnd).
		Find(&lastWeekClient).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Failed to fetch last week's repairmen",
			Data: errors.ErrorModel{
				Message:   "Database error",
				IsSuccess: false,
				Error:     err.Error(),
			},
		})
	}

	// Calculate percentage change
	lastWeekCount := len(lastWeekClient)
	thisWeekCount := len(thisWeekClient)

	var percentageChange float64
	if lastWeekCount == 0 {
		if thisWeekCount > 0 {
			percentageChange = 100 // From 0 to something = 100% increase (or "new" growth)
		} else {
			percentageChange = 0 // 0 to 0 = no change
		}
	} else {
		percentageChange = (float64(thisWeekCount-lastWeekCount) / float64(lastWeekCount)) * 100
	}

	// Collect creation dates for this week's repairmen
	var timestamps []string
	for _, r := range thisWeekClient {
		timestamps = append(timestamps, r.Created_at.Format("2006-01-02"))
	}

	return c.JSON(response.ResponseModel{
		RetCode: "200",
		Message: "Success",
		Data: fiber.Map{
			"role":              "Client",
			"week_start":        weekStart.Format("2006-01-02"),
			"week_end":          weekEnd.Format("2006-01-02"),
			"this_week_count":   thisWeekCount,
			"last_week_count":   lastWeekCount,
			"percentage_change": percentageChange,
			"creation_dates":    timestamps,
		},
	})
}

// func PercentageRequests(c *fiber.Ctx) error {
// 	db := middleware.DBConn
// 	now := time.Now()
// 	location := now.Location()

// 	// Determine the start of this week (Monday)
// 	weekday := int(now.Weekday())
// 	if weekday == 0 {
// 		weekday = 7 // Sunday
// 	}
// 	weekStart := time.Date(now.Year(), now.Month(), now.Day()-weekday+1, 0, 0, 0, 0, location)
// 	weekEnd := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, int(time.Second-time.Nanosecond), location)

// 	// Determine the start and end of last week
// 	lastWeekStart := weekStart.AddDate(0, 0, -7)
// 	lastWeekEnd := weekStart.Add(-time.Nanosecond)

// 	// Fetch this week's service requests
// 	var thisWeekRequests []users.ServiceRequest
// 	err := db.Model(&users.ServiceRequest{}).
// 		Where("request_date BETWEEN ? AND ?", weekStart, weekEnd).
// 		Find(&thisWeekRequests).Error
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(response.ResponseModel{
// 			RetCode: "500",
// 			Message: "Failed to fetch this week's service requests",
// 			Data: errors.ErrorModel{
// 				Message:   "Database error",
// 				IsSuccess: false,
// 				Error:     err.Error(),
// 			},
// 		})
// 	}

// 	// Fetch last week's service requests
// 	var lastWeekRequests []users.ServiceRequest
// 	err = db.Model(&users.ServiceRequest{}).
// 		Where("request_date BETWEEN ? AND ?", lastWeekStart, lastWeekEnd).
// 		Find(&lastWeekRequests).Error
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(response.ResponseModel{
// 			RetCode: "500",
// 			Message: "Failed to fetch last week's service requests",
// 			Data: errors.ErrorModel{
// 				Message:   "Database error",
// 				IsSuccess: false,
// 				Error:     err.Error(),
// 			},
// 		})
// 	}

// 	// Calculate percentage change
// 	lastWeekCount := len(lastWeekRequests)
// 	thisWeekCount := len(thisWeekRequests)

// 	var percentageChange float64
// 	if lastWeekCount == 0 {
// 		if thisWeekCount > 0 {
// 			percentageChange = 100
// 		} else {
// 			percentageChange = 0
// 		}
// 	} else {
// 		percentageChange = (float64(thisWeekCount-lastWeekCount) / float64(lastWeekCount)) * 100
// 	}

// 	// Collect request dates for this week
// 	var timestamps []string
// 	for _, r := range thisWeekRequests {
// 		timestamps = append(timestamps, r.RequestDate.Format("2006-01-02"))
// 	}

// 	return c.JSON(response.ResponseModel{
// 		RetCode: "200",
// 		Message: "Success",
// 		Data: fiber.Map{
// 			"role":              "Requests",
// 			"week_start":        weekStart.Format("2006-01-02"),
// 			"week_end":          weekEnd.Format("2006-01-02"),
// 			"this_week_count":   thisWeekCount,
// 			"last_week_count":   lastWeekCount,
// 			"percentage_change": percentageChange,
// 			"creation_dates":    timestamps,
// 		},
// 	})
// }
