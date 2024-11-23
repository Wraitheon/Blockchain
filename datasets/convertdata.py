import os
import pandas as pd

def process_file(file):
    try:
        # Try to read as an actual XLS file
        xls_file = pd.ExcelFile(file, engine='xlrd')  # Replace with 'pyxlsb' if needed
        for sheet_name in xls_file.sheet_names:
            df = xls_file.parse(sheet_name)
            csv_file_name = f"{os.path.splitext(file)[0]}_{sheet_name}.csv"
            df.to_csv(csv_file_name, index=False)
        print(f"Processed as XLS: {file}")
    except Exception:
        try:
            # If reading as XLS fails, assume it's a CSV
            # Try multiple encodings for CSV files
            for encoding in ['utf-8', 'latin1', 'ISO-8859-1']:
                try:
                    df = pd.read_csv(file, encoding=encoding)
                    csv_file_name = f"{os.path.splitext(file)[0]}.csv"
                    df.to_csv(csv_file_name, index=False)
                    print(f"Processed as CSV with encoding {encoding}: {file}")
                    break  # Exit the loop once successful
                except Exception:
                    continue
            else:
                raise ValueError("Could not process as CSV with any encoding.")
        except Exception as e:
            print(f"Failed to process file {file}: {e}")

# Loop through files in the current directory
for file in os.listdir('.'):
    if file.endswith('.xls'):
        process_file(file)
