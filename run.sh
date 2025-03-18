
##! RUN.sh
##! Author: hmza-sfyn
##! Publish: Sat MAR 15 - 2025

# Get the current directory
current_loc=$(pwd)

# Install QuackoLang
echo "alias quackolang='$current_loc/QuackoLang'" >> ~/.bashrc
#echo "alias quackolang='$current_loc/QuackoLang'" >> ~/.zshrc

# Install DuckPM and LangMang
echo "alias quacko='quackolang $current_loc/QuackoLanguageManager/LangMang.qk.cs'" >> ~/.bashrc
#echo "alias quacko='quackolang $current_loc/QuackoLanguageManager/LangMang.qk'" >> ~/.zshrc

echo "alias duckpm='quackolang $current_loc/DuckPM/duckPM.qk.cs'" >> ~/.bashrc
#echo "alias duckpm='quackolang $current_loc/DuckPM/duckPM.qk'" >> ~/.zshrc

mkdir ~/.quacko_lang_important_do_not_delete
mkdir ~/.quacko_lang_important_do_not_delete/universal_libs

# Refresh shell
source ~/.bashrc
#source ~/.zshrc

echo "Done setting up your env!"
echo "Type 'quackolang' to run QuackoLang."
echo "Type 'quacko --version' to check the version."
echo "Happy coding!"

bash