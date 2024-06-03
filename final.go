package main

import (
	"fmt"
)

type Employee struct {
	DailySalary             float64
	MaxRegularHours         int
	WorkdaysPerWeek         int
	StartWorkTime           string
	EndWorkTime             string
	WorkHoursPerDay         map[string]int
	OvertimeHours           map[string]int
	NightShiftHours         map[string]int
	NightShiftOvertimeHours map[string]int
	IsRestDay               map[string]bool
	IsSpecialNonWorkingDay  map[string]bool
	IsRegularHoliday        map[string]bool
}

func main() {
	employee := Employee{
		DailySalary:             500.0,
		MaxRegularHours:         8,
		WorkdaysPerWeek:         5,
		StartWorkTime:           "0900",
		EndWorkTime:             "1800",
		WorkHoursPerDay:         make(map[string]int),
		OvertimeHours:           make(map[string]int),
		NightShiftHours:         make(map[string]int),
		NightShiftOvertimeHours: make(map[string]int),
		IsRestDay:               make(map[string]bool),
		IsSpecialNonWorkingDay:  make(map[string]bool),
		IsRegularHoliday:        make(map[string]bool),
	}

	for i := 0; i < 7; i++ {
		day := getDayName(i)
		employee.WorkHoursPerDay[day] = employee.MaxRegularHours
		if i == 5 || i == 6 {
			employee.IsRestDay[day] = true
		} else {
			employee.IsRestDay[day] = false
			employee.IsSpecialNonWorkingDay[day] = false
			employee.IsRegularHoliday[day] = false
		}
	}

	for {
		fmt.Println("\nPayroll System Menu:")
		fmt.Println("1. Configure Default Settings")
		fmt.Println("2. Input Employee Attendance")
		fmt.Println("3. Generate Payroll")
		fmt.Println("4. Exit")
		fmt.Print("Select an option: ")

		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			configureSettings(&employee)
		case 2:
			inputAttendance(&employee)
		case 3:
			generatePayroll(&employee)
		case 4:
			fmt.Println("Exiting the Payroll System. Goodbye!")
			return
		default:
			fmt.Println("Invalid choice. Please select a valid option.")
		}
	}
}

func configureSettings(employee *Employee) {
	fmt.Println("\nConfigure Default Settings:")
	fmt.Printf("1. Daily Salary (currently $%.2f)\n", employee.DailySalary)
	fmt.Printf("2. Maximum Regular Work Hours per Day (currently %d hours)\n", employee.MaxRegularHours)
	fmt.Println("3. Assign Type of Work Day")

	fmt.Print("Select a setting to change (1-3) or 0 to go back: ")
	var choice int
	fmt.Scanln(&choice)

	switch choice {
	case 1:
		fmt.Print("Enter the new daily salary: $")
		fmt.Scanln(&employee.DailySalary)
	case 2:
		fmt.Print("Enter the new maximum regular work hours per day: ")
		fmt.Scanln(&employee.MaxRegularHours)
	case 3:
		assignTypeOfWorkDay(employee)
	case 0:
		return
	default:
		fmt.Println("Invalid choice. Please select a valid option.")
	}
}

func assignTypeOfWorkDay(employee *Employee) {
	fmt.Println("\nAssign Type of Work Day:")

	for i := 0; i < 7; i++ {
		day := getDayName(i)
		fmt.Printf("\n%s:", day)
		fmt.Println("\n1. Normal Day")
		fmt.Println("2. Rest Day")
		fmt.Println("3. Special Non-Working Day")
		fmt.Println("4. Special Non-Working and Rest Day")
		fmt.Println("5. Regular Holiday")
		fmt.Println("6. Regular Holiday and Rest Day")
		fmt.Print("Select the type for this day (1-6): ")

		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			employee.IsRestDay[day] = false
			employee.IsSpecialNonWorkingDay[day] = false
			employee.IsRegularHoliday[day] = false
		case 2:
			employee.IsRestDay[day] = true
		case 3:
			employee.IsRestDay[day] = false
			employee.IsSpecialNonWorkingDay[day] = true
			employee.IsRegularHoliday[day] = false
		case 4:
			employee.IsRestDay[day] = true
			employee.IsSpecialNonWorkingDay[day] = true
			employee.IsRegularHoliday[day] = false
		case 5:
			employee.IsRestDay[day] = false
			employee.IsSpecialNonWorkingDay[day] = false
			employee.IsRegularHoliday[day] = true
		case 6:
			employee.IsRestDay[day] = true
			employee.IsSpecialNonWorkingDay[day] = false
			employee.IsRegularHoliday[day] = true
		default:
			fmt.Println("Invalid choice. Please select a valid option.")
			i-- // Decrement i to repeat the current day's selection
		}
	}
}

