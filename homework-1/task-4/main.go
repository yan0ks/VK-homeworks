package main

import "fmt"

func main() {
	const (
		secondsInMinute = 60
		secondsInHour   = 60 * secondsInMinute
		secondsInDay    = 24 * secondsInHour
	)

	uptime := 93784

	days := uptime / secondsInDay
	hours := uptime % secondsInDay / secondsInHour
	minutes := uptime % secondsInHour / secondsInMinute
	seconds := uptime % secondsInMinute

	fmt.Printf("%d -> ", uptime)
	separator := ""
	if days != 0 {
		fmt.Printf("%dд", days)
		separator = " "
	}
	if hours != 0 {
		fmt.Printf("%s%dч", separator, hours)
		separator = " "
	}
	if minutes != 0 {
		fmt.Printf("%s%dм", separator, minutes)
		separator = " "
	}
	if seconds != 0 || uptime == 0 {
		fmt.Printf("%s%dс", separator, seconds)
	}
	fmt.Println()
}
