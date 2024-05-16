import re
import toml
import os
import glob
import logging

def load_existing_strings(directory):
    strings = set()
    for filename in glob.glob(os.path.join(directory, '*.toml')):
        with open(filename, 'r') as f:
            data = toml.load(f)
            strings.update(data.keys())
    return strings

def extract_strings_from_go_files(directory, existing_strings):
    strings = set()
    unfiltered_strings = set()
    constants = {}
    for root, dirs, files in os.walk(directory):
        for file in files:
            if file.endswith(".go"):
                logging.info(f"Checking file: {os.path.join(root, file)}")
                with open(os.path.join(root, file), 'r') as f:
                    content = f.read()
                    const_matches = re.findall(r'(\w+)\s+:=\s+"([^"]*)"', content)
                    for match in const_matches:
                        constants[match[0]] = match[1]
                    matches = re.findall(r'T\(([^)]*)\)', content)
                    for match in matches:
                        if match.startswith("\"") and match.endswith("\""):
                            match_string = match[1:-1]
                        elif match in constants:
                            match_string = constants[match]
                        if match_string not in existing_strings:
                            strings.add(match_string)
                        unfiltered_strings.add(match_string)
    warn_unused_translations(existing_strings, unfiltered_strings)
    return strings

def write_strings_to_toml(strings, directory):
    if len(strings) == 0:
        logging.info(f"No new strings to add to {directory}")
        return
    for filename in glob.glob(os.path.join(directory, '*.toml')):
        with open(filename, 'r') as f:
            data = toml.load(f)
        data.update({f'{s}': {"other": s} for s in strings})
        with open(filename, 'w') as f:
            toml.dump(data, f)
    logging.info(f"Added {len(strings)} new strings to {directory} (keys: {strings})")

def warn_unused_translations(existing_strings, extracted_strings):
    unused_translations = existing_strings - extracted_strings
    for unused in unused_translations:
        logging.warning(f"Unused translation: {unused}")

logging.basicConfig(level=logging.INFO)
existing_strings = load_existing_strings('./translation/locals')
extracted_strings = extract_strings_from_go_files('./', existing_strings)
write_strings_to_toml(extracted_strings, './translation/locals')