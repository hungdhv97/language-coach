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
    ├── docker-compose.dev.yml   (Development environment with port mappings)
    └── docker-compose.prod.yml  (Production environment with port mappings)
