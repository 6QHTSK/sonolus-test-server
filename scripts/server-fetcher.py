import json
import os.path
import platform
import queue
import random
import ssl
import threading
import time
import requests

target_server = "https://servers.sonolus.com/bestdori/"
localization = "zhs"
requests.packages.urllib3.disable_warnings()

exit_flag = False
queueLock = threading.Lock()
queueSemaphore = threading.Semaphore(20)
workQueue = queue.Queue()
threads = []
thread_cnt = 10

headers = {'User-Agent': "Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/112.0"}


def get_json(url, params=None):
    for i in range(3):
        resp = requests.get(url, params=params, headers=headers, verify=False)
        if resp.status_code == 200 and len(resp.content) != 0:
            return resp.json()
    raise "Status Code Not 200 for 3 times"


class DownloadThread(threading.Thread):
    def __init__(self, q: queue.Queue):
        threading.Thread.__init__(self)
        self.q = q

    def run(self):
        while not exit_flag:
            queueLock.acquire()
            if not self.q.empty():
                remote_path = self.q.get()
                queueLock.release()
                download_file(remote_path)
                queueSemaphore.release()
            else:
                queueLock.release()


# only start with / should be downloaded, others will be view as remote repos
def download_file(remote_path):
    assets_path = os.path.join(os.getcwd(), "../sonolus")
    if remote_path.startswith("/sonolus/"):
        local_path = os.path.join(assets_path, remote_path.lstrip("/sonolus/"))
    else:
        local_path = os.path.join(assets_path, remote_path)
    if platform.system().lower() == "windows":
        local_path = local_path.replace("/", "\\")
    remote_url = os.path.join(target_server, remote_path.lstrip("/"))
    os.makedirs(os.path.dirname(local_path), exist_ok=True)
    if not os.path.exists(local_path) or os.path.getsize(local_path) == 0:
        while True:
            print("Start downloading %s" % remote_url)
            f = requests.get(remote_url, headers=headers, verify=False)
            if len(f.content) == 0 or f.status_code != 200:
                print("Download %s Failed({}), Sleep For 10+ Sec".format(f.status_code))
                time.sleep(random.random()*10 + 10)
            else:
                with open(local_path, "wb") as local_file:
                    local_file.write(f.content)
                print("Downloaded %s" % remote_url)
                time.sleep(random.random())
                break
    else:
        print("Already Downloaded %s" % remote_url)


def traverse_section(section_name, download_srl_list):
    i = 0
    page = 1
    items = []
    while i < page:
        if i == 0:
            print("Section:{} Page:{}".format(section_name, i + 1))
        else:
            print("Section:{} Page:{}/{}".format(section_name, i + 1, page))
        params = {"localization": localization, "page": i}
        url = os.path.join(target_server, "sonolus/%s/list" % section_name)
        info = get_json(url, params)
        page = info["pageCount"]
        items.extend(info["items"])
        for item in info["items"]:
            for srl in download_srl_list:
                if item[srl]["url"].startswith("/"):
                    queueSemaphore.acquire()
                    queueLock.acquire()
                    workQueue.put(item[srl]["url"])
                    queueLock.release()
        i += 1
        time.sleep(random.random())
    with open(os.path.join(os.getcwd(), "../sonolus/%s.json" % section_name), "w", encoding="utf8") as list_json_file:
        json.dump(items, list_json_file, ensure_ascii=False)


for i in range(thread_cnt):
    thread = DownloadThread(workQueue)
    thread.start()
    threads.append(thread)

traverse_section("skins", ["data", "texture", "thumbnail"])
traverse_section("backgrounds", ["configuration", "data", "image", "thumbnail"])
traverse_section("effects", ["audio", "data", "thumbnail"])
traverse_section("particles", ["data", "texture", "thumbnail"])
traverse_section("engines", ["data", "thumbnail", "configuration"])

exit_flag = True

for t in threads:
    t.join()

print("Finish fetching server %s" % target_server)
