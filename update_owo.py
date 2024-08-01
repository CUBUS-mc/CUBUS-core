# requires: uwuifyy (https://lib.rs/install/uwuifyy)
import os
import subprocess
import toml

with open('shared/translation/locals/active.en.toml', 'r', encoding='utf-8') as file:
    en_translations = toml.load(file)

owo_translations = {}

for key, value in en_translations.items():
    result = subprocess.run(['uwuifyy', '--text', value], capture_output=True, text=True)
    if result.returncode == 0:
        owo_translations[key] = result.stdout.strip()
        if owo_translations[key].endswith(':'): owo_translations[key] += ' '
    else:
        print(f'Error transforming value for key "{key}": {result.stderr}')
        owo_translations[key] = value
    print(f'{key}: {value} -> {owo_translations[key]}')

with open('shared/translation/locals/active.ky.toml', 'w') as file:
    toml.dump(owo_translations, file)