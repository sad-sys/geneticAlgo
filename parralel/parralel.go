package main

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

// Define constants for subjects
const (
	MATHS            = "MATHS"
	ENGLISH          = "ENGLISH"
	BIOLOGY          = "BIOLOGY"
	CHEMISTRY        = "CHEMISTRY"
	PHYSICS          = "PHYSICS"
	GEOGRAPHY        = "GEOGRAPHY"
	HISTORY          = "HISTORY"
	ECONOMICS        = "ECONOMICS"
	SPANISH          = "SPANISH"
	FRENCH           = "FRENCH"
	GERMAN           = "GERMAN"
	COMPUTER_SCIENCE = "COMPUTER SCIENCE"
	DT               = "D&T"
	ART              = "ART"
	MUSIC            = "MUSIC"
	FREE             = "FREE"
)

// Map for subject period requirements per week
var subjectPeriodRequirements = map[string]int{
	MATHS:            5,
	ENGLISH:          5,
	BIOLOGY:          4,
	CHEMISTRY:        4,
	PHYSICS:          4,
	GEOGRAPHY:        3,
	HISTORY:          3,
	ECONOMICS:        3,
	SPANISH:          3,
	FRENCH:           3,
	GERMAN:           3,
	COMPUTER_SCIENCE: 2,
	DT:               2,
	ART:              2,
	MUSIC:            2,
}

// StudentSubjects holds each student's subject choices
type StudentSubjects struct {
	student  int
	subjects []string
}

// Timetable represents a student's timetable over two weeks
type Timetable struct {
	schedule [10][8]string // 2 weeks, 5 days per week, 8 periods per day
}

// CandidateSolution represents a potential solution in the population
type CandidateSolution struct {
	timetables []Timetable
	fitness    int
}

// Generates a student's subject choices
func chooseSubjectFunction(student int, subjectChoicesClass *[]StudentSubjects) {
	var subjectChoiceStudent StudentSubjects
	subjectChoiceStudent.student = student

	// Core subjects
	subjectChoiceStudent.subjects = append(subjectChoiceStudent.subjects, MATHS, ENGLISH)

	// Sciences
	subjectChoiceStudent.subjects = append(subjectChoiceStudent.subjects, BIOLOGY, CHEMISTRY, PHYSICS)

	// Humanities (choose 2 out of 3)
	humanities := []string{GEOGRAPHY, HISTORY, ECONOMICS}
	randHumanities := make([]string, len(humanities))
	copy(randHumanities, humanities)
	rand.Shuffle(len(randHumanities), func(i, j int) {
		randHumanities[i], randHumanities[j] = randHumanities[j], randHumanities[i]
	})
	subjectChoiceStudent.subjects = append(subjectChoiceStudent.subjects, randHumanities[:2]...)

	// Languages (choose 1)
	languages := []string{SPANISH, FRENCH, GERMAN}
	language := languages[rand.Intn(len(languages))]
	subjectChoiceStudent.subjects = append(subjectChoiceStudent.subjects, language)

	// Optionals (choose 1 from humanities, languages, or creatives)
	creatives := []string{COMPUTER_SCIENCE, DT, ART, MUSIC}
	options := append(append([]string{}, humanities...), creatives...)
	options = append(options, languages...)
	optional := options[rand.Intn(len(options))]
	subjectChoiceStudent.subjects = append(subjectChoiceStudent.subjects, optional)

	// Add to the class list
	*subjectChoicesClass = append(*subjectChoicesClass, subjectChoiceStudent)
}

// Generates student IDs
func generateStudents(numStudents int) []int {
	students := make([]int, numStudents)
	for i := 0; i < numStudents; i++ {
		students[i] = i + 1
	}
	return students
}

// Assigns subject choices to each student
func assignSubjectChoices(students []int) []StudentSubjects {
	subjectChoicesClass := []StudentSubjects{}
	for _, student := range students {
		chooseSubjectFunction(student, &subjectChoicesClass)
	}
	return subjectChoicesClass
}

