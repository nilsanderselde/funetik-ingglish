# Nils Elde
# https://gitlab.com/nilsanderselde
"""Sort words in a text file by custom alphabetical order"""

# Define order in which to sort words
ALPHABET = {c: i for i, c in enumerate('-.aábdeéfghiíklmnoóøprstuúvz')}

# Load words_to_sort.txt as a string
# Dictionary must separate word rows by new lines and values by tabs, including all 
# related fields in order to preserve data integrity.
DICTIONARY_FILE = open('words_to_sort.txt', encoding='utf-8')
DICTIONARY = DICTIONARY_FILE.readlines()

# Sort dictionary according to custom alphabet
DICTIONARY = sorted(DICTIONARY, key=lambda word: [ALPHABET.get(c, ord(c)) for c in word.lower()])

# Save output to file
OUTPUT = open('out.txt', 'w', encoding='utf-8')
for word in DICTIONARY:
    OUTPUT.write(word)

OUTPUT.close()
DICTIONARY_FILE.close()