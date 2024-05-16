# This script is needed because the goi18n tool does not support extracting strings when they are not directly wrapped in i18n.Message()
import re
import toml
import os

def extract_strings_from_go_files(directory):
    strings = set()
    for root, dirs, files in os.walk(directory):
        for file in files:
            if file.endswith(".go"):
                with open(os.path.join(root, file), 'r') as f:
                    content = f.read()
                    matches = re.findall(r'T\("([^"]*)"\)', content)
                    strings.update(matches)
    return strings

def write_strings_to_toml(strings, filename):
    messages = {s: {"other": s} for s in strings}
    with open(filename, 'w') as f:
        toml.dump(messages, f)

strings = extract_strings_from_go_files('./')
write_strings_to_toml(strings, './translation/locals/active.en.toml')

toml_files = [f for f in os.listdir('./translation/locals') if f.startswith('active.') and f.endswith('.toml')]
command = 'goi18n merge -outdir ./translation/locals/' + ' ' + ' '.join(f'./translation/locals/{f}' for f in toml_files)
os.system(command)
