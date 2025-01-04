#!/bin/bash

# Path to the program
PROGRAM="/home/ubuntu/pacshare/pacshare"

# Function to run the program
run_program() {
    echo "Starting pacshare..."
    "$PROGRAM"
    
    # Check the exit status
    if [ $? -ne 0 ]; then
        echo "pacshare crashed with exit code $?. Restarting..."
        return 1
    else
        echo "pacshare exited normally."
        return 0
    fi
}

# Main loop
while true; do
    run_program
    
    # If the program exited normally, break the loop
    if [ $? -eq 0 ]; then
        break
    fi

    rm /home/ubuntu/pacshare/public/src/videos/*.*
    
    # Wait for a short time before restarting
    sleep 5
done

echo "pacshare has stopped running."