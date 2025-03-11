import os, time

while True:

    print("waiting for 60 seconds")
    time.sleep(60)

    print("gonna push now")
    os.system("git add . && git commit -m 'auto push' && git push")