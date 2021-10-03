#!/usr/bin/env python3

from requests import get
from subprocess import run
from threading import Thread
import time
import json

def is_json(myjson):
    try:
        json_object = json.loads(myjson)
    except ValueError as e:
        return False
    return True

def start_server():
    run(["go", "run", "main.go"])

Thread(target=start_server, daemon=True).start()

time.sleep(3)

def test_albums():
    response = get("http://localhost:8080/albums")
    assert response.ok
    assert is_json(response.content)
    data = json.loads(response.content)
    assert type(data) == list
    assert len(data) > 0

if __name__=="__main__":
    test_albums()