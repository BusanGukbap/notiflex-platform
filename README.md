# Notiflex Platform

Notiflex는 B2B 알림 SaaS 플랫폼입니다. 고객사(테넌트)가 API로 알림을 요청하면 이메일·SMS·푸시 등으로 발송합니다.

이 저장소는 Notiflex의 애플리케이션 코드, 쿠버네티스 매니페스트, CI 파이프라인을 **GitOps** 방식으로 관리합니다.

## 구조

```
notiflex-platform/
├── CLAUDE.md            # Claude Code용 프로젝트 컨텍스트와 행동 규칙
├── README.md
├── app/                 # Go 애플리케이션 (notiflex-api)
├── k8s/
│   └── smb/             # 쿠버네티스 매니페스트
└── .github/
    └── workflows/       # CI 파이프라인 (GitHub Actions)
```

## 기술 스택

| 영역 | 사용 기술 |
|------|-----------|
| 애플리케이션 | Go 표준 라이브러리 |
| 컨테이너 | scratch 베이스 이미지 |
| 이미지 저장소 | Google Artifact Registry (`asia-northeast3`) |
| 인프라 | GKE Standard (Zonal), Spot VM |
| GitOps | ArgoCD |
| CI | GitHub Actions |

## 참고

「AI 시대에 개발자가 알아야 하는 인프라 구성 배포 with 클로드 코드」 실습 저장소입니다.
