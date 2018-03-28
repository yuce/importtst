#! /usr/bin/env python3

import sys
import random

def main():
    if len(sys.argv) != 5:
        print("Usage: %s max_row_id max_column_id bit_count output_file.csv" % sys.argv[0])
        sys.exit(1)
    
    max_row_id = int(sys.argv[1])
    max_col_id = int(sys.argv[2])
    bit_count = int(sys.argv[3])
    path = sys.argv[4]

    with open(path, "w") as f:
        for i in range(bit_count):
            row_id = random.randint(0, max_row_id - 1)
            col_id = random.randint(0, max_col_id - 1)
            f.write("%s,%s\n" % (row_id, col_id))

if __name__ == "__main__":
    main()