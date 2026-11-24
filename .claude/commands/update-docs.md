---
description: 저장소 문서를 현재 작업 기준으로 갱신하고 변경 사항 커밋 (각 장 마지막에 실행)
---

# /update-docs

각 장 마지막에 실행한다. 저장소 문서의 **신규 추가**와 **기존 내용 변경**을 모두 감지하고, 그 시점까지의 실제 작업 기준으로 문서를 갱신한 뒤 커밋한다.
특정 파일 목록에 의존하지 않고 아래 탐색 규칙으로 문서를 찾으므로, 이후 장에서 문서가 추가되어도 이 파일을 수정할 필요가 없다.

추가 지시: $ARGUMENTS

## 실행 절차

### 1. 이번 장 작업에서 변경된 내용을 파악한다

- **기준점(BASE)**: 커밋 메시지가 `docs: update-docs`로 시작하는 가장 최근 커밋.
  `git log --grep='^docs: update-docs' -n 1 --format=%H` — 없으면 최초 커밋(`git rev-list --max-parents=0 HEAD`).
- **커밋 기록**: `git log --oneline BASE..HEAD`
- **변경 파일**: `git diff --name-status BASE..HEAD` + `git status --porcelain` (커밋 안 된 작업 포함)
- **도입된 도구·설정**: 변경 파일(매니페스트, Helm values, 워크플로, 앱 코드)과 클러스터 상태
- **결정 사항**: 이번 장에서 독자가 선택한 도구, 구조, 전략 (대화 맥락 포함)
- **규칙·컨텍스트 변경**: CLAUDE.md 행동 규칙, 설정 파일 등

### 2. 문서를 탐색하고 분류한다

**탐색 규칙** (파일 목록을 하드코딩하지 않는다):
- 저장소 루트의 `*.md`
- `docs/`, `claude-context/`, `command-guardrails/` 아래 모든 파일
- 그 밖의 경로라도 이번 장에서 새로 생긴 `*.md` (`.claude/commands/`, `.github/` 아래 템플릿은 제외)

`git diff --name-status BASE..HEAD`와 `git status --porcelain`으로 각 문서를 분류한다:

| 분류 | 처리 |
|------|------|
| **신규 문서** (A, ??) | 내용이 현재 상태와 맞는지 검증하고, 다른 문서(README 구조도, CLAUDE.md 등)에서 참조·소개되도록 연결한다 |
| **변경된 문서** (M) | 변경 내용이 다른 문서와 모순되지 않는지 확인하고, 연관 문서에 전파한다 |
| **변경 없는 문서** | 1단계에서 파악한 작업 내용 때문에 낡은 부분이 있으면 갱신한다 |

### 3. 문서를 파악한 내용에 맞춰 갱신한다

후보 (존재하는 것만 처리):

- `JOURNEY.md` — 진행 현황, 도구 선택, 현재 버전, 리소스 상태
  **진행 현황**: 이번 장에서 완료한 서브챕터를 ⬜ → ✅로 변경하고 완료일(YYYY-MM-DD)을 기록한다. 비워두거나 건너뛰지 않는다.
  **도구 선택 기록**: 이번 장의 각 [필수입력] 실행 단계가 만든 결정을 한 줄씩 추가한다. 이 테이블이 곧 ADR의 source이므로, 누락하면 ADR 번호가 어긋난다. 다음 분배가 보장되도록 모든 결정을 빠짐없이 기록한다:
    - ch3 → ADR-001~002 (ArgoCD, GitHub Actions)
    - ch4 → ADR-003~005 (Prometheus+Grafana, Loki+Fluent Bit, PrometheusRule+Alertmanager)
    - ch5 → ADR-006~007 (Gateway API, Blue/Green)
    - ch6 → ADR-008~010 (Valkey, Secret Manager CSI, Canary 전환)
    - **ch7 → ADR-011~013 (노드풀 분리, App of Apps, 멀티테넌시)** — 멀티테넌시는 "도구 비교"가 아닌 "구조 결정"이므로 빠뜨리기 쉬움
    - **ch8 → ADR-014~016 (Kafka, Tempo, CronJob)** — CronJob은 "스케줄링 결정"이므로 빠뜨리기 쉬움
  **현재 버전**: 이번 장에서 설치하거나 버전이 바뀐 컴포넌트는 추측하지 않고 클러스터에서 직접 조회해 채운다. 비워두지 않는다.
  - Notiflex 이미지 (5장 이전은 Deployment, 이후는 Rollout — 있는 쪽을 조회):
    `kubectl --context gke-sysnet4admin_book_gitaiops get rollout notiflex-api -n notiflex -o jsonpath='{.spec.template.spec.containers[0].image}' 2>/dev/null || kubectl --context gke-sysnet4admin_book_gitaiops get deploy notiflex-api -n notiflex -o jsonpath='{.spec.template.spec.containers[0].image}'`
  - ArgoCD: `kubectl --context gke-sysnet4admin_book_gitaiops get deploy argocd-server -n argocd -o jsonpath='{.spec.template.spec.containers[0].image}'`
  - Kafka: `kubectl --context gke-sysnet4admin_book_gitaiops get kafka -n kafka -o jsonpath='{.items[0].spec.kafka.version}'`
  - 그 외 도구: `kubectl --context gke-sysnet4admin_book_gitaiops get pod -n <namespace> -l <label> -o jsonpath='{.items[0].spec.containers[0].image}'`
  **현재 리소스**: 노드풀이 추가·변경됐으면 `kubectl --context gke-sysnet4admin_book_gitaiops get nodes -L cloud.google.com/gke-nodepool`로 조회해 테이블을 갱신한다.
  **트러블슈팅 이력**: 이번 장에서 가드레일에 없던 문제를 겪었으면 추가한다.

