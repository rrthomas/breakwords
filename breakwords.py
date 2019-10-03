#!/usr/bin/env python3
'''
Generate wordlist:

occurs --nocount 100-0.txt > wordlist.txt

Generate space-less Shakespeare:

cp 100-0.txt 100-0-nospace.txt
rpl --quiet ' ' '' 100-0-nospace.txt
'''

import sys, argparse
try:
    import re2 as re
except ImportError:
    import re


MAX_MATCHES = 30


# Parse command-line arguments
parser = argparse.ArgumentParser(description="Add spaces between words in text.")
parser.add_argument(
    'wordlist_file',
    metavar='WORDLIST-FILE',
    help='word list to use',
)
args = parser.parse_args()


# Read word list, downcase and sort by decreasing word length
with open(args.wordlist_file) as f:
    lexicon = list(set(l.strip().lower() for l in f.readlines()))
lexicon.sort(key=lambda x: len(x), reverse=True)

# FIXME: allow punctuation
lexicon_regex = '|'.join(lexicon)
# It's faster to match with a one-or-more regex:
lexicon_match_regex = '({})+'.format(lexicon_regex)
# Then use a repeated regex to extract the words:
lexicon_extract_regex = '({})?'.format(lexicon_regex) * MAX_MATCHES

def breakwords(text, flags=0):
    m = re.fullmatch(lexicon_match_regex, text, flags=flags)
    if m == None:
        return '(no matches)'
    m = re.fullmatch(lexicon_extract_regex, text, flags=flags)
    return [match for match in m.groups() if match is not None]

for l in [l.strip() for l in sys.stdin.readlines()]:
    print('{}: {}'.format(l, breakwords(l, flags=re.IGNORECASE)))