func inputAttendance(employee *Employee) {
	fmt.Println("\nInput Employee Attendance:")
	for i := 0; i < 7; i++ {
		day := getDayName(i)
		fmt.Printf("Enter IN time for %s (HHMM): ", day)
		var inTime string
		fmt.Scanln(&inTime)
		fmt.Printf("Enter OUT time for %s (HHMM) or 0900 for rest day: ", day)
		var outTime string
		fmt.Scanln(&outTime)
		regularHours, overtimeHours, nightShiftHours, nightShiftOvertimeHours := calculateWorkHours(employee, inTime, outTime, employee.IsRestDay[day], employee.IsSpecialNonWorkingDay[day], employee.IsRegularHoliday[day])
		employee.WorkHoursPerDay[day] = regularHours
		employee.OvertimeHours[day] = overtimeHours
		employee.NightShiftHours[day] = nightShiftHours
		employee.NightShiftOvertimeHours[day] = nightShiftOvertimeHours
	}
	fmt.Println("Attendance recorded successfully.")
}

func calculateWorkHours(employee *Employee, inTime, outTime string, isRestDay, isSpecialNonWorkingDay, isRegularHoliday bool) (int, int, int, int) {
	inHour, inMinute := parseTime(inTime)
	outHour, outMinute := parseTime(outTime)

	if inHour >= outHour {
		outHour += 24 // Handle cases where outTime is on the next day
	}
	if outHour == 0 {
		outHour = 24
	}

	totalHours := outHour - inHour
	totalMinutes := outMinute - inMinute
	if totalMinutes < 0 {
		totalHours--
		totalMinutes += 60
	}

	// Calculate regular work hour
	regularHours := totalHours
	nightShiftHours := 0
	nightShiftOvertimeHours := 0

	// Calculate night shift hours
	fmt.Println("Regular Hours: ", regularHours)
	fmt.Println("Out Hour: ", outHour)

	if regularHours <= (employee.MaxRegularHours + 1) {
		if outHour >= 22 || outHour <= 6 {
			if outHour >= 22 {
				nightShiftHours = outHour - 22
			}
			if outHour <= 6 {
				nightShiftHours = outHour + 2
			}
			regularHours -= nightShiftHours
		}
	}
	fmt.Println("Max Regular Hours: ", employee.MaxRegularHours)
	fmt.Println("Night Shift Hours: ", nightShiftHours)
	// Calculate overtime hours
	overtimeHours := 0
	if regularHours > employee.MaxRegularHours {
		overtimeHours = regularHours - employee.MaxRegularHours
		regularHours = employee.MaxRegularHours
		if outHour >= 22 || outHour <= 6 {
			if outHour >= 22 {
				nightShiftOvertimeHours = outHour - 22
			}
			if outHour <= 6 {
				nightShiftOvertimeHours = outHour
			}
			overtimeHours -= nightShiftOvertimeHours
		}
	}

	if overtimeHours >= 1 {
		overtimeHours -= 1
	}
	fmt.Println("Overtime: ", overtimeHours)
	fmt.Println("Night shift overtime: ", nightShiftOvertimeHours)

	return regularHours, overtimeHours, nightShiftHours, nightShiftOvertimeHours
}

func generatePayroll(employee *Employee) {
	fmt.Println("\nGenerating Payroll:")
	totalSalary := 0.0
	for i := 0; i < 7; i++ {
		day := getDayName(i)
		regularHours := employee.WorkHoursPerDay[day]
		overtimeHours := employee.OvertimeHours[day]
		nightShiftHours := employee.NightShiftHours[day]
		nightShiftOvertimeHours := employee.NightShiftOvertimeHours[day]
		dailySalary := calculateDailySalary(employee, regularHours, overtimeHours, nightShiftHours, nightShiftOvertimeHours, employee.IsRestDay[day], employee.IsSpecialNonWorkingDay[day], employee.IsRegularHoliday[day])
		fmt.Printf("Day: %s, Work Hours: %d, Night Shift Hours: %d, Overtime Hours: %d, Night Shift Overtime Hours: %d, Daily Salary: $%.2f\n", day, regularHours, nightShiftHours, overtimeHours, nightShiftOvertimeHours, dailySalary)
		totalSalary += dailySalary
	}
	fmt.Printf("Total Salary for the Week: $%.2f\n", totalSalary)
}

func calculateDailySalary(employee *Employee, regularHours, overtimeHours, nightShiftHours int, nightShiftOvertimeHours int, isRestDay, isSpecialNonWorkingDay, isRegularHoliday bool) float64 {
	hourlyRate := employee.DailySalary / float64(employee.MaxRegularHours)

	regularSalary := calculateRegularSalary(employee, hourlyRate, regularHours, isRestDay, isSpecialNonWorkingDay, isRegularHoliday)
	nightShiftSalary := calculateNightShiftSalary(employee, nightShiftHours, isRestDay, isSpecialNonWorkingDay, isRegularHoliday)

	totalOvertimeHours := overtimeHours
	overtimeSalary := calculateOvertimeSalary(employee, totalOvertimeHours, isRestDay, isSpecialNonWorkingDay, isRegularHoliday)
	nightShiftOverTimeSalary := calculateNightShiftOvertimeSalary(employee, nightShiftOvertimeHours, isRestDay, isSpecialNonWorkingDay, isRegularHoliday)

	dailySalary := regularSalary + overtimeSalary + nightShiftOverTimeSalary + nightShiftSalary

	return dailySalary
}

