# Nils Elde
# https://gitlab.com/nilsanderselde
# Calculate distance between words stored in text file and save results
import nltk

# Load words_for_distance.txt as a string
# Dictionary must be in the form of [fonetik]\t[traditional]\n
dictionary_file = open('words_for_distance.txt', encoding='utf-8')
dictionary = dictionary_file.readlines()

# Open output file
output = open('out.txt', 'w', encoding='utf-8')

for line in dictionary:
    # Split line into words on tab (since data comes from spreadsheet)
    words = line.split('\t')
    # Calculate the distance between the first word (fonetik) and second word
    # (traditional) with the newline character removed
    distance = nltk.edit_distance(words[0], words[1].replace('\n', ''))
    # Write the words and the Levenshtein distance value into output file
    output.write(''.join([line.replace('\n', ''), '\t', str(distance), '\n']))

output.close()
dictionary_file.close()