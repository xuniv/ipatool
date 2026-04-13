# Embedded Host (Minimum Reference)

이 예제는 임베디드 호스트의 최소 기능 레퍼런스입니다.

포함 내용:

1. `capability_manifest.json`을 읽어 network/storage/secret 권한을 제한.
2. 실패 응답의 `stderr`를 표준 에러(`StandardError`)로 역직렬화.
3. 왕복 지연(`roundtrip-latency`), 타임아웃(`timeout`), 재시도(`retries`)를 호스트 옵션으로 노출.

실행:

```bash
go run ./examples/embedded-host \
  --manifest ./examples/embedded-host/capability_manifest.json \
  --roundtrip-latency=100ms \
  --timeout=1s \
  --retries=1
```
