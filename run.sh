
##! RUN.sh
##! Author: hmza-sfyn
##! Publish: Sat MAR 15 - 2025

# Get the current directory
current_loc=$(pwd)

# Install QuackoLang
echo "alias quackolang='$current_loc/QuackoLang'" >> ~/.bashrc
echo "alias quackolang='$current_loc/QuackoLang'" >> ~/.zshrc

echo "alias retro='$current_loc/DuckPM/retro/retro'" >> ~/.bashrc
echo "alias retro='$current_loc/DuckPM/retro/retro'" >> ~/.zshrc

mkdir ~/.quacko_lang_important_do_not_delete
mkdir ~/.quacko_lang_important_do_not_delete/universal_libs

# Refresh shell
source ~/.bashrc
#source ~/.zshrc

echo "Done setting up your env!"
echo "Type 'quackolang' to run QuackoLang."
echo "Type 'duckpm --langv' to check the version."
echo "Happy coding!"

bash