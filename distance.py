# Nils Elde
# https://gitlab.com/nilsanderselde
"""Calculate distance between words stored in text file and save results"""

import nltk

# Load words_for_distance.txt as a string
# Dictionary must be in the form of [fonetik]\t[traditional]\n
DICTIONARY_FILE = open('words_for_distance.txt', encoding='utf-8')
DICTIONARY = DICTIONARY_FILE.readlines()

# Open output file and write to it
OUTPUT = open('out.txt', 'w', encoding='utf-8')
i = 0
for line in DICTIONARY:
    # Split line into words on tab (since data comes from spreadsheet)
    words = line.split('\t')
    # Calculate the distance between the first word (fonetik) and second word
    # (traditional) with the newline character removed
    distance = nltk.edit_distance(words[0], words[1].replace('\n', ''))
    # Write the words and the Levenshtein distance value into output file
    OUTPUT.write(''.join([line.replace('\n', ''), '\t', str(distance)]))
    # Avoid adding extra line at end of file
    if i < len(DICTIONARY) - 1:
        OUTPUT.write('\n')
    i += 1

OUTPUT.close()
DICTIONARY_FILE.close()