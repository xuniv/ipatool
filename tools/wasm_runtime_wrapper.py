#!/usr/bin/env python3

import json
import os
import shutil
import subprocess
import sys
from typing import Any, Dict, List, Optional, Tuple

EXIT_OK = 0
EXIT_INVALID_INPUT_JSON = 10
EXIT_UNSUPPORTED_RUNTIME = 11
EXIT_RUNTIME_NOT_FOUND = 12
EXIT_MODULE_NOT_FOUND = 13
EXIT_RUNTIME_COMMAND_FAILED = 20
EXIT_STDOUT_NOT_JSON = 21
EXIT_INTERNAL_ERROR = 22


def emit(payload: Dict[str, Any], exit_code: int) -> None:
    sys.stdout.write(json.dumps(payload, separators=(",", ":"), ensure_ascii=False))
    sys.stdout.write("\n")
    raise SystemExit(exit_code)


def debug_log(enabled: bool, message: str) -> None:
    if enabled:
        sys.stderr.write(f"[wasm-wrapper] {message}\n")


def runtime_command(runtime: str, module: str, args: List[str]) -> List[str]:
    if runtime == "wasmtime":
        return [runtime, "run", module, *args]
    if runtime == "wasmer":
        return [runtime, "run", module, "--", *args]
    raise ValueError(f"unsupported runtime: {runtime}")


def resolve_runtime(requested: str) -> Tuple[Optional[str], Optional[str]]:
    if requested == "auto":
        for candidate in ("wasmtime", "wasmer"):
            if shutil.which(candidate):
                return candidate, None
        return None, "no supported runtime found in PATH (tried: wasmtime, wasmer)"

    if requested not in ("wasmtime", "wasmer"):
        return None, f"unsupported runtime '{requested}'"

    if not shutil.which(requested):
        return None, f"runtime '{requested}' not found in PATH"

    return requested, None


def parse_input(raw: str) -> Dict[str, Any]:
    payload = json.loads(raw)
    if not isinstance(payload, dict):
        raise ValueError("input must be a JSON object")

    runtime = payload.get("runtime", "auto")
    if not isinstance(runtime, str):
        raise ValueError("runtime must be a string")

    module = payload.get("module")
    if not isinstance(module, str) or module.strip() == "":
        raise ValueError("module is required and must be a non-empty string")

    args = payload.get("args", [])
    if not isinstance(args, list) or not all(isinstance(item, str) for item in args):
        raise ValueError("args must be an array of strings")

    env = payload.get("env", {})
    if not isinstance(env, dict) or not all(isinstance(k, str) and isinstance(v, str) for k, v in env.items()):
        raise ValueError("env must be an object with string key/value pairs")

    cwd = payload.get("cwd")
    if cwd is not None and not isinstance(cwd, str):
        raise ValueError("cwd must be a string when provided")

    stdin_data = payload.get("stdin", "")
    if not isinstance(stdin_data, str):
        raise ValueError("stdin must be a string when provided")

    debug = bool(payload.get("debug", False))

    return {
        "runtime": runtime,
        "module": module,
        "args": args,
        "env": env,
        "cwd": cwd,
        "stdin": stdin_data,
        "debug": debug,
    }


def main() -> None:
    raw_input = sys.stdin.read()
    try:
        request = parse_input(raw_input)
    except json.JSONDecodeError as exc:
        emit({"ok": False, "error": {"code": "INVALID_INPUT_JSON", "message": str(exc)}}, EXIT_INVALID_INPUT_JSON)
    except ValueError as exc:
        emit({"ok": False, "error": {"code": "INVALID_INPUT", "message": str(exc)}}, EXIT_INVALID_INPUT_JSON)

    debug = request["debug"]

    runtime, runtime_error = resolve_runtime(request["runtime"])
    if runtime is None:
        code = "UNSUPPORTED_RUNTIME" if request["runtime"] not in ("auto", "wasmtime", "wasmer") else "RUNTIME_NOT_FOUND"
        exit_code = EXIT_UNSUPPORTED_RUNTIME if code == "UNSUPPORTED_RUNTIME" else EXIT_RUNTIME_NOT_FOUND
        emit({"ok": False, "error": {"code": code, "message": runtime_error}}, exit_code)

    module = request["module"]
    if not os.path.isfile(module):
        emit(
            {"ok": False, "error": {"code": "MODULE_NOT_FOUND", "message": f"module not found: {module}"}},
            EXIT_MODULE_NOT_FOUND,
        )

    cmd = runtime_command(runtime, module, request["args"])
    debug_log(debug, f"runtime={runtime}")
    debug_log(debug, f"command={' '.join(cmd)}")

    child_env = os.environ.copy()
    child_env.update(request["env"])

    try:
        completed = subprocess.run(
            cmd,
            input=request["stdin"],
            text=True,
            capture_output=True,
            cwd=request["cwd"] or None,
            env=child_env,
        )
    except Exception as exc:
        emit({"ok": False, "error": {"code": "RUNTIME_EXEC_ERROR", "message": str(exc)}}, EXIT_INTERNAL_ERROR)

    if completed.stderr:
        sys.stderr.write(completed.stderr)

    stdout_raw = completed.stdout.strip()

    parsed_stdout: Optional[Any] = None
    if stdout_raw:
        try:
            parsed_stdout = json.loads(stdout_raw)
        except json.JSONDecodeError as exc:
            emit(
                {
                    "ok": False,
                    "runtime": runtime,
                    "error": {
                        "code": "STDOUT_NOT_JSON",
                        "message": f"runtime stdout was not valid JSON: {exc}",
                    },
                    "raw_stdout": completed.stdout,
                    "runtime_exit_code": completed.returncode,
                },
                EXIT_STDOUT_NOT_JSON,
            )

    if completed.returncode != 0:
        emit(
            {
                "ok": False,
                "runtime": runtime,
                "error": {"code": "RUNTIME_COMMAND_FAILED", "message": "wasm runtime returned a non-zero exit code"},
                "runtime_exit_code": completed.returncode,
                "result": parsed_stdout,
            },
            EXIT_RUNTIME_COMMAND_FAILED,
        )

    emit(
        {
            "ok": True,
            "runtime": runtime,
            "runtime_exit_code": completed.returncode,
            "result": parsed_stdout,
        },
        EXIT_OK,
    )


if __name__ == "__main__":
    main()
