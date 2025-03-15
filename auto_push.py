import os, time

while True:

    
    print("gonna push now")
    os.system("git add . && git commit -m 'auto push' && git push")

    print("waiting for 60 seconds")
    time.sleep(60)
