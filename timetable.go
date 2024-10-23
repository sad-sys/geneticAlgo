package main

import "fmt"

func makeStudents() [] int {

	students := make([]int, 30)

	for i := range students {
		students[i] = i + 1
	}
	return students
}

func makeStudentChoices() [] int {
	
}

func main() {
	
	students := makeStudents()
	fmt.Println(students)
}
