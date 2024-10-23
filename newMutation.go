package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Student struct {
	ID      int
	Choices []string
}

func generateStudents(num int) []Student {
	students := make([]Student, num)
	for i := 0; i < num; i++ {
		students[i] = Student{ID: i + 1}
	}
	return students
}

func chooseSubjectFunction(student Student) Student {
	subjectChoiceStudent := []string{}

	// Core subjects
	subjectChoiceStudent = append(subjectChoiceStudent, "MATHS", "ENGLISH")

	// Sciences
	subjectChoiceStudent = append(subjectChoiceStudent, "BIOLOGY", "CHEMISTRY", "PHYSICS")

	// One out of humanities
	humanities := []string{"GEOGRAPHY", "HISTORY", "ECONOMICS"}
	subjectChoiceStudent = append(subjectChoiceStudent, humanities[rand.Intn(len(humanities))])

	// One out of languages
	languages := []string{"SPANISH", "FRENCH", "GERMAN"}
	subjectChoiceStudent = append(subjectChoiceStudent, languages[rand.Intn(len(languages))])

	// One out of any of the optionals
	creatives := []string{"COMPUTER SCIENCE", "D&T", "ART", "MUSIC"}
	optionals := append(append([]string{}, humanities...), languages...)
	optionals = append(optionals, creatives...)
	subjectChoiceStudent = append(subjectChoiceStudent, optionals[rand.Intn(len(optionals))])

	student.Choices = subjectChoiceStudent
	return student
}

func fitnessFunc(subjectChoicesClass []Student, timetables [][]string) int {
	individualFitnessPoints := 0 // Points for fulfilling individual student's requirements
	alignmentFitnessPoints := 0  // Points for subject alignment across students

	// Check if each subject appears exactly twice for each student in 35 periods
	for i, studentTimetable := range timetables {
		subjectCounts := make(map[string]int)
		// Count occurrences of each subject in the student's timetable
		for _, subject := range studentTimetable {
			subjectCounts[subject]++
		}
		// Award points if each subject appears exactly twice
		for _, count := range subjectCounts {
			if count == 2 {
				individualFitnessPoints++
			}
		}
		// Additional fitness for matching student's choices
		for _, choice := range subjectChoicesClass[i].Choices {
			if subjectCounts[choice] >= 2 {
				individualFitnessPoints += 2
			}
		}
	}

	// Check alignment of subjects across all students at each period (35 periods)
	for period := 0; period < 35; period++ {
		subjectCounts := make(map[string]int)
		// Count how many students have each subject at the current time slot
		for _, studentTimetable := range timetables {
			subject := studentTimetable[period]
			subjectCounts[subject]++
		}
		// Award points for alignment
		for _, count := range subjectCounts {
			if count > 1 {
				alignmentFitnessPoints += count - 1
			}
		}
	}

	totalFitness := individualFitnessPoints + alignmentFitnessPoints
	return totalFitness
}

func genStartTimetables(numStudents int) [][]string {
	timetables := make([][]string, numStudents)

	core := []string{"MATHS", "ENGLISH"}
	sciences := []string{"BIOLOGY", "CHEMISTRY", "PHYSICS"}
	humanities := []string{"GEOGRAPHY", "HISTORY", "ECONOMICS"}
	languages := []string{"SPANISH", "FRENCH", "GERMAN"}
	creatives := []string{"COMPUTER SCIENCE", "D&T", "ART", "MUSIC"}

	allSubjects := append(append(append(append(core, sciences...), humanities...), languages...), creatives...)

	for i := 0; i < numStudents; i++ {
		studentTimetable := make([]string, 35) // Change to 35 periods
		for j := 0; j < 35; j++ {
			studentTimetable[j] = allSubjects[rand.Intn(len(allSubjects))]
		}
		timetables[i] = studentTimetable
	}

	return timetables
}

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

func tournamentSelection(subjectChoicesClass []Student, population [][][]string, k int, parentSize int) [][][]string {
	parents := make([][][]string, parentSize)
	for i := 0; i < parentSize; i++ {
		best := population[rand.Intn(len(population))]
		bestFitness := fitnessFunc(subjectChoicesClass, best)
		for j := 1; j < k; j++ {
			choice := population[rand.Intn(len(population))]
			choiceFitness := fitnessFunc(subjectChoicesClass, choice)
			if choiceFitness > bestFitness {
				best = choice
				bestFitness = choiceFitness
			}
		}
		parents[i] = best
	}
	return parents
}

