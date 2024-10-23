package main

import (
	"fmt"
	"math/rand"
	"time"
)

// Define a struct to hold each student's subjects
type StudentSubjects struct {
	student  int
	subjects []string
}

// Function to generate a student's subject choices
func chooseSubjectFunction(student int, subjectChoicesClass *[]StudentSubjects) {
	var subjectChoiceStudent StudentSubjects
	subjectChoiceStudent.student = student

	// Core subjects
	subjectChoiceStudent.subjects = append(subjectChoiceStudent.subjects, "MATHS", "ENGLISH")

	// Sciences
	subjectChoiceStudent.subjects = append(subjectChoiceStudent.subjects, "BIOLOGY", "CHEMISTRY", "PHYSICS")

	// Humanities
	humanities := []string{"GEOGRAPHY", "HISTORY", "ECONOMICS"}
	randHumanities := make([]string, len(humanities))
	copy(randHumanities, humanities)
	rand.Shuffle(len(randHumanities), func(i, j int) { randHumanities[i], randHumanities[j] = randHumanities[j], randHumanities[i] })
	subjectChoiceStudent.subjects = append(subjectChoiceStudent.subjects, randHumanities[:2]...)

	// Languages
	languages := []string{"SPANISH", "FRENCH", "GERMAN"}
	language := languages[rand.Intn(len(languages))]
	subjectChoiceStudent.subjects = append(subjectChoiceStudent.subjects, language)

	// Optionals
	creatives := []string{"COMPUTER SCIENCE", "D&T", "ART", "MUSIC"}
	options := append(append([]string{}, humanities...), creatives...)
	options = append(options, languages...)
	optional := options[rand.Intn(len(options))]
	subjectChoiceStudent.subjects = append(subjectChoiceStudent.subjects, optional)

	// Add to the class list
	*subjectChoicesClass = append(*subjectChoicesClass, subjectChoiceStudent)
}

// Function to calculate the fitness score
func fitnessFunc(subjectChoicesClass []StudentSubjects, timetables [][]string) int {
	points := 0
	for i, student := range subjectChoicesClass {
		subjectCounted := make(map[string]bool)
		for _, subject := range student.subjects {
			if !subjectCounted[subject] && countOccurrences(timetables[i], subject) == 2 {
				points++
				subjectCounted[subject] = true
			}
		}
	}
	return points
}

// Helper function to count occurrences of a subject
func countOccurrences(slice []string, item string) int {
	count := 0
	for _, v := range slice {
		if v == item {
			count++
		}
	}
	return count
}

// Function to generate initial random timetables
func genStartTimetables(numStudents int) [][]string {
	timetables := make([][]string, numStudents)

	core := []string{"MATHS", "ENGLISH"}
	sciences := []string{"BIOLOGY", "CHEMISTRY", "PHYSICS"}
	humanities := []string{"GEOGRAPHY", "HISTORY", "ECONOMICS"}
	languages := []string{"SPANISH", "FRENCH", "GERMAN"}
	creatives := []string{"COMPUTER SCIENCE", "D&T", "ART", "MUSIC"}
	allSubjects := append(append(append(append(core, sciences...), humanities...), languages...), creatives...)

	for i := 0; i < numStudents; i++ {
		for j := 0; j < 78; j++ { // 8 periods * 5 days * 2 weeks - 2 periods
			subject := allSubjects[rand.Intn(len(allSubjects))]
			timetables[i] = append(timetables[i], subject)
		}
	}
	return timetables
}

// Function to mutate the timetable
func mutateTable(table [][]string, prob float64) [][]string {
	core := []string{"MATHS", "ENGLISH"}
	sciences := []string{"BIOLOGY", "CHEMISTRY", "PHYSICS"}
	humanities := []string{"GEOGRAPHY", "HISTORY", "ECONOMICS"}
	languages := []string{"SPANISH", "FRENCH", "GERMAN"}
	creatives := []string{"COMPUTER SCIENCE", "D&T", "ART", "MUSIC"}
	allSubjects := append(append(append(append(core, sciences...), humanities...), languages...), creatives...)

	for i := range table {
		for j := range table[i] {
			if rand.Float64() < prob {
				table[i][j] = allSubjects[rand.Intn(len(allSubjects))]
			}
		}
	}
	return table
}

// Function for tournament selection
func tournamentSelection(subjectChoicesClass []StudentSubjects, bitStrings [][][]string, k, parentSize int) [][][]string {
	parents := make([][][]string, 0, parentSize)
	for i := 0; i < parentSize; i++ {
		var best [][]string
		for j := 0; j < k; j++ {
			choice := bitStrings[rand.Intn(len(bitStrings))]
			if j == 0 || fitnessFunc(subjectChoicesClass, choice) > fitnessFunc(subjectChoicesClass, best) {
				best = choice
			}
		}
		parents = append(parents, best)
	}
	return parents
}

// Function to print the fortnightly timetable
func printFortnightlyTimetable(timetable [][]string, students []int) {
	days := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday"}

	for studentNum, studentTimetable := range timetable {
		fmt.Printf("\nStudent %d Timetable:\n\n", students[studentNum])

		week1 := studentTimetable[:40]
		week2 := studentTimetable[40:]

		fmt.Println("Week 1:")
		for day := 0; day < 5; day++ {
			fmt.Printf("%s: ", days[day])
			fmt.Println(week1[day*8 : (day+1)*8])
		}

		fmt.Println("\nWeek 2:")
		for day := 0; day < 5; day++ {
			fmt.Printf("%s: ", days[day])
			fmt.Println(week2[day*8 : (day+1)*8])
		}

		fmt.Println("\n========================================")
	}
}

// Main function
func main() {
	rand.Seed(time.Now().UnixNano())

	// 1) Generate students
	students := make([]int, 100)
	for i := 0; i < 100; i++ {
		students[i] = i + 1
	}
	fmt.Println(len(students))

	// Attach each student to their subjects
	subjectChoicesClass := []StudentSubjects{}
	for _, student := range students {
		chooseSubjectFunction(student, &subjectChoicesClass)
	}

	// Generate initial timetables
	startTimetables := genStartTimetables(len(students))

	// Calculate initial fitness
	currentFitness := fitnessFunc(subjectChoicesClass, startTimetables)
	fmt.Println("Starting Fitness:", currentFitness)

	// Evolutionary algorithm
	timetable := startTimetables
	for i := 0; i < 1000; i++ {
		fmt.Println("Iteration", i)
		pool := [][][]string{}
		currentFitness := fitnessFunc(subjectChoicesClass, timetable)
		fmt.Println("Current Fitness:", currentFitness)

		// Generate mutations
		for x := 0; x < 1000; x++ {
			tableToMutate := make([][]string, len(timetable))
			for idx := range timetable {
				tableToMutate[idx] = make([]string, len(timetable[idx]))
				copy(tableToMutate[idx], timetable[idx])
			}
			newTable := mutateTable(tableToMutate, 0.01)
			pool = append(pool, newTable)
		}

		// Selection
		parent1 := tournamentSelection(subjectChoicesClass, pool, 50, 1)[0]
		newFitness := fitnessFunc(subjectChoicesClass, parent1)
		if newFitness > currentFitness {
			timetable = parent1
		} else {
			continue
		}
	}

	// Final fitness
	finalFitness := fitnessFunc(subjectChoicesClass, timetable)
	fmt.Println("Final fitness", finalFitness)

	// Print timetables
	printFortnightlyTimetable(timetable, students)
}
