import subprocess
from concurrent.futures import ProcessPoolExecutor, as_completed

# Function to run a shell script with parameters
def run_shell_script_with_params(script_path, param):
    try:
        # Run the shell script with the given parameter
        result = subprocess.run([script_path, param], capture_output=True, text=True, check=True)
        return f"Output from {param}: {result.stdout}"
    except subprocess.CalledProcessError as e:
        return f"Error occurred while running the script with {param}: {e}"

# List of parameters to pass to the script
params = ["param1", "param2", "param3", "param4", "param5"]

# Path to your shell script
script_path = "./your_script.sh"

# Use ProcessPoolExecutor to run the script in parallel with different parameters
with ProcessPoolExecutor(max_workers=len(params)) as executor:
    # Submit multiple jobs with different parameters
    futures = [executor.submit(run_shell_script_with_params, script_path, param) for param in params]

    # Collect results as they are completed
    for future in as_completed(futures):
        result = future.result()
        print(result)