- `CLAUDE.md` — 규칙이나 컨텍스트가 바뀐 경우 (디렉터리 구조, GCP/클러스터 정보, 행동 규칙)

- `README.md` — 디렉터리 구조, 기술 스택이 바뀐 경우. 신규 문서·디렉터리가 생기면 구조도에 추가한다.

- `docs/architecture-decisions.md` — 결정 사항 누적 (도입: 5장)
  **파일이 존재하는 경우에만**: `JOURNEY.md`의 도구 선택 기록 테이블에서 직전 장에 새로 추가된 결정을 **모두** 찾아 ADR 항목으로 추가한다. 같은 장에 여러 결정이 있으면 전부 포함한다 (예: ch7에 nodeSelector + App of Apps가 있으면 둘 다 ADR로 변환). 번호는 기존 마지막 ADR 번호 +1부터 순서대로 매긴다. 파일이 없으면 건너뛴다 (신설은 ch5 가드레일이 담당).
  **ADR 형식 준수 필수**: 각 ADR은 반드시 아래 형식으로 작성한다. `**이유**:` 뒤에 `-` 불릿 리스트로 3~4가지 이유를 작성한다. 이유를 단일 줄이나 쉼표로 압축하지 않는다.
  ```
  ## ADR-NNN: 짧은 결정 제목 (N장)
  **시점**: YYYY-MM / **결정**: 무엇을 채택하고 무엇을 안 쓰는지 한 문장
  **이유**:
  - 이유 1 (가장 결정적인 근거)
  - 이유 2
  - 이유 3
  - 이유 4 (필요 시)
  ```
  **ADR 삽입 위치 (필수)**: 파일 제목(`# Architecture Decision Records`)을 anchor로 쓰면 새 ADR이 파일 상단에 삽입된다. 반드시 아래 순서를 따른다:
  1. `tail -5 docs/architecture-decisions.md` 로 파일 끝을 먼저 확인한다
  2. 마지막 줄(이전 ADR의 마지막 이유 불릿)을 Edit의 `old_string` anchor로 사용한다
  3. 추가 후 `grep "^## ADR-" docs/architecture-decisions.md` 로 번호 오름차순 검증한다 — 순서가 깨져 있으면 파일을 재정렬한 뒤 커밋한다

- `claude-context/` — 현재 아키텍처 스냅샷 (도입: 6장)
  **파일이 존재하는 경우**: 이번 장에서 아키텍처가 변경됐으면 반드시 갱신한다. 갱신 대상 예시: 노드풀 추가(ch7.2), 새 네임스페이스(ch7.4), Kafka/Tempo/CronJob 추가(ch8). 파일 맨 위의 "N장 완료 시점" 텍스트도 현재 장으로 업데이트한다. 건너뛰면 AI가 오래된 아키텍처를 참조하게 된다.

- `command-guardrails/` — 위험 명령 절차 (도입: 8장)
  **파일이 존재하는 경우**: 이번 장에서 새로 도입된 위험 작업(삭제, 수동 실행, 데이터 손실 가능 작업)이 있으면 절차를 추가하고, 기존 절차의 명령·리소스 이름이 현재와 다르면 고친다.

- **그 밖에 2단계에서 찾은 문서** — 위 목록에 없는 신규 문서도 같은 원칙(현재 실제 상태와 일치, 다른 문서와 모순 없음)으로 갱신한다.

### 4. 검증한다

- 문서에 적힌 버전·이미지 태그·리소스 이름이 클러스터/저장소 실제 값과 일치하는지 확인한다.
- 시크릿(토큰, 키, 비밀번호)이 문서에 들어가지 않았는지 확인한다: `git diff | grep -i -E 'BEGIN .*PRIVATE|private_key|ghp_|gho_|password\s*[:=]|secret\s*[:=]'`
- ADR 파일이 있으면 번호 오름차순을 다시 확인한다.

### 5. 커밋한다

새로 추가된 문서와 기존 문서의 내용 변경을 모두 반영한 뒤 Git에 커밋한다.

- 메시지 형식: `docs: update-docs — N장 <이번 장 요약>` (다음 실행 시 기준점이 되므로 접두어 `docs: update-docs`를 반드시 지킨다)
- 본문: 신규 문서 / 변경 문서 / 갱신 이유를 불릿으로 정리
- 문서 외의 커밋 안 된 변경(코드, 매니페스트)이 있으면 함께 커밋하지 말고 독자에게 알린다.
- 커밋 후 `git push origin main`으로 푸시한다.

### 6. 결과를 보고한다

| 문서 | 분류 (신규/변경/갱신) | 갱신 내용 |
|------|----------------------|----------|

존재하지 않는 문서는 건너뛴다. 이후 장에서 새 문서가 추가되면 스킬은 수정 없이 그대로 그 문서까지 포함해 갱신한다.
