import random
import copy

'''
1) First generate a large number of students as a list
2) They get options
3) Requirements for each student 
    a) Each student must have the cores english and maths
    b) Each student must have all three sciences
    c) They must have 2 of three humanities
    d) One of three languages 
    e) One other choice out of creatives, languages, humanities
4) Timetable has 8 periods in a day for 5 days a week
5) Student must have all the subjects every two weeks
'''
'1) First generate a large number of students as a list'
students = []
for i in range (1,101):
    students.append(i)
print(len(students))

subjectChoicesClass = []
'Now attach each student to a particular subject'
def chooseSubjectFunction(student, subjectChoicesClass):
    subjectChoiceStudent = []
    'First add the students'
    subjectChoiceStudent.append(student)

    'Then the core'
    subjectChoiceStudent.append("MATHS")
    subjectChoiceStudent.append("ENGLISH")

    'Then the sciences'
    subjectChoiceStudent.append("BIOLOGY")
    subjectChoiceStudent.append("CHEMISTRY")
    subjectChoiceStudent.append("PHYSICS")

    'Then one out of the humanities'
    humanities = ['GEOGRAPHY','HISTORY','ECONOMICS']
    subjectChoiceStudent.append(random.sample(humanities,1)[0])

    'Then one out of the Languages'
    languages = ['SPANISH','FRENCH','GERMAN']
    subjectChoiceStudent.append(random.sample(languages,1)[0])

    'Then one out of any of the optionals'
    creatives = ['COMPUTER SCIENCE', 'D&T', 'ART','MUSIC']
    subjectChoiceStudent.append(random.sample(humanities + languages + creatives,1)[0])


    'Finish adding to the class'
    subjectChoicesClass.append(subjectChoiceStudent)

    return subjectChoicesClass

for i in range(1,101-1):
    chooseSubjectFunction(students[i],subjectChoicesClass)



def fitnessFunc(subjectChoicesClass, timetables, points):
    for i in range(len(subjectChoicesClass)):
        subject_counted = set()  # Track which subjects have already been counted for the student
        for j in range(1, len(subjectChoicesClass[i])):  # Loop through the subjects the student chose
            subject = subjectChoicesClass[i][j]
            # Check if the subject appears exactly twice in the timetable and hasn't been counted before
            if subject not in subject_counted and timetables[i].count(subject) == 2:
                points += 1
                subject_counted.add(subject)  # Ensure the subject is only counted once per student
    return points

'8 X 5 X 2 - 2 periods in a week '


def genStartTimetables():
    timetables = []

    for i in range (0,101-1):
        timetables.append([])

    core = ['MATHS','ENGLISH']

    sciences = ['BIOLOGY','CHEMISTRY','PHYSICS']

    humanities = ['GEOGRAPHY','HISTORY','ECONOMICS']

    'Then one out of the Languages'
    languages = ['SPANISH','FRENCH','GERMAN']

    'Then one out of any of the optionals'
    creatives = ['COMPUTER SCIENCE', 'D&T', 'ART','MUSIC']

    for studentNum in range(0,len(students)):
        for i in range (0,78):
            timetables[studentNum].append(random.sample(humanities + languages + creatives + core + sciences,1)[0])
    return timetables

startTimetables = genStartTimetables()
fitnessFunc(subjectChoicesClass, startTimetables, 0)
def tournament_selection(timetables, subjectChoicesClass, tournament_size=5):
    selected = random.sample(timetables, tournament_size)
    best_timetable = None
    best_fitness = -1
    for table in selected:
        fitness = fitnessFunc(subjectChoicesClass, [table], 0)
        if fitness > best_fitness:
            best_fitness = fitness
            best_timetable = table
    return best_timetable


def mutateTable(table, prob):
    core = ['MATHS','ENGLISH']

    sciences = ['BIOLOGY','CHEMISTRY','PHYSICS']

    humanities = ['GEOGRAPHY','HISTORY','ECONOMICS']

    'Then one out of the Languages'
    languages = ['SPANISH','FRENCH','GERMAN']

    'Then one out of any of the optionals'
    creatives = ['COMPUTER SCIENCE', 'D&T', 'ART','MUSIC']

    for i in range(0,len(table)):
        for j in range(0,len(table[i])):
            if random.random()<prob:
                table[i][j] = random.sample(humanities + languages + creatives + core + sciences,1)[0]
    return table
def tournamentSelection(bitStrings, k, parentSize):
    parents = []
    for i in range(parentSize):
        for j in range(k):
            if j == 0:
                best = random.choice(bitStrings)
            else:
                choice = random.choice(bitStrings)
                if fitnessFunc(choice) > fitnessFunc(best):
                    best = choice
        parents.append(best)
    return parents
timetable = startTimetables
print("Starting Fitness:", fitnessFunc(subjectChoicesClass,timetable,0))
for i in range(1000):  # Iterating over 10 rounds
    pool = []
    poolFitness = []
    current_fitness = fitnessFunc(subjectChoicesClass, timetable, 0)  # Fitness of current timetable


    for x in range(0,100):
        tableToMutate = copy.deepcopy(timetable)
        newTable = mutateTable(tableToMutate,0.01)  # Generate a new random timetable
        new_fitness = fitnessFunc(subjectChoicesClass, newTable, 0)  # Fitness of new timetabl

        pool.append(newTable)
        poolFitness.append(new_fitness)

    bestFitnessFromPool = max(poolFitness)
    bestIndex = poolFitness.index(bestFitnessFromPool)

    if new_fitness > current_fitness:  # If the new timetable is better
        timetable = pool[bestIndex]  # Accept the new timetable
    else:
        # If the new fitness is not better, do nothing (keep the old timetable)
        continue

print("Final fitness", fitnessFunc(subjectChoicesClass, timetable, 0))

print (timetable[5])
def print_fortnightly_timetable(timetable, students):
    days = ['Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday']
    
    for student_num, student_timetable in enumerate(timetable):
        print(f"\nStudent {students[student_num]} Timetable:\n")
        
        # Organize the timetable into a 2-week schedule (each week having 5 days x 8 periods)
        week1 = student_timetable[:40]  # First week's timetable
        week2 = student_timetable[40:]  # Second week's timetable
        
        # Print Week 1
        print("Week 1:")
        for day in range(5):  # 5 days
            print(f"{days[day]}: ", end="")
            print(week1[day*8:(day+1)*8])  # Print 8 periods for each day

        # Print Week 2
        print("\nWeek 2:")
        for day in range(5):  # 5 days
            print(f"{days[day]}: ", end="")
            print(week2[day*8:(day+1)*8])  # Print 8 periods for each day
        
        print("\n" + "="*40)  # Separator between student timetables

# Example usage
print_fortnightly_timetable(timetable, students)