func calculateNightShiftSalary(employee *Employee, nightShiftHours int, isRestDay, isSpecialNonWorkingDay, isRegularHoliday bool) float64 {
	fmt.Println(nightShiftHours)
	if nightShiftHours == 0 {
		return 0.0
	}

	hourlyRate := employee.DailySalary / float64(employee.MaxRegularHours)
	nightShiftMultiplier := 1.10 // Night shift hourly rate multiplier

	if isRestDay || isSpecialNonWorkingDay {
		nightShiftMultiplier += nightShiftMultiplier * 0.3
	}
	if isSpecialNonWorkingDay && isRegularHoliday {
		nightShiftMultiplier += nightShiftMultiplier * 0.5
	}
	if isRegularHoliday {
		nightShiftMultiplier += nightShiftMultiplier * 1.0
	}
	if isRegularHoliday && isRestDay {
		nightShiftMultiplier += nightShiftMultiplier * 1.6
	}

	nightShiftSalary := float64(nightShiftHours) * hourlyRate * nightShiftMultiplier
	fmt.Println(nightShiftSalary)
	return nightShiftSalary
}

func calculateNightShiftOvertimeSalary(employee *Employee, nightShiftOvertimeHours int, isRestDay, isSpecialNonWorkingDay, isRegularHoliday bool) float64 {
	if nightShiftOvertimeHours == 0 {
		return 0.0
	}

	overtimeRate := employee.DailySalary / float64(employee.MaxRegularHours)
	overtimeMultiplier := 1.375 // Regular Day Nightshift Overtime

	if isRestDay || isSpecialNonWorkingDay {
		overtimeMultiplier = 1.859
	}
	if isSpecialNonWorkingDay && isRestDay {
		overtimeMultiplier = 2.145
	}
	if isRegularHoliday {
		overtimeMultiplier = 2.86
	}
	if isRegularHoliday && isRestDay {
		overtimeMultiplier = 3.718
	}

	nightShiftOvertimeSalary := float64(nightShiftOvertimeHours) * overtimeRate * overtimeMultiplier
	fmt.Println("Night Shift Overtime Salary: ", nightShiftOvertimeSalary)
	return nightShiftOvertimeSalary
}

func calculateRegularSalary(employee *Employee, hourlyRate float64, regularHours int, isRestDay, isSpecialNonWorkingDay, isRegularHoliday bool) float64 {
	dayRateMultiplier := 1.0
	if regularHours < employee.MaxRegularHours {
		regularHours = employee.MaxRegularHours
	}
	if isRestDay || isSpecialNonWorkingDay {
		dayRateMultiplier = 1.3
	}
	if isSpecialNonWorkingDay {
		dayRateMultiplier = 1.5
	}
	if isRegularHoliday {
		dayRateMultiplier = 2.0
	}
	if isRegularHoliday && isRestDay {
		dayRateMultiplier = 2.6
	}

	regularSalary := hourlyRate * float64(regularHours) * dayRateMultiplier
	//if regularHours <= employee.MaxRegularHours {
	//return employee.DailySalary
	//}
	fmt.Println("Regular Salary: ", regularSalary)
	return regularSalary
}

func calculateOvertimeSalary(employee *Employee, totalOvertimeHours int, isRestDay, isSpecialNonWorkingDay, isRegularHoliday bool) float64 {
	overtimeMultiplier := 1.25
	overtimeRate := employee.DailySalary / float64(employee.MaxRegularHours)

	if isRegularHoliday && isSpecialNonWorkingDay {
		overtimeMultiplier = 3.38
	}

	if isRestDay || isSpecialNonWorkingDay {
		overtimeMultiplier = 1.69
		{
			if isSpecialNonWorkingDay && isRestDay {
				overtimeMultiplier = 1.95
			}
		}
	}

	if isRegularHoliday {
		overtimeMultiplier = 2.6
	}

	totalOvertimeSalary := float64(totalOvertimeHours) * overtimeRate * overtimeMultiplier
	fmt.Println("Overtime Salary: ", totalOvertimeSalary)
	return totalOvertimeSalary
}

func parseTime(timeStr string) (int, int) {
	hour := 0
	minute := 0

	fmt.Sscanf(timeStr, "%02d%02d", &hour, &minute)
	return hour, minute
}

func getDayName(dayIndex int) string {
	days := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
	return days[dayIndex]
}
