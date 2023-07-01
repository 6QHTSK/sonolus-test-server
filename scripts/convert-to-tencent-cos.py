import json

files = ["skins.json", "backgrounds.json", "effects.json", "particles.json",
         "engines.json"]
subKeys = [["data", "texture", "thumbnail"],
           ["configuration", "data", "image", "thumbnail"],
           ["audio", "data", "thumbnail"],
           ["data", "texture", "thumbnail"],
           ["data", "thumbnail", "configuration"]]
pwd = "../sonolus/"

for index in range(len(files)):
    with open(pwd + files[index], "r", encoding="utf-8") as f:
        data = json.load(f)
    for obj in data:
        for subkey in subKeys[index]:
            if obj[subkey]["url"].startswith("/sonolus/repository/"):
                obj[subkey]["url"] = "https://repository.ayachan.fun/sonolus/" + obj[subkey]["url"][20:]
    with open(pwd + "tencentCos-" + files[index], "w", encoding="utf-8") as target:
        json.dump(data, target, ensure_ascii=False)
