# Nils Elde
# https://gitlab.com/nilsanderselde
# Sort words in Fonetik Ingliš dictionary according to custom alphabetical order

# Define order in which to sort words
alphabet = {c: i for i, c in enumerate('-.aábcdeéfgiíjklmnoóøprsštuúvzžh')}

# Load dictionary.txt as a string
dictionary = open('dictionary.txt', encoding='utf-8').readlines()

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