##! RUN.sh
##! Author: hmza-sfyn
##! Publish: Sat MAR 15 - 2025

current_loc=$(pwd)

# Install QuackoLang
echo "\nalias quackolang='$current_loc/QuackoLang'" >> ~/.bashrc
echo "\nalias quackolang='$current_loc/QuackoLang'" >> ~/.zshrc

# Refeash BASH
bash

echo "Done settingup your env!"
echo "type: `quackolang` to run QuackoLang"

