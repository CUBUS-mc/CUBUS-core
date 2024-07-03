import os
import subprocess
import toml

with open('shared/translation/locals/active.en.toml', 'r') as file:
    en_translations = toml.load(file)

owo_translations = {}

for key, value in en_translations.items():
    result = subprocess.run(['uwuifyy', '--text', value], capture_output=True, text=True)
    if result.returncode == 0:
        owo_translations[key] = result.stdout.strip()
    else:
        print(f'Error transforming value for key "{key}": {result.stderr}')
        owo_translations[key] = value

with open('shared/translation/locals/active.ky.toml', 'w') as file:
    toml.dump(owo_translations, file)