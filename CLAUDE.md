# Notiflex Platform

## 프로젝트 개요

Notiflex — B2B 알림 SaaS 플랫폼. 고객사(테넌트)가 API로 알림을 요청하면 이메일·SMS·푸시 등으로 발송한다.
이 저장소는 Notiflex의 애플리케이션 코드, 쿠버네티스 매니페스트, CI 파이프라인을 GitOps 방식으로 관리한다.

## 기술 스택

- **언어**: Go 표준 라이브러리 (외부 프레임워크 없음)
- **컨테이너**: scratch 베이스 이미지 (멀티 스테이지 빌드)
- **인프라**: GKE Standard (Zonal), Spot VM
- **GitOps**: ArgoCD (3장부터)
- **CI**: GitHub Actions (3장부터)

## GCP 설정

| 항목 | 값 |
|------|-----|
| 프로젝트 ID | `project-c758a5d0-9fc9-4cc0-b0f` |
| 프로젝트 이름 | `my-gitaiops-project` |
| 리전 | `asia-northeast3` (서울) |
| 존 | `asia-northeast3-a` |
| GKE 클러스터 | `notiflex-cluster` |
| kubectl 컨텍스트 | `gke-sysnet4admin_book_gitaiops` |
| Artifact Registry | `asia-northeast3-docker.pkg.dev/project-c758a5d0-9fc9-4cc0-b0f/notiflex` |

> 책의 예시에 나오는 `my-gitaiops-project`는 프로젝트 **이름**이다. gcloud 명령과 이미지 경로에는 위의 **프로젝트 ID**를 사용한다.

## 디렉터리 구조

```
notiflex-platform/
├── CLAUDE.md
├── app/                 # Go 애플리케이션 (notiflex-api)
├── k8s/
│   └── smb/             # K8s 매니페스트
└── .github/
    └── workflows/       # CI 파이프라인
```

## 행동 규칙

1. **항상 확인 후 실행한다.** 클러스터·GCP 리소스를 생성·변경·삭제하는 명령은 무엇을 바꾸는지 먼저 설명하고 실행한다.
2. **변경 전 현재 상태를 확인한다.** `kubectl get`, `gcloud ... describe/list`로 현재 상태를 본 뒤 변경한다.
3. **모든 kubectl 명령에 `--context gke-sysnet4admin_book_gitaiops`를 지정한다.** 현재 컨텍스트가 다른 클러스터일 수 있다.
4. **매니페스트는 Git이 원본이다.** 클러스터에 직접 적용한 변경도 반드시 `k8s/` 아래 파일에 반영한다.
5. **시크릿을 커밋하지 않는다.** 토큰·키·비밀번호는 저장소에 넣지 않는다.

## 실행 환경 메모 (Windows)

- Git Bash에서 `gcloud`, `gh`는 `~/bin`의 래퍼 스크립트로 실행된다 (gcloud는 번들 Python을 쓰는 `gcloud.cmd`로 연결).
