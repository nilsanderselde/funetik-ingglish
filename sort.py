# Nils Elde
# https://gitlab.com/nilsanderselde
# https://docs.google.com/spreadsheets/d/1Y-NClJDkBJsc3roRPA0Mzo04YCKjlAL8J8pJApCd7mQ/edit?usp=sharing
# Sort words in Fonetik Ingliš dictionary according to custom alphabetical order

# Define order in which to sort words
alphabet = {c: i for i, c in enumerate('-.aábcdeéfgiíjklmnoóøprsštuúvzžh')}

# Load words_to_sort.txt as a string
# Dictionary must separate word rows by new lines and values by tabs, including all 
# related fields in order to preserve data integrity.
dictionary_file = open('words_to_sort.txt', encoding='utf-8')
dictionary = dictionary_file.readlines()

# Append newline character to dictionary so all lines will be processed the same,
# without two lines being glued together in output
dictionary += '\n'

# Sort dictionary according to custom alphabet
dictionary = sorted(dictionary, key=lambda word: [alphabet.get(c, ord(c)) for c in word.lower()])

# Save output to file
output = open('out.txt', 'w', encoding='utf-8') 
for word in dictionary:
    output.write(word)

output.close()
dictionary_file.close()