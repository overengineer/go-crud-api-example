#!/usr/bin/env python3

from requests import get
from subprocess import Popen
import time, json, pytest

def is_json(myjson):
    try:
        json_object = json.loads(myjson)
    except ValueError as e:
        return False
    return True


@pytest.fixture
def server():
    p = Popen(["go", "run", "main.go"])
    time.sleep(3)
    yield
    p.terminate()

def test_albums(server):
    response = get("http://localhost:8088/albums")
    assert response.ok
    assert is_json(response.content)
    data = json.loads(response.content)
    assert type(data) == list
    assert len(data) > 0

if __name__=="__main__":
    pytest.main()