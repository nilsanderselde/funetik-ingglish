# Nils Elde
# https://gitlab.com/nilsanderselde
"""Sort words in a text file by custom alphabetical order"""

# Define order in which to sort words
ALPHABET = {c: i for i, c in enumerate('\taäeoøiuywlrmnbpvfgkdtzsžšh')}

# Load words_to_sort.txt as a string
# Dictionary must separate word rows by new lines and values by tabs, including all 
# related fields in order to preserve data integrity.
DICTIONARY_FILE = open('words_to_sort.txt', encoding='utf-8')
DICTIONARY = DICTIONARY_FILE.readlines()

# Sort dictionary according to custom alphabet
DICTIONARY = sorted(DICTIONARY, key=lambda word: [ALPHABET.get(c, ord(c)) for c in word.lower()])

# Open output file and write to it
OUTPUT = open('out.txt', 'w', encoding='utf-8')
i = 0
for line in DICTIONARY:
    # line with the newline character removed
    OUTPUT.write(line.replace('\n', ''))
    # Avoid adding extra line at end of file
    if i < len(DICTIONARY) - 1:
        OUTPUT.write('\n')
    i += 1

OUTPUT.close()
DICTIONARY_FILE.close()