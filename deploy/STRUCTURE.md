.
├── docker/
│   ├── backend/
│   │   └── Dockerfile
│   ├── frontend/
│   │   └── Dockerfile
│   └── base/
│       └── Dockerfile
├── env/
│   ├── .global.env
│   ├── dev/
│   │   ├── backend.env
│   │   ├── frontend.env
│   │   ├── grafana.env
│   │   └── prometheus.env
│   ├── staging/
│   │   ├── backend.env
│   │   ├── frontend.env
│   │   ├── grafana.env
│   │   └── prometheus.env
│   └── prod/
│       ├── backend.env
│       ├── frontend.env
│       ├── grafana.env
│       └── prometheus.env
└── compose/
    └── docker-compose.yml   (UseControl the env using profiles + env_file separated by env)