// Calculates the fitness score of a candidate solution
func fitnessFunc(subjectChoicesClass []StudentSubjects, timetables []Timetable) int {
	points := 0

	for i, student := range subjectChoicesClass {
		subjectCount := make(map[string]int)
		for _, subject := range student.subjects {
			subjectCount[subject] = 0
		}
		for day := 0; day < 10; day++ {
			for period := 0; period < 8; period++ {
				subject := timetables[i].schedule[day][period]
				if _, ok := subjectCount[subject]; ok {
					subjectCount[subject]++
				}
			}
		}
		for subject, count := range subjectCount {
			required := subjectPeriodRequirements[subject] * 2 // For two weeks
			if count == required {
				points += 2 // Perfect match
			} else if count >= required-1 && count <= required+1 {
				points += 1 // Acceptable deviation
			}
		}
	}
	return points
}

// Calculates the maximum possible fitness score
func calculateMaxPossibleFitness(subjectChoicesClass []StudentSubjects) int {
	maxPoints := 0
	for _, student := range subjectChoicesClass {
		maxPoints += len(student.subjects) * 2
	}
	return maxPoints
}

// Generates initial timetables based on students' subject choices
func genStartTimetables(subjectChoicesClass []StudentSubjects) []Timetable {
	numStudents := len(subjectChoicesClass)
	timetables := make([]Timetable, numStudents)

	for i, student := range subjectChoicesClass {
		// Build a list of subjects according to the required periods
		subjectsWithCounts := make([]string, 0)
		for _, subject := range student.subjects {
			requiredPeriods := subjectPeriodRequirements[subject] * 2 // For two weeks
			for p := 0; p < requiredPeriods; p++ {
				subjectsWithCounts = append(subjectsWithCounts, subject)
			}
		}
		// Shuffle the subjects
		rand.Shuffle(len(subjectsWithCounts), func(i, j int) {
			subjectsWithCounts[i], subjectsWithCounts[j] = subjectsWithCounts[j], subjectsWithCounts[i]
		})

		// Assign to timetable
		index := 0
		for day := 0; day < 10; day++ {
			for period := 0; period < 8; period++ {
				if index < len(subjectsWithCounts) {
					timetables[i].schedule[day][period] = subjectsWithCounts[index]
					index++
				} else {
					// Assign free periods if no subjects left
					timetables[i].schedule[day][period] = FREE
				}
			}
		}
	}
	return timetables
}

// Performs mutation on a timetable
func mutateTimetable(timetable Timetable, studentSubjects []string, prob float64) Timetable {
	subjectsSet := make(map[string]bool)
	for _, subject := range studentSubjects {
		subjectsSet[subject] = true
	}
	// Collect positions of the student's subjects
	positions := make([][2]int, 0)
	for day := 0; day < 10; day++ {
		for period := 0; period < 8; period++ {
			subject := timetable.schedule[day][period]
			if subjectsSet[subject] {
				positions = append(positions, [2]int{day, period})
			}
		}
	}
	// Perform swaps with given probability
	numSwaps := int(float64(len(positions)) * prob)
	for i := 0; i < numSwaps; i++ {
		idx1 := rand.Intn(len(positions))
		idx2 := rand.Intn(len(positions))
		pos1 := positions[idx1]
		pos2 := positions[idx2]
		timetable.schedule[pos1[0]][pos1[1]], timetable.schedule[pos2[0]][pos2[1]] = timetable.schedule[pos2[0]][pos2[1]], timetable.schedule[pos1[0]][pos1[1]]
	}
	return timetable
}

// Performs mutation on a candidate solution
func mutateCandidateSolution(candidate CandidateSolution, subjectChoicesClass []StudentSubjects, prob float64) CandidateSolution {
	for i, timetable := range candidate.timetables {
		studentSubjects := subjectChoicesClass[i].subjects
		candidate.timetables[i] = mutateTimetable(timetable, studentSubjects, prob)
	}
	return candidate
}

// Performs tournament selection
func tournamentSelection(population []CandidateSolution, k int) CandidateSolution {
	best := population[rand.Intn(len(population))]
	for i := 1; i < k; i++ {
		competitor := population[rand.Intn(len(population))]
		if competitor.fitness > best.fitness {
			best = competitor
		}
	}
	return best
}

