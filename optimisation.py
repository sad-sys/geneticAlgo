import random
import numpy as np

# Distance matrix for cities
distances = np.array([[0, 2, 9, 10],
                      [100, 0, 6, 4],
                      [15, 7, 0, 8],
                      [6, 3, 12, 0]])

# Number of cities
num_cities = len(distances)

# Function to calculate fitness (inverse of distance)
def calculate_fitness(route):
    total_distance = 0
    for i in range(len(route) - 1):
        total_distance += distances[route[i]][route[i + 1]]
    total_distance += distances[route[-1]][route[0]]  # Return to the start

    if total_distance == 0:
        total_distance = 1e-6  # Small value to prevent division by zero
    return 1 / total_distance  # Fitness is inverse of distance

# Generate random individual (list of cities)
def generate_individual():
    cities = list(range(num_cities))
    random.shuffle(cities)
    return cities

# Correct PMX crossover implementation
def pmx_crossover(parent1, parent2):
    size = len(parent1)
    p1, p2 = parent1[:], parent2[:]

    # Choose crossover points
    cxpoint1 = random.randint(0, size - 2)
    cxpoint2 = random.randint(cxpoint1 + 1, size - 1)

    # Create offspring placeholders
    child1 = [None]*size
    child2 = [None]*size

    # Copy crossover segment from parents to children
    child1[cxpoint1:cxpoint2] = parent1[cxpoint1:cxpoint2]
    child2[cxpoint1:cxpoint2] = parent2[cxpoint1:cxpoint2]

    # Mapping function
    def pmx_mapping(child, parent, cxpoint1, cxpoint2):
        for i in range(cxpoint1, cxpoint2):
            if parent[i] not in child:
                val = parent[i]
                pos = i
                while True:
                    val_in_child = child[pos]
                    pos_in_parent = parent.index(val_in_child)
                    if child[pos_in_parent] is None:
                        child[pos_in_parent] = val
                        break
                    else:
                        pos = pos_in_parent
        for i in range(size):
            if child[i] is None:
                child[i] = parent[i]

    # Apply mapping to both children
    pmx_mapping(child1, parent2, cxpoint1, cxpoint2)
    pmx_mapping(child2, parent1, cxpoint1, cxpoint2)

    return child1, child2

# Mutation operation (swap mutation)
def mutate(individual):
    idx1, idx2 = random.sample(range(num_cities), 2)
    individual = individual[:]
    individual[idx1], individual[idx2] = individual[idx2], individual[idx1]
    return individual

# Create initial population
def generate_population(size):
    return [generate_individual() for _ in range(size)]

# Tournament selection
def select_parents(population, k=3):
    selected = random.sample(population, k)
    selected = sorted(selected, key=lambda ind: calculate_fitness(ind), reverse=True)
    return selected[0]

# Main GA loop with elitism
def genetic_algorithm(population_size=100, generations=500, mutation_rate=0.01, elitism=True):
    # Generate initial population
    population = generate_population(population_size)

    for generation in range(generations):
        new_population = []

        # Elitism: carry over the best individual
        if elitism:
            best_individual = max(population, key=lambda ind: calculate_fitness(ind))
            new_population.append(best_individual)

        while len(new_population) < population_size:
            # Select parents
            parent1 = select_parents(population)
            parent2 = select_parents(population)
            # Perform crossover
            child1, child2 = pmx_crossover(parent1, parent2)
            # Perform mutation
            if random.random() < mutation_rate:
                child1 = mutate(child1)
            if random.random() < mutation_rate:
                child2 = mutate(child2)
            # Add children to new population
            new_population.extend([child1, child2])

        # Trim population to desired size
        population = new_population[:population_size]

        # Find the best solution in the population
        best_individual = max(population, key=lambda ind: calculate_fitness(ind))
        best_distance = 1 / calculate_fitness(best_individual)
        print(f'Generation {generation}: Best distance = {best_distance}, Route = {best_individual}')

    # Return the best solution found
    best_individual = max(population, key=lambda ind: calculate_fitness(ind))
    best_distance = 1 / calculate_fitness(best_individual)
    return best_individual, best_distance

# Run the genetic algorithm
best_route, best_distance = genetic_algorithm()
print(f'Best Route: {best_route}, Best Distance: {best_distance}')
