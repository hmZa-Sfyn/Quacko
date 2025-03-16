#!/bin/bash

# Script to create a VS Code extension for .qk files with syntax highlighting

# Step 1: Install required tools
echo "Installing yo and generator-code..."
npm install -g yo generator-code

# Step 2: Create the extension
echo "Creating the Quacko extension..."
EXTENSION_NAME="quacko-language"
EXTENSION_DIR="$HOME/.vscode/extensions/$EXTENSION_NAME"

# Check if the extension directory already exists
if [ -d "$EXTENSION_DIR" ]; then
    echo "Extension directory already exists. Removing it..."
    rm -rf "$EXTENSION_DIR"
fi

# Generate the extension
yo code --extensionType="New Language Support" --extensionName="$EXTENSION_NAME" --extensionDescription="Syntax highlighting for Quacko (.qk) files" --extensionId="$EXTENSION_NAME" --gitInit="Yes" --packageManager="npm"

# Move the generated extension to the correct directory
echo "Moving extension to VS Code extensions directory..."
mv "$EXTENSION_NAME" "$EXTENSION_DIR"

# Step 3: Configure the extension
echo "Configuring the extension..."

# Update package.json
cat <<EOL > "$EXTENSION_DIR/package.json"
{
    "name": "$EXTENSION_NAME",
    "displayName": "Quacko Language",
    "description": "Syntax highlighting for Quacko (.qk) files",
    "version": "1.0.0",
    "publisher": "your-name",
    "engines": {
        "vscode": "^1.0.0"
    },
    "categories": [
        "Programming Languages"
    ],
    "contributes": {
        "languages": [{
            "id": "quacko",
            "aliases": ["Quacko", "quacko"],
            "extensions": [".qk"],
            "configuration": "./language-configuration.json"
        }],
        "grammars": [{
            "language": "quacko",
            "scopeName": "source.quacko",
            "path": "./syntaxes/quacko.tmLanguage.json"
        }]
    }
}
EOL

# Create syntax highlighting file
mkdir -p "$EXTENSION_DIR/syntaxes"
cat <<EOL > "$EXTENSION_DIR/syntaxes/quacko.tmLanguage.json"
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
EOL

# Create language configuration file
cat <<EOL > "$EXTENSION_DIR/language-configuration.json"
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
EOL

# Step 4: Notify the user
echo "Quacko extension created successfully!"
echo "The extension has been placed in: $EXTENSION_DIR"
echo "Restart VS Code to activate the extension."