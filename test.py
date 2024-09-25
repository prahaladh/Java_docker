import subprocess
from concurrent.futures import ProcessPoolExecutor, as_completed

# Function to run a shell script with two parameters
def run_shell_script_with_two_params(script_path, param1, param2):
    try:
        # Run the shell script with two parameters
        result = subprocess.run([script_path, param1, param2], capture_output=True, text=True, check=True)
        return f"Output from {param1}, {param2}: {result.stdout}"
    except subprocess.CalledProcessError as e:
        return f"Error occurred while running the script with {param1}, {param2}: {e}"

# First parameter (constant)
param1 = "constant_param"

# List of changing second parameters
params2 = ["param2_1", "param2_2", "param2_3", "param2_4", "param2_5"]

# Path to your shell script
script_path = "./your_script.sh"

# Use ProcessPoolExecutor to run the script in parallel with two parameters
with ProcessPoolExecutor(max_workers=len(params2)) as executor:
    # Submit multiple jobs with the same first parameter and varying second parameters
    futures = [executor.submit(run_shell_script_with_two_params, script_path, param1, param2) for param2 in params2]

    # Collect results as they are completed
    for future in as_completed(futures):
        result = future.result()
        print(result)
