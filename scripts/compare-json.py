import json
import sys

from deepdiff import DeepDiff


def compare_json_files(file1, file2):
    with open(file1, 'r') as f:
        json1 = json.load(f)
    with open(file2, 'r') as f:
        json2 = json.load(f)
    diff = DeepDiff(json1, json2,ignore_order=True)
    return diff.to_json()


if __name__ == "__main__":
    if len(sys.argv) != 3:
        print("Compare script need 2 args")
    file1_path = sys.argv[1]
    file2_path = sys.argv[2]
    diff = compare_json_files(file1_path, file2_path)
    print(diff)
