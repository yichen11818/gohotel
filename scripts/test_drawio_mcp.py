from __future__ import annotations

import json
import subprocess
import sys
import time
from pathlib import Path


NODE = r"C:\nvm4w\nodejs\node.exe"
SERVER = Path.home() / "AppData" / "Roaming" / "npm" / "node_modules" / "@drawio" / "mcp" / "src" / "index.js"


def write_message(proc: subprocess.Popen[bytes], msg: dict) -> None:
    body = json.dumps(msg, ensure_ascii=False, separators=(",", ":")).encode("utf-8")
    proc.stdin.write(f"Content-Length: {len(body)}\r\n\r\n".encode("ascii") + body)
    proc.stdin.flush()


def read_message(proc: subprocess.Popen[bytes], timeout: float = 8) -> dict | None:
    deadline = time.time() + timeout
    header = b""
    while time.time() < deadline:
        b = proc.stdout.read(1)
        if b:
            header += b
            if header.endswith(b"\r\n\r\n"):
                break
            continue
        if proc.poll() is not None:
            break
        time.sleep(0.01)
    if not header:
        return None
    length = None
    for line in header.decode("ascii", "replace").splitlines():
        if line.lower().startswith("content-length:"):
            length = int(line.split(":", 1)[1].strip())
    if length is None:
        return {"bad_header": header.decode("ascii", "replace")}
    body = proc.stdout.read(length)
    return json.loads(body.decode("utf-8"))


def main() -> int:
    proc = subprocess.Popen(
        [NODE, str(SERVER)],
        stdin=subprocess.PIPE,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
    )
    assert proc.stdin and proc.stdout and proc.stderr
    write_message(proc, {
        "jsonrpc": "2.0",
        "id": 1,
        "method": "initialize",
        "params": {
            "protocolVersion": "2024-11-05",
            "capabilities": {},
            "clientInfo": {"name": "codex-drawio-test", "version": "1.0.0"},
        },
    })
    init = read_message(proc)
    print("INIT:", json.dumps(init, ensure_ascii=False)[:1000])
    write_message(proc, {"jsonrpc": "2.0", "method": "notifications/initialized", "params": {}})
    write_message(proc, {"jsonrpc": "2.0", "id": 2, "method": "tools/list", "params": {}})
    tools = read_message(proc)
    print("TOOLS:", json.dumps(tools, ensure_ascii=False)[:1000])
    proc.terminate()
    try:
        proc.wait(timeout=2)
    except subprocess.TimeoutExpired:
        proc.kill()
    err = proc.stderr.read().decode("utf-8", "replace")
    if err:
        print("STDERR:", err[:1000])
    if not init or "result" not in init:
        return 1
    if not tools or "result" not in tools:
        return 2
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
