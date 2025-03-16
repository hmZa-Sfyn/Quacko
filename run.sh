
##! RUN.sh
##! Author: hmza-sfyn
##! Publish: Sat MAR 15 - 2025

# Get the current directory
current_loc=$(pwd)

# Install QuackoLang
echo "alias quackolang='$current_loc/QuackoLang'" >> ~/.bashrc
#echo "alias quackolang='$current_loc/QuackoLang'" >> ~/.zshrc

# Install DuckPM and LangMang
echo "alias quacko='quackolang $current_loc/QuackoLanguageManager/LangMang.qk'" >> ~/.bashrc
#echo "alias quacko='quackolang $current_loc/QuackoLanguageManager/LangMang.qk'" >> ~/.zshrc

echo "alias duckpm='quackolang $current_loc/DuckPM/duckPM.qk'" >> ~/.bashrc
#echo "alias duckpm='quackolang $current_loc/DuckPM/duckPM.qk'" >> ~/.zshrc

# Refresh shell
source ~/.bashrc
#source ~/.zshrc

echo "Done setting up your env!"
echo "Type 'quackolang' to run QuackoLang."
echo "Type 'quacko --version' to check the version."
echo "Happy coding!"

bash