func crossover(parent1, parent2 [][]string) [][]string {
	offspring := make([][]string, len(parent1))
	for i := range parent1 {
		crossoverPoint := rand.Intn(len(parent1[i]))
		offspring[i] = append([]string{}, parent1[i][:crossoverPoint]...)
		offspring[i] = append(offspring[i], parent2[i][crossoverPoint:]...)
	}
	return offspring
}
func printFortnightlyTimetable(timetable [][]string, students []Student) {
	days := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday"}

	for studentNum, studentTimetable := range timetable {
		fmt.Printf("\nStudent %d Timetable:\n", students[studentNum].ID)

		// Since there are 35 periods, let's distribute them properly:
		// Week 1: 17 periods (3 periods for Monday-Thursday, 2 for Friday)
		// Week 2: 18 periods (3 periods for each day)

		week1 := studentTimetable[:17]
		week2 := studentTimetable[17:]

		// Printing Week 1 Timetable (17 periods)
		fmt.Println("Week 1:")
		for day := 0; day < 5; day++ {
			if day < 4 { // Monday to Thursday (3 periods)
				fmt.Printf("%s: %v\n", days[day], week1[day*3:(day+1)*3])
			} else { // Friday (2 periods)
				fmt.Printf("%s: %v\n", days[day], week1[day*3:])
			}
		}

		// Printing Week 2 Timetable (18 periods)
		fmt.Println("\nWeek 2:")
		for day := 0; day < 5; day++ {
			// Each day has exactly 3 periods in Week 2
			fmt.Printf("%s: %v\n", days[day], week2[day*3:(day+1)*3])
		}

		fmt.Println("\n========================================")
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// Generate students
	numStudents := 30
	students := generateStudents(numStudents)

	// Assign subject choices to each student
	subjectChoicesClass := make([]Student, numStudents)
	for i := range students {
		subjectChoicesClass[i] = chooseSubjectFunction(students[i])
	}

	// Generate initial timetables
	timetable := genStartTimetables(numStudents)
	fmt.Println("Starting Fitness:", fitnessFunc(subjectChoicesClass, timetable))

	// Evolutionary algorithm parameters
	numGenerations := 500    // Increased number of generations
	populationSize := 2000   // Increased population size
	baseMutationRate := 0.05 // Starting mutation rate
	mutationRate := baseMutationRate
	maxMutationRate := 0.5
	minMutationRate := 0.01

	sampleSizes := []int{10, 20, 30, 40, 50, 60, 70, 80, 90, 100}

	// Variables to track convergence
	bestFitness := fitnessFunc(subjectChoicesClass, timetable)
	stagnationCounter := 0
	stagnationThreshold := 20 // Number of generations to wait before increasing mutation rate

	for gen := 0; gen < numGenerations; gen++ {
		fmt.Printf("Generation %d\n", gen+1)
		pool := [][][]string{}
		currentFitness := fitnessFunc(subjectChoicesClass, timetable)
		fmt.Println("Current Fitness:", currentFitness)

		// Check for stagnation
		if currentFitness <= bestFitness {
			stagnationCounter++
		} else {
			stagnationCounter = 0
			bestFitness = currentFitness
		}

		// Adjust mutation rate based on stagnation
		if stagnationCounter >= stagnationThreshold {
			mutationRate *= 1.05 // Increase mutation rate by 5%
			if mutationRate > maxMutationRate {
				mutationRate = maxMutationRate
			}
			fmt.Printf("Increasing mutation rate to %.4f due to stagnation\n", mutationRate)
		} else {
			// Gradually reduce mutation rate back to base rate
			if mutationRate > baseMutationRate {
				mutationRate *= 0.95 // Decrease mutation rate by 5%
				if mutationRate < minMutationRate {
					mutationRate = minMutationRate
				}
				fmt.Printf("Decreasing mutation rate to %.4f\n", mutationRate)
			}
		}

		// Generate new population
		for x := 0; x < populationSize; x++ {
			tableToMutate := make([][]string, len(timetable))
			for i := range timetable {
				tableToMutate[i] = append([]string{}, timetable[i]...)
			}
			mutatedTable := mutateTable(tableToMutate, mutationRate)
			pool = append(pool, mutatedTable)
		}

		// Tournament selection with dynamic sample size for diversity
		k := sampleSizes[rand.Intn(len(sampleSizes))]
		parent1 := tournamentSelection(subjectChoicesClass, pool, k, 1)[0]
		parent2 := tournamentSelection(subjectChoicesClass, pool, k, 1)[0]

		// Crossover
		offspring := crossover(parent1, parent2)

		newFitness := fitnessFunc(subjectChoicesClass, offspring)
		if newFitness > currentFitness {
			timetable = offspring
			fmt.Println("New better fitness found:", newFitness)
		}
	}

	fmt.Println("Final Fitness:", fitnessFunc(subjectChoicesClass, timetable))
	printFortnightlyTimetable(timetable, students)
}
