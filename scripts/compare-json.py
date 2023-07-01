import json
from deepdiff import DeepDiff


def compare_json_files(file1, file2):
    with open(file1, 'r') as f:
        json1 = json.load(f)
    with open(file2, 'r') as f:
        json2 = json.load(f)
    diff = DeepDiff(json1, json2,ignore_order=True)
    return diff.to_json()


# 示例
diff = compare_json_files('42620.json', 'gen.json')
print(diff)
