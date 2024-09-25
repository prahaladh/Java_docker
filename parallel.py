import subprocess
from concurrent.futures import ProcessPoolExecutor, as_completed

# Function to run a shell script
def run_shell_script(script_path):
    try:
        # Run the shell script
        result = subprocess.run([script_path], capture_output=True, text=True, check=True)
        return result.stdout
    except subprocess.CalledProcessError as e:
        return f"Error occurred while running the script: {e}"

# Number of times to run the script
num_runs = 5

# Path to your shell script
script_path = "./your_script.sh"

# Use ProcessPoolExecutor for true parallelism
with ProcessPoolExecutor(max_workers=num_runs) as executor:
    # Submit multiple jobs to run the script concurrently using separate processes
    futures = [executor.submit(run_shell_script, script_path) for _ in range(num_runs)]

    # Collect results as they are completed
    for future in as_completed(futures):
        result = future.result()
        print(result)