// Performs crossover between two candidate solutions
func crossover(parent1, parent2 CandidateSolution) CandidateSolution {
	numStudents := len(parent1.timetables)
	child := CandidateSolution{
		timetables: make([]Timetable, numStudents),
	}
	for i := 0; i < numStudents; i++ {
		if rand.Intn(2) == 0 {
			child.timetables[i] = parent1.timetables[i]
		} else {
			child.timetables[i] = parent2.timetables[i]
		}
	}
	return child
}

// Initializes the population
func initializePopulation(populationSize int, subjectChoicesClass []StudentSubjects) ([]CandidateSolution, CandidateSolution) {
	population := make([]CandidateSolution, populationSize)
	var bestSolution CandidateSolution
	bestFitness := 0

	for i := 0; i < populationSize; i++ {
		timetables := genStartTimetables(subjectChoicesClass)
		fitness := fitnessFunc(subjectChoicesClass, timetables)
		candidate := CandidateSolution{
			timetables: timetables,
			fitness:    fitness,
		}
		population[i] = candidate

		if fitness > bestFitness {
			bestFitness = fitness
			bestSolution = candidate
		}
	}
	return population, bestSolution
}

// Prints the fortnightly timetable for each student
func printFortnightlyTimetable(timetables []Timetable, students []int) {
	days := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday"}

	for studentNum, timetable := range timetables {
		fmt.Printf("\nStudent %d Timetable:\n\n", students[studentNum])

		fmt.Println("Week 1:")
		for day := 0; day < 5; day++ {
			fmt.Printf("%s: ", days[day])
			fmt.Println(timetable.schedule[day][:])
		}

		fmt.Println("\nWeek 2:")
		for day := 5; day < 10; day++ {
			fmt.Printf("%s: ", days[day-5])
			fmt.Println(timetable.schedule[day][:])
		}

		fmt.Println("\n========================================")
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// Generate students and their subject choices
	students := generateStudents(100)
	subjectChoicesClass := assignSubjectChoices(students)

	// Initialize population
	populationSize := 50
	population, bestSolution := initializePopulation(populationSize, subjectChoicesClass)
	bestFitness := bestSolution.fitness

	// Evolution parameters
	maxGenerations := 100
	mutationProbability := 0.05
	tournamentSize := 5

	maxPossibleFitness := calculateMaxPossibleFitness(subjectChoicesClass)
	fmt.Println("Max possible fitness:", maxPossibleFitness)

	for generation := 0; generation < maxGenerations; generation++ {
		newPopulation := make([]CandidateSolution, 0, populationSize)

		for i := 0; i < populationSize; i++ {
			// Selection
			parent1 := tournamentSelection(population, tournamentSize)
			parent2 := tournamentSelection(population, tournamentSize)

			// Crossover
			child := crossover(parent1, parent2)

			// Mutation
			child = mutateCandidateSolution(child, subjectChoicesClass, mutationProbability)

			// Evaluate fitness
			child.fitness = fitnessFunc(subjectChoicesClass, child.timetables)

			newPopulation = append(newPopulation, child)
		}

		// Elitism: keep the best individual
		population = append(newPopulation, bestSolution)

		// Sort population by fitness
		sort.Slice(population, func(i, j int) bool {
			return population[i].fitness > population[j].fitness
		})

		// Keep the best ones
		if len(population) > populationSize {
			population = population[:populationSize]
		}

		// Update best solution
		if population[0].fitness > bestFitness {
			bestFitness = population[0].fitness
			bestSolution = population[0]
			fmt.Printf("Generation %d: New best fitness = %d\n", generation, bestFitness)
		}

		// Check for convergence
		if bestFitness == maxPossibleFitness {
			fmt.Println("Optimal solution found.")
			break
		}
	}

	// Final fitness
	fmt.Println("Final best fitness:", bestFitness)

	// Output the best timetable
	printFortnightlyTimetable(bestSolution.timetables, students)
}
