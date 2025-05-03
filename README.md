# Image Down Loader

This tool downloads images from a CSV file containing URLs and optionally processes them before saving them to disk. It's designed for parallel downloading and processing to improve efficiency.


## Key Features

*   **CSV Input:** Reads URLs and filenames from a CSV file.
*   **Parallel Downloading:** Downloads images concurrently using multiple workers.
*   **Preprocessing:** Supports image preprocessing pipelines (currently includes resizing).
*   **Skip Existing Files:** Can skip downloading and processing files that already exist in the output directory.


## Configuration

The tool requires a configuration file (details on the config file structure are not provided in this code snippet).  Key configuration options include:

*   `--input`: Path to the input CSV file.
*   `--dir`: Output directory for saving images.
*   `--url-filed`: The column name of URL in CSV file
*   `--file-filed`: The column name of File in CSV file
*   `--timeout`: Download timeout. Format: `10ms`, `60s`, `1m`, etc
*   `--resize`: Resizing parameters (width and height). Format: `WxH`. For example: `512x512`
*   `--workers`: Number of workers for parallel downloading.
*   `--reload`: Whether to reload images if exists

Usage:

```
idl --input <file.csv> --resize 500x500 --dir <directory/to/downloaded/images>
```
