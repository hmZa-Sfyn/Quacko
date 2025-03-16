import os
import subprocess
import shutil
from tqdm import tqdm
from colorama import Fore, Style, init

# Initialize colorama
init(autoreset=True)

# Configuration
EXTENSION_NAME = "quacko-language"
EXTENSION_DIR = os.path.expanduser(f"~/.vscode/extensions/{EXTENSION_NAME}")
AUTHOR_NAME = "hmza-sfyn"

def run_command(command):
    """Run a shell command and check for errors."""
    result = subprocess.run(command, shell=True, text=True, capture_output=True)
    if result.returncode != 0:
        print(Fore.RED + f"Error: {result.stderr}")
        exit(1)
    return result.stdout

def print_header():
    """Print a colorful header."""
    print(Fore.CYAN + Style.BRIGHT + "====================================")
    print(Fore.CYAN + Style.BRIGHT + "   Quacko Extension Setup Script    ")
    print(Fore.CYAN + Style.BRIGHT + "====================================")
    print()

def print_step(step, description):
    """Print a step with a colorful prefix."""
    print(Fore.YELLOW + Style.BRIGHT + f"[Step {step}] " + Fore.WHITE + description)

def create_extension():
    """Create the Quacko extension."""
    print_header()

    # Step 1: Install required tools
    print_step(1, "Installing yo and generator-code...")
    with tqdm(total=100, desc=Fore.GREEN + "Progress", unit="%") as pbar:
        run_command("npm install -g yo generator-code")
        pbar.update(100)
    print(Fore.GREEN + "✓ Tools installed successfully!")

    # Check if yo is installed
    if not shutil.which("yo"):
        print(Fore.RED + "Error: 'yo' command not found. Please ensure npm global path is in your PATH.")
        print(Fore.YELLOW + "Run the following command and try again:")
        print(Fore.WHITE + f"export PATH=$PATH:{run_command('npm config get prefix').strip()}/bin")
        exit(1)

    # Step 2: Remove existing extension directory
    print_step(2, "Checking for existing extension directory...")
    if os.path.exists(EXTENSION_DIR):
        print(Fore.YELLOW + "Extension directory already exists. Removing it...")
        shutil.rmtree(EXTENSION_DIR)
        print(Fore.GREEN + "✓ Old extension directory removed!")
    else:
        print(Fore.GREEN + "✓ No existing extension directory found.")

    # Step 3: Generate the extension
    print_step(3, "Creating the Quacko extension...")
    with tqdm(total=100, desc=Fore.GREEN + "Progress", unit="%") as pbar:
        run_command(f"yo code --extensionType='New Language Support' --extensionName='{EXTENSION_NAME}' "
                    f"--extensionDescription='Syntax highlighting for Quacko (.qk) files' "
                    f"--extensionId='{EXTENSION_NAME}' --gitInit='Yes' --packageManager='npm'")
        pbar.update(100)

    # Verify that the directory was created
    if not os.path.exists(EXTENSION_NAME):
        print(Fore.RED + "Error: The extension directory was not created by 'yo code'.")
        print(Fore.YELLOW + "Please ensure 'yo code' is working correctly.")
        exit(1)

    print(Fore.GREEN + "✓ Extension generated successfully!")

    # Step 4: Move the extension to the correct directory
    print_step(4, "Moving extension to VS Code extensions directory...")
    with tqdm(total=100, desc=Fore.GREEN + "Progress", unit="%") as pbar:
        shutil.move(EXTENSION_NAME, EXTENSION_DIR)
        pbar.update(100)
    print(Fore.GREEN + f"✓ Extension moved to: {EXTENSION_DIR}")

    # Step 5: Configure the extension
    print_step(5, "Configuring the extension...")
    with tqdm(total=100, desc=Fore.GREEN + "Progress", unit="%") as pbar:
        # Update package.json
        package_json = f"""
        {{
            "name": "{EXTENSION_NAME}",
            "displayName": "Quacko Language",
            "description": "Syntax highlighting for Quacko (.qk) files",
            "version": "1.0.0",
            "publisher": "{AUTHOR_NAME}",
            "engines": {{
                "vscode": "^1.0.0"
            }},
            "categories": [
                "Programming Languages"
            ],
            "contributes": {{
                "languages": [{{
                    "id": "quacko",
                    "aliases": ["Quacko", "quacko"],
                    "extensions": [".qk"],
                    "configuration": "./language-configuration.json"
                }}],
                "grammars": [{{
                    "language": "quacko",
                    "scopeName": "source.quacko",
                    "path": "./syntaxes/quacko.tmLanguage.json"
                }}]
            }}
        }}
        """
        with open(os.path.join(EXTENSION_DIR, "package.json"), "w") as f:
            f.write(package_json)
        pbar.update(25)

        # Create syntax highlighting file
        os.makedirs(os.path.join(EXTENSION_DIR, "syntaxes"), exist_ok=True)
        syntax_file = """
        {
            "scopeName": "source.quacko",
            "name": "Quacko",
            "patterns": [
                {
                    "match": "\\b(let|fn|if|else|elif|return|import|and|or|enum|class|property|static|new|spawn|defer|try|catch|finally|throw|using|qw|case|is|in|for|while|do|break|continue|where|grep|map|true|false|nil)\\b",
                    "name": "keyword.control.quacko"
                },
                {
                    "match": "\\b(Int|UInt|Float|Bool|String|Array|Hash|Tuple|Nil)\\b",
                    "name": "storage.type.quacko"
                },
                {
                    "match": "\\b([0-9]+)\\b",
                    "name": "constant.numeric.quacko"
                },
                {
                    "match": "\\b([A-Za-z_][A-Za-z0-9_]*)\\b",
                    "name": "variable.other.quacko"
                },
                {
                    "begin": "\"",
                    "end": "\"",
                    "name": "string.quoted.double.quacko"
                },
                {
                    "begin": "`",
                    "end": "`",
                    "name": "string.quoted.raw.quacko"
                },
                {
                    "begin": "//",
                    "end": "$",
                    "name": "comment.line.double-slash.quacko"
                },
                {
                    "begin": "/\\*",
                    "end": "\\*/",
                    "name": "comment.block.quacko"
                }
            ],
            "repository": {}
        }
        """
        with open(os.path.join(EXTENSION_DIR, "syntaxes", "quacko.tmLanguage.json"), "w") as f:
            f.write(syntax_file)
        pbar.update(25)

        # Create language configuration file
        language_config = """
        {
            "comments": {
                "lineComment": "//",
                "blockComment": ["/*", "*/"]
            },
            "brackets": [
                ["{", "}"],
                ["[", "]"],
                ["(", ")"]
            ],
            "autoClosingPairs": [
                { "open": "{", "close": "}" },
                { "open": "[", "close": "]" },
                { "open": "(", "close": ")" },
                { "open": "\"", "close": "\"", "notIn": ["string"] },
                { "open": "`", "close": "`", "notIn": ["string"] }
            ],
            "surroundingPairs": [
                ["{", "}"],
                ["[", "]"],
                ["(", ")"],
                ["\"", "\""],
                ["`", "`"]
            ]
        }
        """
        with open(os.path.join(EXTENSION_DIR, "language-configuration.json"), "w") as f:
            f.write(language_config)
        pbar.update(25)

        print(Fore.GREEN + "✓ Extension configured successfully!")
        pbar.update(25)

    # Step 6: Notify the user
    print_step(6, "Finalizing...")
    print(Fore.GREEN + Style.BRIGHT + "Quacko extension created successfully!")
    print(Fore.GREEN + f"The extension has been placed in: {EXTENSION_DIR}")
    print(Fore.YELLOW + "Restart VS Code to activate the extension.")

if __name__ == "__main__":
    create_extension()